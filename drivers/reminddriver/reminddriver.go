package reminddriver

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db/conf"
	"github.com/fluffle/sp0rkle/collections/pushes"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"github.com/fluffle/sp0rkle/util/bson"
	"github.com/fluffle/sp0rkle/util/datetime"
	"github.com/fluffle/sp0rkle/util/push"
)

type Driver struct {
	ctx    context.Context
	rc     *reminders.Collection
	pc     *pushes.Collection
	tzNs   conf.Namespace
}


// We need to be able to kill reminder goroutines
var running = map[bson.ObjectId]context.CancelFunc{}

// It's also nice for people to be able to snooze them
var finished = map[string]*reminders.Reminder{}

// And it's useful to index them for deletion per-person
var listed = map[string][]bson.ObjectId{}

func (d *Driver) Zone(nick string, tz ...string) string {
	nick = strings.ToLower(nick)
	if len(tz) > 0 && tz[0] == "" {
		d.tzNs.Delete(nick)
		return ""
	}
	return d.tzNs.String(nick, tz...)
}

func New(b *bot.Bot, rc *reminders.Collection, pc *pushes.Collection, config *conf.Registry) *Driver {
	d := &Driver{ctx: b.Ctx(), rc: rc, pc: pc, tzNs: config.Ns(datetime.ZoneNs)}
	
	if push.Enabled() {
		d.pc = pc
	}

	// Set up the handlers and commands.
	b.Handle(d.load, client.CONNECTED)
	b.Handle(d.unload, client.DISCONNECTED)
	b.Handle(d.tellCheck, client.PRIVMSG, client.ACTION, client.JOIN, client.NICK)

	b.Command(d.tell, "tell", "tell <nick> <msg>  -- "+
		"Stores a message for the (absent) nick.")
	b.Command(d.tell, "ask", "ask <nick> <msg>  -- "+
		"Stores a message for the (absent) nick.")
	b.Command(d.list, "remind list",
		"remind list  -- Lists reminders set by or for your nick.")
	b.Command(d.del, "remind del",
		"remind del <N>  -- Deletes (previously listed) reminder N.")
	b.Command(d.set, "remind", "remind <nick> <msg> "+
		"in|at|on <time>  -- Reminds nick about msg at time.")
	b.Command(d.snooze, "snooze", "snooze [duration]  -- "+
		"Resets the previously-triggered reminder.")
	b.Command(d.zone, "my timezone is", "my timezone is <zone>  -- "+
		"Sets a local timezone for your nick.")
	b.Command(d.unzone, "forget my timezone", "forget my timezone  -- "+
		"Unsets a local timezone for your nick.")

	return d
}

func (d *Driver) Remind(r *reminders.Reminder, ctx *bot.Context) {
	delta := r.RemindAt.Sub(time.Now())
	if delta < 0 {
		return
	}
	c, cancel := context.WithDeadline(d.ctx, r.RemindAt)
	running[r.Id()] = cancel
	go func() {
		<-c.Done()
		if errors.Is(c.Err(), context.DeadlineExceeded) {
			ctx.Privmsg(string(r.Chan), r.Reply())
			// TODO(fluffle): Tie this into state tracking properly.
			ctx.Privmsg(string(r.Target), r.Reply())
			// This is used in snooze to reinstate reminders.
			finished[strings.ToLower(string(r.Target))] = r
			if d.pc != nil {
				if s := d.pc.GetByNick(string(r.Target), true); s.CanPush() {
					push.Push(s, "Reminder from sp0rkle!", r.Reply())
				}
			}
			d.Forget(r.Id(), false)
		}
	}()
}

func (d *Driver) Forget(id bson.ObjectId, stop bool) {
	cancel, ok := running[id]
	if ok {
		// If it's *not* in running, it's probably a Tell.
		delete(running, id)
		if stop {
			cancel()
		}
	}
	r := d.rc.GetById(id)
	if r == nil {
		return
	}
	if err := d.rc.Del(r); err != nil {
		logging.Error("Failure removing reminder %s: %v", id, err)
	}
}
