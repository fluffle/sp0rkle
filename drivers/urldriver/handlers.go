package urldriver

import (
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/urls"
	"github.com/fluffle/sp0rkle/util"
	"strings"
	"time"
)

func urlScan(line *base.Line) {
	words := strings.Split(line.Args[1], " ")
	n, c := line.Storable()
	for _, w := range words {
		if util.LooksURLish(w) {
			if u := uc.GetByUrl(w); u != nil {
				bot.Reply(line, "%s first mentioned by %s at %s",
					w, u.Nick, u.Timestamp.Format(time.RFC1123))
				continue
			}
			u := urls.NewUrl(w, n, c)
			if len(w) > autoShortenLimit {
				u.Shortened = Encode(w)
			}
			if err := uc.Insert(u); err != nil {
				bot.ReplyN(line, "Couldn't insert url '%s': %s", w, err)
				continue
			}
			if u.Shortened != "" {
				bot.Reply(line, "%s's URL shortened as %s%s%s",
					line.Nick, bot.HttpHost(), shortenPath, u.Shortened)
			}
			lastseen[line.Args[0]] = u.Id
		}
	}
}
