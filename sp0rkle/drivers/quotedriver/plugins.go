package quotedriver

import (
	"github.com/fluffle/sp0rkle/lib/quotes"
	"github.com/fluffle/sp0rkle/lib/util"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"strconv"
)

func (qd *quoteDriver) RegisterPlugins(pm base.PluginManager) {
	pm.AddPlugin(&QuotePlugin{qd, qd_plugin_lookup})
}

type QuotePlugin struct {
	provider  *quoteDriver
	processor func(*quoteDriver, string, *base.Line) string
}

func (qp *QuotePlugin) Apply(val string, line *base.Line) string {
	return qp.processor(qp.provider, val, line)
}

func qd_plugin_lookup(qd *quoteDriver, val string, line *base.Line) string {
	f := func(s string) string {
		var quote *quotes.Quote
		if s == "" {
			quote = qd.GetPseudoRand("")
		} else if s[0] == '#' {
			if qid, err := strconv.Atoi(s[1:]); err == nil {
				quote = qd.GetByQID(qid)
			}
		} else {
			quote = qd.GetPseudoRand(s)
		}
		if quote == nil {
			return "<plugin error>"
		}
		return quote.Quote
	}
	return util.ApplyPluginFunction(val, "quote", f)
}
