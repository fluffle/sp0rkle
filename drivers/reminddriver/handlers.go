package reminddriver

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
)

func load(line *base.Line) {
	// We're connected to IRC, so load saved reminders
	r := rc.LoadAndPrune()
	for i := range r {
		if r[i] == nil {
			logging.Warn("Nil reminder %d from LoadAndPrune", i)
			continue
		}
		Remind(r[i])
	}
}

func tellCheck(line *base.Line) {
	nick := line.Nick
	if line.Cmd == "NICK" {
		// We want the destination nick, not the source.
		nick = line.Args[0]
	}
	r := rc.TellsFor(nick)
	for i := range r {
		if line.Cmd == "NICK" {
			bot.Privmsg(string(r[i].Chan), nick+": "+r[i].Reply())
			bot.Reply(line, "%s", r[i].Reply())
		} else {
			bot.Privmsg(line.Nick, r[i].Reply())
			bot.ReplyN(line, "%s", r[i].Reply())
		}
		rc.RemoveId(r[i].Id)
	}
	if len(r) > 0 {
		delete(listed, line.Nick)
	}
}
