package reminddriver

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
)

func load(ctx *bot.Context) {
	// We're connected to IRC, so load saved reminders
	r := rc.LoadAndPrune()
	for i := range r {
		if r[i] == nil {
			logging.Warn("Nil reminder %d from LoadAndPrune", i)
			continue
		}
		Remind(r[i], ctx)
	}
}

func unload(ctx *bot.Context) {
	// We've been disconnected from IRC: stop all remind goroutines
	// since they will be restarted when we reconnect.
	for id, c := range running {
		c <- struct{}{}
		delete(running, id)
	}
}

func tellCheck(ctx *bot.Context) {
	nick := ctx.Nick
	if ctx.Cmd == client.NICK {
		// We want the destination nick, not the source.
		nick = ctx.Target()
	}
	r := rc.TellsFor(nick)
	for i := range r {
		if ctx.Cmd == client.NICK {
			if r[i].Chan != "" {
				ctx.Privmsg(string(r[i].Chan), nick+": "+r[i].Reply())
			}
			ctx.Reply("%s", r[i].Reply())
		} else {
			ctx.Privmsg(ctx.Nick, r[i].Reply())
			if r[i].Chan != "" {
				ctx.ReplyN("%s", r[i].Reply())
			}
		}
		rc.RemoveId(r[i].Id)
	}
	if len(r) > 0 {
		delete(listed, ctx.Nick)
	}
}
