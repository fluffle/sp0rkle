package urls

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/lib/db"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

const collection string = "urls"

type Url struct {
	db.StorableNick
	db.StorableChan
	Timestamp time.Time
	Url       string
	CachedAs  string
	CacheTime time.Time
	Hash      []byte
	MimeType  string
	Shortened string
	Id        bson.ObjectId "_id"
}

func NewUrl(u string, n db.StorableNick, c db.StorableChan) *Url {
	return &Url{
		Url:          u,
		StorableNick: n,
		StorableChan: c,
		Timestamp:    time.Now(),
		Id:           bson.NewObjectId(),
	}
}

type UrlCollection struct {
	*mgo.Collection
	l logging.Logger
}

func Collection(dbh *db.Database, l logging.Logger) *UrlCollection {
	uc := &UrlCollection{dbh.C(collection), l}
	for _, idx := range []string{"url", "cachedas", "shortened"} {
		err := uc.EnsureIndex(mgo.Index{Key: []string{idx}})
		if err != nil {
			l.Error("Couldn't create %s index on sp0rkle.urls: %s", idx, err)
		}
	}
	return uc
}

func (uc *UrlCollection) GetByUrl(u string) *Url {
	var res Url
	if err := uc.Find(bson.M{"url": u}).One(&res); err == nil {
		return &res
	}
	return nil
}

func (uc *UrlCollection) GetCached(c string) *Url {
	var res Url
	if err := uc.Find(bson.M{"cachedas": c}).One(&res); err == nil {
		return &res
	}
	return nil
}

func (uc *UrlCollection) GetShortened(s string) *Url {
	var res Url
	if err := uc.Find(bson.M{"shortened": s}).One(&res); err == nil {
		return &res
	}
	return nil
}


