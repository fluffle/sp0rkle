package bot

import (
	"flag"
	"fmt"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/collections/conf"
	"github.com/fluffle/sp0rkle/util"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	nick *string = flag.String("nick", "sp0rklf",
		"Name of bot, defaults to 'sp0rklf'")
	server   *string = flag.String("server", "", "IRC server to connect to.")
	ssl      *bool   = flag.Bool("ssl", false, "Use SSL when connecting.")
	httpHost *string = flag.String("http_host", "http://sp0rk.ly",
		"Hostname for HTTP paths served by bot.")
)

// These package globals are an experiment. They have radically simplified some
// of the code in the drivers, and so i'm reserving judgement on the usual
// knee-jerk EWW reaction for the moment. Please don't hate me.
var irc *client.Conn
var ignores conf.Namespace
var lock sync.Mutex

func Init() {
	lock.Lock()
	defer lock.Unlock()
	if irc != nil {
		return
	}

	if *server == "" {
		// Don't call logging.Fatal as we don't want a backtrace in this case
		logging.Error("--server option required. \nOptions are:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Configure IRC client
	cfg := client.NewConfig(*nick, "boing", "slowly becoming sp0rkle")
	cfg.Flood = true
	cfg.SSL = *ssl
	cfg.Server = *server
	cfg.Recover = func (conn *client.Conn, line *client.Line) {
		if err := recover(); err != nil {
			_, f, l, _ := runtime.Caller(2)
			i := strings.Index(f, "sp0rkle/")
			logging.Error("panic at %s:%d: %v", f[i:], l, err)
			conn.Privmsg(line.Target(), fmt.Sprintf(
				"panic at %s:%d: %v", f[i:], l, err))
		}
	}

	irc, _ = client.Client(cfg)

	HandleFunc(bot_connected, client.CONNECTED)
	HandleFunc(bot_disconnected, client.DISCONNECTED)

	// This is a special handler that dispatches commands from the command set
	HandleFunc(bot_command, client.PRIVMSG)
	// This is a special handler that triggers a rebuild and re-exec
	HandleFunc(bot_rebuild, client.NOTICE)
	// This is a special handler that triggers a shutdown and disconnect
	HandleFunc(bot_shutdown, client.NOTICE)

	CommandFunc(bot_help, "help", "If you need to ask, you're beyond help.")

	// Ignores contains a list of Nicks to ignore.
	ignores = conf.Ns("ignore")
	CommandFunc(bot_ignore, "ignore", "ignore <nick>  -- "+
		"make the bot ignore <nick> completely.")
	CommandFunc(bot_unignore, "unignore", "unignore <nick>  -- "+
		"make the bot unignore <nick> again.")
}

func Connect() bool {
	lock.Lock()
	defer lock.Unlock()
	if irc == nil {
		logging.Fatal("Called Connect() before Init().")
	}
	return connectLoop()
}

var shutdown, reexec bool
var disconnected = make(chan bool)

func connectLoop() bool {
	var retries uint32
	for {
		if err := irc.Connect(); err != nil {
			logging.Error("Connection error: %s", err)
			retries++
			if retries > 10 {
				logging.Error("Giving up connection after 10 failed retries.")
				return false
			}
			<-time.After(time.Second * 1 << retries)
		} else {
			retries, shutdown, reexec = 0, false, false
			// Wait here for a signal from bot_disconnected
			<-disconnected
			if shutdown {
				return reexec
			}
		}
	}
	panic("unreachable")
}

func Handle(h base.Handler, events ...string) {
	// TODO(fluffle): rework this properly.
	for _, event := range events {
		irc.HandleFunc(event, func(_ *client.Conn, l *client.Line) {
			if ignores.String(strings.ToLower(l.Nick)) == "" {
				line := transformLine(l)
				h.Execute(line)
			}
		})
	}
}

func HandleFunc(fn base.HandlerFunc, event ...string) {
	Handle(fn, event...)
}

var commands = base.NewCommandSet()

func Command(cmd base.Command, prefix string) {
	commands.Add(cmd, prefix)
}

func CommandFunc(fn base.HandlerFunc, prefix, help string) {
	Command(base.NewCommand(fn, help), prefix)
}

var plugins = base.NewPluginSet()

func Plugin(p base.Plugin) {
	plugins.Add(p)
}

func PluginFunc(fn base.PluginFunc) {
	Plugin(fn)
}

func transformLine(line *client.Line) *base.Line {
	// We want line.Args[1] to contain the (possibly) stripped version of itself
	// but modifying the pointer will result in other goroutines seeing the
	// change, so we need to copy line for our own edification.
	nl := &base.Line{Line: line.Copy()}
	if nl.Cmd != client.PRIVMSG {
		return nl
	}
	nl.Args[1], nl.Addressed = util.RemovePrefixedNick(
		strings.TrimSpace(line.Args[1]), Nick())
	// If we're being talked to in private, line.Args[0] will contain our Nick.
	// To ensure the replies go to the right place (without performing this
	// check everywhere) test for this and set line.Args[0] == line.Nick.
	// We should consider this as "addressing" us too, and set Addressed = true
	if nl.Args[0] == Nick() {
		nl.Args[0] = nl.Nick
		nl.Addressed = true
	}
	return nl
}

// ReplyN() adds a prefix of "nick: " to the reply text,
func ReplyN(line *base.Line, fm string, args ...interface{}) {
	args = append([]interface{}{line.Nick}, args...)
	Reply(line, "%s: "+fm, args...)
}

// whereas Reply() does not.
func Reply(line *base.Line, fm string, args ...interface{}) {
	Privmsg(line.Target(), plugins.Apply(fmt.Sprintf(fm, args...), line))
}

func Do(line *base.Line, fm string, args ...interface{}) {
	Action(line.Target(), plugins.Apply(fmt.Sprintf(fm, args...), line))
}

// Hmmm. Fix these later.
func Privmsg(ch, text string) {
	irc.Privmsg(ch, text)
}

func Action(ch, text string) {
	irc.Action(ch, text)
}

func Nick() string {
	return irc.Me().Nick
}

func HttpHost() string {
	return *httpHost
}
