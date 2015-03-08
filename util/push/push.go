package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	confirmSubj = "Pushbullet PIN = "
	confirmBody = "Tell sp0rkle 'push auth <pin>' to complete setup."
)

type pushData struct {
	Iden  string `json:"device_iden"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func confirmPush(pin string, s *oauthState) error {
	return push(s.Token, s.Iden, confirmSubj+pin, confirmBody)
}

func push(token, iden, title, body string) error {
	u := apiURL + "/v2/pushes"
	p := &pushData{iden, "note", title, body}
	enc, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("POST %s JSON encode failed: %v", u, err)
	}
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(enc))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if err := checkResponseOK(resp); err != nil {
		return fmt.Errorf("POST %s: %v", u, err)
	}
	return nil
}
