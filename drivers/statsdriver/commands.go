package statsdriver

import (
	"fmt"
	"github.com/fluffle/sp0rkle/bot"
	"strings"
)

func statsCmd(ctx *bot.Context) {
	n := ctx.Nick
	if len(ctx.Text()) > 0 {
		n = ctx.Text()
	}
	ns := sc.StatsFor(n, ctx.Target())
	if ns != nil {
		ctx.ReplyN("%s", ns)
	}
}

func topten(ctx *bot.Context) {
	top := sc.TopTen(ctx.Target())
	s := make([]string, 0, 10)
	for i, n := range top {
		s = append(s, fmt.Sprintf("#%d: %s - %d", i+1, n.Nick, n.Lines))
	}
	ctx.Reply("%s", strings.Join(s, ", "))
}
