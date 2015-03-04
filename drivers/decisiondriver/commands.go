package decisiondriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"math/rand"
	"strings"
)

func decideCmd(ctx *bot.Context) {
	opts := splitDelimitedString(ctx.Text())
	chosen := strings.TrimSpace(opts[rand.Intn(len(opts))])
	ctx.ReplyN("%s", chosen)
}

func randCmd(ctx *bot.Context) {
	ctx.ReplyN("%s", randomFloatAsString(ctx.Text()))
}
