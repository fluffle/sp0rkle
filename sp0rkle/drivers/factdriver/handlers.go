package factdriver

import (
	"fmt"
	"github.com/fluffle/goirc/event"
	"launchpad.net/gobson/bson"
	"lib/db"
	"lib/factoids"
	"os"
	"rand"
	"sp0rkle/bot"
	"sp0rkle/base"
	"strings"
	"strconv"
	"time"
)

type FactoidHandler func(*bot.Sp0rkle, *factoidDriver, *base.Line)

// Unboxer for FactoidDriver handlers
func FDHandler(f FactoidHandler) event.Handler {
	return event.NewHandler(func(ev ...interface{}) {
		f(ev[0].(*bot.Sp0rkle), ev[1].(*factoidDriver), ev[2].(*base.Line))
	})
}

func (fd *factoidDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler("bot_privmsg", bot.NewHandler(fd_privmsg))
	r.AddHandler("bot_action", bot.NewHandler(fd_action))
	r.AddHandler("fd_lookup", FDHandler(fd_lookup))
	r.AddHandler("fd_add", FDHandler(fd_add))
	r.AddHandler("fd_delete", FDHandler(fd_delete))
	r.AddHandler("fd_replace", FDHandler(fd_replace))
	r.AddHandler("fd_chance", FDHandler(fd_chance))
	r.AddHandler("fd_literal", FDHandler(fd_literal))
	r.AddHandler("fd_search", FDHandler(fd_search))
	r.AddHandler("fd_info", FDHandler(fd_info))
}

func fd_privmsg(bot *bot.Sp0rkle, line *base.Line) {
	fd := bot.GetDriver(driverName).(*factoidDriver)

	// If we're not being addressed directly, short-circuit to lookup.
	if !line.Addressed {
		bot.Dispatch("fd_lookup", fd, line)
		return
	}

	l := strings.ToLower(line.Args[1])
	// Test for various possible courses of action.
	switch {
	// Factoid add: 'key := value' or 'key :is value'
	case strings.Index(l, ":=") != -1:
		fallthrough
	case strings.Index(l, ":is") != -1:
		bot.Dispatch("fd_add", fd, line)

	// Factoid delete: 'forget|delete that' => deletes fd.lastseen[chan]
	case strings.HasPrefix(l, "forget that"):
		fallthrough
	case strings.HasPrefix(l, "delete that"):
		bot.Dispatch("fd_delete", fd, line)

	// Factoid replace: 'replace that with' => updates fd.lastseen[chan]
	case strings.HasPrefix(l, "replace that with "):
		// chop off the "known" bit to leave just the replacement
		nl := line.Copy()
		nl.Args[1] = nl.Args[1][18:]
		bot.Dispatch("fd_replace", fd, line)

	// Factoid chance: 'chance of that is' => sets chance of fd.lastseen[chan]
	case strings.HasPrefix(l, "chance of that is "):
		// chop off the "known" bit to leave just the replacement
		nl := line.Copy()
		nl.Args[1] = nl.Args[1][18:]
		bot.Dispatch("fd_chance", fd, nl)

	// Factoid literal: 'literal key' => info about factoid
	case strings.HasPrefix(l, "literal "):
		// chop off the "known" bit to leave just the key
		nl := line.Copy()
		nl.Args[1] = nl.Args[1][8:]
		bot.Dispatch("fd_literal", fd, nl)

	// Factoid search: 'fact search regexp' => list of possible key matches
	case strings.HasPrefix(l, "fact search "):
		nl := line.Copy()
		nl.Args[1] = nl.Args[1][12:]
		bot.Dispatch("fd_search", fd, nl)

	// Factoid info: 'fact info key' => some information about key
	case strings.HasPrefix(l, "fact info"):
		nl := line.Copy()
		nl.Args[1] = nl.Args[1][9:]
		bot.Dispatch("fd_info", fd, nl)

	// If we get to here, none of the other FD command possibilities
	// have matched, so try a lookup...
	default:
		bot.Dispatch("fd_lookup", fd, line)
	}
}

func fd_action(bot *bot.Sp0rkle, line *base.Line) {
	fd := bot.GetDriver(driverName).(*factoidDriver)
	// Actions just trigger a lookup.
	bot.Dispatch("fd_lookup", fd, line)
}

