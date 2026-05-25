package quotedriver

import (
	"strconv"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/quotes"
)

func (d *Driver) add(ctx *bot.Context) {
	n, c := ctx.Storable()
	quote := quotes.NewQuote(ctx.Text(), n, c)
	var err error
	if quote.QID, err = d.qc.NewQID(); err != nil {
		ctx.ReplyN("Retrieving new quote ID failed: %v", err)
		return
	}
	if err = d.qc.Put(quote); err == nil {
		ctx.ReplyN("Quote added succesfully, id #%d.", quote.QID)
	} else {
		ctx.ReplyN("Error adding quote: %s.", err)
	}
}

func (d *Driver) del(ctx *bot.Context) {
	txt := ctx.Text()
	// Strip optional # before qid
	if len(txt) > 0 && txt[0] == '#' {
		txt = txt[1:]
	}
	qid, err := strconv.Atoi(txt)
	if err != nil {
		ctx.ReplyN("'%s' doesn't look like a quote id.", ctx.Text())
		return
	}
	if quote := d.qc.GetByQID(qid); quote != nil {
		if err := d.qc.Del(quote); err == nil {
			ctx.ReplyN("I forgot quote #%d: %s", qid, quote.Quote)
		} else {
			ctx.ReplyN("I failed to forget quote #%d: %s", qid, err)
		}
	} else {
		ctx.ReplyN("No quote found for id %d", qid)
	}
}

func (d *Driver) fetch(ctx *bot.Context) {
	if RateLimit(ctx.Nick) {
		return
	}
	qid, err := strconv.Atoi(ctx.Text())
	if err != nil {
		ctx.ReplyN("'%s' doesn't look like a quote id.", ctx.Text())
		return
	}
	quote := d.qc.GetByQID(qid)
	if quote != nil {
		ctx.Reply("#%d: %s", quote.QID, quote.Quote)
	} else {
		ctx.ReplyN("No quote found for id %d", qid)
	}
}

func (d *Driver) lookup(ctx *bot.Context) {
	if RateLimit(ctx.Nick) {
		return
	}
	quote := d.qc.GetPseudoRand(ctx.Text())
	if quote == nil {
		ctx.ReplyN("No quotes matching '%s' found.", ctx.Text())
		return
	}

	// TODO(fluffle): qd should take care of updating Accessed internally
	quote.Accessed++
	if err := d.qc.Put(quote); err != nil {
		ctx.ReplyN("I failed to update quote #%d: %s", quote.QID, err)
	}
	ctx.Reply("#%d: %s", quote.QID, quote.Quote)
}
