package netdriver

import (
	"github.com/fluffle/goirc/client"
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
	bot.Handle(mcStartPoller, client.CONNECTED)
	bot.Handle(mcStopPoller, client.DISCONNECTED)
	bot.Handle(mcChanTopic, "332")
}
