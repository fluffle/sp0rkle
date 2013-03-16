package karmadriver

import (
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
)

func karmaCmd(line *base.Line) {
	if k := kc.KarmaFor(line.Args[1]); k != nil {
		bot.ReplyN(line, "%s", k)
	} else {
		bot.ReplyN(line, "No karma found for '%s'", line.Args[1])
	}
}
