package seendriver

import (
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/seen"
	"regexp"
	"time"
)

var smokeRx *regexp.Regexp = regexp.MustCompile(`(?i)^(?:->\s*?)?(?:s(?:c?h)?m[o0]keh?|cig|fag|spliff|ch[o0]ng|t[o0]ke?)(?:s|z?[0o]r)?\W*?(\?)?$`)

var milestones = []int{100, 500, 1000, 5000, 10000, 25000, 50000, 75000, 100000}

type stupidQuestion struct {
	re string
	rx *regexp.Regexp
	resp string
}

var wittyComebacks []stupidQuestion = []stupidQuestion{
	{`^my (?:arse|ass)$`, nil,
	 "Pull your pants down and hit me with the view, big boy."},
	{`^my (?:penis|cock|dick|wang)$`, nil,
	 "No, thank god... Now put it away, no-one else wants to see it either."},
	{`^(?:yo(?:'|ur)?|\w+'?s) (?:momma|mother|mum)$`, nil,
	 "Yeah, she gives me a discount cos I see her so regularly \\o/"},
	{`^\w+'?s (?:arse|ass|penis|cock|dick|wang)$`, nil,
	 "Unfortunately not... I asked nicely but they're a bit shy :/"},
	{`^me$`, nil, "You're right there, fool."},
}

func init() {
	for i, w := range wittyComebacks {
		// all regex matches for comebacks should be case-insensitive
		wittyComebacks[i].rx = regexp.MustCompile("(?i)" + w.re)
	}
}

var sc *seen.Collection

func Init() {
	sc = seen.Init()

	bot.HandleFunc(smoke, "privmsg", "action")
	bot.HandleFunc(recordLines, "privmsg", "action")
	bot.HandleFunc(recordPrivmsg, "privmsg", "action")
	bot.HandleFunc(recordJoin, "join", "part")
	bot.HandleFunc(recordNick, "nick", "quit")
	bot.HandleFunc(recordKick, "kick")

	bot.CommandFunc(seenCmd, "seen", "seen <nick> [action]  -- " +
		"display the last time <nick> was seen on IRC [doing action]")
	bot.CommandFunc(lines, "lines", "lines [nick]  -- " +
		"display how many lines you [or nick] has said in the channel")
	bot.CommandFunc(topten, "topten", "topten  -- " +
		"display the nicks who have said the most in the channel")
	bot.CommandFunc(topten, "top10", "top10  -- " +
		"display the nicks who have said the most in the channel")
}

// Look up or create a "seen" entry for the line.
// Explicitly don't handle updating line.Text or line.OtherNick
func seenNickFromLine(line *base.Line) *seen.Nick {
	sn := sc.LastSeenDoing(line.Nick, line.Cmd)
	n, c := line.Storable()
	if sn == nil {
		sn = seen.SawNick(n, c, line.Cmd, "")
	} else {
		sn.Nick, sn.Chan = n, c
		sn.Timestamp = time.Now()
	}
	return sn
}
