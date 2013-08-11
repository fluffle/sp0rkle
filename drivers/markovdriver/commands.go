package markovdriver

import (
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/conf"
	chain "github.com/fluffle/sp0rkle/util/markov"
	"strings"
)

func enableMarkov(ctx *bot.Context) {
	conf.Ns(markovNs).String(strings.ToLower(ctx.Nick), "markov")
	ctx.ReplyN("I'll markov you like I markov'd your mum last night.")
}

func disableMarkov(ctx *bot.Context) {
	key := strings.ToLower(ctx.Nick)
	conf.Ns(markovNs).Delete(key)
	if err := mc.ClearTag("user:"+key); err != nil {
		ctx.ReplyN("Failed to clear tag: %s", err)
		return
	}
	ctx.ReplyN("Sure, bro, I'll stop.")
}

func randomCmd(ctx *bot.Context) {
	if len(ctx.Text()) == 0 {
		ctx.ReplyN("Be who? Your mum?")
		return
	}
	whom := strings.ToLower(strings.Fields(ctx.Text())[0])
	if whom == strings.ToLower(ctx.Me()) {
		ctx.ReplyN("Ha, you're funny. No, wait. Retarded... I meant retarded.")
		return
	}
	if !shouldMarkov(whom) {
		if whom == strings.ToLower(ctx.Nick) {
			ctx.ReplyN("You're not recording markov data. "+
				"Use 'markov me' to enable collection.")
		} else {
			ctx.ReplyN("Not recording markov data for %s.", ctx.Text())
		}
		return
	}
	source := mc.Source("user:" + whom)
	if out, err := chain.Sentence(source); err == nil {
		ctx.Reply("%s would say: %s", ctx.Text(), out)
	} else {
		ctx.ReplyN("markov error: %v", err)
	}
}

func insult(ctx *bot.Context) {
	source := mc.Source("tag:insult")
	whom, lc := ctx.Text(), strings.ToLower(ctx.Text())
	if lc == strings.ToLower(ctx.Me()) || lc == "yourself" {
		ctx.ReplyN("Ha, you're funny. No, wait. Retarded... I meant retarded.")
		return
	}
	if lc == "me" {
		whom = ctx.Nick
	}
	if out, err := chain.Sentence(source); err == nil {
		if len(whom) > 0 {
			ctx.Reply("%s: %s", whom, out)
		} else {
			ctx.Reply("%s", out)
		}
	} else {
		ctx.ReplyN("markov error: %v", err)
	}
}

func learn(ctx *bot.Context) {
	s := strings.SplitN(ctx.Text(), " ", 2)
	if len(s) != 2 {
		ctx.ReplyN("I can't learn from you, you're an idiot.")
		return
	}

	// Prepending "tag:" prevents people from learning as "user:foo".
	mc.AddSentence(s[1], "tag:"+s[0])
	if ctx.Public() {
		// Allow large-scale learning via privmsg by not replying there.
		ctx.ReplyN("Ta. You're a fount of knowledge, you are.")
	}
}
