package karma

import (
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

const COLLECTION = "karma"
const TimeFormat = "15:04:05, Monday 2 January 2006"

type Karma struct {
	Subject   string
	Key       string
	Score     int
	Votes     int
	Upvoter   bot.Nick
	Upvtime   time.Time
	Downvoter bot.Nick
	Downvtime time.Time
}

func New(thing string) *Karma {
	return &Karma{
		Subject: thing,
		Key:     strings.ToLower(thing),
	}
}

func (k *Karma) Plus(who bot.Nick) {
	k.Score++
	k.Votes++
	k.Upvoter, k.Upvtime = who, time.Now()
}

func (k *Karma) Minus(who bot.Nick) {
	k.Score--
	k.Votes++
	k.Downvoter, k.Downvtime = who, time.Now()
}

func (k *Karma) String() string {
	s := fmt.Sprintf("'%s' has a karma of %d after %d votes.",
		k.Subject, k.Score, k.Votes)
	if k.Upvoter != "" {
		s += fmt.Sprintf(" Last upvoted by %s at %s.",
			k.Upvoter, k.Upvtime.Format(TimeFormat))
	}
	if k.Downvoter != "" {
		s += fmt.Sprintf(" Last downvoted by %s at %s.",
			k.Downvoter, k.Downvtime.Format(TimeFormat))
	}
	return s
}

func (k *Karma) Id() bson.M {
	return bson.M{"key": k.Key}
}

type Collection struct {
	*mgo.Collection
}

func Init() *Collection {
	kc := &Collection{db.Init().C(COLLECTION)}
	if err := kc.EnsureIndex(mgo.Index{
		Key:    []string{"key"},
		Unique: true,
	}); err != nil {
		logging.Error("Couldn't create index on karma.key: %s", err)
	}
	for _, key := range []string{"score", "votes"} {
		if err := kc.EnsureIndexKey(key); err != nil {
			logging.Error("Couldn't create index on karma.%s: %s", key, err)
		}
	}
	return kc
}

func (kc *Collection) KarmaFor(sub string) *Karma {
	var res Karma
	q := kc.Find(bson.M{"key": strings.ToLower(sub)})
	if err := q.One(&res); err == nil {
		return &res
	}
	return nil
}
