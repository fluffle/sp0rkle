package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"fmt"
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/sp0rkle/lib/util"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
	"strings"
)

func (dd *decisionDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(dd_privmsg), "bot_privmsg")
}

func dd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	if !line.Addressed {
		return
	}

	switch {
	case strings.HasPrefix(line.Args[1], "decide "):
		opts := splitDelimitedString(line.Args[1][8:])
		chosen := strings.TrimSpace(opts[util.RNG.Intn(len(opts))])
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf("%s: %s", line.Nick, chosen))
	case strings.HasPrefix(line.Args[1], "rand "):
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: %s", line.Nick, randomFloatAsString(line.Args[1][5:])
	}
}
