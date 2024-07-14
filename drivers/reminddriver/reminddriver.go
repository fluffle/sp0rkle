package reminddriver

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/pushes"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"github.com/fluffle/sp0rkle/util/push"
	"gopkg.in/mgo.v2/bson"
)

// We use the reminders collection
var rc *reminders.Collection
var pc *pushes.Collection

// We need to be able to kill reminder goroutines
var running = map[bson.ObjectId]context.CancelFunc{}

// It's also nice for people to be able to snooze them
var finished = map[string]*reminders.Reminder{}

// And it's useful to index them for deletion per-person
var listed = map[string][]bson.ObjectId{}

func Init() {
	rc = reminders.Init()
	if push.Enabled() {
		pc = pushes.Init()
	}

	// Set up the handlers and commands.
	bot.Handle(load, client.CONNECTED)
	bot.Handle(unload, client.DISCONNECTED)
	bot.Handle(tellCheck,
		client.PRIVMSG, client.ACTION, client.JOIN, client.NICK)

	bot.Command(tell, "tell", "tell <nick> <msg>  -- "+
		"Stores a message for the (absent) nick.")
	bot.Command(tell, "ask", "ask <nick> <msg>  -- "+
		"Stores a message for the (absent) nick.")
	bot.Command(list, "remind list",
		"remind list  -- Lists reminders set by or for your nick.")
	bot.Command(del, "remind del",
		"remind del <N>  -- Deletes (previously listed) reminder N.")
	bot.Command(set, "remind", "remind <nick> <msg> "+
		"in|at|on <time>  -- Reminds nick about msg at time.")
	bot.Command(snooze, "snooze", "snooze [duration]  -- "+
		"Resets the previously-triggered reminder.")
	bot.Command(zone, "my timezone is", "my timezone is <zone>  -- "+
		"Sets a local timezone for your nick.")
	bot.Command(unzone, "forget my timezone", "forget my timezone  -- "+
		"Unsets a local timezone for your nick.")
}

func Remind(r *reminders.Reminder, ctx *bot.Context) {
	delta := r.RemindAt.Sub(time.Now())
	if delta < 0 {
		return
	}
	c, cancel := context.WithDeadline(bot.Ctx(), r.RemindAt)
	running[r.Id()] = cancel
	go func() {
		<-c.Done()
		if errors.Is(c.Err(), context.DeadlineExceeded) {
			ctx.Privmsg(string(r.Chan), r.Reply())
			// TODO(fluffle): Tie this into state tracking properly.
			ctx.Privmsg(string(r.Target), r.Reply())
			// This is used in snooze to reinstate reminders.
			finished[strings.ToLower(string(r.Target))] = r
			if pc != nil {
				if s := pc.GetByNick(string(r.Target), true); s.CanPush() {
					push.Push(s, "Reminder from sp0rkle!", r.Reply())
				}
			}
			Forget(r.Id(), false)
		}
	}()
}

func Forget(id bson.ObjectId, stop bool) {
	cancel, ok := running[id]
	if ok {
		// If it's *not* in running, it's probably a Tell.
		delete(running, id)
		if stop {
			cancel()
		}
	}
	r := rc.GetById(id)
	if r == nil {
		return
	}
	if err := rc.Del(r); err != nil {
		logging.Error("Failure removing reminder %s: %v", id, err)
	}
}
