package netdriver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func get(req string) (string, bool) {
	res, err := http.Get(req)
	defer res.Body.Close()
	if err != nil {
		return fmt.Sprintf("HTTP error: %v", err), false
	}
	txt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Sprintf("Read error: %v", err), false
	}
	return string(txt), true
}

type httpService struct {
	uri string
}

func (hs *httpService) LookupResult(req string) string {
	ret, _ := get(fmt.Sprintf(hs.uri, url.QueryEscape(req)))
	return ret
}
