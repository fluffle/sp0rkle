package conf

import (
	"fmt"

	"github.com/fluffle/goirc/logging"
	"github.com/fluffle/sp0rkle/db"
	"gopkg.in/mgo.v2"
)

const COLLECTION string = "conf"

var mongo db.C

func mongoIndexes(c db.Collection) {
	err := c.Mongo().EnsureIndex(mgo.Index{Key: []string{"ns", "key"}, Unique: true})
	if err != nil {
		logging.Error("Couldn't create index on sp0rkle.conf: %s", err)
	}
	// Dirty hacks to migrate data.
	var all []Entry
	mongo.All(db.K{}, &all)
	for _, e := range all {
		logging.Debug("Migrating entry %s.", e)
		Bolt(e.Ns).Value(e.Key, e.Value)
	}
}

func Mongo(ns string) *namespace {
	mongo.Init(db.Mongo, COLLECTION, mongoIndexes)
	return &namespace{ns: ns, Collection: &mongo}
}

var bolt db.C

func Bolt(ns string) *namespace {
	bolt.Init(db.Bolt, COLLECTION, nil)
	return &namespace{ns: ns, Collection: &bolt}
}

func Ns(ns string) *both {
	return &both{bolt: Bolt(ns), mongo: Mongo(ns)}
}

type Entry struct {
	Ns, Key string
	Value   interface{}
}

func (e Entry) String() string {
	return fmt.Sprintf("%s<%s: %v>", e.Ns, e.Key, e.Value)
}
