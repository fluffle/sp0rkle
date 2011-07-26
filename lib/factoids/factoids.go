package factoids

// This might get ODM-ish in the future.

import (
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"lib/db"
	"log"
	"os"
	"rand"
	"time"
)

const COLLECTION string = "factoids"

type FactoidType int

const (
	// Factoids can be of these types
	F_FACT FactoidType = iota
	F_ACTION
	F_REPLY
	F_URL
)

// A factoid maps a key to a value, and keeps some stats about it
type Factoid struct {
	Key, Value                  string
	Type                        FactoidType
	Created, Modified, Accessed *FactoidStat
	Perms                       *FactoidPerms
	Id                          bson.ObjectId "_id"
}

// Represent info about things that happened to the factoid
type FactoidStat struct {
	// When <thing> happened
	Timestamp *time.Time
	// Who did <thing>
	Nick, Ident, Host string
	// Where they did <thing>
	Chan string
	// How many times <thing> has been done before
	Count int
}

// Represent info about things that can be done to the factoid
type FactoidPerms struct {
	ReadOnly bool
	Owner    string
}

// Factoids are stored in a mongo collection of Factoid structs
type FactoidCollection struct {
	// We're wrapping mgo.Collection so we can provide our own methods.
	mgo.Collection

	// cache of objectIds for PseudoRand
	seen map[string][]bson.ObjectId
}

// Wrapper to get hold of a factoid collection handle
func Collection(dbh *db.Database) (*FactoidCollection, os.Error) {
	fc := &FactoidCollection{
		Collection: dbh.C(COLLECTION),
		seen: make(map[string][]bson.ObjectId),
	}
	err := fc.EnsureIndex(mgo.Index{Key: []string{"key"}})
	if err != nil {
		log.Printf("Couldn't create index on sp0rkle.factoids: %v", err)
		return nil, err
	}
	return fc, nil
}

func (fc *FactoidCollection) GetFirst(key string) (*Factoid) {
	var res Factoid
	if err := fc.Find(bson.M{"key": key}).One(&res); err != nil {
		return nil
	}
	return &res
}

func (fc *FactoidCollection) GetPseudoRand(key string) (*Factoid) {
	lookup := bson.M{"key": key}
	ids, ok := fc.seen[key]
	if ok && len(ids) > 0 {
		log.Printf("Seen %s before, %d stored id's\n", key, len(ids))
		lookup["_id"] = bson.M{"$nin": ids}
	}
	query := fc.Find(lookup)
	count, err := query.Count()
	if err != nil {
		log.Printf("Counting for key failed: %v", err)
		return nil
	}
	log.Printf("Got %d results back from query.\n", count)
	if count == 0 {
		return nil
	}
	var res Factoid
	if count > 1 {
		query = query.Skip(rand.Intn(count))
	}
	if err = query.One(&res); err != nil {
		log.Printf("Fetching for key failed: %v", err)
		return nil
	}
	if count != 1 {
		if !ok {
			// only store seen for keys that have more than one factoid
			log.Printf("Creating seen data for key %s.\n", key)
			fc.seen[key] = make([]bson.ObjectId, 0, 1)
		}
		log.Printf("Storing id %v for key %s.\n", res.Id, key)
		fc.seen[key] = append(fc.seen[key], res.Id)
	} else if ok {
		// if the count of results is 1 and we're storing seen data for key
		// then we've exhausted the possible results and should wipe it
		log.Printf("Zeroing seen data for key %s.\n", key)
		fc.seen[key] = make([]bson.ObjectId, 0, 1)
	}
	return &res
}
