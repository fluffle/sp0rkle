package quotedriver

import (
//	"fmt"
	"github.com/fluffle/goevent/event"
//	"launchpad.net/gobson/bson"
//	"lib/db"
//	"lib/quotes"
//	"rand"
	"sp0rkle/bot"
	"sp0rkle/base"
//	"strings"
//	"strconv"
//	"time"
)

type QuoteHandler func(*bot.Sp0rkle, *quoteDriver, *base.Line)

// Unboxer for QuoteDriver handlers
func QDHandler(f QuoteHandler) event.Handler {
	return event.NewHandler(func(ev ...interface{}) {
		f(ev[0].(*bot.Sp0rkle), ev[1].(*quoteDriver), ev[2].(*base.Line))
	})
}

func (qd *quoteDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(qd_privmsg), "bot_privmsg")
//	r.AddHandler(bot.NewHandler(qd_action), "bot_action")
	r.AddHandler(QDHandler(qd_lookup), "qd_lookup")
}

func qd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	qd := bot.GetDriver(driverName).(*quoteDriver)
	
	if !line.Addressed {
		bot.Dispatch("qd_lookup", qd, line)
		return
	}
}

func qd_lookup(bot *bot.Sp0rkle, qd *quoteDriver, line *base.Line) {
	bot.Conn.Privmsg(line.Args[0], "Not implemented")
}
