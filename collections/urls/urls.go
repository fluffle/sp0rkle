package urls

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/bson"
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
	idxs := []db.Key{db.K{db.S{"url", u.Url}}}
	// Only add cachedas and shortened keys when the fields have values.
	if u.CachedAs != "" {
		idxs = append(idxs, db.K{db.S{"cachedas", u.CachedAs}})
	}
	if u.Shortened != "" {
		idxs = append(idxs, db.K{db.S{"shortened", u.Shortened}})
	}
	return idxs
}

func (u *Url) Id() bson.ObjectId {
	return u.Id_
}

func (u *Url) Exists() bool {
	return u != nil && len(u.Id_) > 0
}

func (u *Url) byId() db.K {
	return db.K{db.ID{u.Id_}}
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

type Collection struct {
	db.C
	seen map[string]map[bson.ObjectId]bool
}

func Init() *Collection {
	uc := &Collection{
		seen: make(map[string]map[bson.ObjectId]bool),
	}
	uc.Init(db.Bolt.Indexed(), COLLECTION, nil)
	if err := uc.Fsck(&Url{}); err != nil {
		logging.Fatal("urls fsck: %v", err)
	}
	return uc
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
			delete(uc.seen, regex)
		}
		return filtered[0]
	}
	// case count > 1, effectively
	// only store seen for regex that match more than one quote
	if !ok {
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
