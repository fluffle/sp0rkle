package seen

import (
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util"
	"github.com/fluffle/sp0rkle/util/datetime"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

const COLLECTION string = "seen"

type Nick struct {
	Nick      bot.Nick
	Chan      bot.Chan
	OtherNick bot.Nick
	Timestamp time.Time
	Key       string
	Action    string
	Text      string
}

type seenMsg func(*Nick) string

var actionMap map[string]seenMsg = map[string]seenMsg{
	"PRIVMSG": func(n *Nick) string {
		return fmt.Sprintf("in %s, saying '%s'", n.Chan, n.Text)
	},
	"ACTION": func(n *Nick) string {
		return fmt.Sprintf("in %s, saying '%s %s'", n.Chan, n.Nick, n.Text)
	},
	"JOIN": func(n *Nick) string {
		return fmt.Sprintf("joining %s", n.Chan)
	},
	"PART": func(n *Nick) string {
		return fmt.Sprintf("parting %s with the message '%s'", n.Chan, n.Text)
	},
	"KICKING": func(n *Nick) string {
		return fmt.Sprintf("kicking %s from %s with the message '%s'",
			n.OtherNick, n.Chan, n.Text)
	},
	"KICKED": func(n *Nick) string {
		return fmt.Sprintf("being kicked from %s by %s with the message '%s'",
			n.Chan, n.OtherNick, n.Text)
	},
	"QUIT": func(n *Nick) string {
		return fmt.Sprintf("quitting with the message '%s'", n.Text)
	},
	"NICK": func(n *Nick) string {
		return fmt.Sprintf("changing their nick to '%s'", n.Text)
	},
	"SMOKE": func(n *Nick) string { return "going for a smoke." },
}

func SawNick(nick bot.Nick, ch bot.Chan, act, txt string) *Nick {
	return &Nick{
		Nick:      nick,
		Chan:      ch,
		OtherNick: "",
		Timestamp: time.Now(),
		Key:       nick.Lower(),
		Action:    act,
		Text:      txt,
	}
}

func (n *Nick) String() string {
	if act, ok := actionMap[n.Action]; ok {
		return fmt.Sprintf("I last saw %s on %s (%s ago), %s.",
			n.Nick, datetime.Format(n.Timestamp),
			util.TimeSince(n.Timestamp), act(n))
	}
	// No specific message format for the action seen.
	return fmt.Sprintf("I last saw %s at %s (%s ago).",
		n.Nick, datetime.Format(n.Timestamp),
		util.TimeSince(n.Timestamp))
}

func (n *Nick) Id() bson.M {
	return bson.M{"key": n.Key, "action": n.Action}
}

type Collection struct {
	// Wrap mgo.Collection
	*mgo.Collection
}

func Init() *Collection {
	sc := &Collection{db.Init().C(COLLECTION)}
	indexes := [][]string{
		{"key", "action"}, // For searching ...
		{"timestamp"},     // ... and ordering seen entries.
	}
	for _, key := range indexes {
		if err := sc.EnsureIndex(mgo.Index{Key: key}); err != nil {
			logging.Error("Couldn't create %v index on sp0rkle.seen: %v", key, err)
		}
	}
	return sc
}

func (sc *Collection) LastSeen(nick string) *Nick {
	var res Nick
	q := sc.Find(bson.M{"key": strings.ToLower(nick)}).Sort("-timestamp")
	if err := q.One(&res); err == nil {
		return &res
	}
	return nil
}

func (sc *Collection) LastSeenDoing(nick, act string) *Nick {
	var res Nick
	q := sc.Find(bson.M{"key": strings.ToLower(nick), "action": act}).Sort("-timestamp")
	if err := q.One(&res); err == nil {
		return &res
	}
	return nil
}

func (sc *Collection) SeenAnyMatching(rx string) []string {
	var res []string
	q := sc.Find(bson.M{"key": bson.M{"$regex": rx, "$options": "i"}}).Sort("-timestamp")
	if err := q.Distinct("key", &res); err != nil {
		logging.Warn("SeenAnyMatching Find error: %v", err)
		return []string{}
	}
	logging.Debug("Looked for matches, found %#v", res)
	return res
}
