package seen

import (
	"fmt"
	"sort"
	"time"

	"github.com/fluffle/golog/logging"
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
		// NOTE: bson serialization truncates to millisecond precision!
		db.K{db.S{"key", n.Nick.Lower()}, db.I{"ts", uint64(n.Timestamp.UnixMilli())}},
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

type Nicks []*Nick

func (ns Nicks) Strings() []string {
	s := make([]string, len(ns))
	for i, n := range ns {
		s[i] = fmt.Sprintf("%#v", n)
	}
	return s
}

// Implement sort.Interface to sort by descending timestamp.
func (ns Nicks) Len() int           { return len(ns) }
func (ns Nicks) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }
func (ns Nicks) Less(i, j int) bool { return ns[i].Timestamp.After(ns[j].Timestamp) }

type Collection struct {
	db.C
}

func Init() *Collection {
	sc := &Collection{}
	sc.Init(db.Bolt.Indexed(), COLLECTION, nil)

	// Between July 14-September 14 2024 the live sp0rkle instance was not
	// correctly cleaning up/replacing seen Nick instances, instead adding
	// new ones. This has left a bunch of detritus in boltdb, which we can
	// clear up by enforcing some invariants. Some of this has to happen
	// within the db layer, some is dependent on invariants inherent to
	// seen behaviour.
	// This problem was magnified by bson truncating timestamps to ms
	// precision, invalidating indexes.
	if err := sc.Fsck(); err != nil {
		logging.Fatal("seen fsck failed: %v", err)
	}
	return sc
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
	return
}

func (sc *Collection) Fsck() error {
	// First, enforce seen-specific invariants on the stored values.
	var all Nicks
	if err := sc.All(db.K{}, &all); err != nil {
		return fmt.Errorf("seen fsck: fetching all: %w", err)
	}
	rc := &refCheck{}
	for _, n := range all {
		rc.Add(n)
	}
	if len(rc.del) > 0 {
		logging.Warn("seen fsck: removing %d of %d nick values", len(rc.del), len(all))
		for _, n := range rc.del {
			logging.Debug("seen fsck: deleting %#v", n)
			sc.Del(n)
		}
	}
	// Once the values are tidied up, ask db to groom indexes.
	return sc.Collection.Fsck(&Nick{})
}

func (sc *Collection) LastSeen(nick string) *Nick {
	var bAll Nicks
	n := &Nick{Nick: bot.Nick(nick)}
	if err := sc.All(n.byNick(), &bAll); err != nil {
		logging.Error("LastSeen error: %v", err)
		return nil
	}
	if len(bAll) == 0 {
		return nil
	}
	return bAll[len(bAll)-1]
}

func (sc *Collection) LastSeenDoing(nick, act string) *Nick {
	n := &Nick{Nick: bot.Nick(nick), Action: act}
	if err := sc.Get(n.byNickAction(), n); err == nil && n.Exists() {
		return n
	}
	return nil
}

func (sc *Collection) SeenAnyMatching(rx string) []string {
	var ns Nicks
	if err := sc.Match("Nick", rx, &ns); err != nil {
		return nil
	}
	sort.Sort(ns)
	seen := make(map[string]bool)
	res := make([]string, 0, len(ns))
	for _, n := range ns {
		if !seen[n.Nick.Lower()] {
			res = append(res, string(n.Nick))
			seen[n.Nick.Lower()] = true
		}
	}
	return res
}
