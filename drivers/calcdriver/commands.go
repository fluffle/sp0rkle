package calcdriver

import (
	"fmt"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util/calc"
	"github.com/fluffle/sp0rkle/util/datetime"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func calculate(ctx *bot.Context) {
	maths := ctx.Text()
	if num, err := calc.Calc(maths); err == nil {
		ctx.ReplyN("%s = %g", maths, num)
	} else {
		ctx.ReplyN("%s error while parsing %s", err, maths)
	}
}

func date(ctx *bot.Context) {
	tstr, zone := ctx.Text(), ""
	if idx := strings.Index(tstr, "in "); idx != -1 {
		tstr, zone = tstr[:idx], strings.TrimSpace(tstr[idx+3:])
	}
	tm, ok := time.Now(), true
	if tstr != "" {
		if tm, ok = datetime.Parse(tstr); !ok {
			ctx.ReplyN("Couldn't parse time string '%s'.", tstr)
			return
		}
	}
	if loc := datetime.Zone(zone); zone != "" && loc != nil {
		tm = tm.In(loc)
	} else {
		tm = tm.In(time.Local)
	}
	ctx.ReplyN("%s", tm.Format(DateTimeFormat))
}

func netmask(ctx *bot.Context) {
	s := strings.Split(ctx.Text(), " ")
	if len(s) == 1 && strings.Index(s[0], "/") != -1 {
		// Assume we have netmask ip/cidr
		ctx.ReplyN("%s", parseCIDR(s[0]))
	} else if len(s) == 2 {
		// Assume we have netmask ip nm
		ctx.ReplyN("%s", parseMask(s[0], s[1]))
	} else {
		ctx.ReplyN("bad netmask args: %s", ctx.Text())
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
	if (ip.To4() == nil && nmip.To4() != nil ||
		ip.To4() != nil && nmip.To4() == nil) {
		return fmt.Sprintf("can't mix v4 and v6 ip / netmask specifications")
	}
	v4 := ip.To4() != nil
	if v4 {
		// Ensure we're working with 32-bit addrs.
		ip = ip.To4()
		nmip = nmip.To4()
	}
	// this is a bit of a hack, because using ParseIP to parse
	// something that's actually a v4 netmask doesn't quite work
	nm := net.IPMask(nmip)
	cidr, bits := nm.Size()
	if v4 && bits != 32 {
		return fmt.Sprintf("%s doesn't look like a valid IPv4 netmask", nms)
	} else if !v4 && bits != 128 {
		return fmt.Sprintf("%s doesn't look like a valid IPv6 netmask", nms)
	}
	btm, top := maskRange(ip, nm)
	if v4 {
		// Ditto.
		btm = btm.To4()
		top = top.To4()
	}
	return fmt.Sprintf("%s/%d is in the range %s-%s and has the netmask %s",
		ip, cidr, btm, top, nmip)
}

func chr(ctx *bot.Context) {
	chr := strings.ToLower(ctx.Text())
	if strings.HasPrefix(chr, "u+") {
		// Allow "unicode" syntax by translating it to 0x...
		chr = "0x" + chr[2:]
	}
	// handles decimal, hex, and octal \o/
	i, err := strconv.ParseInt(chr, 0, 0)
	if err != nil {
		ctx.ReplyN("Couldn't parse %s as an integer: %s", chr, err)
		return
	}
	ctx.ReplyN("chr(%s) is %c, %U, '%s'", chr, i, i, utf8repr(rune(i)))
}

func ord(ctx *bot.Context) {
	ord := ctx.Text()
	r, _ := utf8.DecodeRuneInString(ord)
	if r == utf8.RuneError {
		ctx.ReplyN("Couldn't parse a utf8 rune from %s", ord)
		return
	}
	ctx.ReplyN("ord(%c) is %d, %U, '%s'", r, r, r, utf8repr(r))
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

func convertBase(ctx *bot.Context) {
	s := strings.Split(ctx.Text(), " ")
	fromto := strings.Split(s[0], "to")
	if len(fromto) != 2 {
		ctx.ReplyN("Specify base as: <from base>to<to base>")
		return
	}
	from, errf := strconv.Atoi(fromto[0])
	to, errt := strconv.Atoi(fromto[1])
	if errf != nil || errt != nil ||
		from < 2 || from > 36 || to < 2 || to > 36 {
		ctx.ReplyN("Either %s or %s is a bad base, must be in range 2-36",
			fromto[0], fromto[1])

		return
	}
	i, err := strconv.ParseInt(s[1], from, 64)
	if err != nil {
		ctx.ReplyN("Couldn't parse %s as a base %d integer", s[1], from)
		return
	}
	ctx.ReplyN("%s in base %d is %s in base %d",
		s[1], from, strconv.FormatInt(i, to), to)

}

func length(ctx *bot.Context) {
	ctx.ReplyN("'%s' is %d characters long",
		ctx.Text(), len(ctx.Text()))

}
