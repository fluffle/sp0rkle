package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"fmt"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Init() {
	bot.Rewrite(randPlugin)
	bot.Rewrite(decidePlugin)

	bot.Command(randCmd, "rand", "rand <range>  -- "+
		"choose a random number in range [lo-]hi")
	bot.Command(decideCmd, "decide", "decide <options>  -- "+
		"choose one of the (space, pipe, quote) delimited options at random")
}

func randomFloatAsString(val string, r *rand.Rand) string {
	// val should be in the format: [lo-]hi[ format]
	var lo, hi float64
	var err error
	format := "%.0f"
	// If there's a space in the string, we also have a format.
	sp := strings.Index(val, " ")
	if sp != -1 {
		format = strings.TrimSpace(val[sp:])
		val = val[:sp]
	}
	// If there's a dash in var, we have a range lo-hi, rather than just 0-hi.
	if dash := strings.Index(val, "-"); dash != -1 {
		if lo, err = strconv.ParseFloat(val[:dash], 32); err != nil {
			lo = 0
		}
		if hi, err = strconv.ParseFloat(val[dash+1:], 32); err != nil {
			hi = 0
		}
	} else {
		lo = 0
		if hi, err = strconv.ParseFloat(val, 32); err != nil {
			hi = 0
		}
	}
	rnd := r.Float64()*(hi-lo) + lo
	return fmt.Sprintf(format, rnd)
}

func splitDelimitedString(val string) []string {
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
	if idx == 0 || val[idx-1] == ' ' {
		// Locate closing quote
		idy := strings.Index(val[idx+1:], string(val[idx]))
		if idy == -1 {
			// No matching quote char found, so try a simple split instead
			return simpleSplit(val)
		}
		// Reindex idy to start of val, as it's currently relative to idx+1
		idy += idx + 1
		// Unicode rune after first quote char
		rx, _ := utf8.DecodeRuneInString(val[idx+1:])
		// Unicode rune before second quote char
		ry, _ := utf8.DecodeLastRuneInString(val[:idy])
		// Check heuristic outlined in 1 above
		if !unicode.IsSpace(rx) && !unicode.IsSpace(ry) &&
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
	ret := make([]string, 0, 10)
	l := &util.Lexer{Input: val}
	for {
		l.Scan(unicode.IsSpace)
		c := l.Peek()
		switch c {
		case '"', '\'':
			// Consume the opening quote
			sep := l.Next()
			// Scan the string until the closing quote
			ret = append(ret, l.Find(c))
			// Advance over closing quote and test for mismatched quotes
			if l.Next() != sep {
				// If we don't find the closing quote, something is broken
				return []string{}
			}
		case 0:
			return ret
		default:
			// It's not a quote or a space, so scan until the next space char.
			// Hopefully on IRC the mismatch between unicode.IsSpace and ' '
			// won't be *too* apparent...
			ret = append(ret, l.Find(' '))
		}
	}
	// Shouldn't ever be reached, but required for the go compiler
	return ret
}
