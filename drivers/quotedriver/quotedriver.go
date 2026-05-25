package quotedriver

import (
	"time"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/quotes"
)

type Driver struct {
	qc *quotes.Collection
}

func New(b *bot.Bot, qc *quotes.Collection) *Driver {
	d := &Driver{qc: qc}

	b.Rewrite(d.quotePlugin)
	b.Command(d.add, "qadd", "qadd <quote>  -- Adds a quote to the db.")
	b.Command(d.add, "quote add",
		"quote add <quote>  -- Adds a quote to the db.")
	b.Command(d.add, "add quote",
		"add quote <quote>  -- Adds a quote to the db.")
	b.Command(d.del, "qdel", "qdel #<qID>  -- Deletes a quote from the db.")
	b.Command(d.del, "quote del",
		"quote del #<qID>  -- Deletes a quote from the db.")
	b.Command(d.del, "del quote",
		"del quote #<qID>  -- Deletes a quote from the db.")
	b.Command(d.fetch, "quote #", "quote #<qID>  -- Displays quote <qID>.")
	b.Command(d.lookup, "quote",
		"quote <regex>  -- Displays quotes matching <regex>")

	return d
}

// Data for rate limiting quote lookups per-nick
type rateLimit struct {
	badness  time.Duration
	lastsent time.Time
}

var limits = map[string]*rateLimit{}

func RateLimit(nick string) bool {
	lim, ok := limits[nick]
	if !ok {
		lim = new(rateLimit)
		limits[nick] = lim
	}
	// limit to 1 quote every 15 seconds, burst to 4 quotes
	elapsed := time.Now().Sub(lim.lastsent)
	if lim.badness += 15*time.Second - elapsed; lim.badness < 0 {
		lim.badness = 0
	}
	if lim.badness > 60*time.Second {
		return true
	}
	lim.lastsent = time.Now()
	return false
}
