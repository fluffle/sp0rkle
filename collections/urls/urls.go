package urls

import (
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

const collection string = "urls"

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
	Id        bson.ObjectId `bson:"_id,omitempty"`
}

func NewUrl(u string, n bot.Nick, c bot.Chan) *Url {
	return &Url{
		Url:       u,
		Nick:      n,
		Chan:      c,
		Timestamp: time.Now(),
		Id:        bson.NewObjectId(),
	}
}

func (u Url) String() string {
	if u.CachedAs != "" {
		return fmt.Sprintf("%s (cached as %s at %s)",
			u.Url, u.CachedAs, u.CacheTime)
	} else if u.Shortened != "" {
		return fmt.Sprintf("%s (shortened as %s)", u.Url, u.Shortened)
	}
	return u.Url
}

type Collection struct {
	*mgo.Collection
}

func Init() *Collection {
	uc := &Collection{db.Init().C(collection)}
	err := uc.EnsureIndex(mgo.Index{Key: []string{"url"}, Unique: true})
	if err != nil {
		logging.Error("Couldn't create url index on sp0rkle.urls: %s", err)
	}
	for _, idx := range []string{"cachedas", "shortened"} {
		err := uc.EnsureIndex(mgo.Index{Key: []string{idx}})
		if err != nil {
			logging.Error("Couldn't create %s index on sp0rkle.urls: %s", idx, err)
		}
	}
	return uc
}

func (uc *Collection) GetById(id bson.ObjectId) *Url {
	var res Url
	if err := uc.Find(bson.M{"_id": id}).One(&res); err == nil {
		return &res
	}
	return nil
}

func (uc *Collection) GetByUrl(u string) *Url {
	var res Url
	if err := uc.Find(bson.M{"url": u}).One(&res); err == nil {
		return &res
	}
	return nil
}

// TODO(fluffle): thisisn't quite PseudoRand but still ...
func (uc *Collection) GetRand(regex string) *Url {
	lookup := bson.M{}
	if regex != "" {
		// Perform a regex lookup if we have one
		lookup["url"] = bson.M{"$regex": regex, "$options": "i"}
	}
	query := uc.Find(lookup)
	count, err := query.Count()
	if err != nil {
		logging.Warn("Count for URL lookup '%s' failed: %s", regex, err)
		return nil
	}
	if count == 0 {
		return nil
	}
	var res Url
	if count > 1 {
		query.Skip(util.RNG.Intn(count))
	}
	if err = query.One(&res); err != nil {
		logging.Warn("Fetch for URL lookup '%s' failed: %s", regex, err)
		return nil
	}
	return &res
}

func (uc *Collection) GetCached(c string) *Url {
	var res Url
	if err := uc.Find(bson.M{"cachedas": c}).One(&res); err == nil {
		return &res
	}
	return nil
}

func (uc *Collection) GetShortened(s string) *Url {
	var res Url
	if err := uc.Find(bson.M{"shortened": s}).One(&res); err == nil {
		return &res
	}
	return nil
}
