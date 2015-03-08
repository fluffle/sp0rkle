package push

import (
	"errors"
	"flag"
	"net/http"
	"net/url"
	"strings"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/conf"
)

var (
	pushNs conf.Namespace

	authURL    = bot.HttpHost() + "/oauth/auth"
	deviceURL  = bot.HttpHost() + "/oauth/device"
	successURL = bot.HttpHost() + "/oauth/success"
	failureURL = bot.HttpHost() + "/oauth/failure"

	pushClientID = flag.String("push_client_id", "",
		"Pushbullet client ID.")
	pushClientSecret = flag.String("push_client_secret", "",
		"Pushbullet client secret.")
	pushNotEnabledError = errors.New("pushbullet is not enabled")
)

// must be public for JSON decode of above, feh.
type Device struct {
	Iden         string  `json:"iden"`
	PushToken    string  `json:"push_token"`
	AppVersion   int     `json:"app_version"`
	Fingerprint  string  `json:"fingerprint"`
	Active       bool    `json:"active"`
	Nickname     string  `json:"nickname"`
	Manufacturer string  `json:"manufacturer"`
	Type         string  `json:"type"`
	Created      float32 `json:"created"`
	Modified     float32 `json:"modified"`
	Model        string  `json:"model"`
	Pushable     bool    `json:"pushable"`
}

func Enabled() bool {
	return !(*pushClientID == "" || *pushClientSecret == "")
}

func Init() {
	if !Enabled() {
		return
	}
	pushNs = conf.Ns("push")
	http.HandleFunc("/oauth/auth", authTokenRedirect)
	http.HandleFunc("/oauth/device", chooseDevice)
	http.HandleFunc("/oauth/success", youAreTehWinnar)
	http.HandleFunc("/oauth/failure", youAreTehLosar)
}

func StartFor(nick, pin string) error {
	s := getState(pin)
	if s == nil || nick != s.Nick {
		return errors.New("No authentication state found.")
	}
	setToken(nick, s.Token)
	setIden(nick, s.Iden)
	return nil
}

func StopFor(nick string) {
	pushNs.Delete("token:" + strings.ToLower(nick))
	pushNs.Delete("iden:" + strings.ToLower(nick))
}

// GenAuthURL generates a URL to start the OAuth dance for a nick.
// Once the user visits this URL and approves the access, they
// get redirected back to /oauth/auth (see http.go). There, we
// complete the OAuth dance and request a list of devices. The
// user chooses the target device for push notifications. Once
// this is done, Push() can push messages to the nick.
func GenAuthURL(nick string) string {
	return "https://www.pushbullet.com/authorize?" + url.Values{
		"client_id":     []string{bot.GetSecret(*pushClientID)},
		"redirect_uri":  []string{authURL},
		"response_type": []string{"code"},
		"state":         []string{newState(nick)},
	}.Encode()
}

func Push(nick, title, body string) error {
	token, iden := getToken(nick), getIden(nick)
	if !Enabled() || token == "" || iden == "" {
		return pushNotEnabledError
	}
	return push(token, iden, title, body)
}
