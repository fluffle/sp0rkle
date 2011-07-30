package factoids

// This might get ODM-ish in the future.

import (
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"lib/db"
	"lib/util"
	"log"
	"os"
	"rand"
	"strings"
	"time"
)

const COLLECTION string = "factoids"

type FactoidType int

const (
	// Factoids can be of these types
	F_FACT FactoidType = iota
	F_ACTION
	F_URL
)

// A factoid maps a key to a value, and keeps some stats about it
type Factoid struct {
	Key, Value                  string
	Chance                      float32
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
	db.StorableNick
	// Where they did <thing>
	db.StorableChan
	// How many times <thing> has been done before
	Count int
}

// Represent info about things that can be done to the factoid
type FactoidPerms struct {
	ReadOnly bool
	db.StorableNick
}

// Helper to make the work of putting together a completely new *Factoid easier
func NewFactoid(key, value string, n db.StorableNick, c db.StorableChan) *Factoid {
	ts := time.LocalTime()
	ft, fv := ParseValue(value)
	return &Factoid{
		Key: key, Value: fv, Type: ft, Chance: 1.0,
		Created:  &FactoidStat{ts, n, c, 1},
		Modified: &FactoidStat{ts, n, c, 0},
		Accessed: &FactoidStat{ts, n, c, 0},
		Perms:    &FactoidPerms{false, n},
		Id:       bson.NewObjectId(),
	}
}

func (f *Factoid) Access(n db.StorableNick, c db.StorableChan) {
	f.Accessed.Timestamp = time.LocalTime()
	f.Accessed.StorableNick = n
	f.Accessed.StorableChan = c
	f.Accessed.Count++
}

func (f *Factoid) Modify(n db.StorableNick, c db.StorableChan) {
	f.Modified.Timestamp = time.LocalTime()
	f.Modified.StorableNick = n
	f.Modified.StorableChan = c
	f.Modified.Count++
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
		seen:       make(map[string][]bson.ObjectId),
	}
	err := fc.EnsureIndex(mgo.Index{Key: []string{"key"}})
	if err != nil {
		log.Printf("Couldn't create index on sp0rkle.factoids: %v", err)
		return nil, err
	}
	return fc, nil
}

// Can't call this Count because that'd override mgo.Collection.Count()
func (fc *FactoidCollection) GetCount(key string) int {
	if num, err := fc.Find(bson.M{"key": key}).Count(); err == nil {
		return num
	}
	return 0
}

func (fc *FactoidCollection) GetById(id bson.ObjectId) *Factoid {
	var res Factoid
	if err := fc.Find(bson.M{"_id": id}).One(&res); err == nil {
		return &res
	}
	return nil
}

func (fc *FactoidCollection) GetFirst(key string) *Factoid {
	var res Factoid
	if err := fc.Find(bson.M{"key": key}).One(&res); err == nil {
		return &res
	}
	return nil
}

func (fc *FactoidCollection) GetPseudoRand(key string) *Factoid {
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
	if count == 0 {
		if ok {
			// we've seen this before, but people have deleted it since.
			fc.seen[key] = nil, false
		}
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
			fc.seen[key] = make([]bson.ObjectId, 0, count)
		}
		log.Printf("Storing id %v for key %s.\n", res.Id, key)
		fc.seen[key] = append(fc.seen[key], res.Id)
	} else if ok {
		// if the count of results is 1 and we're storing seen data for key
		// then we've exhausted the possible results and should wipe it
		log.Printf("Zeroing seen data for key %s.\n", key)
		fc.seen[key] = nil, false
	}
	return &res
}

func ParseValue(v string) (ft FactoidType, fv string) {
	// Assume v is a normal factoid
	ft = F_FACT

	// Check for perlfu prefixes and strip them
	if strings.HasPrefix(v, "<me>") {
		// <me>does something
		ft, fv = F_ACTION, v[4:]
	} else if strings.HasPrefix(v, "<reply>") {
		// <reply> is treated the same as F_FACT now,
		// Factoid.Key is not used except for searching.
		// NOTE: careful with this -- it's used in factimporter too...
		fv = v[7:]
	} else {
		fv = v
	}
	if util.LooksURLish(fv) {
		// Quite a few factoids are just <reply>http://some.url/
		// it's helpful to detect this so we can do useful things
		ft = F_URL
	}
	return
}
