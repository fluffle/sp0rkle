package bot

import (
	"flag"
	"fmt"
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/lib/util"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"strings"
)

var (
	rebuilder *string = flag.String("rebuilder", "",
		"Nick[:password] to accept rebuild command from.")
	prefix *string = flag.String("http_prefix", "http://sp0rk.ly",
		"Prefix for HTTP paths served by bot.")
)

// The bot is called sp0rkle...
const botName string = "sp0rkle"

type Sp0rkle struct {
	// Embed a connection to IRC.
	Conn *client.Conn

	// And an event registry / dispatcher
	ED event.EventDispatcher
	ER event.EventRegistry

	// And a logger.
	l logging.Logger

	// it's got a bunch of drivers that register event handlers
	drivers map[string]base.Driver

	// channel to join on start up
	channels []string

	// nick and password for rebuild command
	rbnick, rbpw string

	// prefix for HTTP paths served
	Prefix string

	// and we need to kill it occasionally.
	reexec, quit bool
	Quit chan bool
}

var bot *Sp0rkle
var irc *client.Conn

func Init(c *client.Conn, l logging.Logger) *Sp0rkle {
	// TODO(fluffle): fix race.
	if bot == nil {
		bot = Bot(c, l)
		irc = c
	}

	HandleFunc(bot_connected, "connected")
	HandleFunc(bot_disconnected, "disconnected")

	// This is a special handler that dispatches commands from the command set
	HandleFunc(bot_command, "privmsg")
	// This is a special handler that triggers a rebuild and re-exec
	HandleFunc(bot_rebuild, "notice")
	// This is a special handler that triggers a shutdown and disconnect
	HandleFunc(bot_shutdown, "notice")

	CommandFunc(bot_help, "help", "If you need to ask, you're beyond help.")
	return bot
}

func Bot(c *client.Conn, l logging.Logger) *Sp0rkle {
	s := strings.Split(*rebuilder, ":")
	bot := &Sp0rkle{
		Conn:     c,
		ER:       c.ER,
		ED:       c.ED,
		l:        l,
		drivers:  make(map[string]base.Driver),
		channels: make([]string, 0, 1),
		rbnick:   s[0],
		Prefix:   *prefix,
		Quit:     make(chan bool),
	}
	if len(s) > 1 {
		bot.rbpw = s[1]
	}
	c.State = bot
	return bot
}

type botFn func(*base.Line)

type botPl func(string, *base.Line) string

type botCommand struct {
	fn botFn
	help string
}

func (bf botFn) Execute(line *base.Line) {
	bf(line)
}

func (bp botPl) Apply(in string, line *base.Line) string {
	return bp(in, line)
}

func (bc *botCommand) Execute(line *base.Line) {
	bc.fn(line)
}

func (bc *botCommand) Help() string {
	return bc.help
}

func Handle(h base.Handler, event ...string) {
	bot.ER.AddHandler(client.NewHandler(func(_ *client.Conn, l *client.Line) {
		h.Execute(Line(l))
	}), event...)
}

func HandleFunc(fn botFn, event ...string) {
	Handle(fn, event...)
}

var commands = base.NewCommandSet()

func Command(cmd base.Command, prefix string) {
	commands.Add(cmd, prefix)
}

func CommandFunc(fn botFn, prefix, help string) {
	Command(&botCommand{fn, help}, prefix)
}

var plugins = base.NewPluginSet()

func Plugin(p base.Plugin) {
	plugins.Add(p)
}

func PluginFunc(fn botPl) {
	Plugin(fn)
}

func Line(line *client.Line) *base.Line {
	// We want line.Args[1] to contain the (possibly) stripped version of itself
	// but modifying the pointer will result in other goroutines seeing the
	// change, so we need to copy line for our own edification.
	nl := &base.Line{Line: line.Copy()}
	if nl.Cmd != "PRIVMSG" {
		return nl
	}
	nl.Args[1], nl.Addressed = util.RemovePrefixedNick(
		strings.TrimSpace(line.Args[1]), irc.Me.Nick)
	// If we're being talked to in private, line.Args[0] will contain our Nick.
	// To ensure the replies go to the right place (without performing this
	// check everywhere) test for this and set line.Args[0] == line.Nick.
	// We should consider this as "addressing" us too, and set Addressed = true
	if nl.Args[0] == irc.Me.Nick {
		nl.Args[0] = nl.Nick
		nl.Addressed = true
	}
	return nl
}

func (bot *Sp0rkle) Name() string {
	return botName
}

func (bot *Sp0rkle) RegisterAll() {
	for _, d := range bot.drivers {
		// Register the driver's event handlers with the event registry.
		d.RegisterHandlers(bot.ER)
		// If the driver provides FactoidPlugins to change factoid output
		// register them with the PluginManager here too.
		if pp, ok := d.(base.PluginProvider); ok {
			pp.RegisterPlugins(plugins)
		}
		// If the driver wants to handle any HTTP paths, register them.
		if hp, ok := d.(base.HttpProvider); ok {
			hp.RegisterHttpHandlers()
		}
	}
}

func (bot *Sp0rkle) Dispatch(name string, ev ...interface{}) {
	// Shim the bot into the parameter list of every event dispatched via it.
	ev = append([]interface{}{bot}, ev...)
	bot.ED.Dispatch(name, ev...)
}

func (bot *Sp0rkle) AddDriver(d base.Driver) {
	bot.drivers[d.Name()] = d
}

func (bot *Sp0rkle) GetDriver(name string) base.Driver {
	// Callers will have to unbox the returned driver themselves
	return bot.drivers[name]
}

func (bot *Sp0rkle) AddChannels(c []string) {
	bot.channels = append(bot.channels, c...)
}

func (bot *Sp0rkle) ReExec() bool {
	return bot.reexec
}

// Currently makes the assumption that we're replying to line.Args[0] in every
// instance. While this is normally the case, it may not be in some cases...
// ReplyN() adds a prefix of "nick: " to the reply text,
func (bot *Sp0rkle) ReplyN(line *base.Line, fm string, args ...interface{}) {
	args = append([]interface{}{line.Nick}, args...)
	bot.Reply(line, "%s: "+fm, args...)
}

// whereas Reply() does not.
func (bot *Sp0rkle) Reply(line *base.Line, fm string, args ...interface{}) {
	bot.Conn.Privmsg(line.Args[0], plugins.Apply(fmt.Sprintf(fm, args...), line))
}

func (bot *Sp0rkle) Do(line *base.Line, fm string, args ...interface{}) {
	bot.Conn.Action(line.Args[0], plugins.Apply(fmt.Sprintf(fm, args...), line))
}

// Currently makes the assumption that we're replying to line.Args[0] in every
// instance. While this is normally the case, it may not be in some cases...
// ReplyN() adds a prefix of "nick: " to the reply text,
func ReplyN(line *base.Line, fm string, args ...interface{}) {
	args = append([]interface{}{line.Nick}, args...)
	Reply(line, "%s: "+fm, args...)
}

// whereas Reply() does not.
func Reply(line *base.Line, fm string, args ...interface{}) {
	Privmsg(line.Args[0], plugins.Apply(fmt.Sprintf(fm, args...), line))
}

func Do(line *base.Line, fm string, args ...interface{}) {
	Action(line.Args[0], plugins.Apply(fmt.Sprintf(fm, args...), line))
}

// Hmmm. Fix these later.
func Privmsg(ch, text string) {
	irc.Privmsg(ch, text)
}

func Action(ch, text string) {
	irc.Action(ch, text)
}

func Flood(f bool) {
	irc.Flood = f
}

func Nick() string {
	return irc.Me.Nick
}
