package reminders

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/datetime"
	"github.com/fluffle/sp0rkle/util/bson"
)

const COLLECTION = "reminders"

type Reminder struct {
	Source   bot.Nick
	Target   bot.Nick
	Chan     bot.Chan
	From, To string
	Reminder string
	Created  time.Time
	RemindAt time.Time
	Tell     bool
	Id_      bson.ObjectId `bson:"_id,omitempty"`
}

var _ db.Indexer = (*Reminder)(nil)

func NewReminder(r string, at time.Time, t, n bot.Nick, c bot.Chan) *Reminder {
	return &Reminder{
		Source:   n,
		Target:   t,
		Chan:     c,
		From:     n.Lower(),
		To:       t.Lower(),
		Reminder: r,
		Created:  time.Now(),
		RemindAt: at,
		Tell:     false,
		Id_:      bson.NewObjectId(),
	}
}

func NewTell(msg string, t, n bot.Nick, c bot.Chan) *Reminder {
	return &Reminder{
		Chan:     c,
		Source:   n,
		Target:   t,
		From:     n.Lower(),
		To:       t.Lower(),
		Reminder: msg,
		Created:  time.Now(),
		Tell:     true,
		Id_:      bson.NewObjectId(),
	}
}

func (r *Reminder) Indexes() []db.Key {
	// Reminders and Tells behave differently and we need to retrieve them
	// separately from each other, so the first level index is on Tell.
	// From and To are not unique so we use a millisecond timestamp from
	// the reminder to differentiate and sort. Tells don't set RemindAt,
	// so we use the create timestamp instead.
	ts := r.RemindAt
	if r.Tell {
		ts = r.Created
	}
	return []db.Key{
		db.K{db.T{"tell", r.Tell}, db.S{"from", r.From}, db.TS{"ts", ts}},
		db.K{db.T{"tell", r.Tell}, db.S{"to", r.To}, db.TS{"ts", ts}},
	}
}

func (r *Reminder) Id() bson.ObjectId {
	return r.Id_
}

func (r *Reminder) byId() db.K {
	return db.K{db.ID{r.Id_}}
}

func tellTo(nick string) db.K {
	return db.K{db.T{"tell", true}, db.S{"to", nick}}
}

func remindFrom(nick string) db.K {
	return db.K{db.T{"tell", false}, db.S{"from", nick}}
}

func remindTo(nick string) db.K {
	return db.K{db.T{"tell", false}, db.S{"to", nick}}
}

func (r *Reminder) At() string {
	return datetime.Format(r.RemindAt)
}

func (r *Reminder) Reply() (s string) {
	switch {
	case r.Tell:
		s = fmt.Sprintf("%s asked me to tell you %s", r.Source, r.Reminder)
	case r.From == r.To:
		s = fmt.Sprintf("%s, you asked me to remind you %s",
			r.Source, r.Reminder)
	default:
		s = fmt.Sprintf("%s, %s asked me to remind you %s",
			r.Target, r.Source, r.Reminder)
	}
	return
}

func (r *Reminder) Acknowledge() (s string) {
	switch {
	case r.Tell:
		s = fmt.Sprintf("okay, i'll tell %s %s when I see them",
			r.Target, r.Reminder)
	case r.From == r.To:
		s = fmt.Sprintf("okay, i'll remind you %s at %s",
			r.Reminder, r.At())
	default:
		s = fmt.Sprintf("okay, i'll remind %s %s at %s",
			r.Target, r.Reminder, r.At())
	}
	return
}

func (r *Reminder) List(nick string) (s string) {
	nick = strings.ToLower(nick)
	switch {
	case r.Tell && nick == r.From:
		s = fmt.Sprintf("you asked me to tell %s %s",
			r.Target, r.Reminder)
	case r.Tell && nick == r.To:
		// this is somewhat unlikely, as it should have triggered already
		s = fmt.Sprintf("%s asked me to tell you %s -- and now I have!",
			r.Source, r.Reminder)
	case nick == r.From && nick == r.To:
		s = fmt.Sprintf("you asked me to remind you %s, at %s",
			r.Reminder, r.At())
	case nick == r.From:
		s = fmt.Sprintf("you asked me to remind %s %s, at %s",
			r.Target, r.Reminder, r.At())
	case nick == r.To:
		s = fmt.Sprintf("%s asked me to remind you %s, at %s",
			r.Source, r.Reminder, r.At())
	default:
		s = fmt.Sprintf("%s asked me to remind %s %s, at %s",
			r.Source, r.Target, r.Reminder, r.At())
	}
	return
}

func (r *Reminder) CollectionName() string { return "reminders" }

func byRemindAt(a, b *Reminder) int {
	return a.RemindAt.Compare(b.RemindAt)
}

type Collection struct {
	collection db.IndexedCollection[*Reminder]
}

func New(collection db.IndexedCollection[*Reminder]) *Collection {
	return &Collection{collection: collection}
}

func (rc *Collection) GetById(id bson.ObjectId) *Reminder {
	r, err := rc.collection.Get(db.K{db.ID{id}})
	if err != nil {
		logging.Error("Reminder GetById(%s) failed: %v", id, err)
		return nil
	}
	return r
}

func (rc *Collection) LoadAndPrune() []*Reminder {
	all, err := rc.collection.All(db.K{db.T{"tell", false}})
	if err != nil {
		logging.Error("Loading all reminders: %v", err)
		return nil
	}
	slices.SortFunc(all, byRemindAt)
	now := time.Now()
	last := -1
	for i, r := range all {
		if r.RemindAt.After(now) {
			last = i
			break
		}
	}

	// If no future reminder was found, all reminders are expired
	if last == -1 && len(all) > 0 {
		last = len(all)
	}

	if last > 0 {
		for _, r := range all[:last] {
			if err := rc.collection.Del(r); err != nil {
				logging.Error("Deleting expired reminder %v (expiry %s): %v", r.Id_, r.At(), err)
			}
		}
		all = all[last:]
		logging.Info("Removed %d old reminders", last)
	}
	return all
}

func (rc *Collection) RemindersFor(nick string) []*Reminder {
	// Note: remindFrom and remindTo queries both filter by Tell=false,
	// so tells are excluded from the results. The dedup below only removes
	// reminders where the sender and receiver are the same nick.
	nick = strings.ToLower(nick)
	from, err := rc.collection.All(remindFrom(nick))
	if err != nil {
		logging.Error("Loading reminders from %s returned error: %v", nick, err)
	}
	to, err := rc.collection.All(remindTo(nick))
	if err != nil {
		logging.Error("Loading reminders to %s returned error: %v", nick, err)
	}
	if len(from) == 0 && len(to) == 0 {
		return nil
	}
	for _, r := range to {
		if r.From != nick {
			from = append(from, r)
		}
	}
	slices.SortFunc(from, byRemindAt)
	return from
}

func (rc *Collection) TellsFor(nick string) []*Reminder {
	tells, err := rc.collection.All(tellTo(strings.ToLower(nick)))
	if err != nil {
		logging.Error("Loading tells for %s returned error: %v", nick, err)
		return nil
	}
	return tells
}

func (rc *Collection) Put(value *Reminder) error {
	return rc.collection.Put(value)
}

func (rc *Collection) Del(value *Reminder) error {
	return rc.collection.Del(value)
}
