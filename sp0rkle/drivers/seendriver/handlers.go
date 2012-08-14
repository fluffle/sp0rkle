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
	r.AddHandler(bot.NewHandler(sd_record_pm), "bot_privmsg", "bot_action")
	r.AddHandler(bot.NewHandler(sd_record_chan), "bot_join", "bot_part")
	r.AddHandler(bot.NewHandler(sd_record_nick), "bot_quit", "bot_nick")
	r.AddHandler(bot.NewHandler(sd_record_kick), "bot_kick")
	r.AddHandler(bot.NewHandler(sd_privmsg), "bot_privmsg")
	r.AddHandler(bot.NewHandler(sd_smoke), "bot_privmsg", "bot_action")
}

func sd_smoke(bot *bot.Sp0rkle, line *base.Line) {
	if ! smokeRx.MatchString(line.Args[1]) {
		return
	}
	sd := bot.GetDriver(driverName).(*seenDriver)
	sn := sd.LastSeenDoing(line.Nick, "SMOKE")
	n, c := line.Storable()
	if sn != nil {
		bot.ReplyN(line, "You last went for a smoke %s ago...",
			time.Since(sn.Timestamp))
		sn.StorableNick, sn.StorableChan = n, c
		sn.Timestamp = time.Now()
	} else {
		sn = seen.SawNick(n, c, "SMOKE", "")
	}
	if _, err := sd.Upsert(sn.Index(), sn); err != nil {
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

func sd_record_pm(bot *bot.Sp0rkle, line *base.Line) {
	sd := bot.GetDriver(driverName).(*seenDriver)
	sn := sd.SeenNickFromLine(line)
	sn.Text = line.Args[1]
	sn.Lines++
	for _, n := range milestones {
		if sn.Lines == n {
			bot.Reply(line, "%s has said %d lines in this channel and" +
				"should now shut the fuck up and do something useful",
				line.Nick, sn.Lines)
		}
	}
	_, err := sd.Upsert(sn.Index(), sn)
	if err != nil {
		bot.Reply(line, "Failed to store seen data: %v", err)
	}
}

func sd_record_chan(bot *bot.Sp0rkle, line *base.Line) {
	sd := bot.GetDriver(driverName).(*seenDriver)
	sn := sd.SeenNickFromLine(line)
	if len(line.Args) > 1 {
		// If we have a PART message
		sn.Text = line.Args[1]
	}
	_, err := sd.Upsert(sn.Index(), sn)
	if err != nil {
		bot.Reply(line, "Failed to store seen data: %v", err)
	}
}

func sd_record_nick(bot *bot.Sp0rkle, line *base.Line) {
	sd := bot.GetDriver(driverName).(*seenDriver)
	sn := sd.SeenNickFromLine(line)
	sn.Chan = ""
	sn.Text = line.Args[0]
	_, err := sd.Upsert(sn.Index(), sn)
	if err != nil {
		// We don't have anyone to reply to in this case, so log instead.
		sd.l.Warn("Failed to store seen data: %v", err)
	}
}

func sd_record_kick(bot *bot.Sp0rkle, line *base.Line) {
	sd := bot.GetDriver(driverName).(*seenDriver)
	n, c := line.Storable()
	kn := db.StorableNick{Nick: line.Args[1]}
	// SeenNickFromLine doesn't work with the hacks for KICKING and KICKED
	// First, handle KICKING
	kr := sd.LastSeenDoing(line.Nick, "KICKING")
	if kr == nil {
		kr = seen.SawNick(n, c, "KICKING", line.Args[2])
	} else {
		kr.StorableNick, kr.StorableChan = n, c
		kr.Timestamp, kr.Text = time.Now(), line.Args[2]
	}
	kr.OtherNick = kn
	_, err := sd.Upsert(kr.Index(), kr)
	if err != nil {
		bot.Reply(line, "Failed to store seen data: %v", err)
	}
	// Now, handle KICKED
	ke := sd.LastSeenDoing(line.Args[1], "KICKED")
	if ke == nil {
		ke = seen.SawNick(kn, c, "KICKED", line.Args[2])
	} else {
		ke.StorableNick, ke.StorableChan = kn, c
		ke.Timestamp, ke.Text = time.Now(), line.Args[2]
	}
	ke.OtherNick = n
	_, err = sd.Upsert(ke.Index(), ke)
	if err != nil {
		bot.Reply(line, "Failed to store seen data: %v", err)
	}
}
