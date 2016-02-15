package quotedriver

import (
	"strconv"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/quotes"
	"github.com/fluffle/sp0rkle/util"
)

func quotePlugin(in string, ctx *bot.Context) string {
	f := func(s string) string {
		var quote *quotes.Quote
		if s == "" {
			quote = qc.GetPseudoRand("")
		} else if s[0] == '#' {
			if qid, err := strconv.Atoi(s[1:]); err == nil {
				quote = qc.GetByQID(qid)
			}
		} else {
			quote = qc.GetPseudoRand(s)
		}
		if quote == nil {
			return "<plugin error>"
		}
		return quote.Quote
	}
	return util.ApplyPluginFunction(in, "quote", f)
}
