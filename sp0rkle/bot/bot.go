package bot

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/goirc/event"
	"sp0rkle/base"
)

// The bot is called sp0rkle...
const botName string = "sp0rkle"

type Sp0rkle struct {
	// Embed a connection to IRC.
	Conn *client.Conn

	// it's got a bunch of drivers that register event handlers
	drivers map[string]base.Driver

	// channel to join on start up
	channels []string

	// and we need to kill it occasionally.
	Quit chan bool
}

func Bot() *Sp0rkle {
	return &Sp0rkle{
		Conn:     nil,
		drivers:  make(map[string]base.Driver),
		channels: make([]string, 0, 1),
		Quit:     make(chan bool),
	}
}

func (bot *Sp0rkle) Name() string {
	return botName
}

func (bot *Sp0rkle) RegisterAll(r event.EventRegistry, pm base.PluginManager) {
	for _, d := range bot.drivers {
		// Register the driver's event handlers with the event registry.
		d.RegisterHandlers(r)
		// If the driver provides FactoidPlugins to change factoid output
		// register them with the PluginManager here too.
		if pp, ok := d.(base.PluginProvider); ok {
			pp.RegisterPlugins(pm)
		}
	}
}

func (bot *Sp0rkle) Dispatch(name string, ev ...interface{}) {
	// Shim the bot into the parameter list of every event dispatched via it.
	ev = append([]interface{}{bot}, ev...)
	bot.Conn.Dispatcher.Dispatch(name, ev...)
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
