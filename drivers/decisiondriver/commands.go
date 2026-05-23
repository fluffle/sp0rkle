package decisiondriver

import (
	"math/rand"
	"strings"

	"github.com/fluffle/sp0rkle/bot"
)

func decideCmd(ctx *bot.Context) {
	opts, err := splitDelimitedString(ctx.Text())
	if err != nil {
		ctx.ReplyN("I can't decide: %v", err)
		return
	}
	if len(opts) == 0 {
		ctx.ReplyN("I need something to choose from...", err)
		return
	}
	chosen := strings.TrimSpace(opts[rand.Intn(len(opts))])
	ctx.ReplyN("%s", chosen)
}

func randCmd(ctx *bot.Context) {
	ctx.ReplyN("%s", randomFloatAsString(ctx.Text()))
}
