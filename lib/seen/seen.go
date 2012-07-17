package seen 

import (
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/lib/db"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

const COLLECTION string = "seen"

type Nick struct {
	db.StorableNick
	db.StorableChan
	Timestamp time.Time
	Key		  string
	Action    string
	Text      string
// Upsert doesn't like it much if we send _id fields when updating:
//   "cannot change _id of a document"
// Since there doesn't appear to be a way to make the _id accessible via
// (de)serialization and *not* send it for an Upsert, just hide it.
//	Id        bson.ObjectId "_id"
}

type seenMsg func(*Nick) string

var actionMap map[string]seenMsg = map[string]seenMsg{
	"PRIVMSG": func(n *Nick) string { 
		return fmt.Sprintf("in %s, saying '%s'", n.Chan, n.Text)},
	"ACTION": func(n *Nick) string { 
		return fmt.Sprintf("in %s, saying '%s %s'", n.Chan, n.Nick, n.Text)},
	"JOIN": func(n *Nick) string {
		return fmt.Sprintf("joining %s", n.Chan)},
	"PART": func(n *Nick) string {
		return fmt.Sprintf("parting %s with the message '%s'", n.Chan, n.Text)},
	"KICK": func(n *Nick) string {
		return fmt.Sprintf("bein kicked from %s with the message '%s'", n.Chan, n.Text)},
	"QUIT": func(n *Nick) string {
		return fmt.Sprintf("quitting with the message '%s'", n.Text)},
	"NICK": func(n *Nick) string {
		return fmt.Sprintf("changing their nick to '%s'", n.Text)},
	"SMOKE": func(n *Nick) string { return "going for a smoke." },
}

func SawNick(nick db.StorableNick, ch db.StorableChan, act, txt string) *Nick {
	return &Nick{nick, ch, time.Now(), strings.ToLower(nick.Nick), act, txt}
}

func (n *Nick) Index() bson.M {
	return bson.M{"key": n.Key, "action": n.Action}
}

func (n *Nick) String() string {
	if act, ok := actionMap[n.Action]; ok {
		return fmt.Sprintf("I last saw %s on %s (%s ago), %s.",
			n.Nick, n.Timestamp.Format(time.RFC1123),
			time.Now().Sub(n.Timestamp), act(n))
	}
	// No specific message format for the action seen.
	return fmt.Sprintf("I last saw %s at %s (%s ago).", 
		n.Nick, n.Timestamp.Format(time.RFC1123),
		time.Now().Sub(n.Timestamp))
}

type SeenCollection struct {
	// Wrap mgo.Collection
	*mgo.Collection

	// logging object
	l logging.Logger
}

func Collection(dbh *db.Database, l logging.Logger) *SeenCollection {
	sc := &SeenCollection{
		Collection: dbh.C(COLLECTION),
		l:          l,
	}
	err := sc.EnsureIndex(mgo.Index{
		Key: []string{"key", "action"},
		Unique: true,
	})
	if err != nil {
		l.Error("Couldn't create index on sp0rkle.seen: %v", err)
	}
	err = sc.EnsureIndexKey([]string{"key"})
	if err != nil {
		l.Error("Couldn't create index on sp0rkle.seen: %v", err)
	}
	return sc
}

func (sc *SeenCollection) LastSeen(nick string) *Nick {
	var res Nick
	q := sc.Find(bson.M{"key": strings.ToLower(nick)}).Sort("-timestamp")
	if err := q.One(&res); err == nil {
		return &res
	}
	return nil
}

func (sc *SeenCollection) LastSeenDoing(nick, act string) *Nick {
	var res Nick
	q := sc.Find(bson.M{"key": strings.ToLower(nick), "action": act}).Sort("-timestamp")
	if err := q.One(&res); err == nil {
		return &res
	}
	return nil
}

func (sc *SeenCollection) SeenAnyMatching(rx string) []string {
	var res []string
	q := sc.Find(bson.M{"key": bson.M{"$regex": rx, "$options": "i"}}).Sort("-timestamp")
	if err := q.Distinct("key", &res); err != nil {
		sc.l.Warn("SeenAnyMatching Find error: %v", err)
		return []string{}
	}
	sc.l.Debug("Looked for matches, found %#v", res)
	return res
}
