package quotes

import (
	"math/rand"
	"regexp"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/bson"
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
		db.K{db.I{"qid", uint64(q.QID)}},
	}
}

func (q *Quote) Id() bson.ObjectId {
	return q.Id_
}

func (q *Quote) byQID() db.K {
	return db.K{db.I{"qid", uint64(q.QID)}}
}

func (q *Quote) CollectionName() string { return "quotes" }

var _ db.Indexer = (*Quote)(nil)

type Quotes []*Quote

type Collection struct {
	collection db.IndexedCollection[*Quote]

	seen map[string]map[bson.ObjectId]bool
}

func New(collection db.IndexedCollection[*Quote]) *Collection {
	return &Collection{
		collection: collection,
		seen:       make(map[string]map[bson.ObjectId]bool),
	}
}

func (qc *Collection) Next(key db.Key, set ...int) (int, error) {
	return qc.collection.Next(key, set...)
}

func (qc *Collection) GetByQID(qid int) *Quote {
	res := &Quote{QID: qid}
	quote, err := qc.collection.Get(res.byQID())
	if err != nil || quote == nil {
		return nil
	}
	return quote
}

func (qc *Collection) NewQID() (int, error) {
	return qc.Next(db.K{})
}

func (qc *Collection) GetPseudoRand(regex string) *Quote {
	rx := regexp.MustCompile("(?i)" + regex)
	quotes, err := qc.collection.Match(func(q *Quote) bool {
		return rx.MatchString(q.Quote)
	})
	if err != nil {
		logging.Warn("Quote Match() failed: %s", err)
		return nil
	}
	filtered := Quotes{}
	ids, ok := qc.seen[regex]
	if ok && len(ids) > 0 {
		logging.Debug("Seen '%s' before, %d stored id's", regex, len(ids))
		for _, q := range quotes {
			if !ids[q.Id_] {
				filtered = append(filtered, q)
			}
		}
	} else {
		filtered = quotes
	}

	count := len(filtered)
	switch count {
	case 0:
		if ok {
			delete(qc.seen, regex)
		}
		return nil
	case 1:
		if ok {
			logging.Debug("Zeroing seen data for key '%s'.", regex)
			delete(qc.seen, regex)
		}
		return filtered[0]
	}
	if !ok {
		logging.Debug("Creating seen data for key '%s'.", regex)
		qc.seen[regex] = make(map[bson.ObjectId]bool)
	}
	res := filtered[rand.Intn(count)]
	logging.Debug("Storing id %v for key '%s'.", res.Id_, regex)
	qc.seen[regex][res.Id_] = true
	return res
}

func (qc *Collection) Put(value *Quote) error {
	return qc.collection.Put(value)
}

func (qc *Collection) Del(value *Quote) error {
	return qc.collection.Del(value)
}
