package seendriver

import (
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/lib/seen"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
	"strings"
	"time"
)

func (sd *seenDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(sd_recordseen), "bot_privmsg", "bot_action",
		"bot_join", "bot_part", "bot_quit", "bot_kick", "bot_nick")
	r.AddHandler(bot.NewHandler(sd_privmsg), "bot_privmsg")
	r.AddHandler(bot.NewHandler(sd_smoke), "bot_privmsg", "bot_action")
}

func sd_smoke(bot *bot.Sp0rkle, line *base.Line) {
	if ! smokeRx.MatchString(line.Args[1]) {
		return
	}
	sd := bot.GetDriver(driverName).(*seenDriver)
	if n := sd.LastSeenDoing(line.Nick, "SMOKE"); n != nil {
		bot.ReplyN(line, "You last went for a smoke %s ago...",
			time.Now().Sub(n.Timestamp))
	}
	n := seen.SawNick(
		db.StorableNick{line.Nick, line.Ident, line.Host},
		db.StorableChan{line.Args[0]},
		"SMOKE",
		"",
	)
	if _, err := sd.Upsert(n.Index(), n); err != nil {
		bot.Reply(line, "Failed to store smoke data: %v", err)
	}
}

func sd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	s := strings.Split(line.Args[1], " ")
	if !line.Addressed || s[0] != "seen" {
		return
	}

	sd := bot.GetDriver(driverName).(*seenDriver)
	if len(s) > 2 {
		// Assume we have "seen <nick> <action>"
		if n := sd.LastSeenDoing(s[1], strings.ToUpper(s[2])); n != nil {
			bot.ReplyN(line, "%s", n)
			return
		}
	}
	// Not specifically asking for that action, or no matching action.
	if n := sd.LastSeen(s[1]); n != nil {
		bot.ReplyN(line, "%s", n)
		return
	}
	// No exact matches for nick found, look for possible partial matches.
	if m := sd.SeenAnyMatching(s[1]); len(m) > 0 {
		if len(m) == 1 {
			if n := sd.LastSeen(m[0]); n != nil {
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
	txt := strings.Join(s[1:], " ")
	for _, w := range wittyComebacks {
		sd.l.Debug("Matching %#v...", w)
		if w.rx.MatchString(txt) {
			bot.ReplyN(line, "%s", w.resp)
			return
		}
	}
	// Ok, probably a genuine query.
	bot.ReplyN(line, "Haven't seen %s before, sorry.", txt)
}

func sd_recordseen(bot *bot.Sp0rkle, line *base.Line) {
	sd := bot.GetDriver(driverName).(*seenDriver)
	text, ch := "", ""
	if len(line.Args) > 1 {
		text = line.Args[1]
		ch = line.Args[0]
	} else if line.Cmd == "NICK" || line.Cmd == "QUIT" {
		// FFUU special cases.
		text = line.Args[0]
	}
	n := seen.SawNick(
		db.StorableNick{line.Nick, line.Ident, line.Host},
		db.StorableChan{ch},
		line.Cmd,
		text,
	)
	_, err := sd.Upsert(n.Index(), n)
	if err != nil && ch != "" {
		bot.ReplyN(line, "Failed to store seen data: %v", err)
	}
}
