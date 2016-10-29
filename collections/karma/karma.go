package karma

import (
	"fmt"
	"strings"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/datetime"
	"gopkg.in/mgo.v2"
)

const COLLECTION = "karma"

type Karma struct {
	Subject   string
	Key       string
	Score     int
	Votes     int
	Upvoter   bot.Nick
	Upvtime   time.Time
	Downvoter bot.Nick
	Downvtime time.Time
}

func New(thing string) *Karma {
	return &Karma{
		Subject: thing,
		Key:     strings.ToLower(thing),
	}
}

func (k *Karma) Plus(who bot.Nick) {
	k.Score++
	k.Votes++
	k.Upvoter, k.Upvtime = who, time.Now()
}

func (k *Karma) Minus(who bot.Nick) {
	k.Score--
	k.Votes++
	k.Downvoter, k.Downvtime = who, time.Now()
}

func (k *Karma) String() string {
	s := fmt.Sprintf("'%s' has a karma of %d after %d votes.",
		k.Subject, k.Score, k.Votes)
	if k.Upvoter != "" {
		s += fmt.Sprintf(" Last upvoted by %s at %s.",
			k.Upvoter, datetime.Format(k.Upvtime))
	}
	if k.Downvoter != "" {
		s += fmt.Sprintf(" Last downvoted by %s at %s.",
			k.Downvoter, datetime.Format(k.Downvtime))
	}
	return s
}

func (k *Karma) K() db.Key {
	return db.K{db.S{"key", k.Key}}
}

var _ db.Keyer = (*Karma)(nil)

type Karmas []*Karma

func (ks Karmas) Strings() []string {
	s := make([]string, len(ks))
	for i, k := range ks {
		s[i] = fmt.Sprintf("%#v", k)
	}
	return s
}

type migrator struct {
	mongo, bolt db.Collection
}

func (m *migrator) Migrate() error {
	var all []*Karma
	if err := m.mongo.All(db.K{}, &all); err != nil {
		return err
	}
	var fail error
	for _, k := range all {
		logging.Debug("Migrating karma entry for %s.", k.Subject)
		if err := m.bolt.Put(k); err != nil {
			logging.Error("Inserting karma entry failed: %v", err)
			fail = err
		}
	}
	return fail
}

func (m *migrator) Diff() ([]string, []string, error) {
	var mAll, bAll Karmas
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
	kc := &Collection{db.Both{}}
	kc.Both.MongoC.Init(db.Mongo, COLLECTION, mongoIndexes)
	kc.Both.BoltC.Init(db.Bolt, COLLECTION, nil)
	m := &migrator{
		mongo: kc.Both.MongoC,
		bolt:  kc.Both.BoltC,
	}
	kc.Both.Checker.Init(m, COLLECTION)
	return kc
}

func mongoIndexes(c db.Collection) {
	if err := c.Mongo().EnsureIndex(mgo.Index{
		Key:    []string{"key"},
		Unique: true,
	}); err != nil {
		logging.Error("Couldn't create index on karma.key: %s", err)
	}
	for _, key := range []string{"score", "votes"} {
		if err := c.Mongo().EnsureIndexKey(key); err != nil {
			logging.Error("Couldn't create index on karma.%s: %s", key, err)
		}
	}
}

func (kc *Collection) KarmaFor(sub string) *Karma {
	res := &Karma{Key: strings.ToLower(sub)}
	if err := kc.Get(res.K(), res); err == nil {
		return res
	}
	return nil
}
