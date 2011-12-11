package quotedriver

import (
	"fmt"
	"github.com/fluffle/goevent/event"
	"launchpad.net/gobson/bson"
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
	r.AddHandler(QDHandler(qd_fetch), "qd_fetch")
	r.AddHandler(QDHandler(qd_lookup), "qd_lookup")
}

func qd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	qd := bot.GetDriver(driverName).(*quoteDriver)
	
	if !line.Addressed {
		return
	}

	l := strings.ToLower(line.Args[1])
	switch {
	// Quote lookup: quote | quote #QID | quote regex
	case strings.HasPrefix(l, "quote"):
		nl := line.Copy()
		handler := "qd_lookup"
		if len(l) < 6 {
			// The line is just "quote" so look up a random quote
			nl.Args[1] = ""
		} else if l[7] == '#' {
			// The remainder of the string should be a quote ID
			nl.Args[1] = nl.Args[1][7:]
			handler = "qd_fetch"
		} else {
			// The remainder of the string is a regex to lookup
			nl.Args[1] = nl.Args[1][6:]
		}
		bot.Dispatch(handler, qd, nl)
	}
}

func qd_fetch(bot *bot.Sp0rkle, qd *quoteDriver, line *base.Line) {
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

func qd_lookup(bot *bot.Sp0rkle, qd *quoteDriver, line *base.Line) {
	quote := qd.GetPseudoRand(line.Args[1])
	if quote == nil {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: No quotes matching '%s' found.", line.Nick, line.Args[1]))
		return
	}

	// TODO(fluffle): qd should take care of updating Accessed internally
	quote.Accessed++
	if err := qd.Update(bson.M{"_id": quote.Id}, quote); err != nil {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: I failed to update quote #%d: %s",
			line.Nick, quote.QID, err))
	}
	bot.Conn.Privmsg(line.Args[0], fmt.Sprintf("#%d: %s",
		quote.QID, quote.Quote))
}
