package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"fmt"
	"os"
	"rand"
	"sp0rkle/base"
	"strings"
	"strconv"
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
	return rand_replacer(val, dd.rng)
}

// Split this out so we can inject a deterministic rand.Rand for testing.
// It's at times like this I miss easy number -> string conversion
// and first-class regex constructs. Doing without is fun!
func rand_replacer(val string, r *rand.Rand) string {
	for {
		var lo, hi float32
		var err os.Error
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
			if lo, err = strconv.Atof32(val[mid:dash]); err != nil {
				lo = 0
			}
			if hi, err = strconv.Atof32(val[dash+1 : sp]); err != nil {
				hi = 0
			}
		} else {
			lo = 0
			if hi, err = strconv.Atof32(val[mid:sp]); err != nil {
				hi = 0
			}
		}
		rnd := r.Float32()*(hi-lo) + lo
		val = val[:ps] + fmt.Sprintf(format, rnd) + val[pe+1:]
	}
	return val
}

func dd_decider(dd *decisionDriver, val string, line *base.Line) string {
	return rand_decider(val, dd.rng)
}

// Split this out so we can inject a deterministic rand.Rand for testing.
func rand_decider(val string, r *rand.Rand) string {
	i := 0
	for {
		i ++
		// TODO(bob-smith): The following block could be refactored out
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
		fmt.Printf("\n\nA: '%s'\n", val[mid:pe])
		// options := strings.SplitN(val[mid:pe]," ", -1)
		options := choices(val[mid:pe])
		rnd := r.Intn(len(options))
		fmt.Printf("length: %d, rnd: %d\n",len(options), rnd)
		chosenone := strings.TrimSpace(options[rnd])
		fmt.Printf("%s\n",strings.Join(options,"|"))
		fmt.Printf("chosen: %s\n",chosenone)
		val = val[:ps] + chosenone  + val[pe+1:]
	}
	return val
}


func choices(val string) []string{
	if strings.IndexAny(val, "\"'") != -1{
		d := strings.IndexAny(val, "\"'")
    var delim string
		delim = string(val[d])
		fmt.Printf("d: %d\ndelim: %q\n", d, delim)
		if strings.Count(val, delim) % 2 == 1{
			return []string{"Unbalanced quotes"}
		}
		fmt.Printf("FOO\n")
		tmp := strings.Split(val,delim)
		var ret []string
		for i := 1; i < len(tmp) ; i += 2{
			ret = append(ret, tmp[i])
		}
		return ret

	} else if strings.IndexAny(val, "|") != -1{
     fmt.Printf("Splitting on |\n")
     fmt.Printf("FOO\n")
     tmp := strings.Split(val,"|")
     var ret []string
     for i := 1; i < len(tmp) ; i ++{
       ret = append(ret, tmp[i])
     }
     return ret
  } else {
	fmt.Printf("Splitting on whitespace\n")
		//String doesn't contains ' or ", so is just a list of words
		return strings.SplitN(val," ", -1)
	}
	fmt.Printf("FAILED\n")
	return []string{"BROKEN", val}
}

