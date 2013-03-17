package seendriver

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/seen"
	"github.com/fluffle/sp0rkle/util"
	"time"
)

func smoke(ctx *bot.Context) {
	if !smokeRx.MatchString(ctx.Text()) {
		return
	}
	sn := sc.LastSeenDoing(ctx.Nick, "SMOKE")
	n, c := ctx.Storable()
	if sn != nil {
		ctx.ReplyN("You last went for a smoke %s ago...",
			util.TimeSince(sn.Timestamp))

		sn.Nick, sn.Chan = n, c
		sn.Timestamp = time.Now()
	} else {
		sn = seen.SawNick(n, c, "SMOKE", "")
	}
	if _, err := sc.Upsert(sn.Id(), sn); err != nil {
		ctx.Reply("Failed to store smoke data: %v", err)
	}
}

func recordPrivmsg(ctx *bot.Context) {
	sn := seenNickFromLine(ctx)
	sn.Text = ctx.Text()
	if _, err := sc.Upsert(sn.Id(), sn); err != nil {
		ctx.Reply("Failed to store seen data: %v", err)
	}
}

func recordJoin(ctx *bot.Context) {
	sn := seenNickFromLine(ctx)
	if len(ctx.Args) > 1 {
		// If we have a PART message
		sn.Text = ctx.Text()
	}
	if _, err := sc.Upsert(sn.Id(), sn); err != nil {
		ctx.Reply("Failed to store seen data: %v", err)
	}
}

func recordNick(ctx *bot.Context) {
	sn := seenNickFromLine(ctx)
	sn.Chan = ""
	sn.Text = ctx.Target()
	if _, err := sc.Upsert(sn.Id(), sn); err != nil {
		// We don't have anyone to reply to in this case, so log instead.
		logging.Warn("Failed to store seen data: %v", err)
	}
}

func recordKick(ctx *bot.Context) {
	n, c := ctx.Storable()
	kn := bot.Nick(ctx.Text())
	// seenNickFromLine doesn't work with the hacks for KICKING and KICKED
	// First, handle KICKING
	kr := sc.LastSeenDoing(ctx.Nick, "KICKING")
	if kr == nil {
		kr = seen.SawNick(n, c, "KICKING", ctx.Args[2])
	} else {
		kr.Nick, kr.Chan = n, c
		kr.Timestamp, kr.Text = time.Now(), ctx.Args[2]
	}
	kr.OtherNick = kn
	_, err := sc.Upsert(kr.Id(), kr)
	if err != nil {
		ctx.Reply("Failed to store seen data: %v", err)
	}
	// Now, handle KICKED
	ke := sc.LastSeenDoing(ctx.Text(), "KICKED")
	if ke == nil {
		ke = seen.SawNick(kn, c, "KICKED", ctx.Args[2])
	} else {
		ke.Nick, ke.Chan = kn, c
		ke.Timestamp, ke.Text = time.Now(), ctx.Args[2]
	}
	ke.OtherNick = n
	_, err = sc.Upsert(ke.Id(), ke)
	if err != nil {
		ctx.Reply("Failed to store seen data: %v", err)
	}
}
