package factoids

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/base"
	_ "github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
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
	Chance                      float64
	Type                        FactoidType
	Created, Modified, Accessed *FactoidStat
	Perms                       *FactoidPerms
	Id                          bson.ObjectId `bson:"_id,omitempty"`
}

// Represent info about things that happened to the factoid
type FactoidStat struct {
	// When <thing> happened
	Timestamp time.Time
	// Who did <thing>
	Nick base.Nick
	// Where they did <thing>
	Chan base.Chan
	// How many times <thing> has been done before
	Count int
}

// Represent info about things that can be done to the factoid
type FactoidPerms struct {
	ReadOnly bool
	Nick base.Nick
}

// Represent info returned from the Info MapReduce
type FactoidInfo struct {
	Created, Modified, Accessed int
}

// Helper to make the work of putting together a completely new *Factoid easier
func NewFactoid(key, value string, n base.Nick, c base.Chan) *Factoid {
	ts := time.Now()
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

func (f *Factoid) Access(n base.Nick, c base.Chan) {
	f.Accessed.Timestamp = time.Now()
	f.Accessed.Nick = n
	f.Accessed.Chan = c
	f.Accessed.Count++
}

func (f *Factoid) Modify(n base.Nick, c base.Chan) {
	f.Modified.Timestamp = time.Now()
	f.Modified.Nick = n
	f.Modified.Chan = c
	f.Modified.Count++
}

// Factoids are stored in a mongo collection of Factoid structs
type FactoidCollection struct {
	// We're wrapping mgo.Collection so we can provide our own methods.
	*mgo.Collection

	// cache of objectIds for PseudoRand
	seen map[string][]bson.ObjectId
}

// Wrapper to get hold of a factoid collection handle
func Collection(dbh *mgo.Database) *FactoidCollection {
	fc := &FactoidCollection{
		Collection: dbh.C(COLLECTION),
		seen:       make(map[string][]bson.ObjectId),
	}
	err := fc.EnsureIndex(mgo.Index{Key: []string{"key"}})
	if err != nil {
		logging.Error("Couldn't create index on sp0rkle.factoids: %v", err)
	}
	return fc
}

// Can't call this Count because that'd override mgo.Collection.Count()
func (fc *FactoidCollection) GetCount(key string) int {
	if num, err := fc.Find(lookup(key)).Count(); err == nil {
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
	if err := fc.Find(lookup(key)).One(&res); err == nil {
		return &res
	}
	return nil
}

func (fc *FactoidCollection) GetPseudoRand(key string) *Factoid {
	lookup := lookup(key)
	ids, ok := fc.seen[key]
	if ok && len(ids) > 0 {
		logging.Debug("Seen '%s' before, %d stored id's", key, len(ids))
		lookup["_id"] = bson.M{"$nin": ids}
	}
	query := fc.Find(lookup)
	count, err := query.Count()
	if err != nil {
		logging.Debug("Counting for key failed: %v", err)
		return nil
	}
	if count == 0 {
		if ok {
			// we've seen this before, but people have deleted it since.
			delete(fc.seen, key)
		}
		return nil
	}
	var res Factoid
	if count > 1 {
		query = query.Skip(util.RNG.Intn(count))
	}
	if err = query.One(&res); err != nil {
		logging.Warn("Fetching factoid for key failed: %v", err)
		return nil
	}
	if count != 1 {
		if !ok {
			// only store seen for keys that have more than one factoid
			logging.Debug("Creating seen data for key '%s'.", key)
			fc.seen[key] = make([]bson.ObjectId, 0, count)
		}
		logging.Debug("Storing id %v for key '%s'.", res.Id, key)
		fc.seen[key] = append(fc.seen[key], res.Id)
	} else if ok {
		// if the count of results is 1 and we're storing seen data for key
		// then we've exhausted the possible results and should wipe it
		logging.Debug("Zeroing seen data for key '%s'.", key)
		delete(fc.seen, key)
	}
	return &res
}

func (fc *FactoidCollection) GetKeysMatching(regex string) []string {
	var res []string
	query := fc.Find(bson.M{"key": bson.M{"$regex": regex}})
	if err := query.Distinct("key", &res); err != nil {
		logging.Warn("Distinct regex query for '%s' failed: %v\n", regex, err)
		return nil
	}
	return res
}

func (fc *FactoidCollection) GetLast(op, key string) *Factoid {
	var res Factoid
	// op == "modified", "accessed", "created"
	op = op + ".timestamp"
	q := fc.Find(lookup(key)).Sort("-op")
	if err := q.One(&res); err == nil {
		return &res
	}
	return nil
}

func (fc *FactoidCollection) InfoMR(key string) *FactoidInfo {
	mr := &mgo.MapReduce{
		Map: `function() { emit("count", {
			accessed: this.accessed.count,
			modified: this.modified.count,
			created: this.created.count,
		})}`,
		Reduce: `function(k,l) {
			var sum = { accessed: 0, modified: 0, created: 0 };
			for each (var v in l) {
				sum.accessed += v.accessed;
				sum.modified += v.modified;
				sum.created  += v.created;
			}
			return sum;
		}`,
	}
	var res []struct {
		Id    int `bson:"_id"`
		Value FactoidInfo
	}
	info, err := fc.Find(lookup(key)).MapReduce(mr, &res)
	if err != nil || len(res) == 0 {
		logging.Warn("Info MR for '%s' failed: %v", key, err)
		return nil
	} else {
		logging.Debug("Info MR mapped %d, emitted %d, produced %d in %d ms.",
			info.InputCount, info.EmitCount, info.OutputCount, info.Time/1e6)
	}
	return &res[0].Value
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

// Shortcut to create correct lookup struct for mgo.Collection.Find().
// Returning an empty bson.M means key == "" can operate on all factoids.
func lookup(key string) bson.M {
	if key == "" {
		return bson.M{}
	}
	return bson.M{"key": key}
}
