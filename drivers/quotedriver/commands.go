package quotedriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/quotes"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

func add(ctx *bot.Context) {
	n, c := ctx.Storable()
	quote := quotes.NewQuote(ctx.Text(), n, c)
	quote.QID = qc.NewQID()
	if err := qc.Insert(quote); err == nil {
		ctx.ReplyN("Quote added succesfully, id #%d.", quote.QID)
	} else {
		ctx.ReplyN("Error adding quote: %s.", err)
	}
}

func del(ctx *bot.Context) {
	txt := ctx.Text()
	// Strip optional # before qid
	if txt[0] == '#' {
		txt = txt[1:]
	}
	qid, err := strconv.Atoi(txt)
	if err != nil {
		ctx.ReplyN("'%s' doesn't look like a quote id.", ctx.Text())
		return
	}
	if quote := qc.GetByQID(qid); quote != nil {
		if err := qc.RemoveId(quote.Id); err == nil {
			ctx.ReplyN("I forgot quote #%d: %s", qid, quote.Quote)
		} else {
			ctx.ReplyN("I failed to forget quote #%d: %s", qid, err)
		}
	} else {
		ctx.ReplyN("No quote found for id %d", qid)
	}
}

func fetch(ctx *bot.Context) {
	if RateLimit(ctx.Nick) {
		return
	}
	qid, err := strconv.Atoi(ctx.Text())
	if err != nil {
		ctx.ReplyN("'%s' doesn't look like a quote id.", ctx.Text())
		return
	}
	quote := qc.GetByQID(qid)
	if quote != nil {
		ctx.Reply("#%d: %s", quote.QID, quote.Quote)
	} else {
		ctx.ReplyN("No quote found for id %d", qid)
	}
}

func lookup(ctx *bot.Context) {
	if RateLimit(ctx.Nick) {
		return
	}
	quote := qc.GetPseudoRand(ctx.Text())
	if quote == nil {
		ctx.ReplyN("No quotes matching '%s' found.", ctx.Text())
		return
	}

	// TODO(fluffle): qd should take care of updating Accessed internally
	quote.Accessed++
	if err := qc.Update(bson.M{"_id": quote.Id}, quote); err != nil {
		ctx.ReplyN("I failed to update quote #%d: %s", quote.QID, err)
	}
	ctx.Reply("#%d: %s", quote.QID, quote.Quote)
}
