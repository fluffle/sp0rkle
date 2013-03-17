package karmadriver

import (
	"github.com/fluffle/sp0rkle/bot"
)

func karmaCmd(ctx *bot.Context) {
	if k := kc.KarmaFor(ctx.Text()); k != nil {
		ctx.ReplyN("%s", k)
	} else {
		ctx.ReplyN("No karma found for '%s'", ctx.Text())
	}
}
