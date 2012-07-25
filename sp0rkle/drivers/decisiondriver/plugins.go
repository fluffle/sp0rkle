package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"github.com/fluffle/sp0rkle/lib/util"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"math/rand"
	"strings"
)

func (dd *decisionDriver) RegisterPlugins(pm base.PluginManager) {
	pm.AddPlugin(&DecisionPlugin{dd, dd_rand})
	pm.AddPlugin(&DecisionPlugin{dd, dd_decider})
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
// It's at times like this I miss easy number -> string conversion
// and first-class regex constructs. Doing without is fun!
func rand_replacer(val string, r *rand.Rand) string {
	for {
		// Work out the indices of the plugin start and end.
		ps := strings.Index(val, "<plugin=rand ")
		if ps == -1 {
			break
		}
		pe := strings.Index(val[ps:], ">")
		if pe == -1 {
			// WTF!?
			break
		}
		pe += ps
		// Mid is where the plugin args start.
		mid := ps + 13
		val = val[:ps] + randomFloatAsString(val[mid:pe], r) + val[pe+1:]
	}
	return val
}

// Split this out so we can inject a deterministic rand.Rand for testing.
func rand_decider(val string, r *rand.Rand) string {
	i := 0
	for {
		i++
		// Work out the indices of the plugin start and end.
		ps := strings.Index(val, "<plugin=decide ")
		if ps == -1 {
			break
		}
		pe := strings.Index(val[ps:], ">")
		if pe == -1 {
			// No closing '>', so abort
			break
		}
		pe += ps
		// Mid is where the plugin args start.
		mid := ps + 15
		if options := splitDelimitedString(val[mid:pe]); len(options) > 0 {
			rnd := r.Intn(len(options))
			chosenone := strings.TrimSpace(options[rnd])
			val = val[:ps] + chosenone + val[pe+1:]
		} else {
			// Previously r.Intn() was called even for errors so this makes
			// tests pass. I'll remove it in a future commit.
			r.Intn(4)
			val = val[:ps] + "<plugin error>" + val[pe+1:]
		}
	}
	return val
}
