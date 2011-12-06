package quotedriver

import (
	"fmt"
	"github.com/fluffle/goevent/event"
//	"launchpad.net/gobson/bson"
//	"lib/db"
//	"lib/quotes"
//	"rand"
	"sp0rkle/bot"
	"sp0rkle/base"
	"strings"
	"strconv"
//	"time"
)

type QuoteHandler func(*bot.Sp0rkle, *quoteDriver, *base.Line)

// Unboxer for QuoteDriver handlers
func QDHandler(f QuoteHandler) event.Handler {
	return event.NewHandler(func(ev ...interface{}) {
		f(ev[0].(*bot.Sp0rkle), ev[1].(*quoteDriver), ev[2].(*base.Line))
	})
}

func (qd *quoteDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(qd_privmsg), "bot_privmsg")
//	r.AddHandler(bot.NewHandler(qd_action), "bot_action")
	r.AddHandler(QDHandler(qd_lookup), "qd_lookup")
}

func qd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	qd := bot.GetDriver(driverName).(*quoteDriver)
	
	if !line.Addressed {
		return
	}

	l := strings.ToLower(line.Args[1])
	switch {
	// quote lookup by QID
	case strings.HasPrefix(l, "quote #"):
		// The remainder of the string should be a quote ID
		nl := line.Copy()
		nl.Args[1] = nl.Args[1][7:]
		bot.Dispatch("qd_lookup", qd, nl)
	}
}

func qd_lookup(bot *bot.Sp0rkle, qd *quoteDriver, line *base.Line) {
	qid, err :=	strconv.Atoi(line.Args[1])
	if err != nil {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: '%s' doesn't look like a quote id.", line.Nick, line.Args[1]))
		return
	}
	quote := qd.GetByQID(qid)
	if quote != nil {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"#%d: %s", quote.QID, quote.Quote))
	} else {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: No quote found for id %d", line.Nick, qid))
	}
}
