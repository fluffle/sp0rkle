package seendriver

import (
	"fmt"
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
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: You last went for a smoke %s ago...",
			line.Nick, time.Now().Sub(n.Timestamp)))
	}
	n := seen.SawNick(
		db.StorableNick{line.Nick, line.Ident, line.Host},
		db.StorableChan{line.Args[0]},
		"SMOKE",
		"",
	)
	if _, err := sd.Upsert(n.Index(), n); err != nil {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"Failed to store smoke data: %v", err))
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
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: %s", line.Nick, n))
			return
		}
	}
	// Not specifically asking for that action, or no matching action.
	if n := sd.LastSeen(s[1]); n != nil {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: %s", line.Nick, n))
		return
	}
	// No exact matches for nick found, look for possible partial matches.
	if m := sd.SeenAnyMatching(s[1]); len(m) > 0 {
		if len(m) == 1 {
			if n := sd.LastSeen(m[0]); n != nil {
				bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
					"%s: 1 possible match: %s", line.Nick, n))
			}
		} else if len(m) > 10 {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: %d possible matches, first 10 are: %s.",
				line.Nick, len(m), strings.Join(m[:9], ", ")))
		} else {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: %d possible matches: %s.",
				line.Nick, len(m), strings.Join(m, ", ")))
		}
		return
	}
	// No partial matches found. Check for people playing silly buggers.
	txt := strings.Join(s[1:], " ")
	for _, w := range wittyComebacks {
		sd.l.Debug("Matching %#v...", w)
		if w.rx.MatchString(txt) {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: %s", line.Nick, w.resp))
			return
		}
	}
	// Ok, probably a genuine query.
	bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
		"%s: Haven't seen %s before, sorry.", line.Nick, txt))
}

func sd_recordseen(bot *bot.Sp0rkle, line *base.Line) {
	sd := bot.GetDriver(driverName).(*seenDriver)
	text := ""
	if len(line.Args) > 1 {
		text = line.Args[1]
	} else if line.Cmd == "NICK" || line.Cmd == "QUIT" {
		// FFUU special cases.
		text = line.Args[0]
	}
	n := seen.SawNick(
		db.StorableNick{line.Nick, line.Ident, line.Host},
		db.StorableChan{line.Args[0]},
		line.Cmd,
		text,
	)
	_, err := sd.Upsert(n.Index(), n)
	if err != nil {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"Failed to store seen data: %v", err))
	}
}
