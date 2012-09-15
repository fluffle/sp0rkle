package reminddriver

import (
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/sp0rkle/lib/datetime"
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/lib/reminders"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
//	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

func (rd *remindDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(rd_load), "bot_connected")
	r.AddHandler(bot.NewHandler(rd_privmsg), "bot_privmsg")
}

func rd_load(bot *bot.Sp0rkle, line *base.Line) {
	rd := bot.GetDriver(driverName).(*remindDriver)
	// We're connected to IRC, so load saved reminders
	r := rd.LoadAndPrune()
	for i := range r {
		go rd.Remind(r[i])(bot)
	}
}

func rd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	rd := bot.GetDriver(driverName).(*remindDriver)

	if !line.Addressed {
		return
	}

	switch {
	case strings.HasPrefix(line.Args[1], "remind "):
		rd_setremind(bot, rd, line)
	}
}

func rd_setremind(bot *bot.Sp0rkle, rd *remindDriver, line *base.Line) {
	// s == remind <target> <reminder> in|at|on <time>
	s := strings.Fields(line.Args[1])
	target := s[1]
	i := len(s)-1
	for i > 0 {
		lc := strings.ToLower(s[i])
		if lc == "in" || lc == "at" || lc == "on" {
			break
		}
		i--
	}
	if i == 0 {
		bot.ReplyN(line, "Invalid remind syntax. Sucka.")
		return
	}
	reminder := strings.Join(s[2:i], " ")
	timestr := strings.ToLower(strings.Join(s[i+1:], " "))
	// TODO(fluffle): surface better errors from datetime.Parse
	at, ok := datetime.Parse(timestr)
	if !ok {
		bot.ReplyN(line, "Couldn't parse time string '%s'", timestr)
		return
	}
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if at.Before(now) && at.After(start) {
		// Perform some basic hacky corrections before giving up
		if strings.Contains(timestr, "am") || strings.Contains(timestr, "pm") {
			at.Add(24 * time.Hour)
		} else {
			at.Add(12 * time.Hour)
		}
	}
	if at.Before(now) {
		bot.ReplyN(line, "Time '%s' is in the past.", timestr)
		return
	}
	n, c := line.Storable()
	// TODO(fluffle): Use state tracking!
	t := db.StorableNick{Nick: target}
	if strings.ToLower(target) == strings.ToLower(line.Nick) {
		t = n
	}
	r := reminders.NewReminder(reminder, at, t, n, c)
	if err := rd.Insert(r); err != nil {
		bot.ReplyN(line, "Error saving reminder: %v", err)
		return
	}
	bot.ReplyN(line, r.Acknowledge())
	go rd.Remind(r)(bot)
}
