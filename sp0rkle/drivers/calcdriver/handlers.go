package calcdriver

import (
	"fmt"
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/sp0rkle/lib/calc"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
	"strings"
)

func (cd *calcDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(cd_privmsg), "bot_privmsg")
}

func cd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	if !line.Addressed {
		return
	}
	cd := bot.GetDriver(driverName).(*calcDriver)

	if strings.HasPrefix(line.Args[1], "calc ") {
		cd.l.Info("calculating %s", line.Args[1][5:])
		if num, err := calc.Calc(line.Args[1][5:]); err == nil {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: %g", line.Nick, num))
		} else {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: %s", line.Nick, err))
		}
	}
}
