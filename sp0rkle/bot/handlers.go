package bot

import (
	"fmt"
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/lib/util"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"os/exec"
	"strings"
)

type BotHandler func(*Sp0rkle, *base.Line)

// NOTE: Nothing but the bot should register for IRC events!
func (bot *Sp0rkle) RegisterHandlers(r event.EventRegistry) {
	// Generic shim to wrap an irc event into a bot event.
	forward_event := func(name string) event.Handler {
		return client.NewHandler(func(irc *client.Conn, line *client.Line) {
			getState(irc).Dispatch("bot_"+name, &base.Line{Line: *line.Copy()})
		})
	}

	r.AddHandler(client.NewHandler(bot_connected), "connected")
	r.AddHandler(client.NewHandler(bot_disconnected), "disconnected")
	r.AddHandler(client.NewHandler(bot_privmsg), "privmsg")
	r.AddHandler(forward_event("action"), "action")
	// These are mostly for the seen plugin.
	r.AddHandler(forward_event("join"), "join")
	r.AddHandler(forward_event("part"), "part")
	r.AddHandler(forward_event("kick"), "kick")
	r.AddHandler(forward_event("quit"), "quit")
	r.AddHandler(forward_event("nick"), "nick")

	// This is a special handler that triggers a rebuild and re-exec
	r.AddHandler(client.NewHandler(bot_rebuild), "notice")
	// This is a special handler that triggers a shutdown and disconnect
	r.AddHandler(client.NewHandler(bot_shutdown), "notice")

	CmdFunc(bot_help, "help", "If you need to ask, you're beyond help.")
}

// Unboxer for bot handlers.
func NewHandler(f BotHandler) event.Handler {
	return event.NewHandler(func(ev ...interface{}) {
		f(ev[0].(*Sp0rkle), ev[1].(*base.Line))
	})
}

func bot_connected(irc *client.Conn, line *client.Line) {
	bot := getState(irc)
	for _, c := range bot.channels {
		bot.l.Info("Joining %s on startup.\n", c)
		irc.Join(c)
	}
	bot.Dispatch("bot_connected", &base.Line{Line: *line.Copy()})
}

func bot_disconnected(irc *client.Conn, line *client.Line) {
	bot := getState(irc)
	bot.Quit <- bot.quit
	bot.l.Info("Disconnected...")
}

// Do some standard processing on incoming lines and dispatch a bot_privmsg
func bot_privmsg(irc *client.Conn, line *client.Line) {
	bot := getState(irc)

	l, p := util.RemovePrefixedNick(strings.TrimSpace(line.Args[1]), irc.Me.Nick)
	// We want line.Args[1] to contain the (possibly) stripped version of itself
	// but modifying the pointer will result in other goroutines seeing the
	// change, so we need to copy line for our own edification.
	nl := &base.Line{Line: *line.Copy()}
	nl.Args[1] = l
	nl.Addressed = p
	// If we're being talked to in private, line.Args[0] will contain our Nick.
	// To ensure the replies go to the right place (without performing this
	// check everywhere) test for this and set line.Args[0] == line.Nick.
	// We should consider this as "addressing" us too, and set Addressed = true
	if nl.Args[0] == irc.Me.Nick {
		nl.Args[0] = nl.Nick
		nl.Addressed = true
	}
	bot.Dispatch("bot_privmsg", nl)
	if cmd := commandMatch(nl.Args[1]); nl.Addressed && cmd != nil {
		cmd.Execute(bot, nl)
	}
}

// Retrieve the bot from irc.State.
func getState(irc *client.Conn) *Sp0rkle {
	return irc.State.(*Sp0rkle)
}

func bot_rebuild(irc *client.Conn, line *client.Line) {
	bot := getState(irc)
	if bot.rbnick == "" || bot.rbnick != line.Nick { return }
	if !strings.HasPrefix(line.Args[1], "rebuild") { return }
	if bot.rbpw != "" && line.Args[1] != "rebuild "+bot.rbpw { return }

	// Ok, we should be good to rebuild now.
	irc.Notice(line.Nick, "Beginning rebuild")
	cmd := exec.Command("go", "get", "-u", "github.com/fluffle/sp0rkle/sp0rkle")
	out, err := cmd.CombinedOutput()
	if err != nil {
		irc.Notice(line.Nick, fmt.Sprintf("Rebuild failed: %s", err))
		for _, l := range strings.Split(string(out), "\n") {
			irc.Notice(line.Nick, l)
		}
		return
	}
	bot.quit = true
	bot.reexec = true
	bot.Conn.Quit("Restarting with new build.")
}

func bot_shutdown(irc *client.Conn, line *client.Line) {
	bot := getState(irc)
	if bot.rbnick == "" || bot.rbnick != line.Nick { return }
	if !strings.HasPrefix(line.Args[1], "shutdown") { return }
	if bot.rbpw != "" && line.Args[1] != "shutdown "+bot.rbpw { return }
	bot.quit = true
	bot.Conn.Quit("Shutting down.")
}

func bot_help(bot *Sp0rkle, line *base.Line) {
	s := strings.Join(strings.Fields(line.Args[1])[1:], " ")
	if cmd := commandMatch(s); cmd != nil {
		bot.ReplyN(line, cmd.Help())
	} else if len(s) == 0 {
		bot.ReplyN(line, "https://github.com/fluffle/sp0rkle/wiki " +
			"-- pull requests welcome ;-)")
	} else {
		bot.ReplyN(line, "Unrecognised command '%s'.", s)
	}
}
