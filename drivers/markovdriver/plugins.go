package markovdriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util"
	chain "github.com/fluffle/sp0rkle/util/markov"
)

func insultPlugin(in string, ctx *bot.Context) string {
	f := func(string) string {
		source := mc.Source("tag:insult")
		if insult, err := chain.Sentence(source); err == nil {
			return insult
		}
		return "<plugin error>"
	}
	return util.ApplyPluginFunction(in, "insult", f)
}
