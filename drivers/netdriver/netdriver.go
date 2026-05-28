package netdriver

import (
	"context"
	"net/http"
	"sync"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db/conf"
	"github.com/fluffle/sp0rkle/collections/pushes"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"github.com/fluffle/sp0rkle/util/push"
	"github.com/fluffle/sp0rkle/util/flights"
	"github.com/google/go-github/github"
)

type Driver struct {
	ctx    context.Context
	pc     *pushes.Collection
	rc     *reminders.Collection
	gh     *github.Client
	fp     *flightPoller
	confNs conf.Namespace
}


func New(b *bot.Bot, rc *reminders.Collection, pc *pushes.Collection, config *conf.Registry) *Driver {
	d := &Driver{ctx: b.Ctx(), rc: rc, pc: pc, confNs: config.Ns("mc")}

	b.Command(d.urbanDictionary, "ud", "ud <term>  -- "+
		"Look up <term> on UrbanDictionary.")

	srv := d.confNs.String(mcServer)
	if srv != "" {
		if st, err := pollServer(srv); err == nil {
			st.confNs = d.confNs
			logging.Info("Starting MC poller for '%s'", srv)
			b.Poll(st)
			b.HandleBG(st.Topic, "332")
		} else {
			logging.Error("Not starting MC poller: %v", err)
		}
	}
	b.Command(d.mcSet, "mc set", "mc set <key> <value>  -- "+
		"Set minecraft server polling config vars.")

	if *githubToken != "" {
		d.rc = rc
		d.gh = githubClient(b)

		b.Handle(d.githubWatcher, client.PRIVMSG)

		b.Command(d.githubCreateIssue, "file bug:", "file bug: <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
		b.Command(d.githubCreateIssue, "file bug", "file bug <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
		b.Command(d.githubCreateIssue, "report bug", "report bug <title>. "+
			"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
		b.Command(d.githubUpdateIssue, "update bug #", "update bug #<number> "+
			"<comment>  -- Adds a comment to bug <number>. Abusers will be hurt.")
	}

	if push.Enabled() {
		d.pc = pc
		b.Command(d.pushEnable, "push enable", "push enable  -- "+
			"Start the OAuth flow to enable pushbullet notifications.")
		b.Command(d.pushDisable, "push disable", "push disable  -- "+
			"Disable pushbullet notifications and delete tokens.")
		b.Command(d.pushConfirm, "push auth", "push auth <pin>  -- "+
			"Confirm pushed PIN to finish pushbullet auth dance.")
		b.Command(d.pushAddAlias, "push add alias", "push add alias  -- "+
			"Add a push alias for your nick.")
		b.Command(d.pushDelAlias, "push del alias", "push del alias  -- "+
			"Delete a push alias for your nick.")

		http.HandleFunc("/oauth/auth", d.pushAuthHTTP)
		http.HandleFunc("/oauth/device", d.pushDeviceHTTP)
		http.HandleFunc("/oauth/success", d.pushSuccessHTTP)
		http.HandleFunc("/oauth/failure", d.pushFailureHTTP)
	}

	if flights.Enabled() {
		logging.Info("Starting flight poller")
		d.fp = &flightPoller{
			Mutex: &sync.Mutex{},
			tracking: make(map[string]*flightInfo),
			aliases: make(map[string]string),
		}
		b.Command(d.stalk, "stalk", "stalk <flight>  -- trackers a flight via AviationStack")
		b.Command(d.stalk, "flight tracking start", "flight tracking start <flight>  -- trackers a flight via AviationStack")
		b.Command(d.unstalk, "unstalk", "unstalk <flight>  -- stops tracking a flight")
		b.Command(d.unstalk, "flight tracking stop", "flight tracking stop <flight>  -- stops tracking a flight")
		b.Command(d.status, "flight status", "flight status <flight> -- returns the cached status of the flight if it is being stalked")
		b.Command(d.alias, "flight alias", "flight alias <flight> <alias> -- adds a memorable alias for the flight")
		b.Poll(d.fp)
	}
	return d
}
