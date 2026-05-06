package bot

import (
	"strings"

	"github.com/fluffle/sp0rkle/collections/conf"
)

const ignoreNs = "ignore"

func ignore(ctx *Context) {
	fields := strings.Fields(ctx.Text())
	if len(fields) == 0 {
		return
	}
	nick := strings.ToLower(fields[0])
	conf.Ns(ignoreNs).String(nick, "ignore")
	ctx.ReplyN("I'll ignore '%s'.", nick)
}

func unignore(ctx *Context) {
	fields := strings.Fields(ctx.Text())
	if len(fields) == 0 {
		return
	}
	nick := strings.ToLower(fields[0])
	conf.Ns(ignoreNs).Delete(nick)
	ctx.ReplyN("No longer ignoring '%s'.", nick)
}
