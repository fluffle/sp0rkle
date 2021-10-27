package factoids

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	Id_                         bson.ObjectId `bson:"_id,omitempty"`
}

var _ db.Indexer = (*Factoid)(nil)

// Represent info about things that happened to the factoid
type FactoidStat struct {
	// When <thing> happened
	Timestamp time.Time
	// Who did <thing>
	Nick bot.Nick
	// Where they did <thing>
	Chan bot.Chan
	// How many times <thing> has been done before
	Count int
}

func (fs *FactoidStat) String() string {
	return fmt.Sprintf("%s,%s@%s#%d", fs.Nick, fs.Chan,
		fs.Timestamp.Format(time.RFC3339), fs.Count)
}

// Represent info about things that can be done to the factoid
type FactoidPerms struct {
	ReadOnly bool
	Nick     bot.Nick
}

func (fp *FactoidPerms) String() string {
	if fp.ReadOnly {
		return string(fp.Nick) + "(ro)"
	}
	return string(fp.Nick)
}

// Represent info returned from the Info MapReduce
type FactoidInfo struct {
	Created, Modified, Accessed int
}

// Helper to make the work of putting together a completely new *Factoid easier
func NewFactoid(key, value string, n bot.Nick, c bot.Chan) *Factoid {
	ts := time.Now()
	ft, fv := ParseValue(value)
	return &Factoid{
		Key: key, Value: fv, Type: ft, Chance: 1.0,
		Created:  &FactoidStat{ts, n, c, 1},
		Modified: &FactoidStat{ts, n, c, 0},
		Accessed: &FactoidStat{ts, n, c, 0},
		Perms:    &FactoidPerms{false, n},
		Id_:      bson.NewObjectId(),
	}
}

func (f *Factoid) String() string {
	return fmt.Sprintf("<%s/%d>=%q (%g%%) c=%s/m=%s/a=%s owner=%s",
		f.Key, f.Type, f.Value, f.Chance,
		f.Created, f.Modified, f.Accessed, f.Perms)
}

func (f *Factoid) Access(n bot.Nick, c bot.Chan) {
	f.Accessed.Timestamp = time.Now()
	f.Accessed.Nick = n
	f.Accessed.Chan = c
	f.Accessed.Count++
}

func (f *Factoid) Modify(n bot.Nick, c bot.Chan) {
	f.Modified.Timestamp = time.Now()
	f.Modified.Nick = n
	f.Modified.Chan = c
	f.Modified.Count++
}

func (f *Factoid) Id() bson.ObjectId {
	return f.Id_
}

func (f *Factoid) Indexes() []db.Key {
	// Factoids are indexed because f.Key is not unique and looking up
	// all the factoids with a given key is a common operation.
	// It's a bit wasteful sticking the ObjectId in the key when its
	// also the value, but we need a unique key name inside the "key" bucket.
	// A more optimal and less lazy solution might involve using bucket
	// sequences to provide the keys inside the "key" bucket, but meh.
	return []db.Key{
		db.K{db.S{"key", f.Key}, db.S{"v", string(f.Id_)}},
	}
}

func (f *Factoid) byId() db.K {
	return db.K{db.S{"_id", string(f.Id_)}}
}

func byKey(key string) db.K {
	k := db.K{}
	if key != "" {
		k = append(k, db.S{"key", key})
	}
	return k
}

type Factoids []*Factoid

func (fs Factoids) Strings() []string {
	s := make([]string, len(fs))
	for i, f := range fs {
		// We can't use %#v here because a Factoid struct contains pointers.
		s[i] = f.String()
	}
	return s
}

type migrator struct {
	mongo, bolt db.Collection
}

func (m migrator) MigrateTo(newState db.MigrationState) error {
	if newState != db.MONGO_PRIMARY {
		return nil
	}
	var all Factoids
	if err := m.mongo.All(db.K{}, &all); err != nil {
		return err
	}
	if err := m.bolt.BatchPut(all); err != nil {
		logging.Error("Migrating factoids: %v", err)
		return err
	}
	logging.Info("Migrated %d factoid entries.", len(all))
	return nil
}

func (m migrator) Diff() ([]string, []string, error) {
	var mAll, bAll Factoids
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

	// cache of objectIds for PseudoRand
	seen map[string]map[bson.ObjectId]bool
}

// Wrapper to get hold of a factoid collection handle
func Init() *Collection {
	fc := &Collection{
		Both: db.Both{},
		seen: make(map[string]map[bson.ObjectId]bool),
	}
	fc.Both.MongoC.Init(db.Mongo, COLLECTION, mongoIndexes)
	fc.Both.BoltC.Init(db.Bolt.Indexed(), COLLECTION, nil)
	fc.Both.Debug(true)
	m := &migrator{
		mongo: fc.Both.MongoC,
		bolt:  fc.Both.BoltC,
	}
	fc.Both.Checker.Init(m, COLLECTION)
	return fc
}

func mongoIndexes(c db.Collection) {
	err := c.Mongo().EnsureIndex(mgo.Index{Key: []string{"key"}})
	if err != nil {
		logging.Error("Couldn't create index on sp0rkle.factoids: %v", err)
	}
}

// Can't call this Count because that'd override mgo.Collection.Count()
func (fc *Collection) GetCount(key string) int {
	// TODO(fluffle): less-wasteful GetCount()
	return len(fc.GetAll(key))
}

