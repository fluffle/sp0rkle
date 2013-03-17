package urldriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/urls"
	"github.com/fluffle/sp0rkle/util"
	"strings"
	"time"
)

func find(ctx *bot.Context) {
	if u := uc.GetRand(ctx.Text()); u != nil {
		ctx.ReplyN("%s", u)
	}
}

func shorten(ctx *bot.Context) {
	var u *urls.Url
	if ctx.Text() == "" {
		// assume we have been given "shorten that"
		if u = uc.GetById(lastseen[ctx.Target()]); u == nil {
			ctx.ReplyN("I seem to have forgotten what to shorten")
			return
		}
		if u.Shortened != "" {
			ctx.ReplyN("That was already shortened as %s%s%s",
				bot.HttpHost(), shortenPath, u.Shortened)

			return
		}
	} else {
		url := strings.TrimSpace(ctx.Text())
		if idx := strings.Index(url, " "); idx != -1 {
			url = url[:idx]
		}
		if !util.LooksURLish(url) {
			ctx.ReplyN("'%s' doesn't look URLish", url)
			return
		}
		if u = uc.GetByUrl(url); u == nil {
			n, c := ctx.Storable()
			u = urls.NewUrl(url, n, c)
		} else if u.Shortened != "" {
			ctx.ReplyN("That was already shortened as %s%s%s",
				bot.HttpHost(), shortenPath, u.Shortened)

			return
		}
	}
	if err := Shorten(u); err != nil {
		ctx.ReplyN("Failed to store shortened url: %s", err)
		return
	}
	ctx.ReplyN("%s shortened to %s%s%s",
		u.Url, bot.HttpHost(), shortenPath, u.Shortened)

}

func cache(ctx *bot.Context) {
	var u *urls.Url
	if ctx.Text() == "" {
		// assume we have been given "cache that"
		if u = uc.GetById(lastseen[ctx.Target()]); u == nil {
			ctx.ReplyN("I seem to have forgotten what to cache")
			return
		}
		if u.CachedAs != "" {
			ctx.ReplyN("That was already cached as %s%s%s at %s",
				bot.HttpHost(), cachePath, u.CachedAs,
				u.CacheTime.Format(time.RFC1123))

			return
		}
	} else {
		url := strings.TrimSpace(ctx.Text())
		if idx := strings.Index(url, " "); idx != -1 {
			url = url[:idx]
		}
		if !util.LooksURLish(url) {
			ctx.ReplyN("'%s' doesn't look URLish", url)
			return
		}
		if u = uc.GetByUrl(url); u == nil {
			n, c := ctx.Storable()
			u = urls.NewUrl(url, n, c)
		} else if u.CachedAs != "" {
			ctx.ReplyN("That was already cached as %s%s%s at %s",
				bot.HttpHost(), cachePath, u.CachedAs,
				u.CacheTime.Format(time.RFC1123))

			return
		}
	}
	if err := Cache(u); err != nil {
		ctx.ReplyN("Failed to store cached url: %s", err)
		return
	}
	ctx.ReplyN("%s cached as %s%s%s",
		u.Url, bot.HttpHost(), cachePath, u.CachedAs)

}
