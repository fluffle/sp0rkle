package karmadriver

import (
	"strings"

	"github.com/fluffle/sp0rkle/bot"
)

func karmaCmd(ctx *bot.Context) {
	if strings.TrimSpace(ctx.Text()) == "" {
		ctx.ReplyN("karma karma karma karma, karma chameeleeoooooonnn")
		return
	}
	if k := kc.KarmaFor(ctx.Text()); k != nil {
		ctx.ReplyN("%s", k)
		return
	}
	ctx.ReplyN("No karma found for '%s'", ctx.Text())
}
