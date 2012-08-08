package bot

import (
	"flag"
	"fmt"
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"strings"
)

var rebuilder *string = flag.String("rebuilder", "",
		"Nick[:password] to accept rebuild command from.")

// The bot is called sp0rkle...
const botName string = "sp0rkle"

type Sp0rkle struct {
	// Embed a connection to IRC.
	Conn *client.Conn

	// And an event registry / dispatcher
	ED event.EventDispatcher
	ER event.EventRegistry

	// And a plugin manager.
	PM base.PluginManager

	// And a logger.
	l logging.Logger

	// it's got a bunch of drivers that register event handlers
	drivers map[string]base.Driver

	// channel to join on start up
	channels []string

	// nick and password for rebuild command
	rbnick, rbpw string

	// and we need to kill it occasionally.
	reexec, quit bool
	Quit chan bool
}

func Bot(c *client.Conn, pm base.PluginManager, l logging.Logger) *Sp0rkle {
	s := strings.Split(*rebuilder, ":")
	bot := &Sp0rkle{
		Conn:     c,
		ER:       c.ER,
		ED:       c.ED,
		PM:       pm,
		l:        l,
		drivers:  make(map[string]base.Driver),
		channels: make([]string, 0, 1),
		rbnick:   s[0],
		Quit:     make(chan bool),
	}
	if len(s) > 1 {
		bot.rbpw = s[1]
	}
	c.State = bot
	return bot
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
			pp.RegisterPlugins(bot.PM)
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
	bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(fm, args...))
}
