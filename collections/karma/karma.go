package karma

import (
	"fmt"
	"strings"
	"time"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/datetime"
)

const COLLECTION = "karma"

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

func NewKarma(thing string) *Karma {
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
			k.Upvoter, datetime.Format(k.Upvtime))
	}
	if k.Downvoter != "" {
		s += fmt.Sprintf(" Last downvoted by %s at %s.",
			k.Downvoter, datetime.Format(k.Downvtime))
	}
	return s
}

func (k *Karma) K() db.Key {
	return db.K{db.S{"key", k.Key}}
}

func (k *Karma) CollectionName() string { return "karma" }

var _ db.Keyer = (*Karma)(nil)

type Collection struct {
	collection db.KeyedCollection[*Karma]
}

func New(collection db.KeyedCollection[*Karma]) *Collection {
	return &Collection{collection: collection}
}

func (kc *Collection) KarmaFor(sub string) *Karma {
	karma, err := kc.collection.Get(db.K{db.S{"key", strings.ToLower(sub)}})
	if err != nil || karma == nil {
		return nil
	}
	return karma
}

func (kc *Collection) Put(k *Karma) error {
	return kc.collection.Put(k)
}

func (kc *Collection) Del(k *Karma) error {
	return kc.collection.Del(k)
}
