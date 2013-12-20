package netdriver

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"io/ioutil"
	"net/http"
)

func get(req string) ([]byte, error) {
	res, err := http.Get(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func Init() {
	bot.Command(urbanDictionary, "ud", "ud <term>  -- "+
		"Look up <term> on UrbanDictionary.")
	bot.Command(createGitHubIssue, "file bug:", "file bug: <title>. "+
		"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
	bot.Command(createGitHubIssue, "file bug", "file bug <title>. "+
		"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
	bot.Command(createGitHubIssue, "report bug", "file bug: <title>. "+
		"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")

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
}
