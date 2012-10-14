package quotedriver

import (
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/quotes"
	"labix.org/v2/mgo/bson"
	"strconv"
)

func add(line *base.Line) {
	n, c := line.Storable()
	quote := quotes.NewQuote(line.Args[1], n, c)
	quote.QID = qc.NewQID()
	if err := qc.Insert(quote); err == nil {
		bot.ReplyN(line, "Quote added succesfully, id #%d.", quote.QID)
	} else {
		bot.ReplyN(line, "Error adding quote: %s.", err)
	}
}

func del(line *base.Line) {
	// Strip optional # before qid
	if line.Args[1][0] == '#' {
		line.Args[1] = line.Args[1][1:]
	}
	qid, err := strconv.Atoi(line.Args[1])
	if err != nil {
		bot.ReplyN(line, "'%s' doesn't look like a quote id.", line.Args[1])
		return
	}
	if quote := qc.GetByQID(qid); quote != nil {
		if err := qc.RemoveId(quote.Id); err == nil {
			bot.ReplyN(line, "I forgot quote #%d: %s", qid, quote.Quote)
		} else {
			bot.ReplyN(line, "I failed to forget quote #%d: %s", qid, err)
		}
	} else {
		bot.ReplyN(line, "No quote found for id %d", qid)
	}
}

func fetch(line *base.Line) {
	if RateLimit(line.Nick) {
		return
	}
	qid, err := strconv.Atoi(line.Args[1])
	if err != nil {
		bot.ReplyN(line, "'%s' doesn't look like a quote id.", line.Args[1])
		return
	}
	quote := qc.GetByQID(qid)
	if quote != nil {
		bot.Reply(line, "#%d: %s", quote.QID, quote.Quote)
	} else {
		bot.ReplyN(line, "No quote found for id %d", qid)
	}
}

func lookup(line *base.Line) {
	if RateLimit(line.Nick) {
		return
	}
	quote := qc.GetPseudoRand(line.Args[1])
	if quote == nil {
		bot.ReplyN(line, "No quotes matching '%s' found.", line.Args[1])
		return
	}

	// TODO(fluffle): qd should take care of updating Accessed internally
	quote.Accessed++
	if err := qc.Update(bson.M{"_id": quote.Id}, quote); err != nil {
		bot.ReplyN(line, "I failed to update quote #%d: %s", quote.QID, err)
	}
	bot.Reply(line, "#%d: %s", quote.QID, quote.Quote)
}
