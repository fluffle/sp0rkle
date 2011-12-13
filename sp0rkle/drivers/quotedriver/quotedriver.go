package quotedriver

import (
	"github.com/fluffle/golog/logging"
	//	"launchpad.net/gobson/bson"
	"lib/db"
	"lib/quotes"
	//	"lib/util"
	//	"sp0rkle/base"
	//	"strings"
	"time"
)

const driverName string = "quotes"
const second int64 = 1e9

type rateLimit struct {
	badness  int64
	lastsent int64
}

type quoteDriver struct {
	*quotes.QuoteCollection

	// Data for rate limiting quote lookups per-nick
	limits map[string]*rateLimit

	// logging object
	l logging.Logger
}

func QuoteDriver(db *db.Database, l logging.Logger) *quoteDriver {
	qc := quotes.Collection(db, l)
	return &quoteDriver{
		QuoteCollection: qc,
		limits:          make(map[string]*rateLimit),
		l:               l,
	}
}

func (qd *quoteDriver) Name() string {
	return driverName
}

func (qd *quoteDriver) rateLimit(nick string) bool {
	lim, ok := qd.limits[nick]
	if !ok {
		lim = new(rateLimit)
		qd.limits[nick] = lim
	}
	// limit to 1 quote every 15 seconds, burst to 4 quotes
	elapsed := time.Nanoseconds() - lim.lastsent
	if lim.badness += 15*second - elapsed; lim.badness < 0 {
		lim.badness = 0
	}
	if lim.badness > 60*second {
		return true
	}
	lim.lastsent = time.Nanoseconds()
	return false
}
