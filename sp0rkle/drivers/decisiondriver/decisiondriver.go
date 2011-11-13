package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"github.com/fluffle/goevent/event"
	"lib/util"
	"rand"
	"sp0rkle/base"
	"sp0rkle/bot"
	"time"
)

const driverName string = "decisions"

type decisionDriver struct {
	rng *rand.Rand
}

func DecisionDriver() *decisionDriver {
	return &decisionDriver{util.NewRand(time.Nanoseconds())}
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
