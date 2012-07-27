package quotedriver

import (
	"fmt"
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/lib/quotes"
	"github.com/fluffle/sp0rkle/lib/util"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
	"labix.org/v2/mgo/bson"
	//	"rand"
	"strconv"
	"strings"
	//	"time"
)

func (qd *quoteDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(qd_privmsg), "bot_privmsg")
	//	r.AddHandler(bot.NewHandler(qd_action), "bot_action")
}

func qd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	qd := bot.GetDriver(driverName).(*quoteDriver)

	if !line.Addressed {
		return
	}

	nl := line.Copy()
	switch {
	// Quote add: qadd | quote add | add quote
	case util.StripAnyPrefix(&nl.Args[1], []string{"quote add ", "qadd ", "add quote "}):
		qd_add(bot, qd, nl)
	// Quote delete: qdel | quote del | del quote  #?QID
	case util.StripAnyPrefix(&nl.Args[1], []string{"quote del ", "qdel ", "del quote "}):
		// Strip optional # before qid
		if nl.Args[1][0] == '#' {
			nl.Args[1] = nl.Args[1][1:]
		}
		qd_delete(bot, qd, nl)
	// Quote lookup: quote #QID
	case util.StripAnyPrefix(&nl.Args[1], []string{"quote #"}):
		qd_fetch(bot, qd, nl)
	// Quote lookup: quote | quote regex
	case strings.ToLower(nl.Args[1]) == "quote":
		nl.Args[1] = ""
		fallthrough
	// This needs to come after the other cases as it will strip just "quote "
	case util.StripAnyPrefix(&nl.Args[1], []string{"quote "}):
		qd_lookup(bot, qd, nl)
	}
}

func qd_add(bot *bot.Sp0rkle, qd *quoteDriver, line *base.Line) {
	n := db.StorableNick{line.Nick, line.Ident, line.Host}
	c := db.StorableChan{line.Args[0]}
	quote := quotes.NewQuote(line.Args[1], n, c)
	quote.QID = qd.NewQID()
	if err := qd.Insert(quote); err == nil {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: Quote added succesfully, id #%d.", line.Nick, quote.QID))
	} else {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf("Oh no! %s.", err))
	}
}

func qd_delete(bot *bot.Sp0rkle, qd *quoteDriver, line *base.Line) {
	qid, err := strconv.Atoi(line.Args[1])
	if err != nil {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: '%s' doesn't look like a quote id.", line.Nick, line.Args[1]))
		return
	}
	if quote := qd.GetByQID(qid); quote != nil {
		if err := qd.Remove(bson.M{"_id": quote.Id}); err == nil {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I forgot quote #%d: %s", line.Nick, qid, quote.Quote))
		} else {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I failed to forget quote #%d: %s", line.Nick, qid, err))
		}
	} else {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: No quote found for id %d", line.Nick, qid))
	}
}

func qd_fetch(bot *bot.Sp0rkle, qd *quoteDriver, line *base.Line) {
	if qd.rateLimit(line.Nick) {
		return
	}
	qid, err := strconv.Atoi(line.Args[1])
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
	if qd.rateLimit(line.Nick) {
		return
	}
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
