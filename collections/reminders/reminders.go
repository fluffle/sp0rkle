package reminders

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/datetime"
	"gopkg.in/mgo.v2/bson"
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
	// From and To are not unique so we use a nanosecond timestamp from
	// the reminder to differentiate and sort. Tells don't set RemindAt,
	// so we use the create timestamp instead.
	ts := uint64(r.RemindAt.UnixNano())
	if r.Tell {
		ts = uint64(r.Created.UnixNano())
	}
	return []db.Key{
		db.K{db.T{"tell", r.Tell}, db.S{"from", r.From}, db.I{"ts", ts}},
		db.K{db.T{"tell", r.Tell}, db.S{"to", r.To}, db.I{"ts", ts}},
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

type Reminders []*Reminder

func (rs Reminders) Strings() []string {
	s := make([]string, len(rs))
	for i, r := range rs {
		s[i] = fmt.Sprintf("%#v", r)
	}
	return s
}

func (rs Reminders) sortByRemindAt() {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].RemindAt.Before(rs[j].RemindAt)
	})
}

type migrator struct {
	mongo, bolt db.Collection
}

func (m *migrator) MigrateTo(newState db.MigrationState) error {
	if newState != db.MONGO_PRIMARY {
		return nil
	}
	var all Reminders
	if err := m.mongo.All(db.K{}, &all); err != nil {
		return err
	}
	if err := m.bolt.BatchPut(all); err != nil {
		logging.Error("Migrating reminders: %v", err)
		return err
	}
	logging.Info("Migrated %d reminders.", len(all))
	return nil
}

func (m *migrator) Diff() ([]string, []string, error) {
	var mAll, bAll Reminders
	if err := m.mongo.All(db.K{}, &mAll); err != nil {
		return nil, nil, err
	}
	if err := m.bolt.All(db.K{}, &bAll); err != nil {
		return nil, nil, err
	}
	return mAll.Strings(), bAll.Strings(), nil
}

type Collection struct {
	db.Both
}

func Init() *Collection {
	rc := &Collection{db.Both{}}
	rc.Both.MongoC.Init(db.Mongo, COLLECTION, mongoIndexes)
	rc.Both.BoltC.Init(db.Bolt.Indexed(), COLLECTION, nil)
	m := &migrator{
		mongo: rc.Both.MongoC,
		bolt:  rc.Both.BoltC,
	}
	rc.Both.Checker.Init(m, COLLECTION)
	return rc
}

func mongoIndexes(c db.Collection) {
	for _, k := range []string{"remindat", "from", "to", "tell"} {
		if err := c.Mongo().EnsureIndexKey(k); err != nil {
			logging.Error("Couldn't create %s index on sp0rkle.reminders: %v", k, err)
		}
	}
}

func (rc *Collection) GetById(id bson.ObjectId) *Reminder {
	r := &Reminder{Id_: id}
	if err := rc.Get(r.byId(), r); err != nil {
		logging.Error("Reminder GetById(%s) failed: %v", id, err)
		return nil
	}
	return r
}

func (rc *Collection) LoadAndPrune() Reminders {
	// Can't delete from Bolt without loading everything and sorting.
	// This will work fine in Mongo too, so let's just do that.
	var all Reminders
	if err := rc.All(db.K{db.T{"tell", false}}, &all); err != nil {
		logging.Error("Loading all reminders: %v", err)
		return nil
	}
	all.sortByRemindAt()
	now := time.Now()
	var last int
	for i, r := range all {
		if r.RemindAt.After(now) {
			last = i
			break
		}
	}

	if last > 0 {
		for _, r := range all[:last] {
			if err := rc.Del(r); err != nil {
				logging.Error("Deleting expired reminder %v (expiry %s): %v", r.Id_, r.At(), err)
			}
		}
		all = all[last:]
		logging.Info("Removed %d old reminders", last)
	}
	return all
}

func (rc *Collection) RemindersFor(nick string) Reminders {
	nick = strings.ToLower(nick)
	var from, to Reminders
	if err := rc.All(remindFrom(nick), &from); err != nil {
		logging.Error("Loading reminders from %s returned error: %v", nick, err)
	}
	if err := rc.All(remindTo(nick), &to); err != nil {
		logging.Error("Loading reminders to %s returned error: %v", nick, err)
	}
	if len(from) == 0 && len(to) == 0 {
		return nil
	}
	// A reminder that is both from nick and to nick will appear in
	// both lists, so we can't just append one to the other...
	for _, r := range to {
		if r.From != nick {
			from = append(from, r)
		}
	}
	from.sortByRemindAt()
	for _, r := range from {
		logging.Debug("REMINDER: %s %v", r.At(), r.Id_)
	}
	return from
}

func (rc *Collection) TellsFor(nick string) Reminders {
	var tells Reminders
	if err := rc.All(tellTo(strings.ToLower(nick)), &tells); err != nil {
		logging.Error("Loading tells for %s returned error: %v", nick, err)
		return nil
	}
	for _, r := range tells {
		logging.Debug("TELL: %s %v", r.At(), r.Id_)
	}
	return tells
}
