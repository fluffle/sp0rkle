package main

import (
	"fmt"
	"github.com/fluffle/goirc/event"
	"github.com/fluffle/goirc/client"
	"lib/db"
	"lib/factoids"
	"lib/util"
	"log"
	"strings"
)

type factoidDriver struct {
	*factoids.FactoidCollection
}

func FactoidDriver(db *db.Database) Driver {
	fc, err := factoids.Collection(db)
	if err != nil {
		log.Fatalf("factoid collection failed: %v\n", err)
	}
	return &factoidDriver{fc}
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
}

func fd_privmsg(irc *client.Conn, line *client.Line) {
	fd := getFD(irc)
	l, p := util.RemovePrefixedNick(strings.TrimSpace(line.Args[1]), irc.Me.Nick)
	// We want line.Args[1] to contain the (possibly) stripped version of itself
	// but modifying the pointer will result in other goroutines seeing the
	// change, so we need to copy line for our own edification.
	nl := line.Copy()
	nl.Args[1] = l

	if p && strings.Index(l, ":=") != -1 {
		// We're being addressed directly, this could be a factoid add.
		// Currently, just support := for adds. English parsing is hard.
		irc.Dispatcher.Dispatch("fd_add", irc, nl, fd)
		return
	}
	// If we get to here, none of the other FD command possibilities
	// have matched, so try a lookup...
	irc.Dispatcher.Dispatch("fd_lookup", irc, nl, fd)
}

func fd_action(irc *client.Conn, line *client.Line) {
	fd := getFD(irc)
	// Actions just trigger a lookup.
	irc.Dispatcher.Dispatch("fd_lookup", irc, line, fd)
}

func fd_add(irc *client.Conn, line *client.Line, fd *factoidDriver) {
	kv := strings.Split(line.Args[1], ":=", 2)
	key := strings.ToLower(strings.TrimSpace(kv[0]))
	val := strings.TrimSpace(kv[1])
	n := db.StorableNick{line.Nick, line.Ident, line.Host}
	c := db.StorableChan{line.Args[0]}
	fact := factoids.NewFactoid(key, val, n, c)
	if err := fd.Insert(fact); err == nil {
		count := fd.GetCount(key)
		irc.Privmsg(line.Args[0],
			fmt.Sprintf("Woo, I now know %d things about '%s'.", count, key))
	} else {
		irc.Privmsg(line.Args[0], fmt.Sprintf("Oh no! %s.", err))
	}
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
		switch fact.Type {
		case factoids.F_ACTION:
			irc.Action(line.Args[0], fact.Value)
		default:
			irc.Privmsg(line.Args[0], fact.Value)
		}
	}
}

func getFD(irc *client.Conn) *factoidDriver {
	return getState(irc).drivers["factoids"].(*factoidDriver)
}
