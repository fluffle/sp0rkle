package reminddriver

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/lib/reminders"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
	"labix.org/v2/mgo/bson"
	"time"
)

type remindFn func(*remindDriver, *base.Line)

// A remindCommand fulfils base.Handler and base.Command
type remindCommand struct {
	rd *remindDriver
	fn remindFn
	help string
}

func (rc *remindCommand) Execute(l *base.Line) {
	rc.fn(rc.rd, l)
}

func (rc *remindCommand) Help() string {
	return rc.help
}

// These two shim the remind driver into the command / handler
func (rd *remindDriver) Cmd(fn remindFn, prefix, help string) {
	bot.Cmd(&remindCommand{rd,fn,help}, prefix)
}

func (rd *remindDriver) Handle(fn remindFn, event ...string) {
	bot.Handle(&remindCommand{rd, fn, ""}, event...)
}

type remindDriver struct {
	*reminders.ReminderCollection

	// We need to be able to kill reminder goroutines
	kill map[bson.ObjectId]chan bool

	// And it's useful to index them for deletion per-person
	list map[string][]bson.ObjectId
}

func Init(db *db.Database) *remindDriver {
	rd := &remindDriver{
		ReminderCollection: reminders.Collection(db),
		kill: make(map[bson.ObjectId]chan bool),
		list: make(map[string][]bson.ObjectId),
	}

	// Set up the handlers and commands.
	rd.Handle((*remindDriver).Load, "connected")
	rd.Handle((*remindDriver).TellCheck,
		"privmsg", "action", "join", "nick")

	rd.Cmd((*remindDriver).Tell, "tell", "tell <nick> <msg>  -- " +
		"Stores a message for the (absent) nick.")
	rd.Cmd((*remindDriver).List, "remind list",
		"remind list  -- Lists reminders set by or for your nick.")
	rd.Cmd((*remindDriver).Del, "remind del",
		"remind del <N>  -- Deletes (previously listed) reminder N.")
	rd.Cmd((*remindDriver).Set, "remind", "remind <nick> <msg> " +
		"in|at|on <time>  -- Reminds nick about msg at time.")
	return rd
}

func (rd *remindDriver) Remind(r *reminders.Reminder) {
	delta := r.RemindAt.Sub(time.Now())
	if delta < 0 {
		return
	}
	c := make(chan bool)
	rd.kill[r.Id] = c
	go func() {
		select {
		case <-time.After(delta):
			bot.Privmsg(r.Chan, r.Reply())
			if r.Target.Host != "" {
				// At the time of the reminder being created, target existed
				// TODO(fluffle): Tie this into state tracking properly.
				bot.Privmsg(r.Target.Nick, r.Reply())
			}
			rd.Forget(r.Id, false)
		case <-c:
			return
		}
	}()
}

func (rd *remindDriver) Forget(id bson.ObjectId, kill bool) {
	c, ok := rd.kill[id]
	if !ok { return }
	delete(rd.kill, id)
	if kill {
		c <- true
	}
	if err := rd.RemoveId(id); err != nil {
		logging.Error("Failure removing reminder %s: %v", id, err)
	}
}
