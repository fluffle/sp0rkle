package factdriver

import (
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"strings"
	"time"
)

// Replicate perlfu's $<stuff> identifiers
func replaceIdentifiers(in string, line *base.Line) string {
	return id_replacer(in, line, time.Now())
}

// Split this out so we can inject a deterministic time for testing.
func id_replacer(val string, line *base.Line, ts time.Time) string {
	val = strings.Replace(val, "$nick", line.Nick, -1)
	val = strings.Replace(val, "$chan", line.Args[0], -1)
	val = strings.Replace(val, "$username", line.Ident, -1)
	val = strings.Replace(val, "$user", line.Ident, -1)
	val = strings.Replace(val, "$host", line.Host, -1)
	val = strings.Replace(val, "$date", ts.Format(time.ANSIC), -1)
	val = strings.Replace(val, "$time", ts.Format("15:04:05"), -1)
	return val
}
