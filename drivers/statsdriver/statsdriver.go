package statsdriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/stats"
)

var sc *stats.Collection

func Init() {
	sc = stats.Init()

	bot.HandleFunc(recordStats, "privmsg", "action")

	bot.CommandFunc(statsCmd, "lines", "lines [nick]  -- "+
		"display how many lines you [or nick] has said in the channel")
	bot.CommandFunc(statsCmd, "stats", "stats [nick]  -- "+
		"display how many lines you [or nick] has said in the channel")
	bot.CommandFunc(topten, "topten", "topten  -- "+
		"display the nicks who have said the most in the channel")
	bot.CommandFunc(topten, "top10", "top10  -- "+
		"display the nicks who have said the most in the channel")
}
