package reminddriver

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/conf"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"github.com/fluffle/sp0rkle/util/datetime"
	"github.com/fluffle/sp0rkle/util/push"
	"gopkg.in/mgo.v2/bson"
)

// remind del
func del(ctx *bot.Context) {
	list, ok := listed[ctx.Nick]
	if !ok {
		ctx.ReplyN("Please use 'remind list' first, " +
			"to be sure of what you're deleting.")

		return
	}
	idx, err := strconv.Atoi(ctx.Text())
	if err != nil || idx > len(list) || idx <= 0 {
		ctx.ReplyN("Invalid reminder index '%s'", ctx.Text())
		return
	}
	idx--
	Forget(list[idx], true)
	delete(listed, ctx.Nick)
	ctx.ReplyN("I'll forget that one, then...")
}

// remind list
func list(ctx *bot.Context) {
	r := rc.RemindersFor(ctx.Nick)
	c := len(r)
	if c == 0 {
		ctx.ReplyN("You have no reminders set.")
		return
	}
	if c > 5 && ctx.Public() {
		ctx.ReplyN("You've got lots of reminders, ask me privately.")
		return
	}
	// Save an ordered list of ObjectIds for easy reminder deletion
	ctx.ReplyN("You have %d reminders set:", c)
	list := make([]bson.ObjectId, c)
	for i := range r {
		ctx.Reply("%d: %s", i+1, r[i].List(ctx.Nick))
		list[i] = r[i].Id
	}
	listed[ctx.Nick] = list
}

// remind
func set(ctx *bot.Context) {
	// s == <target> <reminder> in|at|on <time>
	s := strings.Fields(ctx.Text())
	if len(s) < 4 {
		ctx.ReplyN("You asked me to remind %s.", ctx.Text())
		return
	}
	// Look up a per-user timezone if one is set.
	z := datetime.ZoneOrLocal(conf.Zone(ctx.Nick))
	// Parse the reminder time from the input.
	at, err, reminder, timestr := time.Now(), error(nil), "", ""
	for i := 1; i+1 < len(s); i++ {
		lc := strings.ToLower(s[i])
		if lc == "in" || lc == "at" || lc == "on" {
			timestr = strings.Join(s[i+1:], " ")
		} else if i+2 == len(s) {
			// Hack to test the last word for e.g. "tomorrow"
			i++
			timestr = strings.ToLower(s[i])
		} else {
			continue
		}
		at, err = datetime.ParseZ(timestr, z)
		if err == nil {
			reminder = strings.Join(s[1:i], " ")
			break
		}
	}
	if reminder == "" {
		ctx.ReplyN("You asked me to remind %s.", ctx.Text())
		return
	}
	if err != nil {
		ctx.ReplyN("Couldn't parse time string %q: %v", timestr, err)
		return
	}
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, datetime.TZ())
	if at.Before(now) && at.After(start) {
		at = at.Add(24 * time.Hour)
	}
	if at.Before(now) {
		ctx.ReplyN("Time %q is in the past.", timestr)
		return
	}
	n, c := ctx.Storable()
	// TODO(fluffle): Use state tracking! And do this better.
	t := bot.Nick(s[0])
	if t.Lower() == strings.ToLower(ctx.Nick) ||
		t.Lower() == "me" {
		t = n
	}
	r := reminders.NewReminder(reminder, at, t, n, c)
	if err := rc.Insert(r); err != nil {
		ctx.ReplyN("Error saving reminder: %v", err)
		return
	}
	// Any previously-generated list of reminders is now obsolete.
	delete(listed, ctx.Nick)
	ctx.ReplyN("%s", r.Acknowledge())
	Remind(r, ctx)
}

// snooze
func snooze(ctx *bot.Context) {
	r, ok := finished[strings.ToLower(ctx.Nick)]
	if !ok {
		ctx.ReplyN("No record of an expired reminder for you, sorry!")
		return
	}
	// Look up a per-user timezone if one is set.
	z := datetime.ZoneOrLocal(conf.Zone(ctx.Nick))
	now := time.Now().In(z)
	at := now.Add(30 * time.Minute)
	if ctx.Text() != "" {
		var err error
		at, err = datetime.ParseZ(ctx.Text(), z)
		if err != nil {
			ctx.ReplyN("Couldn't parse time string %q: %v.", ctx.Text(), err)
			return
		}
		if at.Before(now) {
			ctx.ReplyN("You can't snooze reminder into the past, fool.")
			return
		}
	}
	r.Created = now
	r.RemindAt = at
	if _, err := rc.UpsertId(r.Id, r); err != nil {
		ctx.ReplyN("Error saving reminder: %v", err)
		return
	}
	delete(listed, ctx.Nick)
	ctx.ReplyN("%s", r.Acknowledge())
	Remind(r, ctx)
}

// tell
func tell(ctx *bot.Context) {
	// s == <target> <stuff>
	txt := ctx.Text()
	idx := strings.Index(txt, " ")
	if idx == -1 {
		ctx.ReplyN("Tell who what?")
		return
	}
	tell := txt[idx+1:]
	n, c := ctx.Storable()
	t := bot.Nick(txt[:idx])
	if t.Lower() == strings.ToLower(ctx.Nick) ||
		t.Lower() == "me" {
		ctx.ReplyN("You're a dick. Oh, wait, that wasn't *quite* it...")
		return
	}
	r := reminders.NewTell(tell, t, n, c)
	if err := rc.Insert(r); err != nil {
		ctx.ReplyN("Error saving tell: %v", err)
		return
	}
	if s := pc.GetByNick(txt[:idx], true); s.CanPush() {
		push.Push(s, fmt.Sprintf("%s in %s asked me to tell you:",
			ctx.Nick, ctx.Target()), tell)
	}
	// Any previously-generated list of reminders is now obsolete.
	delete(listed, ctx.Nick)
	ctx.ReplyN("%s", r.Acknowledge())
}

// zone
func zone(ctx *bot.Context) {
	fields := strings.Fields(ctx.Text())
	if len(fields) == 0 {
		ctx.ReplyN("Your timezone is ... fat? Like your mum?")
		return
	}
	if z := datetime.Zone(fields[0]); z != nil {
		conf.Zone(ctx.Nick, fields[0])
		ctx.ReplyN("Reminders will now be in %q.", z)
	} else {
		ctx.ReplyN("Don't recognise %q as a valid timezone, sorry.", fields[0])
	}
}

// unzone
func unzone(ctx *bot.Context) {
	conf.Zone(ctx.Nick, "")
	ctx.ReplyN("I've forgotten where you live... honest!")
}
