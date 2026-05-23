package statsdriver

import (
	"fmt"
	"strings"

	"github.com/fluffle/sp0rkle/bot"
)

func statsCmd(ctx *bot.Context) {
	n := ctx.Nick
	if len(ctx.Text()) > 0 {
		n = ctx.Text()
	}
	ns := sc.StatsFor(n, ctx.Target())
	if ns != nil {
		ctx.ReplyN("%s", ns)
	} else {
		ctx.ReplyN("I've not seen %q before, sorry.", n)
	}
}

func topten(ctx *bot.Context) {
	top := sc.TopTen(ctx.Target())
	s := make([]string, len(top))
	for i, n := range top {
		s[i] = fmt.Sprintf("#%d: %s - %d", i+1, n.Nick, n.Lines)
	}
	ctx.Reply("%s", strings.Join(s, ", "))
}
