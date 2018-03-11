package push

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/pushes"
	"golang.org/x/oauth2"
)

// TODO(fluffle):
//  - End to end encryption: https://docs.pushbullet.com/#end-to-end-encryption
//  - Log and surface push errors
//  - Allow users to view and update the devices sp0rkle knows about / pushes to
//  - Read api docs and add features

var (
	pushClientID = flag.String("push_client_id", "",
		"Pushbullet client ID.")
	pushClientSecret = flag.String("push_client_secret", "",
		"Pushbullet client secret.")
)

func pushAPI(path string) string {
	return "https://api.pushbullet.com" + path
}

func config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     bot.GetSecret(*pushClientID),
		ClientSecret: bot.GetSecret(*pushClientSecret),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.pushbullet.com/authorize",
			TokenURL: pushAPI("/oauth2/token"),
		},
		RedirectURL: bot.HttpHost() + "/oauth/auth",
		Scopes:      []string{"everything"},
	}
}

func client(s *pushes.State) *http.Client {
	return config().Client(oauth2.NoContext, s.Token)
}

func checkResponseOK(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		errmsg := &struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Param   string `json:"param,omitempty"`
				Cat     string `json:"cat"`
			} `json:"error"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(errmsg); err != nil {
			return errors.New(resp.Status)
		}
		return fmt.Errorf("%s: %s", resp.Status, errmsg.Error.Message)
	}
	return nil
}

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

func AuthCodeURL(s *pushes.State) string {
	return config().AuthCodeURL(s.State())
}

func Exchange(code string) (*oauth2.Token, error) {
	// Pushbullet don't support passing client secret via http basic auth headers.
	return config().Exchange(oauth2.NoContext, code)
}

func GetDevices(s *pushes.State) ([]*Device, error) {
	u := pushAPI("/v2/devices")
	resp, err := client(s).Get(u)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if err := checkResponseOK(resp); err != nil {
		return nil, fmt.Errorf("GET %s: %v", u, err)
	}
	devs := &struct {
		Devices []*Device `json:"devices"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(devs); err != nil {
		return nil, fmt.Errorf("GET %s JSON decode failed: %v", u, err)
	}
	return devs.Devices, nil
}

func Confirm(s *pushes.State) error {
	if s.CanConfirm() {
		return push(s, "Pushbullet PIN = "+s.Pin,
			"Tell sp0rkle 'push auth <pin>' to complete setup.")
	}
	return errors.New("Not in correct state to send confirmation push.")
}

func Push(s *pushes.State, title, body string) error {
	if s.CanPush() {
		return push(s, title, body)
	}
	return errors.New("Push not enabled.")
}

func push(s *pushes.State, title, body string) error {
	u := pushAPI("/v2/pushes")
	enc, err := json.Marshal(&struct {
		Iden  string `json:"device_iden"`
		Type  string `json:"type"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}{s.Iden, "note", title, body})
	if err != nil {
		return fmt.Errorf("POST %s JSON encode failed: %v", u, err)
	}
	resp, err := client(s).Post(u, "application/json", bytes.NewBuffer(enc))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if err := checkResponseOK(resp); err != nil {
		return fmt.Errorf("POST %s: %v", u, err)
	}
	return nil
}
