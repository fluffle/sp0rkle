package factdriver

import (
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"strings"
	"time"
)

func (fd *factoidDriver) RegisterPlugins(pm base.PluginManager) {
	// pm == fd in this case, but meh.
	pm.Add(&FactoidPlugin{fd, fd_identifiers})
}

type FactoidPlugin struct {
	provider  *factoidDriver
	processor func(*factoidDriver, string, *base.Line) string
}

func (fp *FactoidPlugin) Apply(val string, line *base.Line) string {
	return fp.processor(fp.provider, val, line)
}

// Replicate perlfu's $<stuff> identifiers
func fd_identifiers(fd *factoidDriver, val string, line *base.Line) string {
	return id_replacer(val, line, time.Now())
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
