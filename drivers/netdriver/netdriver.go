package netdriver

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"io/ioutil"
	"net/http"
)

// We store 'tell' notices for github updates
var rc *reminders.Collection


func get(req string) ([]byte, error) {
	res, err := http.Get(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func Init() {
	rc = reminders.Init()

	bot.Command(urbanDictionary, "ud", "ud <term>  -- "+
		"Look up <term> on UrbanDictionary.")

	if *mcServer != "" {
		if st, err := pollServer(*mcServer); err == nil {
			bot.Poll(st)
			bot.Handle(func(ctx *bot.Context) {
				st.Topic(ctx)
			}, "332")
		} else {
			logging.Error("Not starting MC poller: %v", err)
		}
	}

	if *githubToken != "" {
		gh := githubClient()
		bot.Poll(githubPoller(gh))
		wrap := func (ctx *bot.Context) { githubCreateIssue(ctx, gh) }

		bot.Command(wrap, "file bug:", "file bug: <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
		bot.Command(wrap, "file bug", "file bug <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
		bot.Command(wrap, "report bug", "file bug: <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
	}
}
