package quotes

import (
	"github.com/fluffle/golog/logging"
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"lib/db"
	"rand"
	"time"
)

const COLLECTION string = "quotes"

type Quote struct {
	Quote     string
	QID       int
	db.StorableNick
	db.StorableChan
	Accessed  int
	Timestamp *time.Time
	Id        bson.ObjectId "_id"
}

func NewQuote(q string, n db.StorableNick, c db.StorableChan) *Quote {
	ts := time.LocalTime()
	return &Quote{q, 0, n, c, 0, ts, bson.NewObjectId()}
}

type QuoteCollection struct {
	// Wrap mgo.Collection
	mgo.Collection

	// logging object
	l logging.Logger
}

func Collection(dbh *db.Database, l logging.Logger) *QuoteCollection {
	qc := &QuoteCollection{
		Collection: dbh.C(COLLECTION),
		l:          l,
	}
	err := qc.EnsureIndex(mgo.Index{Key: []string{"qid"}})
	if err != nil {
		l.Error("Couldn't create index on sp0rkle.quotes: %v", err)
	}
	return qc
}

func (qc *QuoteCollection) GetByQID(qid int) *Quote {
	var res Quote
	if err := qc.Find(bson.M{"qid": qid}).One(&res); err == nil {
		return &res
	}
	return nil
}

func (qc *QuoteCollection) GetQuoteMatching(regex string) *Quote {
	query := qc.Find(bson.M{"quote": bson.M{"$regex": regex}})
	count, err := query.Count()
	if err != nil {
		qc.l.Warn("Count for quote lookup '%s' failed: %s", regex, err)
		return nil
	}
	if count == 0 {
		return nil
	}
	var res Quote
	if count > 1 {
		query = query.Skip(rand.Intn(count))
	}
	if err = query.One(&res); err != nil {
		qc.l.Warn("Fetch for quote lookup '%s' failed: %s", regex, err)
		return nil
	}
	return &res
}