func fd_add(bot *bot.Sp0rkle, fd *factoidDriver, line *base.Line) {
	var key, val string
	if strings.Index(line.Args[1], ":=") != -1 {
		kv := strings.SplitN(line.Args[1], ":=", 2)
		key = ToKey(kv[0], false)
		val = strings.TrimSpace(kv[1])
	} else {
		// we use :is to add val = "key is val"
		kv := strings.SplitN(line.Args[1], ":is", 2)
		key = ToKey(kv[0], false)
		val = strings.Join([]string{strings.TrimSpace(kv[0]),
			"is", strings.TrimSpace(kv[1])}, " ")
	}
	n := db.StorableNick{line.Nick, line.Ident, line.Host}
	c := db.StorableChan{line.Args[0]}
	fact := factoids.NewFactoid(key, val, n, c)
	if err := fd.Insert(fact); err == nil {
		count := fd.GetCount(key)
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: Woo, I now know %d things about '%s'.",
			line.Nick, count, key))
	} else {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf("Oh no! %s.", err))
	}
}

func fd_chance(bot *bot.Sp0rkle, fd *factoidDriver, line *base.Line) {
	str := strings.TrimSpace(line.Args[1])
	var chance float32

	if strings.HasSuffix(str, "%") {
		// Handle 'chance of that is \d+%'
		if i, err := strconv.Atoi(str[:len(str)-1]); err != nil {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: '%s' didn't look like a % chance to me.",
				line.Nick, str))
			return
		} else {
			chance = float32(i) / 100
		}
	} else {
		// Assume the chance is a floating point number.
		if c, err := strconv.Atof32(str); err != nil {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: '%s' didn't look like a chance to me.",
				line.Nick, str))
			return
		} else {
			chance = c
		}
	}

	// Make sure the chance we've parsed lies in (0.0,1.0]
	if chance > 1.0 || chance <= 0.0 {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: '%s' was outside possible chance ranges.",
			line.Nick, str))
		return
	}

	// Retrieve last seen ObjectId, replace with ""
	ls := fd.Lastseen(line.Args[0], "")
	// ok, we're good to update the chance.
	if fact := fd.GetById(ls); fact != nil {
		// Store the old chance, update with the new
		old := fact.Chance
		fact.Chance = chance
		// Update the Modified field
		n := db.StorableNick{line.Nick, line.Ident, line.Host}
		c := db.StorableChan{line.Args[0]}
		fact.Modify(n, c)
		// And store the new factoid data
		if err := fd.Update(bson.M{"_id": ls}, fact); err == nil {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: '%s' was at %.0f%% chance, now is at %.0f%%.",
				line.Nick, fact.Key, old*100, chance*100))
		} else {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I failed to replace '%s': %s",
				line.Nick, fact.Key, err))
		}
	} else {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: Whatever that was, I've already forgotten it.", line.Nick))
	}
}

func fd_delete(bot *bot.Sp0rkle, fd *factoidDriver, line *base.Line) {
	// Get fresh state on the last seen factoid.
	ls := fd.Lastseen(line.Args[0], "")
	if fact := fd.GetById(ls); fact != nil {
		if err := fd.Remove(bson.M{"_id": ls}); err == nil {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I forgot that '%s' was '%s'.",
				line.Nick, fact.Key, fact.Value))
		} else {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I failed to forget '%s': %s",
				line.Nick, fact.Key, err))
		}
	} else {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: Whatever that was, I've already forgotten it.", line.Nick))
	}
}

func fd_info(bot *bot.Sp0rkle, fd *factoidDriver, line *base.Line) {
	key := ToKey(line.Args[1], false)
	count := fd.GetCount(key);
	if count == 0 {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: I don't know anything about '%s'.",
			line.Nick, key))
		return
	}
	msgs := make([]string, 0, 10)
	if key == "" {
		msgs = append(msgs, fmt.Sprintf("%s: In total, I know %d things.",
			line.Nick, count))
	} else {
		msgs = append(msgs, fmt.Sprintf("%s: I know %d things about '%s'.",
			line.Nick, count, key))
	}
	if created := fd.GetLast("created", key); created != nil {
		c := created.Created
		msgs = append(msgs, "A factoid")
		if key != "" {
			msgs = append(msgs, fmt.Sprintf("for '%s'", key))
		}
		msgs = append(msgs, fmt.Sprintf("was last created on %s by %s,",
			c.Timestamp.Format(time.ANSIC), c.Nick))
	}
	if modified := fd.GetLast("modified", key); modified != nil {
		m := modified.Modified
		msgs = append(msgs, fmt.Sprintf("modified on %s by %s,",
			m.Timestamp.Format(time.ANSIC), m.Nick))
	}
	if accessed := fd.GetLast("accessed", key); accessed != nil {
		a := accessed.Accessed
		msgs = append(msgs, fmt.Sprintf("and accessed on %s by %s.",
			a.Timestamp.Format(time.ANSIC), a.Nick))
	}
	if info := fd.InfoMR(key); info != nil {
		if key == "" {
			msgs = append(msgs, "These factoids have")
		} else {
			msgs = append(msgs, fmt.Sprintf("'%s' has", key))
		}
		msgs = append(msgs, fmt.Sprintf(
			"been modified %d times and accessed %d times.",
			info.Modified, info.Accessed))
	}
	bot.Conn.Privmsg(line.Args[0], strings.Join(msgs, " "))
}

