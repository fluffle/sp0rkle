package netdriver

import (
	"fmt"
	"github.com/fluffle/golog/logging"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
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
	l   logging.Logger
}

func (hs *httpService) LookupResult(req string) string {
	ret, _ := get(fmt.Sprintf(hs.uri, url.QueryEscape(req)))
	return ret
}

type iGoogleCalcService httpService

func IGoogleCalcService(l logging.Logger) *iGoogleCalcService {
	return &iGoogleCalcService{
		uri: "http://www.google.com/ig/calculator?hl=en&q=%s",
		l: l,
	}
}

func (is *iGoogleCalcService) LookupResult(q string) string {
	data, ok := get(fmt.Sprintf(is.uri, url.QueryEscape(q)))
	is.l.Info("Got '%s' for query '%s'", data, q)
	if !ok {
		return data
	}
	// Irritatingly, the iGoogle calculator sends back malformed JSON
	// where the object keys aren't quoted, so we do this manually.
	// The format is: {lhs: "2 + 2",rhs: "4",error: "",icc: false}
	rx := regexp.MustCompile(`rhs: "([^"]*)",error: "([^"]*)"`)
	res := rx.FindStringSubmatch(data)
	if res == nil {
		return fmt.Sprintf("Regex error: no match in '%s'", data)
	}
	if res[2] != "" {
		return res[2] // err
	}
	return res[1] // rhs
}
