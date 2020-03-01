package seen

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util"
	"github.com/fluffle/sp0rkle/util/datetime"
	"github.com/fluffle/sp0rkle/util/diff"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
		db.K{db.S{"key", n.Nick.Lower()}, db.I{"ts", uint64(n.Timestamp.UnixNano())}},
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

type migrator struct {
	mongo, bolt db.Collection
}

func (m *migrator) Migrate() error {
	var all Nicks
	if err := m.mongo.All(db.K{}, &all); err != nil {
		return err
	}
	if err := m.bolt.BatchPut(all); err != nil {
		logging.Error("Migrating seen: %v", err)
		return err
	}
	logging.Info("Migrated %d seen entries.", len(all))
	return nil
}

func (m *migrator) Diff() ([]string, []string, error) {
	var mAll, bAll Nicks
	if err := m.mongo.All(db.K{}, &mAll); err != nil {
		return nil, nil, err
	}
	if err := m.bolt.All(db.K{}, &bAll); err != nil {
		return nil, nil, err
	}
	return mAll.Strings(), bAll.Strings(), nil
}

type Collection struct {
	db.Both
}

func Init() *Collection {
	sc := &Collection{db.Both{}}
	sc.Both.MongoC.Init(db.Mongo, COLLECTION, mongoIndexes)
	sc.Both.BoltC.Init(db.Bolt.Indexed(), COLLECTION, nil)
	m := &migrator{
		mongo: sc.Both.MongoC,
		bolt:  sc.Both.BoltC,
	}
	sc.Both.Checker.Init(m, COLLECTION)
	return sc
}

func mongoIndexes(c db.Collection) {
	indexes := [][]string{
		{"key", "action"}, // For searching ...
		{"timestamp"},     // ... and ordering seen entries.
	}
	for _, key := range indexes {
		if err := c.Mongo().EnsureIndex(mgo.Index{Key: key}); err != nil {
			logging.Error("Couldn't create %v index on sp0rkle.seen: %v", key, err)
		}
	}
}

func (sc *Collection) LastSeen(nick string) *Nick {
	var mAll, bAll Nicks
	n := &Nick{Nick: bot.Nick(nick)}

	// Not using Both here because it's a useful test of BoltDB ordering.
	q := sc.Mongo().Find(bson.M{"key": strings.ToLower(nick)}).Sort("timestamp")
	mErr := q.All(&mAll)
	bErr := sc.BoltC.All(n.byNick(), &bAll)
	if mErr != bErr {
		logging.Warn("LastSeen errors differ: %v != %v", mErr, bErr)
	}
	mStr := mAll.Strings()
	bStr := bAll.Strings()
	unified, err := diff.Unified(mStr, bStr)
	if err != nil {
		logging.Debug("LastSeen: %v\n%s", err, strings.Join(unified, "\n"))
	}
	if sc.Migrated() {
		if len(bAll) == 0 {
			return nil
		}
		return bAll[len(bAll)-1]
	}
	if len(mAll) == 0 {
		return nil
	}
	return mAll[len(mAll)-1]
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
