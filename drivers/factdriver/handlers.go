package factdriver

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/factoids"
	"github.com/fluffle/sp0rkle/util"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strings"
)

// Factoid add: 'key := value' or 'key :is value'
func insert(ctx *bot.Context) {
	if !ctx.Addressed || !util.IsFactoidAddition(ctx.Text()) {
		return
	}

	var key, val string
	if strings.Index(ctx.Text(), ":=") != -1 {
		kv := strings.SplitN(ctx.Text(), ":=", 2)
		key = ToKey(kv[0], false)
		val = strings.TrimSpace(kv[1])
	} else {
		// we use :is to add val = "key is val"
		kv := strings.SplitN(ctx.Text(), ":is", 2)
		key = ToKey(kv[0], false)
		val = strings.Join([]string{strings.TrimSpace(kv[0]),
			"is", strings.TrimSpace(kv[1])}, " ")
	}
	n, c := ctx.Storable()
	fact := factoids.NewFactoid(key, val, n, c)

	// The "randomwoot" factoid contains random positive phrases for success.
	joy := "Woo"
	if rand := fc.GetPseudoRand("randomwoot"); rand != nil {
		joy = rand.Value
	}

	if err := fc.Insert(fact); err == nil {
		count := fc.GetCount(key)
		LastSeen(ctx.Target(), fact.Id)
		ctx.ReplyN("%s, I now know %d things about '%s'.", joy, count, key)
	} else {
		ctx.ReplyN("Error storing factoid: %s.", err)
	}
}

func lookup(ctx *bot.Context) {
	// Only perform extra prefix removal if we weren't addressed directly
	key := ToKey(ctx.Text(), !ctx.Addressed)
	var fact *factoids.Factoid

	if fact = fc.GetPseudoRand(key); fact == nil && ctx.Cmd == client.ACTION {
		// Support sp0rkle's habit of stripping off it's own nick
		// but only for actions, not privmsgs.
		if strings.HasSuffix(key, ctx.Me()) {
			key = strings.TrimSpace(key[:len(key)-len(ctx.Me())])
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
		LastSeen(ctx.Target(), fact.Id)
		// Update the Accessed field
		// TODO(fluffle): fd should take care of updating Accessed internally
		fact.Access(ctx.Storable())
		// And store the new factoid data
		if err := fc.Update(bson.M{"_id": fact.Id}, fact); err != nil {
			ctx.ReplyN("I failed to update '%s' (%s): %s ",
				fact.Key, fact.Id, err)

		}
		recurse(fact, map[string]bool{key: true})
		switch fact.Type {
		case factoids.F_ACTION:
			ctx.Do("%s", fact.Value)
		default:
			ctx.Reply("%s", fact.Value)
		}
	}
}

// Recursively resolve pointers to other factoids
func recurse(fact *factoids.Factoid, keys map[string]bool) {
	val := fact.Value
	key, start, end := util.FactPointer(val)
	if key == "" { return }
	if _, ok := keys[key]; ok || len(keys) > 20 {
		fact.Value = val[:start] + "[circular reference]" + val[end:]
		return
	}
	keys[key] = true
	if f2 := fc.GetPseudoRand(key); f2 != nil {
		fact.Value = val[:start] + f2.Value + val[end:]
		if start == 0 && fact.Type != f2.Type {
			// Propagate change of factoid type when the pointer
			// is at the beginning of the string.
			fact.Type = f2.Type
		}
		recurse(fact, keys)
		return
	}
	// if we get here, we found a pointer key but no matching factoid
	// so recurse on the stuff after that key *only* to avoid loops.
	fact.Value = val[end:]
	recurse(fact, keys)
	fact.Value = val[:end] + fact.Value
}
