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

var _ db.Keyer = (*Karma)(nil)

type Collection struct {
	db.C
}

func Init() *Collection {
	kc := &Collection{}
	kc.Init(db.Bolt.Keyed(), COLLECTION, nil)
	return kc
}

func (kc *Collection) KarmaFor(sub string) *Karma {
	res := &Karma{Key: strings.ToLower(sub)}
	if err := kc.Get(res.K(), res); err == nil {
		return res
	}
	return nil
}
