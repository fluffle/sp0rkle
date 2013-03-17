package statsdriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/stats"
)

func recordStats(ctx *bot.Context) {
	ns := sc.StatsFor(ctx.Nick, ctx.Target())
	if ns == nil {
		n, c := ctx.Storable()
		ns = stats.NewStat(n, c)
	}
	ns.Update(ctx.Text())
	if ns.Lines%10000 == 0 {
		ctx.Reply("%s has said %d lines in this channel and "+
			"should now shut the fuck up and do something useful",
			ctx.Nick, ns.Lines)

	}
	if _, err := sc.Upsert(ns.Id(), ns); err != nil {
		ctx.Reply("Failed to store stats data: %v", err)
	}
}
