package urldriver

import (
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/sp0rkle/lib/urls"
	"github.com/fluffle/sp0rkle/lib/util"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
	"strings"
	"time"
)

func (ud *urlDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(ud_privmsg), "bot_privmsg")
}

func ud_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	ud := bot.GetDriver(driverName).(*urlDriver)

	// If we're not being addressed directly, short-circuit to scan. 
	if !line.Addressed {
		ud_scan(bot, ud, line)
		return
	}

	nl := line.Copy()
	switch {
		case util.StripAnyPrefix(&nl.Args[1],
			[]string{"url find ", "urlfind ", "url search ", "urlsearch "}):
			ud_find(bot, ud, nl)
		case util.HasAnyPrefix(nl.Args[1], []string{"random url", "randurl"}):
			nl.Args[1] = ""
			ud_find(bot, ud, nl)
		case util.StripAnyPrefix(&nl.Args[1],
			[]string{"shorten that", "shorten"}):
			ud_shorten(bot, ud, nl)
		case util.StripAnyPrefix(&nl.Args[1],
			[]string{"cache that", "cache "}):
			ud_cache(bot, ud, nl)
		default:
			ud_scan(bot, ud, line)
	}
}

func ud_scan(bot *bot.Sp0rkle, ud *urlDriver, line *base.Line) {
	words := strings.Split(line.Args[1], " ")
	n, c := line.Storable()
	for _, w := range words {
		if util.LooksURLish(w) {
			if u := ud.GetByUrl(w); u != nil {
				bot.Reply(line, "%s first mentioned by %s at %s",
					w, u.Nick, u.Timestamp.Format(time.RFC1123))
				continue
			}
			u := urls.NewUrl(w, n, c)
			if len(w) > autoShortenLimit {
				u.Shortened = ud.Encode(w)
			}
			if err := ud.Insert(u); err != nil {
				bot.ReplyN(line, "Couldn't insert url '%s': %s", w, err)
				continue
			}
			if u.Shortened != "" {
				bot.Reply(line, "%s's URL shortened as %s%s%s",
					line.Nick, bot.Prefix, shortenPath, u.Shortened)
			}
			ud.lastseen[line.Args[0]] = u.Id
		}
	}
}

func ud_find(bot *bot.Sp0rkle, ud *urlDriver, line *base.Line) {
	if u := ud.GetRand(line.Args[1]); u != nil {
		bot.ReplyN(line, "%s", u)
	}
}

func ud_shorten(bot *bot.Sp0rkle, ud *urlDriver, line *base.Line) {
	var u *urls.Url
	if line.Args[1] == "" {
		// assume we have been given "shorten that"
		if u = ud.GetById(ud.lastseen[line.Args[0]]); u == nil {
			bot.ReplyN(line, "I seem to have forgotten what to shorten")
			return
		}
		if u.Shortened != "" {
			bot.ReplyN(line, "That was already shortened as %s%s%s",
				bot.Prefix, shortenPath, u.Shortened)
		}
	} else {
		url := strings.TrimSpace(line.Args[1])
		if idx := strings.Index(line.Args[1], " "); idx != -1 {
			url = url[:idx]
		}
		if !util.LooksURLish(url) {
			bot.ReplyN(line, "'%s' doesn't look URLish", url)
			return
		}
		n, c := line.Storable()
		u = urls.NewUrl(url, n, c)
	}
	if err := ud.Shorten(u); err != nil {
		bot.ReplyN(line, "Failed to store shortened url: %s", err)
		return
	}
	bot.ReplyN(line, "%s shortened to %s%s%s",
		u.Url, bot.Prefix, shortenPath, u.Shortened)
}

func ud_cache(bot *bot.Sp0rkle, ud *urlDriver, line *base.Line) {
	var u *urls.Url
	if line.Args[1] == "" {
		// assume we have been given "cache that"
		if u = ud.GetById(ud.lastseen[line.Args[0]]); u == nil {
			bot.ReplyN(line, "I seem to have forgotten what to cache")
			return
		}
		if u.CachedAs != "" {
			bot.ReplyN(line, "That was already cached as %s%s%s at %s",
			bot.Prefix, cachePath, u.CachedAs,
			u.CacheTime.Format(time.RFC1123))
		}
	} else {
		url := strings.TrimSpace(line.Args[1])
		if idx := strings.Index(line.Args[1], " "); idx != -1 {
			url = url[:idx]
		}
		if !util.LooksURLish(url) {
			bot.ReplyN(line, "'%s' doesn't look URLish", url)
			return
		}
		n, c := line.Storable()
		u = urls.NewUrl(url, n, c)
	}
	if err := ud.Cache(u); err != nil {
		bot.ReplyN(line, "Failed to store cached url: %s", err)
		return
	}
	bot.ReplyN(line, "%s cached as %s%s%s",
		u.Url, bot.Prefix, cachePath, u.CachedAs)
}
