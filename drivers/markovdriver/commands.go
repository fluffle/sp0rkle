package markovdriver

import (
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/bot"
	umarkov "github.com/fluffle/sp0rkle/util/markov"
	"time"
)

func randomCmd(line *base.Line) {
	source := mc.CreateSourceForTag("user:" + line.Args[1])
	seed := time.Now().UTC().UnixNano()
	first_word := umarkov.SENTENCE_START

	out, err := umarkov.Generate(source, first_word, seed, 150)
	if err == nil {
		bot.ReplyN(line, "%s would say: %s", line.Args[1], out)
	} else {
		bot.ReplyN(line, "error: %v", err)
	}
}
