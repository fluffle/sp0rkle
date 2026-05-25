package urls

import (
	"fmt"
	"math/rand"
	"regexp"
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
	return u != nil && len(u.Id_) > 0 && u.Url != ""
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

func (u *Url) CollectionName() string { return "urls" }

type Collection struct {
	collection db.IndexedCollection[*Url]
	seen       map[string]map[bson.ObjectId]bool
}

func New(collection db.IndexedCollection[*Url]) *Collection {
	return &Collection{
		collection: collection,
		seen:       make(map[string]map[bson.ObjectId]bool),
	}
}

func (uc *Collection) GetById(id bson.ObjectId) *Url {
	res := &Url{Id_: id}
	url, err := uc.collection.Get(res.byId())
	if err != nil || !url.Exists() {
		return nil
	}
	return url
}

func (uc *Collection) GetByUrl(u string) *Url {
	res := &Url{Url: u}
	url, err := uc.collection.Get(res.byUrl())
	if err != nil || !url.Exists() {
		return nil
	}
	return url
}

func (uc *Collection) GetRand(regex string) *Url {
	var urls []*Url
	var err error
	if regex == "" {
		urls, err = uc.collection.All(db.K{})
	} else {
		rx := regexp.MustCompile("(?i)" + regex)
		urls, err = uc.collection.Match(func(u *Url) bool {
			return rx.MatchString(u.Url)
		})
	}
	if err != nil {
		logging.Warn("URL Match() failed: %v", err)
		return nil
	}

	filtered := []*Url{}
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
			delete(uc.seen, regex)
		}
		return nil
	case 1:
		if ok {
			delete(uc.seen, regex)
		}
		return filtered[0]
	}
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
	url, err := uc.collection.Get(res.byCachedAs())
	if err != nil || !url.Exists() {
		return nil
	}
	return url
}

func (uc *Collection) GetShortened(s string) *Url {
	res := &Url{Shortened: s}
	url, err := uc.collection.Get(res.byShortened())
	if err != nil || !url.Exists() {
		return nil
	}
	return url
}

func (uc *Collection) Put(value *Url) error {
	return uc.collection.Put(value)
}

func (uc *Collection) Del(value *Url) error {
	return uc.collection.Del(value)
}
