package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"github.com/fluffle/sp0rkle/util"
	"github.com/fluffle/sp0rkle/base"
	"strings"
)

func randPlugin(val string, line *base.Line) string {
	f := func(s string) string {
		return randomFloatAsString(s, util.RNG)
	}
	return util.ApplyPluginFunction(val, "rand", f)
}

func decidePlugin(val string, line *base.Line) string {
	f := func(s string) string {
		if options := splitDelimitedString(s); len(options) > 0 {
			return strings.TrimSpace(options[util.RNG.Intn(len(options))])
		}
		return "<plugin error>"
	}
	return util.ApplyPluginFunction(val, "decide", f)
}
