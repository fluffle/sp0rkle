package karmadriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/karma"
)

func recordKarma(ctx *bot.Context) {
	// Karma can look like some.text.string++ or (text with spaces)--
	// and there could be multiple occurrences of it in a string.
	nick, _ := ctx.Storable()
	for _, kt := range karmaThings(ctx.Text()) {
		k := kc.KarmaFor(kt.thing)
		if k == nil {
			k = karma.New(kt.thing)
		}
		if kt.plus {
			k.Plus(nick)
		} else {
			k.Minus(nick)
		}
		if _, err := kc.Upsert(k.Id(), k); err != nil {
			ctx.Reply("Failed to insert Karma: %s", err)
		}
	}
}
