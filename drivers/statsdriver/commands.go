package statsdriver

import (
	"fmt"
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"strings"
)

func statsCmd(line *base.Line) {
	n := line.Nick
	if len(line.Args[1]) > 0 {
		n = line.Args[1]
	}
	ns := sc.StatsFor(n, line.Args[0])
	if ns != nil {
		bot.ReplyN(line, "%s", ns)
	}
}

func topten(line *base.Line) {
	top := sc.TopTen(line.Args[0])
	s := make([]string, 0, 10)
	for i, n := range top {
		s = append(s, fmt.Sprintf("#%d: %s - %d", i+1, n.Nick, n.Lines))
	}
	bot.Reply(line, "%s", strings.Join(s, ", "))
}
