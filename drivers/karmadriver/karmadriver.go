package karmadriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/karma"
	"strings"
)

var kc *karma.Collection

func Init() {
	kc = karma.Init()

	bot.HandleFunc(recordKarma, "privmsg", "action")

	bot.CommandFunc(karmaCmd, "karma", "karma <thing>  -- "+
		"Retrieve the karma score of <thing>.")
}

type kt struct {
	thing string
	plus  bool
}

func karmaThings(s string) []kt {
	ret := make([]kt, 0)
	start, end, endp, endm, index, plus := 0, 0, 0, 0, 0, true
	for {
		endp = strings.Index(s[index:], "++")
		endm = strings.Index(s[index:], "--")
		if endp == endm {
			// the only time this can be true is when both are -1
			break
		}
		plus = true
		end = endp
		if endp == -1 || (endm >= 0 && endm < endp) {
			plus = false
			end = endm
		}
		end += index
		index = end + 2
		if end == 0 {
			// String begins with ++ or --, ignore it
			continue
		}
		start = strings.LastIndex(s[:end], " ") + 1
		if s[end-1] == ')' {
			end--
			start = strings.LastIndex(s[:end], "(") + 1
			if start == 0 {
				// Missing opening bracket for )++
				continue
			}
		}
		if start == end {
			// Either ' --' or '()++'
			continue
		}
		ret = append(ret, kt{s[start:end], plus})
	}
	return ret
}
