package urls

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION string = "urls"

type Url struct {
	Nick      bot.Nick
	Chan      bot.Chan
	Timestamp time.Time
	Url       string
	CachedAs  string
	CacheTime time.Time
	Hash      []byte
	MimeType  string
	Shortened string
	Id_       bson.ObjectId `bson:"_id,omitempty"`
}

var _ db.Indexer = (*Url)(nil)

func NewUrl(u string, n bot.Nick, c bot.Chan) *Url {
	return &Url{
		Url:       u,
		Nick:      n,
		Chan:      c,
		Timestamp: time.Now(),
		Id_:       bson.NewObjectId(),
	}
}

func (u *Url) String() string {
	if u.CachedAs != "" {
		return fmt.Sprintf("%s (cached as %s at %s)",
			u.Url, u.CachedAs, u.CacheTime)
	} else if u.Shortened != "" {
		return fmt.Sprintf("%s (shortened as %s)", u.Url, u.Shortened)
	}
	return u.Url
}

func (u *Url) Indexes() []db.Key {
	return []db.Key{
		db.K{db.S{"url", u.Url}},
		db.K{db.S{"cachedas", u.CachedAs}},
		db.K{db.S{"shortened", u.Shortened}},
	}
}

func (u *Url) Id() bson.ObjectId {
	return u.Id_
}

func (u *Url) Exists() bool {
	return u != nil && len(u.Id_) > 0
}

func (u *Url) byId() db.K {
	return db.K{db.S{"_id", string(u.Id_)}}
}

func (u *Url) byUrl() db.K {
	return db.K{db.S{"url", u.Url}}
}

func (u *Url) byCachedAs() db.K {
	return db.K{db.S{"cachedas", u.CachedAs}}
}

func (u *Url) byShortened() db.K {
	return db.K{db.S{"shortened", u.Shortened}}
}

type Urls []*Url

func (us Urls) Strings() []string {
	s := make([]string, len(us))
	for i, u := range us {
		s[i] = fmt.Sprintf("%#v", u)
	}
	return s
}

type migrator struct {
	mongo, bolt db.Collection
}

func (m *migrator) Migrate() error {
	var all []*Url
	if err := m.mongo.All(db.K{}, &all); err != nil {
		return err
	}
	if err := m.bolt.BatchPut(all); err != nil {
		logging.Error("Migrating urls: %v", err)
		return err
	}
	logging.Debug("Migrated %d urls.", len(all))
	return nil
}

func (m *migrator) Diff() ([]string, []string, error) {
	var mAll, bAll Urls
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

	// Cache of ObjectId's for GetRand.
	seen map[string]map[bson.ObjectId]bool
}

func Init() *Collection {
	uc := &Collection{
		Both: db.Both{},
		seen: make(map[string]map[bson.ObjectId]bool),
	}
	uc.Both.MongoC.Init(db.Mongo, COLLECTION, mongoIndexes)
	uc.Both.BoltC.Init(db.Bolt, COLLECTION, nil)
	m := &migrator{
		mongo: uc.Both.MongoC,
		bolt:  uc.Both.BoltC,
	}
	uc.Both.Checker.Init(m, COLLECTION)
	return uc
}

func mongoIndexes(c db.Collection) {
	err := c.Mongo().EnsureIndex(mgo.Index{Key: []string{"url"}, Unique: true})
	if err != nil {
		logging.Error("Couldn't create url index on sp0rkle.urls: %s", err)
	}
	for _, idx := range []string{"cachedas", "shortened"} {
		err := c.Mongo().EnsureIndex(mgo.Index{Key: []string{idx}})
		if err != nil {
			logging.Error("Couldn't create %s index on sp0rkle.urls: %s", idx, err)
		}
	}
}

func (uc *Collection) GetById(id bson.ObjectId) *Url {
	res := &Url{Id_: id}
	if err := uc.Get(res.byId(), res); err == nil && res.Exists() {
		return res
	}
	return nil
}

func (uc *Collection) GetByUrl(u string) *Url {
	res := &Url{Url: u}
	if err := uc.Get(res.byUrl(), res); err == nil && res.Exists() {
		return res
	}
	return nil
}

// TODO(fluffle): Dedupe with quotes and other pseudo-rand implementations.
// Comments in quotes collection about efficiency apply here too.
func (uc *Collection) GetRand(regex string) *Url {
	urls := Urls{}
	if regex == "" {
		if err := uc.All(db.K{}, &urls); err != nil {
			logging.Warn("URL All() failed: %v", err)
			return nil
		}
	} else {
		if err := uc.Match("Url", regex, &urls); err != nil {
			logging.Warn("URL Match(%q) failed: %v", regex, err)
			return nil
		}
	}

	filtered := Urls{}
	ids, ok := uc.seen[regex]
	if ok && len(ids) > 0 {
		logging.Debug("Looked for URLs matching %q before, %d stored id's", regex, len(ids))
		for _, url := range urls {
			if !ids[url.Id_] {
				filtered = append(filtered, url)
			}
		}
	} else {
		filtered = urls
	}

	count := len(filtered)
	switch count {
	case 0:
		if ok {
			// Looked for this regex before, but nothing matches now
			delete(uc.seen, regex)
		}
		return nil
	case 1:
		if ok {
			// if the count of results is 1 and we're storing seen data for regex
			// then we've exhausted the possible results and should wipe it
			logging.Debug("Zeroing seen data for regex %q.", regex)
			delete(uc.seen, regex)
		}
		return filtered[0]
	}
	// case count > 1:
	if !ok {
		// only store seen for regex that match more than one quote
		logging.Debug("Creating seen data for regex %q.", regex)
		uc.seen[regex] = map[bson.ObjectId]bool{}
	}
	url := filtered[rand.Intn(count)]
	logging.Debug("Storing id %v for regex %q.", url.Id_, regex)
	uc.seen[regex][url.Id_] = true
	return url
}

func (uc *Collection) GetCached(c string) *Url {
	res := &Url{CachedAs: c}
	if err := uc.Get(res.byCachedAs(), res); err == nil && res.Exists() {
		return res
	}
	return nil
}

func (uc *Collection) GetShortened(s string) *Url {
	res := &Url{Shortened: s}
	if err := uc.Get(res.byShortened(), res); err == nil && res.Exists() {
		return res
	}
	return nil
}
