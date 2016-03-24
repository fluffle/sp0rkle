package karmadriver

import (
	"unicode"
	"unicode/utf8"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/karma"
	"github.com/fluffle/sp0rkle/util"
)

var kc *karma.Collection

func Init() {
	kc = karma.Init()

	bot.Handle(recordKarma, client.PRIVMSG, client.ACTION)

	bot.Command(karmaCmd, "karma", "karma <thing>  -- "+
		"Retrieve the karma score of <thing>.")
}

type kt struct {
	thing string
	plus  bool
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func isPlusMinus(r rune) bool {
	if r == '+' || r == '-' {
		return true
	}
	return false
}

func notPlusMinus(r rune) bool { return !isPlusMinus(r) }

func isAlphanumeric(r rune) bool {
	return unicode.In(r, unicode.L, unicode.N)
}

func karmaThings(s string) []kt {
	ret := make([]kt, 0)
	// Reversing the input is easier than making the lexer go backwards.
	l := &util.Lexer{Input: reverse(s)}
	for {
		prefix := l.Scan(notPlusMinus)
		if l.Peek() == 0 {
			break
		} else if len(prefix) != 0 {
			// Require a space or end of string after karma identifier.
			if r, _ := utf8.DecodeLastRuneInString(prefix); !unicode.IsSpace(r) {
				l.Next()
				continue
			}
		}

		thing := kt{}
		pm := l.Scan(isPlusMinus)
		if pm == "++" {
			thing.plus = true
		} else if pm != "--" {
			continue
		}

		if l.Peek() == ')' {
			// TODO: Still doesn't handle nested brackets. Do I care?
			l.Next()
			thing.thing = reverse(l.Find('('))
			if l.Peek() == 0 {
				// Hit EOF while looking for the opening bracket.
				break
			}
		} else {
			// TODO: This may be unexpected behaviour to some.
			thing.thing = reverse(l.Scan(isAlphanumeric))
		}
		if len(thing.thing) == 0 {
			continue
		}
		ret = append(ret, thing)
	}
	return ret
}
