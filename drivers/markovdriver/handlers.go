package markovdriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/conf"
)

func shouldMarkov(nick string) bool {
	return conf.Ns(markovNs).String(nick) != ""
}

func recordMarkov(ctx *bot.Context) {
	if !ctx.Addressed && ctx.Public() && shouldMarkov(ctx.Nick) {
		// Only markov lines that are public, not addressed to us,
		// and from markov-enabled nicks
		mc.AddSentence(ctx.Text(), "user:"+ctx.Nick)
	}
}
