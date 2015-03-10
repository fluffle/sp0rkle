package netdriver

import (
	"io/ioutil"
	"net/http"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/conf"
	"github.com/fluffle/sp0rkle/collections/pushes"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"github.com/fluffle/sp0rkle/util/push"
	"github.com/google/go-github/github"
)

var pc *pushes.Collection
var rc *reminders.Collection
var gh *github.Client

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

	mcConf = conf.Ns("mc")
	srv := mcConf.String(mcServer)
	if srv != "" {
		if st, err := pollServer(srv); err == nil {
			logging.Info("Starting MC poller for '%s'", srv)
			bot.Poll(st)
			bot.Handle(func(ctx *bot.Context) {
				st.Topic(ctx)
			}, "332")
		} else {
			logging.Error("Not starting MC poller: %v", err)
		}
	}
	bot.Command(mcSet, "mc set", "mc set <key> <value>  -- "+
		"Set minecraft server polling config vars.")
	// TODO(fluffle): Polling can only be en/disabled at reconnect.
	//	bot.Command(mcPoll, "mc poll", "mc poll start|stop  -- "+
	//		"Enable or disable minecraft server polling.")

	if *githubToken != "" {
		rc = reminders.Init()
		gh = githubClient()

		bot.Handle(githubWatcher, client.PRIVMSG)

		bot.Command(githubCreateIssue, "file bug:", "file bug: <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
		bot.Command(githubCreateIssue, "file bug", "file bug <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
		bot.Command(githubCreateIssue, "report bug", "report bug <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
		bot.Command(githubUpdateIssue, "update bug #", "update bug #<number> "+
			"<comment>  -- Adds a comment to bug <number>. Abusers will be hurt.")
	}

	if push.Enabled() {
		pc = pushes.Init()
		bot.Command(pushEnable, "push enable", "push enable  -- "+
			"Start the OAuth flow to enable pushbullet notifications.")
		bot.Command(pushDisable, "push disable", "push disable  -- "+
			"Disable pushbullet notifications and delete tokens.")
		bot.Command(pushConfirm, "push auth", "push auth <pin>  -- "+
			"Confirm pushed PIN to finish pushbullet auth dance.")

		http.HandleFunc("/oauth/auth", pushAuthHTTP)
		http.HandleFunc("/oauth/device", pushDeviceHTTP)
		http.HandleFunc("/oauth/success", pushSuccessHTTP)
		http.HandleFunc("/oauth/failure", pushFailureHTTP)
	}
}
