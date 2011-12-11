package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"github.com/fluffle/goevent/event"
	"sp0rkle/bot"
	"sp0rkle/base"
)


func (dd *decisionDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(dd_privmsg), "bot_privmsg")
}

func dd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	 //dd := bot.GetDriver(driverName).(*decisionDriver)
}
