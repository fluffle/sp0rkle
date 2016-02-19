package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"math/rand"
	"strings"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util"
)

func randPlugin(val string, ctx *bot.Context) string {
	f := func(s string) string {
		return randomFloatAsString(s)
	}
	return util.ApplyPluginFunction(val, "rand", f)
}

func decidePlugin(val string, ctx *bot.Context) string {
	f := func(s string) string {
		if options := splitDelimitedString(s); len(options) > 0 {
			return strings.TrimSpace(options[rand.Intn(len(options))])
		}
		return "<plugin error>"
	}
	return util.ApplyPluginFunction(val, "decide", f)
}
