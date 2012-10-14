package netdriver

import (
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"strings"
)

func (nd *netDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(nd_privmsg), "bot_privmsg")
}

func nd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	nd := bot.GetDriver(driverName).(*netDriver)

	idx := strings.Index(line.Args[1], " ")
	if !line.Addressed || idx == -1 {
		return
	}
	svc, query := line.Args[1][:idx], line.Args[1][idx+1:]
	if s, ok := nd.services[svc]; ok {
		bot.ReplyN(line, "%s", s.LookupResult(query))
	}
}
