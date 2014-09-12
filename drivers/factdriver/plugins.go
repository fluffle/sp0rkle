package factdriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util/datetime"
	"strings"
	"time"
)

// Replicate perlfu's $<stuff> identifiers
func replaceIdentifiers(in string, ctx *bot.Context) string {
	return id_replacer(in, ctx, time.Now())
}

// Split this out so we can inject a deterministic time for testing.
func id_replacer(val string, ctx *bot.Context, ts time.Time) string {
	val = strings.Replace(val, "$nick", ctx.Nick, -1)
	val = strings.Replace(val, "$chan", ctx.Target(), -1)
	val = strings.Replace(val, "$username", ctx.Ident, -1)
	val = strings.Replace(val, "$user", ctx.Ident, -1)
	val = strings.Replace(val, "$host", ctx.Host, -1)
	val = strings.Replace(val, "$date", datetime.Format(ts), -1)
	val = strings.Replace(val, "$time", datetime.Format(ts, "15:04:05"), -1)
	return val
}
