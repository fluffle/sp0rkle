package markovdriver

import (
	"strings"

	"github.com/fluffle/sp0rkle/apis/llama"
	"github.com/fluffle/sp0rkle/bot"
	chain "github.com/fluffle/sp0rkle/util/markov"
)

func (d *Driver) enableMarkov(ctx *bot.Context) {
	d.confNs.String(strings.ToLower(ctx.Nick), "markov")
	ctx.ReplyN("I'll markov you like I markov'd your mum last night.")
}

func (d *Driver) disableMarkov(ctx *bot.Context) {
	key := strings.ToLower(ctx.Nick)
	d.confNs.Delete(key)
	if err := d.mc.ClearTag("user:" + key); err != nil {
		ctx.ReplyN("Failed to clear tag: %s", err)
		return
	}
	ctx.ReplyN("Sure, bro, I'll stop.")
}

func (d *Driver) randomCmd(ctx *bot.Context) {
	if len(ctx.Text()) == 0 {
		ctx.ReplyN("Be who? Your mum?")
		return
	}
	whom := strings.ToLower(strings.Fields(ctx.Text())[0])
	if whom == strings.ToLower(ctx.Me()) {
		ctx.ReplyN("Ha, you're funny. No, wait. Retarded... I meant retarded.")
		return
	}
	if !d.shouldMarkov(whom) {
		if whom == strings.ToLower(ctx.Nick) {
			ctx.ReplyN("You're not recording markov data. " +
				"Use 'markov me' to enable collection.")
		} else {
			ctx.ReplyN("Not recording markov data for %s.", ctx.Text())
		}
		return
	}
	source := d.mc.Source("user:" + whom)
	if out, err := chain.Sentence(source); err == nil {
		ctx.Reply("%s would say: %s", ctx.Text(), out)
	} else {
		ctx.ReplyN("markov error: %v", err)
	}
}

func (d *Driver) insult(ctx *bot.Context) {
	whom, lc := ctx.Text(), strings.ToLower(ctx.Text())
	if lc == strings.ToLower(ctx.Me()) || lc == "yourself" {
		ctx.ReplyN("Ha, you're funny. No, wait. Retarded... I meant retarded.")
		return
	}
	if lc == "me" {
		whom = ctx.Nick
	}
	out, err := llama.Complete(randomPrompt())
	if err == nil {
		if len(whom) > 0 {
			ctx.Reply("%s: %s", whom, out)
		} else {
			ctx.Reply("%s", out)
		}
	} else {
		ctx.ReplyN("The LLM couldn't be bothered: %v", err)
	}
}

func (d *Driver) learn(ctx *bot.Context) {
	s := strings.SplitN(ctx.Text(), " ", 2)
	if len(s) != 2 {
		ctx.ReplyN("I can't learn from you, you're an idiot.")
		return
	}

	// Prepending "tag:" prevents people from learning as "user:foo".
	d.mc.AddSentence(s[1], "tag:"+s[0])
	if ctx.Public() {
		// Allow large-scale learning via privmsg by not replying there.
		ctx.ReplyN("Ta. You're a fount of knowledge, you are.")
	}
}
