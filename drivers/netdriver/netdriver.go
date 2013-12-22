package netdriver

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"github.com/google/go-github/github"
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

		bot.Handle(wrap(githubWatcher, gh), client.PRIVMSG)

		bot.Command(wrap(githubCreateIssue, gh), "file bug:", "file bug: <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
		bot.Command(wrap(githubCreateIssue, gh), "file bug", "file bug <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
		bot.Command(wrap(githubCreateIssue, gh), "report bug", "file bug: <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
	}
}

func wrap(f func(*bot.Context, *github.Client), gh *github.Client) func (*bot.Context) {
	return func(ctx *bot.Context) {
		f(ctx, gh)
	}
}
