package seendriver

import (
	"strings"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
)

func seenCmd(ctx *bot.Context) {
	s := strings.Fields(ctx.Text())
	if len(s) == 2 {
		// Assume we have "seen <nick> <action>"
		if n := sc.LastSeenDoing(s[0], strings.ToUpper(s[1])); n != nil {
			ctx.ReplyN("%s", n)
			return
		}
	}
	// Not specifically asking for that action, or no matching action.
	if n := sc.LastSeen(s[0]); n != nil {
		ctx.ReplyN("%s", n)
		return
	}
	// No exact matches for nick found, look for possible partial matches.
	if m := sc.SeenAnyMatching(s[0]); len(m) > 0 {
		if len(m) == 1 {
			if n := sc.LastSeen(m[0]); n != nil {
				ctx.ReplyN("1 possible match: %s", n)
			}
		} else if len(m) > 10 {
			ctx.ReplyN("%d possible matches, most recent 10 are: %s.",
				len(m), strings.Join(m[:9], ", "))

		} else {
			ctx.ReplyN("%d possible matches: %s.",
				len(m), strings.Join(m, ", "))

		}
		return
	}
	// No partial matches found. Check for people playing silly buggers.
	for _, w := range wittyComebacks {
		logging.Debug("Matching %#v...", w)
		if w.rx.MatchString(ctx.Text()) {
			ctx.ReplyN("%s", w.resp)
			return
		}
	}
	// Ok, probably a genuine query.
	ctx.ReplyN("Haven't seen %s before, sorry.", ctx.Text())
}
