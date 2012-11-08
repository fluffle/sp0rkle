package urldriver

import (
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/urls"
	"github.com/fluffle/sp0rkle/util"
	"strings"
	"time"
)

func find(line *base.Line) {
	if u := uc.GetRand(line.Args[1]); u != nil {
		bot.ReplyN(line, u.String())
	}
}

func shorten(line *base.Line) {
	var u *urls.Url
	if line.Args[1] == "" {
		// assume we have been given "shorten that"
		if u = uc.GetById(lastseen[line.Args[0]]); u == nil {
			bot.ReplyN(line, "I seem to have forgotten what to shorten")
			return
		}
		if u.Shortened != "" {
			bot.ReplyN(line, "That was already shortened as %s%s%s",
				bot.HttpHost(), shortenPath, u.Shortened)
			return
		}
	} else {
		url := strings.TrimSpace(line.Args[1])
		if idx := strings.Index(url, " "); idx != -1 {
			url = url[:idx]
		}
		if !util.LooksURLish(url) {
			bot.ReplyN(line, "'%s' doesn't look URLish", url)
			return
		}
		if u = uc.GetByUrl(url); u == nil {
			n, c := line.Storable()
			u = urls.NewUrl(url, n, c)
		} else if u.Shortened != "" {
			bot.ReplyN(line, "That was already shortened as %s%s%s",
				bot.HttpHost(), shortenPath, u.Shortened)
			return
		}
	}
	if err := Shorten(u); err != nil {
		bot.ReplyN(line, "Failed to store shortened url: %s", err)
		return
	}
	bot.ReplyN(line, "%s shortened to %s%s%s",
		u.Url, bot.HttpHost(), shortenPath, u.Shortened)
}

func cache(line *base.Line) {
	var u *urls.Url
	if line.Args[1] == "" {
		// assume we have been given "cache that"
		if u = uc.GetById(lastseen[line.Args[0]]); u == nil {
			bot.ReplyN(line, "I seem to have forgotten what to cache")
			return
		}
		if u.CachedAs != "" {
			bot.ReplyN(line, "That was already cached as %s%s%s at %s",
				bot.HttpHost(), cachePath, u.CachedAs,
				u.CacheTime.Format(time.RFC1123))
			return
		}
	} else {
		url := strings.TrimSpace(line.Args[1])
		if idx := strings.Index(url, " "); idx != -1 {
			url = url[:idx]
		}
		if !util.LooksURLish(url) {
			bot.ReplyN(line, "'%s' doesn't look URLish", url)
			return
		}
		if u = uc.GetByUrl(url); u == nil {
			n, c := line.Storable()
			u = urls.NewUrl(url, n, c)
		} else if u.CachedAs != "" {
			bot.ReplyN(line, "That was already cached as %s%s%s at %s",
				bot.HttpHost(), cachePath, u.CachedAs,
				u.CacheTime.Format(time.RFC1123))
			return
		}
	}
	if err := Cache(u); err != nil {
		bot.ReplyN(line, "Failed to store cached url: %s", err)
		return
	}
	bot.ReplyN(line, "%s cached as %s%s%s",
		u.Url, bot.HttpHost(), cachePath, u.CachedAs)
}
