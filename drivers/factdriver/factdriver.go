package factdriver

import (
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/factoids"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util"
	"labix.org/v2/mgo/bson"
	"strings"
)

type factoidFn func(*factoidDriver, *base.Line)

// A factoidCommand fulfills base.Handler and base.Command
type factoidCommand struct {
	fd *factoidDriver
	fn factoidFn
	help string
}

func (fc *factoidCommand) Execute(l *base.Line) {
	fc.fn(fc.fd, l)
}

func (fc *factoidCommand) Help() string {
	return fc.help
}

func (fd *factoidDriver) Command(fn factoidFn, prefix, help string) {
	bot.Command(&factoidCommand{fd, fn, help}, prefix)
}

func (fd *factoidDriver) Handle(fn factoidFn, event ...string) {
	bot.Handle(&factoidCommand{fd, fn, ""}, event...)
}

type factoidDriver struct {
	*factoids.FactoidCollection

	// Keep a reference to the last factoid looked up around
	// for use with 'edit that' and 'delete that' commands.
	// Do this on a per-channel basis to avoid (too much) confusion.
	lastseen map[string]bson.ObjectId
}

func Init(db *db.Database) *factoidDriver {
	fd := &factoidDriver{
		FactoidCollection: factoids.Collection(db),
		lastseen:          make(map[string]bson.ObjectId),
	}

	fd.Handle((*factoidDriver).Add, "privmsg")
	fd.Handle((*factoidDriver).Lookup, "privmsg", "action")

	bot.PluginFunc(replaceIdentifiers)

	fd.Command((*factoidDriver).Chance, "chance of that is",
		"chance  -- Sets trigger chance of the last displayed factoid value.")
	fd.Command((*factoidDriver).Delete, "delete that",
		"delete  -- Forgets the last displayed factoid value.")
	fd.Command((*factoidDriver).Delete, "forget that",
		"forget  -- Forgets the last displayed factoid value.")
	fd.Command((*factoidDriver).Info, "fact info",
		"fact info <key>  -- Displays some stats about factoid <key>.")
	fd.Command((*factoidDriver).Literal, "literal",
		"literal <key>  -- Displays the factoid values stored for <key>.")
	fd.Command((*factoidDriver).Replace, "replace that with",
		"replace  -- Replaces the last displayed factoid value.")
	fd.Command((*factoidDriver).Search, "fact search",
		"fact search <regexp>  -- Searches for factoids matching <regexp>.")
	return fd
}

func (fd *factoidDriver) Lastseen(ch string, id ...bson.ObjectId) bson.ObjectId {
	if len(id) > 0 {
		old, ok := fd.lastseen[ch]
		fd.lastseen[ch] = id[0]
		if ok && old != "" {
			return old
		}
	} else if lastseen, ok := fd.lastseen[ch]; ok {
		return lastseen
	}
	return ""
}

// Does some standard processing on s to make it key-like.
func ToKey(s string, prefixes bool) string {
	// Lowercase and strip leading/trailing spaces and (some) punctuation
	s = strings.ToLower(strings.Trim(s, "!?., "))
	// Remove IRC formatting and colours
	s = util.RemoveColours(util.RemoveFormatting(s))
	if prefixes {
		// Remove commonly-written "prefixes" (see lib/util/prefix.rl)
		// We don't always want this, so guard it with a boolean.
		s = util.RemovePrefixes(s)
	}
	return s
}
