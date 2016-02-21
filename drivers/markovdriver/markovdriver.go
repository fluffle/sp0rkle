package markovdriver

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/markov"
)

const markovNs = "markov"

var mc *markov.Collection

func Init() {
	mc = markov.Init()

	bot.Handle(recordMarkov, client.PRIVMSG, client.ACTION)
	bot.Rewrite(insultPlugin)

	bot.Command(enableMarkov, "markov me", "markov me  -- "+
		"Enable recording of your public messages to generate chains.")
	bot.Command(disableMarkov, "don't markov me", "don't markov me  -- "+
		"Disable (and delete) recording of your public messages.")
	bot.Command(disableMarkov, "don't markov me, bro", "don't markov me  -- "+
		"Disable (and delete) recording of your public messages.")
	bot.Command(randomCmd, "markov", "markov <nick>  -- "+
		"Generate random sentence for given <nick>.")
	bot.Command(insult, "insult", "insult <nick>  -- Insult <nick> at random.")
	bot.Command(learn, "learn", "learn <tag> <sentence>  -- "+
		"Learns a sentence for a particular.")
}
