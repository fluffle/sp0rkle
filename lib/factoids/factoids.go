package factoids

// This might get ODM-ish in the future.

import (
	"github.com/garyburd/go-mongo"
	"log"
	"os"
	"rand"
	"time"
)

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
	Id                          mongo.ObjectId `bson:"_id"`
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
	*mongo.Collection

	// cache of objectIds for PseudoRand
	seen map[string][]mongo.ObjectId
}

// Wrapper to get hold of a factoid collection handle
func Collection(conn mongo.Conn) (*FactoidCollection, os.Error) {
	fc := &FactoidCollection{
		Collection: &mongo.Collection{
			Conn:         conn,
			Namespace:    "sp0rkle.factoids",
			LastErrorCmd: mongo.DefaultLastErrorCmd,
		},
		seen: make(map[string][]mongo.ObjectId),
	}
	err := fc.CreateIndex(mongo.D{mongo.DocItem{Key: "Key", Value: 1}}, nil)
	if err != nil {
		log.Printf("Couldn't create index on sp0rkle.factoids: %v", err)
		return nil, err
	}
	return fc, nil
}

func (fc *FactoidCollection) GetFirst(key string) (*Factoid) {
	var res Factoid
	if err := fc.Find(mongo.M{"Key": key}).One(&res); err != nil {
		return nil
	}
	return &res
}

func (fc *FactoidCollection) GetPseudoRand(key string) (*Factoid) {
	lookup := mongo.M{"Key": key}
	ids, ok := fc.seen[key]
	if ok && len(ids) > 0 {
		log.Printf("Seen %s before, %d stored id's\n", key, len(ids))
		lookup["_id"] = mongo.M{"$nin": ids}
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
		query = query.Skip(rand.Intn(int(count)))
	}
	if err = query.One(&res); err != nil {
		log.Printf("Fetching for key failed: %v", err)
		return nil
	}
	if count != 1 {
		if !ok {
			// only store seen for keys that have more than one factoid
			log.Printf("Creating seen data for key %s.\n", key)
			fc.seen[key] = make([]mongo.ObjectId, 0, 1)
		}
		log.Printf("Storing id %v for key %s.\n", res.Id, key)
		fc.seen[key] = append(fc.seen[key], res.Id)
	} else if ok {
		// if the count of results is 1 and we're storing seen data for key
		// then we've exhausted the possible results and should wipe it
		log.Printf("Zeroing seen data for key %s.\n", key)
		fc.seen[key] = make([]mongo.ObjectId, 0, 1)
	}
	return &res
}
