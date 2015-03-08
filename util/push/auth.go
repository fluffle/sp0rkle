package push

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fluffle/sp0rkle/bot"
)

const apiURL = "https://api.pushbullet.com"

func setToken(nick, token string) {
	pushNs.String("token:"+strings.ToLower(nick), token)
}

func getToken(nick string) string {
	return pushNs.String("token:" + strings.ToLower(nick))
}

func setIden(nick, iden string) {
	pushNs.String("iden:"+strings.ToLower(nick), iden)
}

func getIden(nick string) string {
	return pushNs.String("iden:" + strings.ToLower(nick))
}

// was already using json encoding for api calls
type oauthState struct {
	Nick  string    `json:"nick"`
	Time  time.Time `json:"time"`
	Token string    `json:"token,omitempty"`
	Iden  string    `json:"iden,omitempty"`
}

func newState(nick string) string {
	id := fmt.Sprintf("%016x", rand.Int63())
	setState(id, &oauthState{
		Nick: nick,
		Time: time.Now(),
	})
	return id
}

func getState(id string) *oauthState {
	js := pushNs.Value("state:" + id)
	if js == nil {
		return nil
	}
	s := &oauthState{}
	if err := json.Unmarshal(js.([]byte), s); err != nil {
		return nil
	}
	if time.Now().After(s.Time.Add(time.Hour)) {
		// We have an hour's grace time to complete the auth flow.
		delState(id)
		return nil
	}
	return s
}

func setState(id string, s *oauthState) {
	js, _ := json.Marshal(s)
	pushNs.Value("state:"+id, js)
}

func delState(id string) {
	pushNs.Delete("state:" + id)
}

func getAccessToken(code string) (string, error) {
	if !Enabled() {
		return "", pushNotEnabledError
	}
	u := apiURL + "/oauth2/token"
	resp, err := http.PostForm(u, url.Values{
		"client_id":     []string{bot.GetSecret(*pushClientID)},
		"client_secret": []string{bot.GetSecret(*pushClientSecret)},
		"code":          []string{code},
		"grant_type":    []string{"authorization_code"},
	})
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if err := checkResponseOK(resp); err != nil {
		return "", fmt.Errorf("POST %s: %v", u, err)
	}

	auth := &struct {
		TokenType string `json:"token_type"`
		Token     string `json:"access_token"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(auth); err != nil {
		return "", fmt.Errorf("POST %s JSON decode failed: %v", u, err)
	}
	return auth.Token, nil
}

func getDevices(token string) ([]*Device, error) {
	u := apiURL + "/v2/devices"
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
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
