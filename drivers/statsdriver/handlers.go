package statsdriver

import (
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/stats"
)

func recordStats(line *base.Line) {
	ns := sc.StatsFor(line.Nick, line.Args[0])
	if ns == nil {
		n, c := line.Storable()
		ns = stats.NewStat(n, c)
	}
	ns.Update(line.Args[1])
	if ns.Lines%10000 == 0 {
		bot.Reply(line, "%s has said %d lines in this channel and "+
			"should now shut the fuck up and do something useful",
			line.Nick, ns.Lines)
	}
	if _, err := sc.Upsert(ns.Id(), ns); err != nil {
		bot.Reply(line, "Failed to store stats data: %v", err)
	}
}
