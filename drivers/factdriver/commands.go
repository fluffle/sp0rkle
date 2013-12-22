package factdriver

import (
	"fmt"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util"
	"labix.org/v2/mgo/bson"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Factoid chance: 'chance of that is' => sets chance of lastSeen[chan]
func chance(ctx *bot.Context) {
	str := ctx.Text()
	var chance float64

	if strings.HasSuffix(str, "%") {
		// Handle 'chance of that is \d+%'
		if i, err := strconv.Atoi(str[:len(str)-1]); err != nil {
			ctx.ReplyN("'%s' didn't look like a % chance to me.", str)
			return
		} else {
			chance = float64(i) / 100
		}
	} else {
		// Assume the chance is a floating point number.
		if c, err := strconv.ParseFloat(str, 64); err != nil {
			ctx.ReplyN("'%s' didn't look like a chance to me.", str)
			return
		} else {
			chance = c
		}
	}

	// Make sure the chance we've parsed lies in (0.0,1.0]
	if chance > 1.0 || chance <= 0.0 {
		ctx.ReplyN("'%s' was outside possible chance ranges.", str)
		return
	}

	// Retrieve last seen ObjectId, replace with ""
	ls := LastSeen(ctx.Target(), "")
	// ok, we're good to update the chance.
	if fact := fc.GetById(ls); fact != nil {
		// Store the old chance, update with the new
		old := fact.Chance
		fact.Chance = chance
		// Update the Modified field
		fact.Modify(ctx.Storable())
		// And store the new factoid data
		if err := fc.Update(bson.M{"_id": ls}, fact); err == nil {
			ctx.ReplyN("'%s' was at %.0f%% chance, now is at %.0f%%.",
				fact.Key, old*100, chance*100)

		} else {
			ctx.ReplyN("I failed to replace '%s': %s", fact.Key, err)
		}
	} else {
		ctx.ReplyN("Whatever that was, I've already forgotten it.")
	}
}

// Pulls out regexp or replacement, allowing for escaped delimiters.
func extractRx(l *util.Lexer, delim rune) string {
	ret, i := "", 0
	for {
		ret += l.Find(delim)
		for i = len(ret) - 1; i >= 0 && ret[i] == '\\'; i-- { }
		if l.Peek() == 0 || (len(ret)-i)%2 == 1 {
			// Even number of backslashes at end of string
			// => delimiter isn't escaped. (Or we're at EOF).
			break
		}
		ret += l.Next()
	}
	return ret
}

// Factoid edit: that =~ s/<regex>/<replacement>/
func edit(ctx *bot.Context) {
	// extract regexp and replacement
	l := &util.Lexer{Input: ctx.Text()}
	if l.Next() != "s" {
		ctx.ReplyN("It's 'that =~ s/<regex>/<replacement>/', fool.")
		return
	}
	delim := l.Peek()         // Identify delimiting character
	l.Next()                  // Skip past that delimiter
	re := extractRx(l, delim) // Extract regex from string
	l.Next()                  // Skip past next delimiter
	rp := extractRx(l, delim) // Extract replacement from string
	if l.Next() != string(delim) {
		ctx.ReplyN("Couldn't parse regex: re='%s', rp='%s'.", re, rp)
		return
	}
	rx, err := regexp.Compile(re)
	if err != nil {
		ctx.ReplyN("Couldn't compile regex '%s': %s", re, err)
		return
	}
	// Retrieve last seen ObjectId, replace with ""
	ls := LastSeen(ctx.Target(), "")
	fact := fc.GetById(ls)
	if fact == nil {
		ctx.ReplyN("I've forgotten what we were talking about, sorry!")
		return
	}
	old := fact.Value
	fact.Value = rx.ReplaceAllString(old, rp)
	fact.Modify(ctx.Storable())
	if err := fc.UpdateId(ls, fact); err == nil {
		ctx.ReplyN("'%s' was '%s', is now '%s'.",
			fact.Key, old, fact.Value)

	} else {
		ctx.ReplyN("I failed to replace '%s': %s", fact.Key, err)
	}
}

// Factoid delete: 'forget|delete that' => deletes lastSeen[chan]
func forget(ctx *bot.Context) {
	// Get fresh state on the last seen factoid.
	ls := LastSeen(ctx.Target(), "")
	if fact := fc.GetById(ls); fact != nil {
		if err := fc.Remove(bson.M{"_id": ls}); err == nil {
			ctx.ReplyN("I forgot that '%s' was '%s'.",
				fact.Key, fact.Value)

		} else {
			ctx.ReplyN("I failed to forget '%s': %s", fact.Key, err)
		}
	} else {
		ctx.ReplyN("Whatever that was, I've already forgotten it.")
	}
}

// Factoid info: 'fact info key' => some information about key
func info(ctx *bot.Context) {
	key := ToKey(ctx.Text(), false)
	count := fc.GetCount(key)
	if count == 0 {
		ctx.ReplyN("I don't know anything about '%s'.", key)
		return
	}
	msgs := make([]string, 0, 10)
	if key == "" {
		msgs = append(msgs, fmt.Sprintf("In total, I know %d things.", count))
	} else {
		msgs = append(msgs, fmt.Sprintf("I know %d things about '%s'.",
			count, key))
	}
	if created := fc.GetLast("created", key); created != nil {
		c := created.Created
		msgs = append(msgs, "A factoid")
		if key != "" {
			msgs = append(msgs, fmt.Sprintf("for '%s'", key))
		}
		msgs = append(msgs, fmt.Sprintf("was last created on %s by %s,",
			c.Timestamp.Format(time.ANSIC), c.Nick))
	}
	if modified := fc.GetLast("modified", key); modified != nil {
		m := modified.Modified
		msgs = append(msgs, fmt.Sprintf("modified on %s by %s,",
			m.Timestamp.Format(time.ANSIC), m.Nick))
	}
	if accessed := fc.GetLast("accessed", key); accessed != nil {
		a := accessed.Accessed
		msgs = append(msgs, fmt.Sprintf("and accessed on %s by %s.",
			a.Timestamp.Format(time.ANSIC), a.Nick))
	}
	if info := fc.InfoMR(key); info != nil {
		if key == "" {
			msgs = append(msgs, "These factoids have")
		} else {
			msgs = append(msgs, fmt.Sprintf("'%s' has", key))
		}
		msgs = append(msgs, fmt.Sprintf(
			"been modified %d times and accessed %d times.",
			info.Modified, info.Accessed))
	}
	ctx.ReplyN("%s", strings.Join(msgs, " "))
}

// Factoid literal: 'literal key' => info about factoid
func literal(ctx *bot.Context) {
	key := ToKey(ctx.Text(), false)
	if count := fc.GetCount(key); count == 0 {
		ctx.ReplyN("I don't know anything about '%s'.", key)
		return
	} else if count > 10 && ctx.Public() {
		ctx.ReplyN("I know too much about '%s', ask me privately.", key)
		return
	}

	if facts := fc.GetAll(key); facts != nil {
		for _, fact := range facts {
			// Use Privmsg directly here so that the results aren't output
			// via the plugin system and contain the literal data.
			ctx.Privmsg(ctx.Target(), fmt.Sprintf(
				"[%3.0f%%] %s", fact.Chance*100, fact.Value))
		}
	} else {
		ctx.ReplyN("Something literally went wrong :-(")
	}
}

// Factoid replace: 'replace that with' => updates lastSeen[chan]
func replace(ctx *bot.Context) {
	ls := LastSeen(ctx.Target(), "")
	if fact := fc.GetById(ls); fact != nil {
		// Store the old factoid value
		old := fact.Value
		// Replace the value with the new one
		fact.Value = ctx.Text()
		// Update the Modified field
		fact.Modify(ctx.Storable())
		// And store the new factoid data
		if err := fc.Update(bson.M{"_id": ls}, fact); err == nil {
			ctx.ReplyN("'%s' was '%s', now is '%s'.",
				fact.Key, old, fact.Value)

		} else {
			ctx.ReplyN("I failed to replace '%s': %s", fact.Key, err)
		}
	} else {
		ctx.ReplyN("Whatever that was, I've already forgotten it.")
	}
}

// Factoid search: 'fact search regexp' => list of possible key matches
func search(ctx *bot.Context) {
	keys := fc.GetKeysMatching(ctx.Text())
	if keys == nil || len(keys) == 0 {
		ctx.ReplyN("I couldn't think of anything matching '%s'.",
			ctx.Text())

		return
	}
	// RESULTS.
	count := len(keys)
	if count > 10 {
		res := strings.Join(keys[:10], "', '")
		ctx.ReplyN(
			"I found %d keys matching '%s', here's the first 10: '%s'.",
			count, ctx.Text(), res)

	} else {
		res := strings.Join(keys, "', '")
		ctx.ReplyN(
			"I found %d keys matching '%s', here they are: '%s'.",
			count, ctx.Text(), res)

	}
}
