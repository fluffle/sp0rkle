package markovdriver

import (
	"strings"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
)

func (d *Driver) shouldMarkov(nick string) bool {
	return d.confNs.String(nick) != ""
}

func (d *Driver) recordMarkov(ctx *bot.Context) {
	whom := strings.ToLower(ctx.Nick)
	if !ctx.Addressed && ctx.Public() && d.shouldMarkov(whom) {
		// Only markov lines that are public, not addressed to us,
		// and from markov-enabled nicks
		switch ctx.Cmd {
		case client.PRIVMSG:
			d.mc.AddSentence(ctx.Text(), "user:"+whom)
		case client.ACTION:
			d.mc.AddAction(ctx.Text(), "user:"+whom)
		}
	}
}
