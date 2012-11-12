package factdriver

import (
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/factoids"
	"github.com/fluffle/sp0rkle/util"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"strings"
)

// Factoid add: 'key := value' or 'key :is value'
func insert(line *base.Line) {
	if !line.Addressed || !util.IsFactoidAddition(line.Args[1]) {
		return
	}

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
	n, c := line.Storable()
	fact := factoids.NewFactoid(key, val, n, c)
	if err := fc.Insert(fact); err == nil {
		count := fc.GetCount(key)
		bot.ReplyN(line, "Woo, I now know %d things about '%s'.", count, key)
	} else {
		bot.ReplyN(line, "Error storing factoid: %s.", err)
	}
}

func lookup(line *base.Line) {
	// Only perform extra prefix removal if we weren't addressed directly
	key := ToKey(line.Args[1], !line.Addressed)
	var fact *factoids.Factoid

	if fact = fc.GetPseudoRand(key); fact == nil && line.Cmd == "ACTION" {
		// Support sp0rkle's habit of stripping off it's own nick
		// but only for actions, not privmsgs.
		if strings.HasSuffix(key, bot.Nick()) {
			key = strings.TrimSpace(key[:len(key)-len(bot.Nick())])
			fact = fc.GetPseudoRand(key)
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
	if rand.Float64() < chance {
		// Store this as the last seen factoid
		LastSeen(line.Args[0], fact.Id)
		// Update the Accessed field
		// TODO(fluffle): fd should take care of updating Accessed internally
		fact.Access(line.Storable())
		// And store the new factoid data
		if err := fc.Update(bson.M{"_id": fact.Id}, fact); err != nil {
			bot.ReplyN(line, "I failed to update '%s' (%s): %s ",
				fact.Key, fact.Id, err)
		}
		keys := map[string]bool{key: true}
		val := recurse(fact.Value, keys)
		switch fact.Type {
		case factoids.F_ACTION:
			bot.Do(line, "%s", val)
		default:
			bot.Reply(line, "%s", val)
		}
	}
}

// Recursively resolve pointers to other factoids
func recurse(val string, keys map[string]bool) string {
	key, start, end := util.FactPointer(val)
	if key == "" {
		return val
	}
	if _, ok := keys[key]; ok || len(keys) > 20 {
		val = val[:start] + "[circular reference]" + val[end:]
		return val
	}
	keys[key] = true
	if fact := fc.GetPseudoRand(key); fact != nil {
		val = val[:start] + fact.Value + val[end:]
		return recurse(val, keys)
	}
	// if we get here, we found a pointer key but no matching factoid
	// so recurse on the stuff after that key *only* to avoid loops.
	return val[:end] + recurse(val[end:], keys)
}
