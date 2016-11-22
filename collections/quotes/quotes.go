package quotes

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/datetime"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION string = "quotes"

type Quote struct {
	Quote     string
	QID       int
	Nick      bot.Nick
	Chan      bot.Chan
	Accessed  int
	Timestamp time.Time
	Id_       bson.ObjectId `bson:"_id,omitempty"`
}

var _ db.Indexer = (*Quote)(nil)

func NewQuote(q string, n bot.Nick, c bot.Chan) *Quote {
	return &Quote{q, 0, n, c, 0, time.Now(), bson.NewObjectId()}
}

func (q *Quote) Indexes() []db.Key {
	return []db.Key{
		db.K{db.I{"qid", q.QID}},
	}
}

func (q *Quote) Id() bson.ObjectId {
	return q.Id_
}

func (q *Quote) byQID() db.K {
	return db.K{db.I{"qid", q.QID}}
}

type Quotes []*Quote

func (qs Quotes) Strings() []string {
	s := make([]string, len(qs))
	for i, q := range qs {
		// Explicitly omit QID here since QIDs in Bolt will come from
		// the bucket sequence and probably be != Mongo.
		s[i] = fmt.Sprintf("%s <%s:%s> %s (%d)", datetime.Format(q.Timestamp),
			q.Nick, q.Chan, q.Quote, q.Accessed)
	}
	return s
}

type migrator struct {
	mongo, bolt db.Collection
}

func (m *migrator) Migrate() error {
	var all Quotes
	// Break encapsulation to preserve quote ID ordering.
	if err := m.mongo.Mongo().Find(bson.M{}).Sort("qid").All(&all); err != nil {
		return err
	}
	if err := m.bolt.BatchPut(all); err != nil {
		logging.Error("Migrating quotes: %v", err)
		return err
	}
	logging.Info("Migrated %d quotes.", len(all))
	// Update sequence with current largest QID.
	_, err := m.bolt.Next(db.K{}, uint64(all[len(all)-1].QID))
	return err
}

func (m *migrator) Diff() ([]string, []string, error) {
	var mAll, bAll Quotes
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

	// Cache of ObjectId's for PseudoRand
	seen map[string]map[bson.ObjectId]bool

	// This is a bit of a gratuitous hack to allow for easier numeric quote IDs.
	maxQID int32
}

func Init() *Collection {
	qc := &Collection{
		Both:   db.Both{},
		seen:   make(map[string]map[bson.ObjectId]bool),
		maxQID: 1,
	}
	qc.Both.MongoC.Init(db.Mongo, COLLECTION, mongoIndexes)
	qc.Both.BoltC.Init(db.Bolt, COLLECTION, nil)
	m := &migrator{
		mongo: qc.Both.MongoC,
		bolt:  qc.Both.BoltC,
	}
	qc.Both.Checker.Init(m, COLLECTION)

	// QID incrementing is not in mongodb so we break out here.
	var res Quote
	if err := qc.Mongo().Find(bson.M{}).Sort("-qid").One(&res); err == nil {
		qc.maxQID = int32(res.QID)
	}
	return qc
}

func mongoIndexes(c db.Collection) {
	err := c.Mongo().EnsureIndex(mgo.Index{Key: []string{"qid"}, Unique: true})
	if err != nil {
		logging.Error("Couldn't create index on sp0rkle.quotes: %v", err)
	}
}

func (qc *Collection) GetByQID(qid int) *Quote {
	res := &Quote{QID: qid}
	if err := qc.Get(res.byQID(), res); err == nil {
		return res
	}
	return nil
}

func (qc *Collection) NewQID() (int, error) {
	if qc.Migrated() {
		return qc.Next(db.K{})
	}
	return int(atomic.AddInt32(&qc.maxQID, 1)), nil
}

func (qc *Collection) GetPseudoRand(regex string) *Quote {
	// TODO(fluffle): This implementation of GetPseudoRand is inefficient
	// for either Bolt or Mongo on their own. It's a lowest-common-denominator
	// that should work for both. There are 3 steps: fetch all quotes matching
	// the regex, filter out already-seen ObjectIds, and return a result while
	// updating the ObjectId filters.

	quotes := Quotes{}
	if regex == "" {
		if err := qc.All(db.K{}, &quotes); err != nil {
			logging.Warn("Quote All() failed: %s", err)
			return nil
		}
	} else {
		if err := qc.Match("Quote", regex, &quotes); err != nil {
			logging.Warn("Quote Match(%q) failed: %s", regex, err)
			return nil
		}
	}

	filtered := Quotes{}
	ids, ok := qc.seen[regex]
	if ok && len(ids) > 0 {
		logging.Debug("Looked for quotes matching %q before, %d stored id's",
			regex, len(ids))
		for _, quote := range quotes {
			if !ids[quote.Id_] {
				filtered = append(filtered, quote)
			}
		}
	} else {
		filtered = quotes
	}

	count := len(filtered)
	switch count {
	case 0:
		if ok {
			// Looked for this regex before, but nothing matches now
			delete(qc.seen, regex)
		}
		return nil
	case 1:
		if ok {
			// if the count of results is 1 and we're storing seen data for regex
			// then we've exhausted the possible results and should wipe it
			logging.Debug("Zeroing seen data for regex %q.", regex)
			delete(qc.seen, regex)
		}
		return filtered[0]
	}
	// case count > 1:
	if !ok {
		// only store seen for regex that match more than one quote
		logging.Debug("Creating seen data for regex %q.", regex)
		qc.seen[regex] = map[bson.ObjectId]bool{}
	}
	res := filtered[rand.Intn(count)]
	logging.Debug("Storing id %v for regex %q.", res.Id_, regex)
	qc.seen[regex][res.Id_] = true
	return res
}
