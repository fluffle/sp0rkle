package decisiondriver

import (
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util"
	"strings"
)

func decideCmd(line *base.Line) {
	opts := splitDelimitedString(line.Args[1])
	chosen := strings.TrimSpace(opts[util.RNG.Intn(len(opts))])
	bot.ReplyN(line, chosen)
}

func randCmd(line *base.Line) {
	bot.ReplyN(line, "%s", randomFloatAsString(line.Args[1], util.RNG))
}
