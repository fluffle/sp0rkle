package calcdriver

import (
	"github.com/fluffle/sp0rkle/bot"
)

func Init() {
	bot.CommandFunc(calculate, "calc", "calc <expr>  -- does maths for you")
	bot.CommandFunc(netmask, "netmask", "netmask <ip/cidr>|<ip> <mask>" +
		"  -- calculate IPv4 / IPv6 netmasks")
	bot.CommandFunc(chr, "chr", "chr <int>  -- " +
		"prints the character represented by <int> in various formats")
	bot.CommandFunc(ord, "ord", "ord <char>  -- " +
		"prints the numeric and UTF-8 representations of <char>")
	bot.CommandFunc(convertBase, "base", "base <from>to<to> <num>  -- " +
		"converts <num> from base <from> to base <to>")
	bot.CommandFunc(length, "length", "length <string>  -- " +
		"prints the length of <string>")
}
