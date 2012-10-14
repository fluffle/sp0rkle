package calcdriver

import (
	"fmt"
	"github.com/fluffle/sp0rkle/util/calc"
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"net"
	"strings"
	"strconv"
	"unicode/utf8"
)

func calculate(line *base.Line) {
	maths := line.Args[1]
	if num, err := calc.Calc(maths); err == nil {
		bot.ReplyN(line, "%s = %g", maths, num)
	} else {
		bot.ReplyN(line, "%s error while parsing %s", err, maths)
	}
}

func netmask(line *base.Line) {
	s := strings.Split(line.Args[1], " ")
	if strings.Index(s[1], "/") != -1 {
		// Assume we have netmask ip/cidr
		bot.ReplyN(line, parseCIDR(s[0]))
	} else if len(s) == 2 {
		// Assume we have netmask ip nm
		bot.ReplyN(line, parseMask(s[0], s[1]))
	} else {
		bot.ReplyN(line, "bad netmask args: %s", line.Args[1])
	}
}

func maskRange(ip net.IP, mask net.IPMask) (btm, top net.IP) {
	btm = ip.Mask(mask)
	top = make(net.IP, len(ip))
	copy(top, ip)
	for i, b := range mask {
		top[i] |= ^b
	}
	return
}

func parseCIDR(cidr string) string {
	_, nm, err := net.ParseCIDR(cidr)
	if err == nil {
		btm, top := maskRange(nm.IP, nm.Mask)
		return fmt.Sprintf("%s is in the range %s-%s and has the netmask %s",
			cidr, btm, top, net.IP(nm.Mask))
	}
	return fmt.Sprintf("error parsing ip/cidr %s: %s", cidr, err)
}

func parseMask(ips, nms string) string {
	ip := net.ParseIP(ips)
	nmip := net.ParseIP(nms)
	if ip == nil || nmip == nil {
		return fmt.Sprintf("either %s or %s couldn't be parsed as an IP",
			ips, nms)
	}
	// this is a bit of a hack, because using ParseIP to parse
	// something that's actually a v4 netmask doesn't quite work
	nm := net.IPMask(nmip.To4())
	cidr, bits := nm.Size()
	if ip.To4() != nil && nm != nil {
		if bits != 32 {
			return fmt.Sprintf("%s doesn't look like a valid IPv4 netmask", nms)
		}
	} else {
		// IPv6, hopefully
		nm = net.IPMask(nmip)
		cidr, bits = nm.Size()
		if bits != 128 {
			return fmt.Sprintf("%s doesn't look like a valid IPv6 netmask", nms)
		}
	}
	btm, top := maskRange(ip, nm)
	return fmt.Sprintf("%s/%d is in the range %s-%s and has the netmask %s",
		ip, cidr, btm, top, nmip)
}

func chr(line *base.Line) {
	chr := line.Args[1]
	// handles decimal, hex, and octal \o/
	i, err := strconv.ParseInt(chr, 0, 0)
	if err != nil {
		bot.ReplyN(line, "Couldn't parse %s as an integer: %s", chr, err)
		return
	}
	bot.ReplyN(line, "chr(%s) is %c, %U, '%s'", chr, i, i, utf8repr(rune(i)))
}

func ord(line *base.Line) {
	ord := line.Args[1]
	r, _ := utf8.DecodeRuneInString(ord)
	if r == utf8.RuneError {
		bot.ReplyN(line, "Couldn't parse a utf8 rune from %s", ord)
		return
	}
	bot.ReplyN(line, "ord(%c) is %d, %U, '%s'", r, r, r, utf8repr(r))
}

func utf8repr(r rune) string {
	p := make([]byte, 4)
	n := utf8.EncodeRune(p, r)
	s := make([]string, n)
	for i, c := range p[:n] {
		s[i] = fmt.Sprintf("0x%x", c)
	}
	return strings.Join(s, " ")
}

func convertBase(line *base.Line) {
	s := strings.Split(line.Args[1], " ")
	fromto := strings.Split(s[0], "to")
	if len(fromto) != 2 {
		bot.ReplyN(line, "Specify base as: <from base>to<to base>")
		return
	}
	from, errf := strconv.Atoi(fromto[0])
	to, errt := strconv.Atoi(fromto[1])
	if errf != nil || errt != nil ||
		from < 2 || from > 36 || to < 2 || to > 36 {
		bot.ReplyN(line, "Either %s or %s is a bad base, must be in range 2-36",
			fromto[0], fromto[1])
		return
	}
	i, err := strconv.ParseInt(s[1], from, 64)
	if err != nil {
		bot.ReplyN(line, "Couldn't parse %s as a base %d integer", s[1], from)
		return
	}
	bot.ReplyN(line, "%s in base %d is %s in base %d",
		s[1], from, strconv.FormatInt(i, to), to)
}

func length(line *base.Line) {
	bot.ReplyN(line, "'%s' is %d characters long",
		line.Args[1], len(line.Args[1]))
}
