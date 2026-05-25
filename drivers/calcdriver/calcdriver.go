package calcdriver

import (
	"strings"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db/conf"
	"github.com/fluffle/sp0rkle/util/datetime"
)

type Driver struct {
	tzNs conf.Namespace
}

func (d *Driver) Zone(nick string, tz ...string) string {
	nick = strings.ToLower(nick)
	if len(tz) > 0 && tz[0] == "" {
		d.tzNs.Delete(nick)
		return ""
	}
	return d.tzNs.String(nick, tz...)
}

func New(b *bot.Bot, config *conf.Registry) *Driver {
	d := &Driver{tzNs: config.Ns(datetime.ZoneNs)}
	b.Command(d.calculate, "calc", "calc <expr>  -- does maths for you")
	b.Command(d.date, "date", "date <time/date> [in <zone>] -- "+
		"works out the absolute time for <time/date> [in <zone>]")
	b.Command(d.netmask, "netmask", "netmask <ip/cidr>|<ip> <mask>"+
		"  -- calculate IPv4 / IPv6 netmasks")
	b.Command(d.chr, "chr", "chr <int>  -- "+
		"prints the character represented by <int> in various formats")
	b.Command(d.ord, "ord", "ord <char>  -- "+
		"prints the numeric and UTF-8 representations of <char>")
	b.Command(d.convertBase, "base", "base <from>to<to> <num>  -- "+
		"converts <num> from base <from> to base <to>")
	b.Command(d.length, "length", "length <string>  -- "+
		"prints the length of <string>")
	return d
}
