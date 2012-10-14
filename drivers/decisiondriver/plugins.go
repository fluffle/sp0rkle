package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"github.com/fluffle/sp0rkle/util"
	"github.com/fluffle/sp0rkle/base"
	"math/rand"
	"strings"
)

func (dd *decisionDriver) RegisterPlugins(pm base.PluginManager) {
	pm.Add(&DecisionPlugin{dd, dd_rand})
	pm.Add(&DecisionPlugin{dd, dd_decider})
}

type DecisionPlugin struct {
	provider  *decisionDriver
	processor func(*decisionDriver, string, *base.Line) string
}

func (fp *DecisionPlugin) Apply(val string, line *base.Line) string {
	return fp.processor(fp.provider, val, line)
}

func dd_rand(dd *decisionDriver, val string, line *base.Line) string {
	return rand_replacer(val, util.RNG)
}

func dd_decider(dd *decisionDriver, val string, line *base.Line) string {
	return rand_decider(val, util.RNG)
}

// Split this out so we can inject a deterministic rand.Rand for testing.
func rand_replacer(val string, r *rand.Rand) string {
	f := func(s string) string {
		return randomFloatAsString(s, r)
	}
	return util.ApplyPluginFunction(val, "rand", f)
}

// Split this out so we can inject a deterministic rand.Rand for testing.
func rand_decider(val string, r *rand.Rand) string {
	f := func(s string) string {
		if options := splitDelimitedString(s); len(options) > 0 {
			return strings.TrimSpace(options[r.Intn(len(options))])
		}
		// Previously r.Intn() was called even for errors so this makes
		// tests pass. I'll remove it in a future commit.
		r.Intn(4)
		return "<plugin error>"
	}
	return util.ApplyPluginFunction(val, "decide", f)
}
