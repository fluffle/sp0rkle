package urldriver

import (
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/lib/urls"
	"github.com/fluffle/sp0rkle/lib/util"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
	"strings"
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
		default:
			ud_scan(bot, ud, line)
	}
}

func ud_scan(bot *bot.Sp0rkle, ud *urlDriver, line *base.Line) {
	words := strings.Split(line.Args[1], " ")
	n := db.StorableNick{line.Nick, line.Ident, line.Host}
	c := db.StorableChan{line.Args[0]}
	for _, w := range words {
		if util.LooksURLish(w) {
			u := urls.NewUrl(w, n, c)
			if err := ud.Insert(u); err != nil {
				bot.ReplyN(line, "Couldn't insert url '%s': %s", w, err)
			}
		}
	}
}

func ud_find(bot *bot.Sp0rkle, ud *urlDriver, line *base.Line) {
	if u := ud.GetRand(line.Args[1]); u != nil {
		bot.ReplyN(line, "%s", u)
	}
}
