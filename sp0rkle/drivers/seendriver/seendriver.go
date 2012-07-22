package seendriver

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/lib/seen"
	"regexp"
)

const driverName string = "seen"

var smokeRx *regexp.Regexp = regexp.MustCompile(`(?i)^(?:->\s*?)?(?:s(?:c?h)?m[o0]keh?|cig|fag|spliff|ch[o0]ng|t[o0]ke?)(?:s|z?[0o]r)?\W*?(\?)?$`)

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

type seenDriver struct {
	*seen.SeenCollection
	l logging.Logger
}

func SeenDriver(db *db.Database, l logging.Logger) *seenDriver {
	sc := seen.Collection(db, l)
	return &seenDriver{sc, l}
}

func (sd *seenDriver) Name() string {
	return driverName
}