func fd_literal(bot *bot.Sp0rkle, fd *factoidDriver, line *base.Line) {
	key := ToKey(line.Args[1], false)
	if count := fd.GetCount(key); count == 0 {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: I don't know anything about '%s'.",
			line.Nick, key))
		return
	} else if count > 10 && strings.HasPrefix(line.Args[0], "#") {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: I know too much about '%s', ask me privately.",
			line.Nick, key))
		return
	}

	// Temporarily turn off flood protection cos we could be spamming a bit.
	bot.Conn.Flood = true
	defer func() { bot.Conn.Flood = false }()
	// Passing an anonymous function to For makes it a little hard to abstract
	// away in lib/factoids. Fortunately this is something of a one-off.
	var fact *factoids.Factoid
	f := func() os.Error {
		if fact != nil {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: [%3.0f%%] %s", line.Nick, fact.Chance*100, fact.Value))
		}
		return nil
	}
	if err := fd.Find(bson.M{"key": key}).For(&fact, f); err != nil {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: Something went wrong: %s", line.Nick, err))
	}
}

func fd_lookup(bot *bot.Sp0rkle, fd *factoidDriver, line *base.Line) {
	// Only perform extra prefix removal if we weren't addressed directly
	key := ToKey(line.Args[1], !line.Addressed)
	var fact *factoids.Factoid

	if fact = fd.GetPseudoRand(key); fact == nil && line.Cmd == "ACTION" {
		// Support sp0rkle's habit of stripping off it's own nick
		// but only for actions, not privmsgs.
		if strings.HasSuffix(key, bot.Conn.Me.Nick) {
			key = strings.TrimSpace(key[:len(key)-len(bot.Conn.Me.Nick)])
			fact = fd.GetPseudoRand(key)
		}
	}
	if fact == nil {
		return
	}
	// Chance is used to limit the rate of factoid replies for things
	// people say a lot, like smilies, or 'lol', or 'i love the peen'.
	chance := fact.Chance
	if key == "" {
		// This is doing a "random" lookup, triggered by someone typing in
		// something entirely composed of the chars stripped by ToKey().
		// To avoid making this too spammy, forcibly limit the chance to 40%.
		chance = 0.4
	}
	if rand.Float32() < chance {
		// Store this as the last seen factoid
		fd.Lastseen(line.Args[0], fact.Id)
		// Update the Accessed field
		n := db.StorableNick{line.Nick, line.Ident, line.Host}
		c := db.StorableChan{line.Args[0]}
		fact.Access(n, c)
		// And store the new factoid data
		if err := fd.Update(bson.M{"_id": fact.Id}, fact); err != nil {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I failed to update '%s': %s",
				line.Nick, fact.Key, err))
		}

		// Apply the list of factoid plugins to the factoid value.
		val := fd.ApplyPlugins(fact.Value, line)
		switch fact.Type {
		case factoids.F_ACTION:
			bot.Conn.Action(line.Args[0], val)
		default:
			bot.Conn.Privmsg(line.Args[0], val)
		}
	}
}

func fd_replace(bot *bot.Sp0rkle, fd *factoidDriver, line *base.Line) {
	ls := fd.Lastseen(line.Args[0], "")
	if fact := fd.GetById(ls); fact != nil {
		// Store the old factoid value
		old := fact.Value
		// Replace the value with the new one
		fact.Value = strings.TrimSpace(line.Args[1])
		// Update the Modified field
		n := db.StorableNick{line.Nick, line.Ident, line.Host}
		c := db.StorableChan{line.Args[0]}
		fact.Modify(n, c)
		// And store the new factoid data
		if err := fd.Update(bson.M{"_id": ls}, fact); err == nil {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: '%s' was '%s', now is '%s'.",
				line.Nick, fact.Key, old, fact.Value))
		} else {
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I failed to replace '%s': %s",
				line.Nick, fact.Key, err))
		}
	} else {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: Whatever that was, I've already forgotten it.", line.Nick))
	}
}

func fd_search(bot *bot.Sp0rkle, fd *factoidDriver, line *base.Line) {
	if keys := fd.GetKeysMatching(line.Args[1]); keys == nil || len(keys) == 0 {
		bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: I couldn't think of anything matching '%s'.",
			line.Nick, line.Args[0]))
	} else {
		// RESULTS.
		count := len(keys)
		if count > 10 {
			res := strings.Join(keys[:10], "', '")
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I found %d keys matching '%s', here's the first 10: '%s'.",
				line.Nick, count, line.Args[1], res))
		} else {
			res := strings.Join(keys, "', '")
			bot.Conn.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I found %d keys matching '%s', here they are: '%s'.",
				line.Nick, count, line.Args[1], res))
		}
	}
}
