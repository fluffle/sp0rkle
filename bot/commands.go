package bot

import (
	"strings"
)

func (b *Bot) ignore(ctx *Context) {
	fields := strings.Fields(ctx.Text())
	if len(fields) == 0 {
		return
	}
	nick := strings.ToLower(fields[0])
	b.ignoreNs.String(nick, "ignore")
	ctx.ReplyN("I'll ignore '%s'.", nick)
}

func (b *Bot) unignore(ctx *Context) {
	fields := strings.Fields(ctx.Text())
	if len(fields) == 0 {
		return
	}
	nick := strings.ToLower(fields[0])
	b.ignoreNs.Delete(nick)
	ctx.ReplyN("No longer ignoring '%s'.", nick)
}
