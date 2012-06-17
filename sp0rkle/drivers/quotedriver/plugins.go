package quotedriver

import (
	"github.com/fluffle/sp0rkle/lib/quotes"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"strconv"
	"strings"
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

// TODO(fluffle): add a util function to locate <plugin=SOMETHING [args...]>
// and use it here and in the decision driver...
func qd_plugin_lookup(qd *quoteDriver, val string, line *base.Line) string {
	for {
		var quote *quotes.Quote
		ps := strings.Index(val, "<plugin=quote")
		if ps == -1 {
			break
		}
		pe := strings.Index(val[ps:], ">")
		switch {
		case pe == -1:
			break
		case pe == 13:
			// <plugin=quote>
			quote = qd.GetPseudoRand("")
		case val[ps+14] == '#':
			// <plugin=quote #QID>
			if qid, err := strconv.Atoi(val[ps+15:ps+pe]); err == nil {
				quote = qd.GetByQID(qid)
			}
		default:
			// we have " some key to look up" between ps+14 and ps+pe
			quote = qd.GetPseudoRand(val[ps+14:ps+pe])
		}
		if quote == nil {
			continue
		}
		pe += ps
		val = val[:ps] + quote.Quote + val[pe+1:]
	}
	return val
}
