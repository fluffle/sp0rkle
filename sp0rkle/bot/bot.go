package bot

import (
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
)

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

	// and we need to kill it occasionally.
	Quit chan bool
}

func Bot(c *client.Conn, pm base.PluginManager, l logging.Logger) *Sp0rkle {
	bot := &Sp0rkle{
		Conn:     c,
		ER:       c.ER,
		ED:       c.ED,
		PM:       pm,
		l:        l,
		drivers:  make(map[string]base.Driver),
		channels: make([]string, 0, 1),
		Quit:     make(chan bool),
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

func (bot *Sp0rkle) AddChannel(c string) {
	bot.channels = append(bot.channels, c)
}
