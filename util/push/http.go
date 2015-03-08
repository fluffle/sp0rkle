package push

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/fluffle/golog/logging"
)

func authTokenRedirect(rw http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Redirect(rw, req, failureURL+"?fail=parse", 302)
		return
	}
	if req.FormValue("error") != "" {
		http.Redirect(rw, req, failureURL+"?fail=denied", 302)
		return
	}
	id := req.FormValue("state")
	s := getState(id)
	if id == "" || s == nil {
		http.Redirect(rw, req, failureURL+"?fail=nostate", 302)
		return
	}
	code := req.FormValue("code")
	if code == "" {
		http.Redirect(rw, req, failureURL+"?fail=notoken", 302)
		return
	}
	tok, err := getAccessToken(code)
	if err != nil {
		logging.Error("Failed to get access token for %s: %v", s.Nick, err)
		http.Redirect(rw, req, failureURL+"?fail=exchange", 302)
		return
	}

	s.Token = tok
	setState(id, s)
	http.Redirect(rw, req, deviceURL+"?state="+id, 302)
}

func chooseDevice(rw http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Redirect(rw, req, failureURL+"?fail=parse", 302)
		return
	}
	id := req.FormValue("state")
	s := getState(id)
	if id == "" || s == nil {
		http.Redirect(rw, req, failureURL+"?fail=nostate", 302)
		return
	}
	if req.Method == "POST" {
		if s.Iden = req.FormValue("iden"); s.Iden == "" {
			http.Redirect(rw, req, failureURL+"?fail=noiden", 302)
			return
		}
		pin := fmt.Sprintf("%06x", rand.Intn(1e6))
		if err := confirmPush(pin, s); err != nil {
			http.Redirect(rw, req, failureURL+"?fail=push", 302)
			return
		}
		// Store state under the PIN now (lazy...)
		setState(pin, s)
		delState(id)
		http.Redirect(rw, req, successURL, 302)
		return
	}
	// get device list and print a form
	devs, err := getDevices(s.Token)
	if err != nil || len(devs) == 0 {
		logging.Error("Failed to get devices for %s: %v", s.Nick, err)
		http.Redirect(rw, req, failureURL+"?fail=device", 302)
		return
	}
	if err = deviceTmpl.Execute(rw, &deviceData{id, devs}); err != nil {
		logging.Error("Template execution failed: %v", err)
		// assuming here that failure occured because we couldn't write
		return
	}
}

func youAreTehWinnar(rw http.ResponseWriter, req *http.Request) {
	bytes.NewBufferString(successHtml).WriteTo(rw)
}

func youAreTehLosar(rw http.ResponseWriter, req *http.Request) {
	f := "parse"
	if err := req.ParseForm(); err == nil {
		f = "nofail"
		if _, ok := failures[req.FormValue("fail")]; ok {
			f = req.FormValue("fail")
		}
	}
	failureTmpl.Execute(rw, failureData{failures[f]})
}
