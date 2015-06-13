package reminddriver

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/pushes"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"github.com/fluffle/sp0rkle/util/push"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

// We use the reminders collection
var rc *reminders.Collection
var pc *pushes.Collection

// We need to be able to kill reminder goroutines
var running = map[bson.ObjectId]chan struct{}{}

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
}

func Remind(r *reminders.Reminder, ctx *bot.Context) {
	delta := r.RemindAt.Sub(time.Now())
	if delta < 0 {
		return
	}
	c := make(chan struct{})
	running[r.Id] = c
	go func() {
		select {
		case <-time.After(delta):
			ctx.Privmsg(string(r.Chan), r.Reply())
			// TODO(fluffle): Tie this into state tracking properly.
			ctx.Privmsg(string(r.Target), r.Reply())
			// This is used in snooze to reinstate reminders.
			finished[strings.ToLower(string(r.Target))] = r
			if s := pc.GetByNick(string(r.Target)); s.CanPush() {
				push.Push(s, "Reminder from sp0rkle!", r.Reply())
			}
			Forget(r.Id, false)
		case <-c:
			return
		}
	}()
}

func Forget(id bson.ObjectId, stop bool) {
	c, ok := running[id]
	if ok {
		// If it's *not* in running, it's probably a Tell.
		delete(running, id)
		if stop {
			c <- struct{}{}
		}
	}
	if err := rc.RemoveId(id); err != nil {
		logging.Error("Failure removing reminder %s: %v", id, err)
	}
}
