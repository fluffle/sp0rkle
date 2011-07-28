package main

import (
	"fmt"
	"github.com/fluffle/goirc/event"
	"github.com/fluffle/goirc/client"
	"launchpad.net/gobson/bson"
	"lib/db"
	"lib/factoids"
	"lib/util"
	"log"
	"strings"
)

const _FD_NAME string = "factoids"

type factoidDriver struct {
	*factoids.FactoidCollection

	// Keep a reference to the last factoid looked up around
	// for use with 'edit that' and 'delete that' commands.
	lastseen bson.ObjectId
}

func FactoidDriver(db *db.Database) Driver {
	fc, err := factoids.Collection(db)
	if err != nil {
		log.Fatalf("factoid collection failed: %v\n", err)
	}
	return &factoidDriver{fc, ""}
}

type FactoidHandler func(*client.Conn, *client.Line, *factoidDriver)

// Unboxer for FactoidDriver handlers
func FDHandler(f FactoidHandler) event.Handler {
	return event.NewHandler(func(ev ...interface{}) {
		f(ev[0].(*client.Conn), ev[1].(*client.Line), ev[2].(*factoidDriver))
	})
}

func (fd *factoidDriver) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler("privmsg", client.IRCHandler(fd_privmsg))
	r.AddHandler("action", client.IRCHandler(fd_action))
	r.AddHandler("fd_lookup", FDHandler(fd_lookup))
	r.AddHandler("fd_add", FDHandler(fd_add))
	r.AddHandler("fd_delete", FDHandler(fd_delete))
	r.AddHandler("fd_replace", FDHandler(fd_replace))
}

func (fd *factoidDriver) Name() string {
	return _FD_NAME
}

func fd_privmsg(irc *client.Conn, line *client.Line) {
	fd := getFD(irc)
	l, p := util.RemovePrefixedNick(strings.TrimSpace(line.Args[1]), irc.Me.Nick)
	// We want line.Args[1] to contain the (possibly) stripped version of itself
	// but modifying the pointer will result in other goroutines seeing the
	// change, so we need to copy line for our own edification.
	nl := line.Copy()
	nl.Args[1] = l
	l = strings.ToLower(l)

	if !p {
		// If we're not being addressed directly, short-circuit to lookup.
		irc.Dispatcher.Dispatch("fd_lookup", irc, nl, fd)
		return
	}

	// Test for various possible courses of action.
	switch {
	// Factoid add: 'key := value' or 'key :is value'
	case strings.Index(l, ":=") != -1: fallthrough
	case strings.Index(l, ":is") != -1:
		irc.Dispatcher.Dispatch("fd_add", irc, nl, fd)
	// Factoid delete: 'forget|delete that' => deletes fd.lastseen
	case strings.HasPrefix(l, "forget that"): fallthrough
	case strings.HasPrefix(l, "delete that"):
		irc.Dispatcher.Dispatch("fd_delete", irc, nl, fd)
	// Factoid replace: 'replace that with' => updates fd.lastseen
	case strings.HasPrefix(l, "replace that with "):
		// chop off the "known" bit to leave just the replacement
		nl.Args[1] = nl.Args[1][18:]
		irc.Dispatcher.Dispatch("fd_replace", irc, nl, fd)
	// If we get to here, none of the other FD command possibilities
	// have matched, so try a lookup...
	default:
		irc.Dispatcher.Dispatch("fd_lookup", irc, nl, fd)
	}
}

func fd_action(irc *client.Conn, line *client.Line) {
	fd := getFD(irc)
	// Actions just trigger a lookup.
	irc.Dispatcher.Dispatch("fd_lookup", irc, line, fd)
}

func fd_add(irc *client.Conn, line *client.Line, fd *factoidDriver) {
	var key, val string
	if strings.Index(line.Args[1], ":=") != -1 {
		kv := strings.Split(line.Args[1], ":=", 2)
		key = strings.ToLower(strings.TrimSpace(kv[0]))
		val = strings.TrimSpace(kv[1])
	} else {
		// we use :is to add val = "key is val"
		kv := strings.Split(line.Args[1], ":is", 2)
		key = strings.ToLower(strings.TrimSpace(kv[0]))
		val = strings.Join([]string{strings.TrimSpace(kv[0]),
			"is", strings.TrimSpace(kv[1])}, " ")
	}
	n := db.StorableNick{line.Nick, line.Ident, line.Host}
	c := db.StorableChan{line.Args[0]}
	fact := factoids.NewFactoid(key, val, n, c)
	if err := fd.Insert(fact); err == nil {
		count := fd.GetCount(key)
		irc.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: Woo, I now know %d things about '%s'.",
			line.Nick, count, key))
	} else {
		irc.Privmsg(line.Args[0], fmt.Sprintf("Oh no! %s.", err))
	}
}

func fd_delete(irc *client.Conn, line *client.Line, fd *factoidDriver) {
	// Get fresh state on the last seen factoid.
	if fact := fd.GetById(fd.lastseen); fact != nil {
		if err := fd.Remove(bson.M{"_id": fd.lastseen}); err == nil {
			irc.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I forgot that '%s' was '%s'.",
				line.Nick, fact.Key, fact.Value))
		} else {
			irc.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I failed to forget '%s': %s",
				line.Nick, fact.Key, err))
		}
	} else {
		irc.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: Whatever that was, I've already forgotten it.", line.Nick))
	}
	fd.lastseen = ""
}

func fd_lookup(irc *client.Conn, line *client.Line, fd *factoidDriver) {
	key := strings.ToLower(strings.TrimSpace(line.Args[1]))
	var fact *factoids.Factoid

	if fact = fd.GetPseudoRand(key); fact == nil && line.Cmd == "ACTION" {
		// Support sp0rkle's habit of stripping off it's own nick
		// but only for actions, not privmsgs.
		if strings.HasSuffix(key, irc.Me.Nick) {
			key = strings.TrimSpace(key[:len(key)-len(irc.Me.Nick)])
			fact = fd.GetPseudoRand(key)
		}
	}
	if fact != nil {
		fd.lastseen = fact.Id
		switch fact.Type {
		case factoids.F_ACTION:
			irc.Action(line.Args[0], fact.Value)
		default:
			irc.Privmsg(line.Args[0], fact.Value)
		}
	}
}

func fd_replace(irc *client.Conn, line *client.Line, fd *factoidDriver) {
	if fact := fd.GetById(fd.lastseen); fact != nil {
		old := fact.Value
		fact.Value = line.Args[1]
		if err := fd.Update(bson.M{"_id": fd.lastseen}, fact); err == nil {
			irc.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: '%s' was '%s', now is '%s'.",
				line.Nick, fact.Key, old, fact.Value))
		} else {
			irc.Privmsg(line.Args[0], fmt.Sprintf(
				"%s: I failed to replace '%s': %s",
				line.Nick, fact.Key, err))
		}
	} else {
		irc.Privmsg(line.Args[0], fmt.Sprintf(
			"%s: Whatever that was, I've already forgotten it.", line.Nick))
	}
	fd.lastseen = ""
}

func getFD(irc *client.Conn) *factoidDriver {
	return getState(irc).drivers[_FD_NAME].(*factoidDriver)
}
