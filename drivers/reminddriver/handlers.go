package reminddriver

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
)

func (rd *remindDriver) Load(line *base.Line) {
	// We're connected to IRC, so load saved reminders
	r := rd.LoadAndPrune()
	for i := range r {
		if r[i] == nil {
			logging.Warn("Nil reminder %d from LoadAndPrune", i)
			continue
		}
		rd.Remind(r[i])
	}
}

func (rd *remindDriver) TellCheck(line *base.Line) {
	nick := line.Nick
	if line.Cmd == "NICK" {
		// We want the destination nick, not the source.
		nick = line.Args[0]
	}
	r := rd.TellsFor(nick)
	for i := range r {
		if line.Cmd == "NICK" {
			bot.Privmsg(string(r[i].Chan), nick + ": " + r[i].Reply())
			bot.Reply(line, r[i].Reply())
		} else {
			bot.Privmsg(line.Nick, r[i].Reply())
			bot.ReplyN(line, r[i].Reply())
		}
		rd.RemoveId(r[i].Id)
	}
	if len(r) > 0 {
		delete(rd.list, line.Nick)
	}
}
