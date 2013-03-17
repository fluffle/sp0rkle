package decisiondriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util"
	"strings"
)

func decideCmd(ctx *bot.Context) {
	opts := splitDelimitedString(ctx.Text())
	chosen := strings.TrimSpace(opts[util.RNG.Intn(len(opts))])
	ctx.ReplyN("%s", chosen)
}

func randCmd(ctx *bot.Context) {
	ctx.ReplyN("%s", randomFloatAsString(ctx.Text(), util.RNG))
}
