package quotes

import (
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
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
	Id        bson.ObjectId `bson:"_id,omitempty"`
}

func NewQuote(q string, n bot.Nick, c bot.Chan) *Quote {
	return &Quote{q, 0, n, c, 0, time.Now(), bson.NewObjectId()}
}

type Collection struct {
	// Wrap mgo.Collection
	*mgo.Collection

	// Cache of ObjectId's for PseudoRand
	seen map[string][]bson.ObjectId

	// This is a bit of a gratuitous hack to allow for easier numeric quote IDs.
	maxQID int32
}

func Init() *Collection {
	qc := &Collection{
		Collection: db.Mongo.C(COLLECTION).Mongo(),
		seen:       make(map[string][]bson.ObjectId),
		maxQID:     1,
	}
	err := qc.EnsureIndex(mgo.Index{Key: []string{"qid"}, Unique: true})
	if err != nil {
		logging.Error("Couldn't create index on sp0rkle.quotes: %v", err)
	}

	var res Quote
	if err := qc.Find(bson.M{}).Sort("-qid").One(&res); err == nil {
		qc.maxQID = int32(res.QID)
	}
	return qc
}

func (qc *Collection) GetByQID(qid int) *Quote {
	var res Quote
	if err := qc.Find(bson.M{"qid": qid}).One(&res); err == nil {
		return &res
	}
	return nil
}

func (qc *Collection) NewQID() int {
	return int(atomic.AddInt32(&qc.maxQID, 1))
}

// TODO(fluffle): reduce duplication with lib/factoids?
func (qc *Collection) GetPseudoRand(regex string) *Quote {
	lookup := bson.M{}
	if regex != "" {
		// Only perform a regex lookup if there's a regex to match against,
		// otherwise this just fetches a quote at pseudo-random.
		lookup["quote"] = bson.M{"$regex": regex, "$options": "i"}
	}
	ids, ok := qc.seen[regex]
	if ok && len(ids) > 0 {
		logging.Debug("Looked for quotes matching '%s' before, %d stored id's",
			regex, len(ids))
		lookup["_id"] = bson.M{"$nin": ids}
	}
	query := qc.Find(lookup)
	count, err := query.Count()
	if err != nil {
		logging.Warn("Count for quote lookup '%s' failed: %s", regex, err)
		return nil
	}
	if count == 0 {
		if ok {
			// Looked for this regex before, but nothing matches now
			delete(qc.seen, regex)
		}
		return nil
	}
	var res Quote
	if count > 1 {
		query = query.Skip(rand.Intn(count))
	}
	if err = query.One(&res); err != nil {
		logging.Warn("Fetch for quote lookup '%s' failed: %s", regex, err)
		return nil
	}
	if count != 1 {
		if !ok {
			// only store seen for regex that match more than one quote
			logging.Debug("Creating seen data for regex '%s'.", regex)
			qc.seen[regex] = make([]bson.ObjectId, 0, count)
		}
		logging.Debug("Storing id %v for regex '%s'.", res.Id, regex)
		qc.seen[regex] = append(qc.seen[regex], res.Id)
	} else if ok {
		// if the count of results is 1 and we're storing seen data for regex
		// then we've exhausted the possible results and should wipe it
		logging.Debug("Zeroing seen data for regex '%s'.", regex)
		delete(qc.seen, regex)
	}
	return &res
}
