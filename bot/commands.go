package bot

import (
	"github.com/fluffle/sp0rkle/collections/conf"
	"strings"
)

const ignoreNs = "ignore"

func ignore(ctx *Context) {
	nick := strings.ToLower(strings.Fields(ctx.Text())[0])
	if nick == "" {
		return
	}
	conf.Ns(ignoreNs).String(nick, "ignore")
	ctx.ReplyN("I'll ignore '%s'.", nick)
}

func unignore(ctx *Context) {
	nick := strings.ToLower(strings.Fields(ctx.Text())[0])
	if nick == "" {
		return
	}
	conf.Ns(ignoreNs).Delete(nick)
	ctx.ReplyN("No longer ignoring '%s'.", nick)
}
