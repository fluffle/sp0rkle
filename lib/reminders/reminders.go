package reminders

import (
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/lib/db"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

const COLLECTION = "reminders"
const RemindTimeFormat = "15:04:05, Monday 2 January 2006"

type Reminder struct {
	db.StorableChan
	db.StorableNick
	Target    db.StorableNick
	Reminder  string
	Created   time.Time
	RemindAt  time.Time
	Id        bson.ObjectId `bson:"_id,omitempty"`
}

func NewReminder(r string, at time.Time, t, n db.StorableNick, c db.StorableChan) *Reminder {
	return &Reminder{
		StorableChan: c,
		StorableNick: n,
		Target: t,
		Reminder: r,
		Created: time.Now(),
		RemindAt: at,
		Id: bson.NewObjectId(),
	}
}

func (r *Reminder) Reply() string {
	if r.Nick == r.Target.Nick {
		return fmt.Sprintf("%s, you asked me to remind you %s",
			r.Nick, r.Reminder)
	}
	return fmt.Sprintf("%s, %s asked me to remind you %s",
		r.Target.Nick, r.Nick, r.Reminder)
}

func (r *Reminder) Acknowledge() string {
	if r.Nick == r.Target.Nick {
		return fmt.Sprintf("okay, i'll remind you %s at %s",
			r.Reminder, r.RemindAt.Format(RemindTimeFormat))
	}
	return fmt.Sprintf("okay, i'll remind %s %s at %s",
		r.Target.Nick, r.Reminder, r.RemindAt.Format(RemindTimeFormat))
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
	if err := rc.EnsureIndexKey("remindat"); err != nil {
		l.Error("Couldn't create index on sp0rkle.reminders: %v", err)
	}
	return rc
}

func (rc *ReminderCollection) LoadAndPrune() []*Reminder {
	// First, drop any reminders where RemindAt < time.Now()
	ci, err := rc.RemoveAll(bson.M{"remindat": bson.M{"$lt": time.Now()}})
	if err != nil {
		rc.l.Error("Pruning reminders returned error: %v", err)
	}
	if ci.Removed > 0 {
		rc.l.Info("Removed %d old reminders", ci.Removed)
	}
	// Now, load the remainder; the db is just used for persistence
	q := rc.Find(nil)
	count, err := q.Count()
	ret := make([]*Reminder, count)
	if err != nil {
		rc.l.Error("Counting reminders returned error: %v", err)
	}
	if count == 0 || err != nil {
		return []*Reminder{}
	}
	iter := q.Iter()
	i := 0
	for iter.Next(ret[i]) {
		i++
	}
	return ret
}
