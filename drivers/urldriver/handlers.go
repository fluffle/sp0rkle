package urldriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/urls"
	"github.com/fluffle/sp0rkle/util"
	"strings"
	"time"
)

func urlScan(ctx *bot.Context) {
	words := strings.Split(ctx.Text(), " ")
	n, c := ctx.Storable()
	for _, w := range words {
		if util.LooksURLish(w) {
			if u := uc.GetByUrl(w); u != nil &&
				u.Nick != bot.Nick(ctx.Nick) &&
				time.Since(u.Timestamp) > 2*time.Hour {
				ctx.Reply("that URL first mentioned by %s %s ago",
					u.Nick, util.TimeSince(u.Timestamp))

				continue
			}
			u := urls.NewUrl(w, n, c)
			if len(w) > autoShortenLimit {
				u.Shortened = Encode(w)
			}
			if err := uc.Insert(u); err != nil {
				ctx.ReplyN("Couldn't insert url '%s': %s", w, err)
				continue
			}
			if u.Shortened != "" {
				ctx.Reply("%s's URL shortened as %s%s%s",
					ctx.Nick, bot.HttpHost(), shortenPath, u.Shortened)
			}
			lastseen[ctx.Target()] = u.Id
		}
	}
}
