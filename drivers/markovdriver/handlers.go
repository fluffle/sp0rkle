package markovdriver

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/conf"
	"strings"
)

func shouldMarkov(nick string) bool {
	return conf.Ns(markovNs).String(nick) != ""
}

func recordMarkov(ctx *bot.Context) {
	whom := strings.ToLower(ctx.Nick)
	if !ctx.Addressed && ctx.Public() && shouldMarkov(whom) {
		// Only markov lines that are public, not addressed to us,
		// and from markov-enabled nicks
		switch ctx.Cmd {
		case client.PRIVMSG:
			mc.AddSentence(ctx.Text(), "user:"+whom)
		case client.ACTION:
			mc.AddAction(ctx.Text(), "user:"+whom)
		}
	}
}
