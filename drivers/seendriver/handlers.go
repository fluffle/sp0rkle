package seendriver

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/seen"
	"github.com/fluffle/sp0rkle/util"
	"time"
)

func smoke(line *base.Line) {
	if ! smokeRx.MatchString(line.Args[1]) {
		return
	}
	sn := sc.LastSeenDoing(line.Nick, "SMOKE")
	n, c := line.Storable()
	if sn != nil {
		bot.ReplyN(line, "You last went for a smoke %s ago...",
			util.TimeSince(sn.Timestamp))
		sn.Nick, sn.Chan = n, c
		sn.Timestamp = time.Now()
	} else {
		sn = seen.SawNick(n, c, "SMOKE", "")
	}
	if _, err := sc.Upsert(sn.Id(), sn); err != nil {
		bot.Reply(line, "Failed to store smoke data: %v", err)
	}
}

func recordLines(line *base.Line) {
	sn := sc.LinesFor(line.Nick, line.Args[0])
	if sn == nil {
		n, c := line.Storable()
		sn = seen.SawNick(n, c, "LINES", "")
	}
	sn.Lines++
	for _, n := range milestones {
		if sn.Lines == n {
			bot.Reply(line, "%s has said %d lines in this channel and" +
				"should now shut the fuck up and do something useful",
				line.Nick, sn.Lines)
		}
	}
	if _, err := sc.Upsert(sn.Id(), sn); err != nil {
		bot.Reply(line, "Failed to store seen data: %v", err)
	}
}

func recordPrivmsg(line *base.Line) {
	sn := seenNickFromLine(line)
	sn.Text = line.Args[1]
	if _, err := sc.Upsert(sn.Id(), sn); err != nil {
		bot.Reply(line, "Failed to store seen data: %v", err)
	}
}

func recordJoin(line *base.Line) {
	sn := seenNickFromLine(line)
	if len(line.Args) > 1 {
		// If we have a PART message
		sn.Text = line.Args[1]
	}
	if _, err := sc.Upsert(sn.Id(), sn); err != nil {
		bot.Reply(line, "Failed to store seen data: %v", err)
	}
}

func recordNick(line *base.Line) {
	sn := seenNickFromLine(line)
	sn.Chan = ""
	sn.Text = line.Args[0]
	if _, err := sc.Upsert(sn.Id(), sn); err != nil {
		// We don't have anyone to reply to in this case, so log instead.
		logging.Warn("Failed to store seen data: %v", err)
	}
}

func recordKick(line *base.Line) {
	n, c := line.Storable()
	kn := base.Nick(line.Args[1])
	// seenNickFromLine doesn't work with the hacks for KICKING and KICKED
	// First, handle KICKING
	kr := sc.LastSeenDoing(line.Nick, "KICKING")
	if kr == nil {
		kr = seen.SawNick(n, c, "KICKING", line.Args[2])
	} else {
		kr.Nick, kr.Chan = n, c
		kr.Timestamp, kr.Text = time.Now(), line.Args[2]
	}
	kr.OtherNick = kn
	_, err := sc.Upsert(kr.Id(), kr)
	if err != nil {
		bot.Reply(line, "Failed to store seen data: %v", err)
	}
	// Now, handle KICKED
	ke := sc.LastSeenDoing(line.Args[1], "KICKED")
	if ke == nil {
		ke = seen.SawNick(kn, c, "KICKED", line.Args[2])
	} else {
		ke.Nick, ke.Chan = kn, c
		ke.Timestamp, ke.Text = time.Now(), line.Args[2]
	}
	ke.OtherNick = n
	_, err = sc.Upsert(ke.Id(), ke)
	if err != nil {
		bot.Reply(line, "Failed to store seen data: %v", err)
	}
}
