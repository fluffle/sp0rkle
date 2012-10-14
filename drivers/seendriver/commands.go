package seendriver

import (
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"strings"
)

func seenCmd(line *base.Line) {
	s := strings.Fields(line.Args[1])
	if len(s) == 2 {
		// Assume we have "seen <nick> <action>"
		if n := sc.LastSeenDoing(s[0], strings.ToUpper(s[1])); n != nil {
			bot.ReplyN(line, "%s", n)
			return
		}
	}
	// Not specifically asking for that action, or no matching action.
	if n := sc.LastSeen(s[0]); n != nil {
		bot.ReplyN(line, "%s", n)
		return
	}
	// No exact matches for nick found, look for possible partial matches.
	if m := sc.SeenAnyMatching(s[1]); len(m) > 0 {
		if len(m) == 1 {
			if n := sc.LastSeen(m[0]); n != nil {
				bot.ReplyN(line, "1 possible match: %s", n)
			}
		} else if len(m) > 10 {
			bot.ReplyN(line, "%d possible matches, first 10 are: %s.",
				len(m), strings.Join(m[:9], ", "))
		} else {
			bot.ReplyN(line, "%d possible matches: %s.",
				len(m), strings.Join(m, ", "))
		}
		return
	}
	// No partial matches found. Check for people playing silly buggers.
	for _, w := range wittyComebacks {
		logging.Debug("Matching %#v...", w)
		if w.rx.MatchString(line.Args[1]) {
			bot.ReplyN(line, "%s", w.resp)
			return
		}
	}
	// Ok, probably a genuine query.
	bot.ReplyN(line, "Haven't seen %s before, sorry.", line.Args[1])
}

func lines(line *base.Line) {
	n := line.Nick
	if len(line.Args[1]) > 0 {
		n = line.Args[1]
	}
	sn := sc.LinesFor(n, line.Args[0])
	if sn != nil {
		bot.ReplyN(line, "%s has said %d lines in this channel",
			sn.Nick, sn.Lines)
	}
}

func topten(line *base.Line) {
	top := sc.TopTen(line.Args[0])
	s := make([]string, 0, 10)
	for i, n := range top {
		s = append(s, fmt.Sprintf("#%d: %s - %d", i+1, n.Nick, n.Lines))
	}
	bot.Reply(line, "%s", strings.Join(s, ", "))
}
