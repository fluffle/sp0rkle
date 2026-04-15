package quotes

import (
	"math/rand"
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

type Quotes []*Quote

type Collection struct {
	db.C

	seen map[string]map[bson.ObjectId]bool
}

func Init() *Collection {
	qc := &Collection{
		seen: make(map[string]map[bson.ObjectId]bool),
	}
	qc.Init(db.Bolt.Indexed(), COLLECTION, nil)
	if err := qc.Fsck(&Quote{}); err != nil {
		logging.Fatal("quotes fsck failed: %v", err)
	}
	return qc
}

func (qc *Collection) GetByQID(qid int) *Quote {
	res := &Quote{QID: qid}
	if err := qc.Get(res.byQID(), res); err == nil {
		return res
	}
	return nil
}

func (qc *Collection) NewQID() (int, error) {
	return qc.Next(db.K{})
}

func (qc *Collection) GetPseudoRand(regex string) *Quote {
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
			delete(qc.seen, regex)
		}
		return filtered[0]
	}
	// case count > 1, effectively
	// only store seen for regex that match more than one quote
	if !ok {
		qc.seen[regex] = map[bson.ObjectId]bool{}
	}
	res := filtered[rand.Intn(count)]
	logging.Debug("Storing id %v for regex %q.", res.Id_, regex)
	qc.seen[regex][res.Id_] = true
	return res
}
