package bot

import (
	"flag"
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/util"
	"os/exec"
	"strings"
)

var (
	channels *string = flag.String("channels", "#sp0rklf",
		"Comma-separated list of channels to join, defaults to '#sp0rklf'")
	rebuilder *string = flag.String("rebuilder", "",
		"Nick[:password] to accept rebuild command from.")
	oper *string = flag.String("oper", "",
		"user:password for server OPER command on connect.")
	vhost *string = flag.String("vhost", "",
		"user:password for server VHOST command on connect.")
)

func bot_connected(line *base.Line) {
	for _, c := range strings.Split(*channels, ",") {
		logging.Info("Joining %s on startup.\n", c)
		irc.Join(c)
	}
	if *oper != "" {
		up := strings.SplitN(*oper, ":", 2)
		if len(up) == 2 {
			irc.Oper(up[0], up[1])
		}
	}
	if *vhost != "" {
		up := strings.SplitN(*vhost, ":", 2)
		if len(up) == 2 {
			irc.Raw(fmt.Sprintf("VHOST %s %s", up[0], up[1]))
		}
	}
}

func bot_disconnected(line *base.Line) {
	// The read from this channel is in connectLoop
	disconnected <- true
	logging.Info("Disconnected...")
}

func bot_command(l *base.Line) {
	// This is a dirty hack to treat factoid additions as a special
	// case, since they may begin with command string prefixes.
	if util.IsFactoidAddition(l.Args[1]) {
		return
	}
	if cmd, ln := commands.Match(l.Args[1]); l.Addressed && cmd != nil {
		// Cut command off, trim and compress spaces.
		l.Args[1] = strings.Join(strings.Fields(l.Args[1][ln:]), " ")
		cmd.Execute(l)
	}
}

func bot_rebuild(line *base.Line) {
	if !check_rebuilder("rebuild", line) { return }

	// Ok, we should be good to rebuild now.
	irc.Notice(line.Nick, "Beginning rebuild")
	cmd := exec.Command("go", "get", "-u", "github.com/fluffle/sp0rkle")
	out, err := cmd.CombinedOutput()
	if err != nil {
		irc.Notice(line.Nick, fmt.Sprintf("Rebuild failed: %s", err))
		for _, l := range strings.Split(string(out), "\n") {
			irc.Notice(line.Nick, l)
		}
		return
	}
	shutdown = true
	reexec = true
	irc.Quit("Restarting with new build.")
}

func bot_shutdown(line *base.Line) {
	if check_rebuilder("shutdown", line) {
		shutdown = true
		irc.Quit("Shutting down.")
	}
}

func check_rebuilder(cmd string, line *base.Line) bool {
	s := strings.Split(*rebuilder, ":")
	if s[0] == "" || s[0] != line.Nick || !strings.HasPrefix(line.Args[1], cmd) {
		return false
	}
	if len(s) > 1 && line.Args[1] != fmt.Sprintf("%s %s", cmd, s[1]) {
		return false
	}
	return true
}

func bot_help(line *base.Line) {
	if cmd, _ := commands.Match(line.Args[1]); cmd != nil {
		ReplyN(line, "%s", cmd.Help())
	} else if len(line.Args[1]) == 0 {
		ReplyN(line, "https://github.com/fluffle/sp0rkle/wiki "+
			"-- pull requests welcome ;-)")
	} else {
		ReplyN(line, "Unrecognised command '%s'.", line.Args[1])
	}
}

func bot_ignore(line *base.Line) {
	nick := strings.ToLower(strings.Fields(line.Args[1])[0])
	if nick == "" {
		return
	}
	ignores.String(nick, "ignore")
	ReplyN(line, "I'll ignore '%s'.", nick)
}

func bot_unignore(line *base.Line) {
	nick := strings.ToLower(strings.Fields(line.Args[1])[0])
	if nick == "" {
		return
	}
	ignores.Delete(nick)
	ReplyN(line, "No longer ignoring '%s'.", nick)
}
