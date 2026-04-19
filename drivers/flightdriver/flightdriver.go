package flightdriver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/conf"
)

const (
	flightsNs   = "flights"
	apiKeyKey   = "api_key"
	apiBase     = "http://api.aviationstack.com/v1/flights"
	pollFreq    = 10 * time.Minute
	maxStalkAge = 24 * time.Hour
)

type flightInfo struct {
	FlightNum string
	Target    bot.Chan
	// We need to know which server this flight is being tracked on
	Me        string
	LastState string
	LastStatus string
	StartTime time.Time
}

type flightPoller struct {
	sync.Mutex
	// tracking maps "me:channel:flightNum" to flightInfo
	tracking map[string]*flightInfo
}

func (fp *flightPoller) Poll(ctxs []*bot.Context) {
	if len(ctxs) == 0 {
		return
	}

	key := conf.Ns(flightsNs).String(apiKeyKey)
	if key == "" {
		logging.Error("AviationStack API key not set. Use !stalkkey to set it.")
		return
	}

	fp.Lock()
	// Create a copy of the tracking list to process without holding the lock
	toProcess := make([]*flightInfo, 0, len(fp.tracking))
	for k, info := range fp.tracking {
		if time.Since(info.StartTime) > maxStalkAge {
			delete(fp.tracking, k)
			continue
		}
		toProcess = append(toProcess, info)
	}
	fp.Unlock()

	// To avoid calling the API multiple times for the same flight in the same poll
	type apiResult struct {
		status string
		rawStatus string
		err    error
	}
	cache := make(map[string]apiResult)

	for _, info := range toProcess {
		res, ok := cache[info.FlightNum]
		if !ok {
			status, rawStatus, err := getFlightStatus(info.FlightNum, key)
			res = apiResult{status, rawStatus, err}
			cache[info.FlightNum] = res
		}

		if res.err != nil {
			logging.Error("Error getting flight status for %s: %v", info.FlightNum, res.err)
			continue
		}

		if res.status == "" {
			continue
		}

		if res.status != info.LastState {
			// Find the correct context for this flight
			var targetCtx *bot.Context
			for _, ctx := range ctxs {
				if ctx.Me() == info.Me {
					targetCtx = ctx
					break
				}
			}

			if targetCtx != nil {
				targetCtx.Privmsg(string(info.Target), fmt.Sprintf("Flight %s update: %s", info.FlightNum, res.status))

				fp.Lock()
				// Re-verify it's still in the map before updating
				mapKey := fmt.Sprintf("%s:%s:%s", info.Me, info.Target, info.FlightNum)
				if refreshedInfo, ok := fp.tracking[mapKey]; ok {
					refreshedInfo.LastState = res.status
					refreshedInfo.LastStatus = res.rawStatus
					if res.rawStatus == "landed" {
						delete(fp.tracking, mapKey)
					}
				}
				fp.Unlock()
			}
		}
	}
}

func (fp *flightPoller) Start() {}
func (fp *flightPoller) Stop()  {}
func (fp *flightPoller) Tick() time.Duration {
	return pollFreq
}

var fp *flightPoller

func Init() {
	fp = &flightPoller{
		tracking: make(map[string]*flightInfo),
	}
	bot.Command(stalk, "stalk", "stalk <flight>  -- trackers a flight via AviationStack")
	bot.Command(stalkoff, "stalkoff", "stalkoff <flight>  -- stops tracking a flight")
	bot.Command(stalkkey, "stalkkey", "stalkkey <key>  -- sets AviationStack API key")
	bot.Poll(fp)
}

func stalk(ctx *bot.Context) {
	fn := strings.ToUpper(strings.TrimSpace(ctx.Text()))
	if fn == "" {
		ctx.ReplyN("Which flight do you want to stalk?")
		return
	}

	fp.Lock()
	defer fp.Unlock()
	_, ch := ctx.Storable()
	me := ctx.Me()
	key := fmt.Sprintf("%s:%s:%s", me, ch, fn)
	fp.tracking[key] = &flightInfo{
		FlightNum: fn,
		Target:    ch,
		Me:        me,
		StartTime: time.Now(),
	}
	ctx.ReplyN("Now stalking flight %s", fn)
}

func stalkoff(ctx *bot.Context) {
	fn := strings.ToUpper(strings.TrimSpace(ctx.Text()))
	if fn == "" {
		ctx.ReplyN("Which flight do you want to stop stalking?")
		return
	}

	fp.Lock()
	defer fp.Unlock()
	_, ch := ctx.Storable()
	me := ctx.Me()
	key := fmt.Sprintf("%s:%s:%s", me, ch, fn)
	if _, ok := fp.tracking[key]; ok {
		delete(fp.tracking, key)
		ctx.ReplyN("Stopped stalking flight %s", fn)
	} else {
		ctx.ReplyN("I wasn't stalking flight %s in this channel.", fn)
	}
}

func stalkkey(ctx *bot.Context) {
	key := strings.TrimSpace(ctx.Text())
	if key == "" {
		ctx.ReplyN("Please provide an API key.")
		return
	}
	conf.Ns(flightsNs).String(apiKeyKey, key)
	ctx.ReplyN("AviationStack API key updated.")
}

// AviationStack API structs
type apiResponse struct {
	Data []struct {
		FlightStatus string `json:"flight_status"`
		Departure    struct {
			Airport  string      `json:"airport"`
			Delay    interface{} `json:"delay"`
			Scheduled string      `json:"scheduled"`
		} `json:"departure"`
		Arrival struct {
			Airport string      `json:"airport"`
			Delay   interface{} `json:"delay"`
		} `json:"arrival"`
		Airline struct {
			Name string `json:"name"`
		} `json:"airline"`
	} `json:"data"`
}

func getFlightStatus(flightNum, apiKey string) (string, string, error) {
	u, _ := url.Parse(apiBase)
	q := u.Query()
	q.Set("access_key", apiKey)
	q.Set("flight_iata", flightNum)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", "", err
	}

	if len(apiResp.Data) == 0 {
		q.Del("flight_iata")
		q.Set("flight_icao", flightNum)
		u.RawQuery = q.Encode()
		resp, err = http.Get(u.String())
		if err == nil {
			defer resp.Body.Close()
			body, _ = ioutil.ReadAll(resp.Body)
			json.Unmarshal(body, &apiResp)
		}
	}

	if len(apiResp.Data) == 0 {
		return "", "", nil
	}

	data := apiResp.Data[0]
	res := fmt.Sprintf("%s from %s to %s is %s.",
		data.Airline.Name, data.Departure.Airport, data.Arrival.Airport, data.FlightStatus)

	depDelay := formatDelay(data.Departure.Delay)
	arrDelay := formatDelay(data.Arrival.Delay)

	if depDelay != "" {
		res += " Departure delay: " + depDelay
	}
	if arrDelay != "" {
		res += " Arrival delay: " + arrDelay
	}

	return res, data.FlightStatus, nil
}

func formatDelay(d interface{}) string {
	if d == nil {
		return ""
	}
	switch v := d.(type) {
	case float64:
		if v == 0 {
			return ""
		}
		return fmt.Sprintf("%.0f mins", v)
	case int:
		if v == 0 {
			return ""
		}
		return fmt.Sprintf("%d mins", v)
	}
	return ""
}
