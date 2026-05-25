package factdriver

import (
	"strings"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/factoids"
	"github.com/fluffle/sp0rkle/util"
	"github.com/fluffle/sp0rkle/util/bson"
)

// We talk to the factoids collection
type Driver struct {
	fc       *factoids.Collection
	// Keep a reference to the last factoid looked up around
	// for use with 'edit that' and 'delete that' commands.
	// Do this on a per-channel basis to avoid (too much) confusion.
	lastSeen map[string]bson.ObjectId
}


func New(b *bot.Bot, fc *factoids.Collection) *Driver {
	d := &Driver{fc: fc, lastSeen: map[string]bson.ObjectId{}}

	b.Handle(d.insert, client.PRIVMSG)
	b.Handle(d.lookup, client.PRIVMSG, client.ACTION)

	b.Rewrite(d.replaceIdentifiers)

	b.Command(d.chance, "chance of that is",
		"chance  -- Sets trigger chance of the last displayed factoid value.")
	b.Command(d.edit, "that =~",
		"=~ s/regex/replacement/ -- Edits the last factoid value using regex.")
	b.Command(d.forget, "delete that",
		"delete  -- Forgets the last displayed factoid value.")
	b.Command(d.forget, "forget that",
		"forget  -- Forgets the last displayed factoid value.")
	b.Command(d.info, "fact info",
		"fact info <key>  -- Displays some stats about factoid <key>.")
	b.Command(d.literal, "literal",
		"literal <key>  -- Displays the factoid values stored for <key>.")
	b.Command(d.replace, "replace that with",
		"replace  -- Replaces the last displayed factoid value.")
	b.Command(d.search, "fact search",
		"fact search <regexp>  -- Searches for factoids matching <regexp>.")

	return d
}

func (d *Driver) LastSeen(ch string, id ...bson.ObjectId) bson.ObjectId {
	if len(id) > 0 {
		old, ok := d.lastSeen[ch]
		d.lastSeen[ch] = id[0]
		if ok && old != "" {
			return old
		}
	} else if ls, ok := d.lastSeen[ch]; ok {
		return ls
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
