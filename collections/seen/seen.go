package seen

import (
	"fmt"
	"regexp"
	"slices"
	"time"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util"
	"github.com/fluffle/sp0rkle/util/datetime"
	"github.com/fluffle/sp0rkle/util/bson"
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
	Id_       bson.ObjectId `bson:"_id"`
}

var _ db.Indexer = (*Nick)(nil)

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
		Id_:       bson.NewObjectId(),
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

func (n *Nick) Indexes() []db.Key {
	// Yes, this creates two buckets per nick, but then we don't have to worry
	// about the keys *in* the bucket. Using "nick" for both keys would mean an
	// All() lookup for "nick" would resolve both action and ts pointers.
	// This way either we look up nick + action or key (implicitly ordered by ts).
	//
	// This could *theoretically* be reduced to one bucket by taking into
	// account implementation details of All() and boltdb key ordering --
	// if the timestamp key sorts lexographically before the action key then
	// those pointers will be resolved first (in timestamp order), and
	// the action pointers *should* be deduped and ignored by All().
	// This means the results of All() would still be in timestamp order.
	return []db.Key{
		db.K{db.S{"nick", n.Nick.Lower()}, db.S{"action", n.Action}},
		db.K{db.S{"key", n.Nick.Lower()}, db.TS{"ts", n.Timestamp}},
	}
}

func (n *Nick) Id() bson.ObjectId {
	return n.Id_
}

func (n *Nick) Exists() bool {
	return n != nil && len(n.Id_) > 0
}

func (n *Nick) byNick() db.K {
	// Uses "key" not "nick" bucket, so that results are ordered by timestamp.
	return db.K{db.S{"key", n.Nick.Lower()}}
}

func (n *Nick) byNickAction() db.K {
	return db.K{db.S{"nick", n.Nick.Lower()}, db.S{"action", n.Action}}
}

func (n *Nick) CollectionName() string { return "seen" }

func byTimestamp(a, b *Nick) int {
	return a.Timestamp.Compare(b.Timestamp)
}

type Collection struct {
	collection db.IndexedCollection[*Nick]
}

func New(collection db.IndexedCollection[*Nick]) *Collection {
	return &Collection{collection: collection}
}

// actMap keys are Actions
type actMap map[string]*Nick

type refCheck struct {
	del []*Nick
	// seen is a two-level map that tracks the hierarchy in boltdb
	// the invariant we want to enforce is that a given IRC nick must only
	// have one stored *Nick per action type, and that this is the newest
	// of the available ones.
	seen map[bot.Nick]actMap
}

func (rc *refCheck) Add(n *Nick) {
	if rc.seen == nil {
		rc.seen = map[bot.Nick]actMap{}
	}
	am, ok := rc.seen[n.Nick]
	if !ok {
		am = actMap{}
		rc.seen[n.Nick] = am
	}
	prev, ok := am[n.Action]
	if !ok {
		am[n.Action] = n
		return
	}
	if prev.Timestamp.Before(n.Timestamp) {
		am[n.Action] = n
		rc.del = append(rc.del, prev)
	} else {
		rc.del = append(rc.del, n)
	}
}

func (n *Nick) ValidateValues(all []*Nick) ([]*Nick, error) {
	rc := &refCheck{}
	for _, nick := range all {
		rc.Add(nick)
	}
	return rc.del, nil
}

func (sc *Collection) LastSeen(nick string) *Nick {
	n := &Nick{Nick: bot.Nick(nick)}
	nicks, err := sc.collection.All(n.byNick())
	if err != nil || len(nicks) == 0 {
		return nil
	}
	return nicks[len(nicks)-1]
}

func (sc *Collection) LastSeenDoing(nick, act string) *Nick {
	n := &Nick{Nick: bot.Nick(nick), Action: act}
	saw, err := sc.collection.Get(n.byNickAction())
	if err != nil || !saw.Exists() {
		return nil
	}
	return saw
}

func (sc *Collection) SeenAnyMatching(rx string) []string {
	r := regexp.MustCompile("(?i)" + rx)
	nicks, err := sc.collection.Match(func(n *Nick) bool {
		return r.MatchString(string(n.Nick))
	})
	if err != nil || len(nicks) == 0 {
		return nil
	}
	slices.SortFunc(nicks, byTimestamp)
	seen := make(map[string]bool)
	res := make([]string, 0, len(nicks))
	for _, n := range nicks {
		if !seen[n.Nick.Lower()] {
			res = append(res, string(n.Nick))
			seen[n.Nick.Lower()] = true
		}
	}
	return res
}

func (sc *Collection) Put(value *Nick) error {
	return sc.collection.Put(value)
}

func (sc *Collection) Del(value *Nick) error {
	return sc.collection.Del(value)
}
