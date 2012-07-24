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
		if options := choices(val[mid:pe]); options != nil
			rnd := r.Intn(len(options))
			chosenone := strings.TrimSpace(options[rnd])
			val = val[:ps] + chosenone + val[pe+1:]
		} else {
			val = val[:ps] + "<plugin error>" + val[pe+1:]
		}
	}
	return val
}

func choices(val string) []string {
	// We accept three different delimiter types in the input string, and we use
	// the following heuristics to determine what type of parsing style to use.

	// 1. Look for \s(["'])\S *and* \S\1\s in the input string. If we find this,
	//    use a mixed parsing style that accepts space separated bare words and
	//    ' or " delimited strings that may contain spaces *and* |
	// 2. Look for an occurrence of the | character. If we find it, split on it.
	// 3. Split on spaces.

	idx := strings.IndexAny(val, `"'`)
	// Careful to make sure that a string where the only quote is the last
	// character doesn't cause a panic.
	if idx == len(val)-1 || idx == -1 {
		return simpleSplit(val)
	}
	// This should all be safe as `'` `"` and ` ` are all one byte in UTF-8
	if idx == 0 || val[idx-1] != ' ' {
		idy := strings.Index(val[idx+1:], string(val[idx]))
		if idy == -1 {
			// No matching quote char found, so try a simple split instead
			return simpleSplit(val)
		}
		// Unicode rune after first quote
		rx, _ := utf8.DecodeRuneInString(val[idx+1:])
		// Unicode rune before second quote
		ry, _ := utf8.DecodeLastRuneInString(val[:idy])
		// 
		if (unicode.IsLetter(rx) || unicode.IsNumber(rx)) &&
			(unicode.IsLetter(ry) || unicode.IsNumber(ry)) &&
			(idy == len(val)-1 || val[idy+1] == ' ') {
			return quoteSplit(val)
		}
	}
	return simpleSplit(val)
}

func simpleSplit(val string) []string {
	if strings.Index(val, "|") != -1 {
		// | is a simple delimiter
		// NOTE: spaces either side of the | are taken care of by the caller
		return strings.Split(val, "|")
	}
	// String doesn't contain any seperator chars,
	// so is just a list of options to choose from
	return strings.Split(val, " ")
}

func quoteSplit(val string) []string {
	
