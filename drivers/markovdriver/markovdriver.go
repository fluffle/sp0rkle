package markovdriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/markov"
)

var mc *markov.Collection

func Init() {
	mc = markov.Init()

	bot.HandleFunc(recordMarkov, "privmsg")

	bot.CommandFunc(randomCmd, "markov", "markov <user>  -- "+
		"Generate random sentece for given <user>.")
}
