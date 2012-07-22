package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
)

const driverName string = "decisions"

type decisionDriver struct {
	// Nothing needed here, yet.
}

func DecisionDriver() *decisionDriver {
	return &decisionDriver{}
}

type DecisionHandler func(*bot.Sp0rkle, *decisionDriver, *base.Line)

// Unboxer for DecisionDriver handlers
func DDHandler(f DecisionHandler) event.Handler {
	return event.NewHandler(func(ev ...interface{}) {
		f(ev[0].(*bot.Sp0rkle), ev[1].(*decisionDriver), ev[2].(*base.Line))
	})
}

func (dd *decisionDriver) Name() string {
	return driverName
}
