package reminddriver

import (
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/sp0rkle/lib/datetime"
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/lib/reminders"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
	"labix.org/v2/mgo/bson"
	"strconv"
	"strings"
	"time"
)

type remindFn func(*remindDriver, *bot.Sp0rkle, *base.Line)

type remindCommand struct {
	rd *remindDriver
	fn remindFn
	help string
}

func (rc *remindCommand) Execute(b *bot.Sp0rkle, l *base.Line) {
	rc.fn(rc.rd, b, l)
}

func (rc *remindCommand) Help() string {
	return rc.help
}

func (rd *remindDriver) Cmd(fn remindFn, prefix, help string) {
	bot.Cmd(prefix, &remindCommand{rd,fn,help})
}

func (rd *remindDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler(bot.NewHandler(rd_load), "bot_connected")
	r.AddHandler(bot.NewHandler(rd_privmsg), "bot_privmsg")
	r.AddHandler(bot.NewHandler(rd_tell_check),
		"bot_privmsg", "bot_action", "bot_join", "bot_nick")
	rd.Cmd((*remindDriver).Tell, "tell", "tell <nick> <msg>  -- " +
		"Stores a message for the (absent) nick.")
}

func rd_load(bot *bot.Sp0rkle, line *base.Line) {
	rd := bot.GetDriver(driverName).(*remindDriver)
	// We're connected to IRC, so load saved reminders
	r := rd.LoadAndPrune()
	for i := range r {
		if r[i] == nil {
			rd.l.Warn("Nil reminder %d from LoadAndPrune", i)
			continue
		}
		go rd.Remind(r[i])(bot)
	}
}

func rd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	rd := bot.GetDriver(driverName).(*remindDriver)

	if !line.Addressed {
		return
	}

	switch {
//	case strings.HasPrefix(line.Args[1], "tell "):
//		rd_tell(bot, rd, line)
	case strings.HasPrefix(line.Args[1], "remind list"):
		rd_list(bot, rd, line)
	case strings.HasPrefix(line.Args[1], "remind del"):
		rd_del(bot, rd, line)
	case strings.HasPrefix(line.Args[1], "remind "):
		rd_set(bot, rd, line)
	}
}

func rd_del(bot *bot.Sp0rkle, rd *remindDriver, line *base.Line) {
	list, ok := rd.list[line.Nick]
	if !ok {
		bot.ReplyN(line, "Please use 'remind list' first, " +
			"to be sure of what you're deleting.")
		return
	}
	s := strings.Fields(line.Args[1])
	idx, err := strconv.Atoi(s[len(s)-1])
	if err != nil || idx > len(list) || idx <= 0 {
		bot.ReplyN(line, "Invalid reminder index '%s'", s[len(s)-1])
		return
	}
	idx--
	rd.Delete(list[idx], true)
	delete(rd.list, line.Nick)
	bot.ReplyN(line, "I'll forget that one, then...")
}

func rd_list(bot *bot.Sp0rkle, rd *remindDriver, line *base.Line) {
	r := rd.RemindersFor(line.Nick)
	c := len(r)
	if c == 0 {
		bot.ReplyN(line, "You have no reminders set.")
		return
	}
	if c > 5 && line.Args[0][0] == '#' {
		bot.ReplyN(line, "You've got lots of reminders, ask me privately.")
		return
	}
	// Save an ordered list of ObjectIds for easy reminder deletion
	bot.ReplyN(line, "You have %d reminders set:", c)
	list := make([]bson.ObjectId, c)
	for i := range r {
		bot.Reply(line, "%d: %s", i+1, r[i].List(line.Nick))
		list[i] = r[i].Id
	}
	rd.list[line.Nick] = list
}

func rd_set(bot *bot.Sp0rkle, rd *remindDriver, line *base.Line) {
	// s == remind <target> <reminder> in|at|on <time>
	s := strings.Fields(line.Args[1])
	if len(s) < 5 {
		bot.ReplyN(line, "Invalid remind syntax. Sucka.")
		return
	}
	i := len(s)-1
	for i > 0 {
		lc := strings.ToLower(s[i])
		if lc == "in" || lc == "at" || lc == "on" {
			break
		}
		i--
	}
	if i < 2 {
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
			at = at.Add(24 * time.Hour)
		} else {
			at = at.Add(12 * time.Hour)
		}
	}
	if at.Before(now) {
		bot.ReplyN(line, "Time '%s' is in the past.", timestr)
		return
	}
	n, c := line.Storable()
	// TODO(fluffle): Use state tracking!
	t := db.StorableNick{Nick: s[1]}
	if strings.ToLower(t.Nick) == strings.ToLower(line.Nick) ||
		strings.ToLower(t.Nick) == "me" {
		t = n
	}
	r := reminders.NewReminder(reminder, at, t, n, c)
	if err := rd.Insert(r); err != nil {
		bot.ReplyN(line, "Error saving reminder: %v", err)
		return
	}
	// Any previously-generated list of reminders is now obsolete.
	delete(rd.list, line.Nick)
	bot.ReplyN(line, r.Acknowledge())
	go rd.Remind(r)(bot)
}

func (rd *remindDriver) Tell(bot *bot.Sp0rkle, line *base.Line) {
	// s == tell <target> <stuff>
	s := strings.Fields(line.Args[1])
	if len(s) < 3 {
		bot.ReplyN(line, "Tell who what?")
		return
	}
	tell := strings.Join(s[2:], " ")
	n, c := line.Storable()
	t := db.StorableNick{Nick: s[1]}
	if strings.ToLower(t.Nick) == strings.ToLower(line.Nick) ||
		strings.ToLower(t.Nick) == "me" {
		bot.ReplyN(line, "You're a dick. Oh, wait, that wasn't *quite* it...")
		return
	}
	r := reminders.NewTell(tell, t, n, c)
	if err := rd.Insert(r); err != nil {
		bot.ReplyN(line, "Error saving tell: %v", err)
		return
	}
	// Any previously-generated list of reminders is now obsolete.
	delete(rd.list, line.Nick)
	bot.ReplyN(line, r.Acknowledge())
}

func rd_tell_check(bot *bot.Sp0rkle, line *base.Line) {
	rd := bot.GetDriver(driverName).(*remindDriver)
	nick := line.Nick
	if line.Cmd == "NICK" {
		// We want the destination nick, not the source.
		nick = line.Args[0]
	}
	r := rd.TellsFor(nick)
	for i := range r {
		if line.Cmd == "NICK" {
			bot.Conn.Privmsg(r[i].Chan, nick + ": " + r[i].Reply())
			bot.Reply(line, r[i].Reply())
		} else {
			bot.Conn.Privmsg(line.Nick, r[i].Reply())
			bot.ReplyN(line, r[i].Reply())
		}
		rd.RemoveId(r[i].Id)
	}
	if len(r) > 0 {
		delete(rd.list, line.Nick)
	}
}
