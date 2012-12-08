package calcdriver

import (
	"github.com/fluffle/sp0rkle/bot"
)

const DateTimeFormat = "15:04:05, Monday 2 January 2006 -0700"

func Init() {
	bot.CommandFunc(calculate, "calc", "calc <expr>  -- does maths for you")
	bot.CommandFunc(date, "date", "date <time/date> [in <zone>] -- "+
		"works out the absolute time for <time/date> [in <zone>]")
	bot.CommandFunc(netmask, "netmask", "netmask <ip/cidr>|<ip> <mask>"+
		"  -- calculate IPv4 / IPv6 netmasks")
	bot.CommandFunc(chr, "chr", "chr <int>  -- "+
		"prints the character represented by <int> in various formats")
	bot.CommandFunc(ord, "ord", "ord <char>  -- "+
		"prints the numeric and UTF-8 representations of <char>")
	bot.CommandFunc(convertBase, "base", "base <from>to<to> <num>  -- "+
		"converts <num> from base <from> to base <to>")
	bot.CommandFunc(length, "length", "length <string>  -- "+
		"prints the length of <string>")
}
