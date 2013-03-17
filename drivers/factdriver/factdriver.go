package factdriver

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/factoids"
	"github.com/fluffle/sp0rkle/util"
	"labix.org/v2/mgo/bson"
	"strings"
)

// We talk to the factoids collection
var fc *factoids.Collection

// Keep a reference to the last factoid looked up around
// for use with 'edit that' and 'delete that' commands.
// Do this on a per-channel basis to avoid (too much) confusion.
var lastSeen = map[string]bson.ObjectId{}

func Init() {
	fc = factoids.Init()

	bot.Handle(insert, client.PRIVMSG)
	bot.Handle(lookup, client.PRIVMSG, client.ACTION)

	bot.Rewrite(replaceIdentifiers)

	bot.Command(chance, "chance of that is",
		"chance  -- Sets trigger chance of the last displayed factoid value.")
	bot.Command(edit, "that =~",
		"=~ s/regex/replacement/ -- Edits the last factoid value using regex.")
	bot.Command(forget, "delete that",
		"delete  -- Forgets the last displayed factoid value.")
	bot.Command(forget, "forget that",
		"forget  -- Forgets the last displayed factoid value.")
	bot.Command(info, "fact info",
		"fact info <key>  -- Displays some stats about factoid <key>.")
	bot.Command(literal, "literal",
		"literal <key>  -- Displays the factoid values stored for <key>.")
	bot.Command(replace, "replace that with",
		"replace  -- Replaces the last displayed factoid value.")
	bot.Command(search, "fact search",
		"fact search <regexp>  -- Searches for factoids matching <regexp>.")
}

func LastSeen(ch string, id ...bson.ObjectId) bson.ObjectId {
	if len(id) > 0 {
		old, ok := lastSeen[ch]
		lastSeen[ch] = id[0]
		if ok && old != "" {
			return old
		}
	} else if ls, ok := lastSeen[ch]; ok {
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
