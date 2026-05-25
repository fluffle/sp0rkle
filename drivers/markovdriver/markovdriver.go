package markovdriver

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db/conf"
	"github.com/fluffle/sp0rkle/collections/markov"
)

const markovNsName = "markov"

type Driver struct {
	mc     *markov.Collection
	confNs conf.Namespace
}


func New(b *bot.Bot, mc *markov.Collection, config *conf.Registry) *Driver {
	d := &Driver{mc: mc, confNs: config.Ns(markovNsName)}

	b.Handle(d.recordMarkov, client.PRIVMSG, client.ACTION)
	b.Rewrite(d.insultPlugin)

	b.Command(d.enableMarkov, "markov me", "markov me  -- "+
		"Enable recording of your public messages to generate chains.")
	b.Command(d.disableMarkov, "don't markov me", "don't markov me  -- "+
		"Disable (and delete) recording of your public messages.")
	b.Command(d.disableMarkov, "don't markov me, bro", "don't markov me  -- "+
		"Disable (and delete) recording of your public messages.")
	b.Command(d.randomCmd, "markov", "markov <nick>  -- "+
		"Generate random sentence for given <nick>.")
	b.Command(d.insult, "insult", "insult <nick>  -- Insult <nick> at random.")
	b.Command(d.learn, "learn", "learn <tag> <sentence>  -- "+
		"Learns a sentence for a particular.")

	return d
}