func (fc *Collection) GetById(id bson.ObjectId) *Factoid {
	res := &Factoid{Id_: id}
	if err := fc.Get(res.byId(), res); err != nil {
		logging.Warn("Factoid GetById failed: %v", err)
		return nil
	}
	return res
}

func (fc *Collection) GetAll(key string) []*Factoid {
	res := Factoids{}
	if err := fc.All(byKey(key), &res); err != nil {
		logging.Warn("Factoid GetAll failed: %v", err)
		return nil
	}
	return res
}

func (fc *Collection) GetPseudoRand(key string) *Factoid {
	// TODO(fluffle): GetPR implementation in package db.
	facts := fc.GetAll(key)
	filtered := Factoids{}
	ids, ok := fc.seen[key]
	if ok && len(ids) > 0 {
		logging.Debug("Seen '%s' before, %d stored id's", key, len(ids))
		for _, fact := range facts {
			if !ids[fact.Id()] {
				filtered = append(filtered, fact)
			}
		}
	} else {
		filtered = facts
	}

	count := len(filtered)
	switch count {
	case 0:
		if ok {
			// we've seen this before, but people have deleted it since.
			delete(fc.seen, key)
		}
		return nil
	case 1:
		if ok {
			// if the count of results is 1 and we're storing seen data for key
			// then we've exhausted the possible results and should wipe it
			logging.Debug("Zeroing seen data for key '%s'.", key)
			delete(fc.seen, key)
		}
		return filtered[0]
	}
	// case count > 1
	if !ok {
		// only store seen for keys that have more than one factoid
		logging.Debug("Creating seen data for key '%s'.", key)
		fc.seen[key] = make(map[bson.ObjectId]bool)
	}
	res := filtered[rand.Intn(count)]
	logging.Debug("Storing id %v for key '%s'.", res.Id(), key)
	fc.seen[key][res.Id()] = true
	return res
}

func (fc *Collection) GetKeysMatching(regex string) []string {
	facts := Factoids{}
	if err := fc.Match("Key", regex, &facts); err != nil {
		logging.Warn("Factoid GetKeyMatching failed: %v", err)
		return nil
	}
	// Have to dedupe here now :-/
	res := []string{}
	set := map[string]bool{}
	for _, fact := range facts {
		if _, ok := set[fact.Key]; !ok {
			set[fact.Key] = true
			res = append(res, fact.Key)
		}
	}
	return res
}

func (fc *Collection) GetLast(key string) (c *Factoid, m *Factoid, a *Factoid) {
	// Waaay less efficient for MongoDB but works for both.
	facts := fc.GetAll(key)
	for _, fact := range facts {
		if c == nil || c.Created.Timestamp.Before(fact.Created.Timestamp) {
			c = fact
		}
		if m == nil || m.Modified.Timestamp.Before(fact.Modified.Timestamp) {
			m = fact
		}
		if a == nil || a.Accessed.Timestamp.Before(fact.Accessed.Timestamp) {
			a = fact
		}
	}
	return
}

func (fc *Collection) InfoMR(key string) *FactoidInfo {
	// MapReduce has no BoltDB equivalent and building one seems excessive.
	minfo := &FactoidInfo{}
	binfo := &FactoidInfo{}
	state := fc.Check()

	if state < db.BOLT_ONLY {
		// Mongo
		mr := &mgo.MapReduce{
			Map: `function() { emit("count", {
				accessed: this.accessed.count,
				modified: this.modified.count,
				created: this.created.count
			})}`,
			Reduce: `function(k,l) {
				var sum = { accessed: 0, modified: 0, created: 0 };
				for (var i = 0; i < l.length; i++) {
					sum.accessed += l[i].accessed;
					sum.modified += l[i].modified;
					sum.created  += l[i].created;
				}
				return sum;
			}`,
		}
		var res []struct {
			Id    int `bson:"_id"`
			Value FactoidInfo
		}
		k := bson.M{}
		if key != "" {
			k["key"] = key
		}
		info, err := fc.Mongo().Find(k).MapReduce(mr, &res)
		if err != nil || len(res) == 0 {
			logging.Warn("Info MR for '%s' failed: %v", key, err)
		} else {
			logging.Debug("Info MR mapped %d, emitted %d, produced %d in %d ms.",
				info.InputCount, info.EmitCount, info.OutputCount, info.Time/1e6)
			*minfo = res[0].Value
		}
	}
	if state > db.MONGO_ONLY {
		// Bolt, we have to do things manually, which is way easier even if it
		// does involve maybe slurping all the factoids into a slice, *again*.
		// TODO(fluffle): Add a ForEach() to boltdb wrapper once migrated.
		facts := Factoids{}
		if err := fc.Both.BoltC.All(byKey(key), &facts); err != nil {
			logging.Warn("Factoid InfoMR All failed: %v", err)
		}

		for _, fact := range facts {
			binfo.Accessed += fact.Accessed.Count
			binfo.Modified += fact.Modified.Count
			binfo.Created += fact.Created.Count
		}
	}
	if (state == db.MONGO_PRIMARY || state == db.BOLT_PRIMARY) &&
		(binfo.Accessed != minfo.Accessed ||
			binfo.Modified != minfo.Modified ||
			binfo.Created != minfo.Created) {
		logging.Warn("Factoid InfoMR: diff detected!\n\tMongo: %v\n\tBolt: %v", minfo, binfo)
	}
	if state >= db.BOLT_PRIMARY {
		return binfo
	}
	return minfo
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
