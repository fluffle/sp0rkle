package markovdriver

import (
	"github.com/fluffle/sp0rkle/apis/llama"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util"
)

func (d *Driver) insultPlugin(in string, ctx *bot.Context) string {
	f := func(string) string {
		out, err := llama.Complete(randomPrompt())
		if err == nil {
			return out
		}
		return "<plugin error>"
	}
	return util.ApplyPluginFunction(in, "insult", f)
}
