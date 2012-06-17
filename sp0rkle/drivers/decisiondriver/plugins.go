package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"fmt"
	"github.com/fluffle/sp0rkle/lib/util"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"math/rand"
	"strconv"
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

// Split this out so we can inject a deterministic rand.Rand for testing.
// It's at times like this I miss easy number -> string conversion
// and first-class regex constructs. Doing without is fun!
func rand_replacer(val string, r *rand.Rand) string {
	for {
		var lo, hi float64
		var err error
		format := "%.0f"
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
		// If there's a space before the plugin ends, we also have a format.
		sp := strings.Index(val[mid:pe], " ")
		if sp != -1 {
			sp += mid
			format = strings.TrimSpace(val[sp:pe])
		} else {
			sp = pe
		}
		// If there's a dash before the space or the plugin ends, we have a
		// range lo-hi, rather than just 0-hi.
		if dash := strings.Index(val[mid:sp], "-"); dash != -1 {
			dash += mid
			if lo, err = strconv.ParseFloat(val[mid:dash], 32); err != nil {
				lo = 0
			}
			if hi, err = strconv.ParseFloat(val[dash+1:sp], 32); err != nil {
				hi = 0
			}
		} else {
			lo = 0
			if hi, err = strconv.ParseFloat(val[mid:sp], 32); err != nil {
				hi = 0
			}
		}
		rnd := r.Float64()*(hi-lo) + lo
		val = val[:ps] + fmt.Sprintf(format, rnd) + val[pe+1:]
	}
	return val
}

func dd_decider(dd *decisionDriver, val string, line *base.Line) string {
	return rand_decider(val, util.RNG)
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
		// options := strings.SplitN(val[mid:pe]," ", -1)
		options := choices(val[mid:pe])
		rnd := r.Intn(len(options))
		chosenone := strings.TrimSpace(options[rnd])
		val = val[:ps] + chosenone + val[pe+1:]
	}
	return val
}

func choices(val string) []string {
	d := strings.IndexAny(val, `"'|`)
	if d == -1 {
		// String doesn't contain any seperator chars,
		// so is just a list of options to choose from
		return strings.Split(val, " ")
	}
	delim := string(val[d])
	tmp := strings.Split(val, delim)
	if delim == "|" {
		return tmp
	}
	// Make sure we have balanced quotes
	if len(tmp) % 2 == 0 {
		return []string{"Unbalanced quotes"}
	}
	// Copy out the possible choice values
	ret := make([]string, (len(tmp)-1)/2)
	for i, j := 1, 0; i < len(tmp); i, j = i+2, j+1 {
		ret[j] = tmp[i]
	}
	return ret
}
