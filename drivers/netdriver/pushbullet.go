package netdriver

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/fluffle/goirc/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util/push"
)

func pushFailureURL(fail string) string {
	return bot.HttpHost() + "/oauth/failure?fail=" + fail
}

func pushSuccessURL() string {
	return bot.HttpHost() + "/oauth/success"
}

func pushDeviceURL(state string) string {
	return bot.HttpHost() + "/oauth/device?state=" + state
}

// pushEnable generates a URL to start the OAuth dance for a nick.
// Once the user visits this URL and approves the access, they
// get redirected back to /oauth/auth, which is handled by
// authTokenRedirect below. There, we complete the OAuth dance
// and redirect to /oauth/devices (chooseDevice) which lists
// the user's devices and accepts a POST to choose a target
// device for push notifications. Once this is done, we push
// a confirmation notification to the chosen device with a 6
// digit pin and require that they msg that to us via IRC.
func pushEnable(ctx *bot.Context) {
	if s := pc.GetByNick(ctx.Nick); s != nil {
		if s.HasAlias(ctx.Nick) {
			ctx.ReplyN("Your nick is already used as an alias for %s.", s.Nick)
			return
		}
		if s.CanPush() {
			ctx.ReplyN("Pushes already enabled.")
			return
		}
		ctx.Privmsg(ctx.Nick, "Hmm. Deleting partially-complete state...")
		pc.DelState(s)
	}
	s, err := pc.NewState(ctx.Nick)
	if err != nil {
		ctx.ReplyN("Error creating push state: %v", err)
		return
	}
	// Privmsg the URL so randoms don't click it.
	ctx.Privmsg(ctx.Nick, "Hi! Visit the following URL while logged into "+
		"the account you want to use to push to your device.")
	ctx.Privmsg(ctx.Nick, push.AuthCodeURL(s))
}

func pushDisable(ctx *bot.Context) {
	// Do not search by aliases here: it allows someone to change nick
	// to a known alias and then disable pushes for that user.
	s := pc.GetByNick(ctx.Nick, false)
	if s == nil {
		ctx.ReplyN("Pushes not enabled.")
		return
	}
	if err := pc.DelState(s); err != nil {
		ctx.ReplyN("Error deleting push state: %v", err)
		return
	}
	ctx.ReplyN("Ok, pushes disabled.")
}

func pushConfirm(ctx *bot.Context) {
	pin := strings.Fields(ctx.Text())[0]
	s := pc.GetByNick(ctx.Nick, false)
	switch {
	case s == nil:
		ctx.ReplyN("No authentication state found.")
		return
	case s.Done:
		ctx.ReplyN("Pushes already enabled.")
		return
	case pin != s.Pin:
		ctx.ReplyN("Incorrect pin.")
		return
	}
	s.Done = true
	if err := pc.SetState(s); err != nil {
		ctx.ReplyN("Error setting push state: %v", err)
		return
	}
	ctx.ReplyN("Pushes enabled! Yay!")
}

func pushAddAlias(ctx *bot.Context) {
	alias := strings.Fields(ctx.Text())[0]
	s := pc.GetByNick(ctx.Nick, false)
	if s == nil || !s.CanPush() {
		ctx.ReplyN("Pushes not enabled.")
		return
	}
	if s.HasAlias(alias) {
		ctx.ReplyN("Alias %q already exists.", alias)
		return
	}
	if a := pc.GetByNick(alias); a != nil {
		ctx.ReplyN("Alias %q already exists for nick %s.", alias, a.Nick)
		return
	}
	s.AddAlias(alias)
	if err := pc.SetState(s); err != nil {
		ctx.ReplyN("Error setting push state: %v", err)
		return
	}
	ctx.ReplyN("Added alias %q to your push state.", alias)
}

func pushDelAlias(ctx *bot.Context) {
	alias := strings.Fields(ctx.Text())[0]
	s := pc.GetByNick(ctx.Nick, false)
	if s == nil || !s.CanPush() {
		ctx.ReplyN("Pushes not enabled.")
		return
	}
	if !s.HasAlias(alias) {
		ctx.ReplyN("%q is not one of your aliases.", alias)
		return
	}
	s.DelAlias(alias)
	if err := pc.SetState(s); err != nil {
		ctx.ReplyN("Error setting push state: %v", err)
		return
	}
	ctx.ReplyN("Deleted alias %q from your push state.", alias)
}

func pushAuthHTTP(rw http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Redirect(rw, req, pushFailureURL("parse"), 302)
		return
	}
	if req.FormValue("error") != "" {
		http.Redirect(rw, req, pushFailureURL("denied"), 302)
		return
	}
	id := req.FormValue("state")
	s := pc.GetByB64(id)
	if id == "" || s == nil {
		http.Redirect(rw, req, pushFailureURL("nostate"), 302)
		return
	}
	code := req.FormValue("code")
	if code == "" {
		http.Redirect(rw, req, pushFailureURL("notoken"), 302)
		return
	}
	tok, err := push.Exchange(code)
	if err != nil {
		logging.Error("Failed to get access token for %s: %v", s.Nick, err)
		http.Redirect(rw, req, pushFailureURL("exchange"), 302)
		return
	}

	s.Token = tok
	pc.SetState(s)
	http.Redirect(rw, req, pushDeviceURL(id), 302)
}

func pushDeviceHTTP(rw http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Redirect(rw, req, pushFailureURL("parse"), 302)
		return
	}
	id := req.FormValue("state")
	s := pc.GetByB64(id)
	if id == "" || s == nil {
		http.Redirect(rw, req, pushFailureURL("nostate"), 302)
		return
	}
	if req.Method == "POST" {
		if s.Iden = req.FormValue("iden"); s.Iden == "" {
			http.Redirect(rw, req, pushFailureURL("noiden"), 302)
			return
		}
		s.Pin = fmt.Sprintf("%06x", rand.Intn(1e6))
		if err := push.Confirm(s); err != nil {
			logging.Error("Failed to send confirmation push for %s: %v", s.Nick, err)
			http.Redirect(rw, req, pushFailureURL("push"), 302)
			return
		}
		pc.SetState(s)
		http.Redirect(rw, req, pushSuccessURL(), 302)
		return
	}
	// get device list and print a form
	devs, err := push.GetDevices(s)
	if err != nil  {
		logging.Error("Failed to get devices for %s: %v", s.Nick, err)
		http.Redirect(rw, req, pushFailureURL("device"), 302)
		return
	}
	if len(devs) == 0 {
		strings.NewReader(pushNoDeviceHTML).WriteTo(rw)
		return
	}
	if err = pushDeviceTmpl.Execute(rw, &pushDevice{id, devs}); err != nil {
		logging.Error("Template execution failed: %v", err)
		// assuming here that failure occured because we couldn't write
		return
	}
}

func pushSuccessHTTP(rw http.ResponseWriter, req *http.Request) {
	strings.NewReader(pushSuccessHTML).WriteTo(rw)
}

func pushFailureHTTP(rw http.ResponseWriter, req *http.Request) {
	f := "parse"
	if err := req.ParseForm(); err == nil {
		f = "nofail"
		if _, ok := pushFailures[req.FormValue("fail")]; ok {
			f = req.FormValue("fail")
		}
	}
	pushFailureTmpl.Execute(rw, pushFailure{pushFailures[f]})
}
