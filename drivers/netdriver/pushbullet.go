 package netdriver

import (
	"strings"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util/push"
)

func pushEnable(ctx *bot.Context) {
	ctx.Privmsg(ctx.Nick, "Hi! Visit the following URL while logged into "+
		"the account you want to use to push to your device.")
	ctx.Privmsg(ctx.Nick, push.GenAuthURL(ctx.Nick))
}

func pushDisable(ctx *bot.Context) {
	push.StopFor(ctx.Nick)
	ctx.Privmsg(ctx.Nick, "Ok, pushes disabled.")
}

func pushAuth(ctx *bot.Context) {
	s := strings.Fields(ctx.Text())
	if err := push.StartFor(ctx.Nick, s[0]); err != nil {
		ctx.Privmsg(ctx.Nick, err.Error())
		return
	}
	ctx.Privmsg(ctx.Nick, "Pushes enabled! Yay!")
}
