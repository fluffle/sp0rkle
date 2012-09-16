package reminders

import (
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/lib/db"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

const COLLECTION = "reminders"
const RemindTimeFormat = "15:04:05, Monday 2 January 2006"

type Reminder struct {
	db.StorableChan
	db.StorableNick
	Target    db.StorableNick
	From, To  string
	Reminder  string
	Created   time.Time
	RemindAt  time.Time
	Tell      bool
	Id        bson.ObjectId `bson:"_id,omitempty"`
}

func NewReminder(r string, at time.Time, t, n db.StorableNick, c db.StorableChan) *Reminder {
	return &Reminder{
		StorableChan: c,
		StorableNick: n,
		Target: t,
		From: strings.ToLower(n.Nick),
		To: strings.ToLower(t.Nick),
		Reminder: r,
		Created: time.Now(),
		RemindAt: at,
		Tell: false,
		Id: bson.NewObjectId(),
	}
}

func NewTell(msg string, t, n db.StorableNick, c db.StorableChan) *Reminder {
	return &Reminder{
		StorableNick: n,
		StorableChan: c,
		Target: t,
		From: strings.ToLower(n.Nick),
		To: strings.ToLower(t.Nick),
		Reminder: msg,
		Created: time.Now(),
		Tell: true,
		Id: bson.NewObjectId(),
	}
}

func (r *Reminder) Reply() string {
	if r.Tell {
		return fmt.Sprintf("%s asked me to tell you %s at %s",
			r.Nick, r.Reminder, r.Created.Format(RemindTimeFormat))
	} else if r.Nick == r.Target.Nick {
		return fmt.Sprintf("%s, you asked me to remind you %s",
			r.Nick, r.Reminder)
	}
	return fmt.Sprintf("%s, %s asked me to remind you %s",
		r.Target.Nick, r.Nick, r.Reminder)
}

func (r *Reminder) Acknowledge() string {
	if r.Tell {
		return fmt.Sprintf("okay, i'll tell %s %s when I see them",
			r.Target.Nick, r.Reminder)
	} else if r.Nick == r.Target.Nick {
		return fmt.Sprintf("okay, i'll remind you %s at %s",
			r.Reminder, r.RemindAt.Format(RemindTimeFormat))
	}
	return fmt.Sprintf("okay, i'll remind %s %s at %s",
		r.Target.Nick, r.Reminder, r.RemindAt.Format(RemindTimeFormat))
}

func (r *Reminder) List(nick string) string {
	nick = strings.ToLower(nick)
	if r.Tell {
		if nick == r.From {
			return fmt.Sprintf("you asked me to tell %s %s",
				r.Target.Nick, r.Reminder)
		} else if nick == r.To {
			// this is somewhat unlikely, as it should have triggered already
			return fmt.Sprintf("%s asked me to tell you %s -- and now I have!",
				r.Nick, r.Reminder)
		}
	} else if nick == r.From && nick == r.To {
		return fmt.Sprintf("you asked me to remind you %s, at %s",
			r.Reminder, r.RemindAt.Format(RemindTimeFormat))
	} else if nick == r.From {
		return fmt.Sprintf("you asked me to remind %s %s, at %s",
			r.Target.Nick, r.Reminder, r.RemindAt.Format(RemindTimeFormat))
	} else if nick == r.To {
		return fmt.Sprintf("%s asked me to remind you %s, at %s",
			r.Nick, r.Reminder, r.RemindAt.Format(RemindTimeFormat))
	}
	return fmt.Sprintf("%s asked me to remind %s %s, at %s",
		r.Nick, r.Target, r.Reminder, r.RemindAt.Format(RemindTimeFormat))
}

type ReminderCollection struct {
	*mgo.Collection
	l logging.Logger
}

func Collection(dbh *db.Database, l logging.Logger) *ReminderCollection {
	rc := &ReminderCollection{
		Collection: dbh.C(COLLECTION),
		l: l,
	}
	for _, k := range []string{"remindat", "from", "to", "tell"} {
		if err := rc.EnsureIndexKey(k); err != nil {
			l.Error("Couldn't create %s index on sp0rkle.reminders: %v", k, err)
		}
	}
	return rc
}

func (rc *ReminderCollection) LoadAndPrune() []*Reminder {
	// First, drop any reminders where RemindAt < time.Now()
	ci, err := rc.RemoveAll(bson.M{"$and": []bson.M{
		{"remindat": bson.M{"$lt": time.Now()}},
		{"tell": false},
	}})
	if err != nil {
		rc.l.Error("Pruning reminders returned error: %v", err)
	}
	if ci.Removed > 0 {
		rc.l.Info("Removed %d old reminders", ci.Removed)
	}
	// Now, load the remainder; the db is just used for persistence
	q := rc.Find(bson.M{"tell": false})
	ret := make([]*Reminder, 0)
	if err := q.All(&ret); err != nil {
		rc.l.Error("Loading reminders returned error: %v", err)
		return nil
	}
	return ret
}

func (rc *ReminderCollection) RemindersFor(nick string) []*Reminder {
	nick = strings.ToLower(nick)
	q := rc.Find(bson.M{"$or": []bson.M{{"from": nick}, {"to": nick}}})
	q.Sort("remindat")
	ret := make([]*Reminder, 0)
	if err := q.All(&ret); err != nil {
		rc.l.Error("Loading reminders for %s returned error: %v", nick, err)
		return nil
	}
	return ret
}
