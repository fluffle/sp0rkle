package reminddriver

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/lib/reminders"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
	"labix.org/v2/mgo/bson"
	"time"
)

const driverName = "reminders"

type remindDriver struct {
	*reminders.ReminderCollection

	// We need to be able to kill reminder goroutines
	kill map[bson.ObjectId]chan bool

	// And it's useful to index them for deletion per-person
	list map[string][]bson.ObjectId

	l logging.Logger
}

func RemindDriver(db *db.Database, l logging.Logger) *remindDriver {
	rc := reminders.Collection(db, l)
	return &remindDriver{
		ReminderCollection: rc,
		kill: make(map[bson.ObjectId]chan bool),
		list: make(map[string][]bson.ObjectId),
		l: l,
	}
}

func (rd *remindDriver) Name() string {
	return driverName
}

func (rd *remindDriver) Remind(r *reminders.Reminder) func(*bot.Sp0rkle) {
	delta := r.RemindAt.Sub(time.Now())
	if delta < 0 {
		return nil
	}
	c := make(chan bool)
	rd.kill[r.Id] = c
	return func(b *bot.Sp0rkle) {
		select {
		case <-time.After(delta):
			b.Conn.Privmsg(r.Chan, r.Reply())
			if r.Target.Host != "" {
				// At the time of the reminder being created, target existed
				// TODO(fluffle): Tie this into state tracking properly.
				b.Conn.Privmsg(r.Target.Nick, r.Reply())
			}
			rd.Forget(r.Id, false)
		case <-c:
			return
		}
	}
}

func (rd *remindDriver) Forget(id bson.ObjectId, kill bool) {
	c, ok := rd.kill[id]
	if !ok { return }
	delete(rd.kill, id)
	if kill {
		c <- true
	}
	if err := rd.RemoveId(id); err != nil {
		rd.l.Error("Failure removing reminder %s: %v", id, err)
	}
}
