package statsdriver

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/stats"
)

type Driver struct {
	sc *stats.Collection
}

func New(b *bot.Bot, sc *stats.Collection) *Driver {
	d := &Driver{sc: sc}

	b.Handle(d.recordStats, client.PRIVMSG, client.ACTION)

	b.Command(d.statsCmd, "lines", "lines [nick]  -- "+
		"display how many lines you [or nick] has said in the channel")
	b.Command(d.statsCmd, "stats", "stats [nick]  -- "+
		"display how many lines you [or nick] has said in the channel")
	b.Command(d.topten, "topten", "topten  -- "+
		"display the nicks who have said the most in the channel")
	b.Command(d.topten, "top10", "top10  -- "+
		"display the nicks who have said the most in the channel")

	return d
}
