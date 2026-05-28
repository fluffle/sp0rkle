package flights

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/fluffle/goirc/logging"
)

const (
	apiBase = "http://api.aviationstack.com/v1/flights"
	dateFormat = "2006-01-02T15:04:05-07:00"
)

var aviationStackAPIKey = flag.String("aviation_stack_api_key", "",
	"API key for AviationStack.")

func Enabled() bool {
	return *aviationStackAPIKey != ""
}

// AviationStack API structs
type apiResponse struct {
	Data []Data `json:"data"`
}

type Data struct {
	FlightStatus string `json:"flight_status"`
	Departure    struct {
		Airport  string  `json:"airport"`
		Code     string  `json:"iata"`
		Terminal string  `json:"terminal"`
		Gate     string  `json:"gate"`
		Delay    float64 `json:"delay"`
	} `json:"departure"`
	Arrival struct {
		Airport   string  `json:"airport"`
		Code      string  `json:"iata"`
		Terminal  string  `json:"terminal"`
		Gate      string  `json:"gate"`
		Baggage   string  `json:"baggage"`
		Delay     float64 `json:"delay"`
	} `json:"arrival"`
	Flight struct {
		IATA string `json:"iata"`
		ICAO string `json:"icao"`
	} `json:"flight"`
	Airline struct {
		Name string `json:"name"`
	} `json:"airline"`
	Aircraft struct {
		Reg  string `json:"registration"`
		Code string `json:"icao24"`
	} `json:"aircraft"`
	Live struct {
		Updated string `json:"updated"`
		Latitude float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Altitude float64 `json:"altitude"`
		Direction float64 `json:"direction"`
		Speed float64 `json:"speed_horizontal"`
		Ground bool `json:"is_ground"`
	} `json:"live"`
}

func (data *Data) String() string {
	return data.status(false)
}

func (data *Data) Details() string {
	return data.status(true)
}

func (data *Data) status(details bool) string {
	b := &strings.Builder{}
	// United Airlines flight UA2402
	fmt.Fprintf(b, "%s flight %s", data.Airline.Name, data.Flight.IATA)
	// from George Bush Intercontinental (IAH) gate E9
	fmt.Fprintf(b, " from %s (%s)", data.Departure.Airport, data.Departure.Code)
	if details && (data.Departure.Terminal != "" || data.Departure.Gate != "") {
		b.WriteString(", departing")
		if data.Departure.Terminal != "" {
			fmt.Fprintf(b, " T", data.Departure.Terminal)
		}
		if data.Departure.Gate != "" {
			fmt.Fprintf(b, " gate %s", data.Departure.Gate)
		}
	}
	// to Newark Liberty International (EWR) gate C109, baggage carousel 2
	fmt.Fprintf(b, " to %s (%s)", data.Arrival.Airport, data.Arrival.Code)
	if details && (data.Arrival.Terminal != "" || data.Arrival.Gate != "") {
		b.WriteString(", arriving")
		if data.Arrival.Terminal != "" {
			fmt.Fprintf(b, " T%s", data.Arrival.Terminal)
		}
		if details && data.Arrival.Gate != "" {
			fmt.Fprintf(b, " gate %s", data.Arrival.Gate)
		}
	}
	if details && data.Arrival.Baggage != "" {
		fmt.Fprintf(b, ", baggage carousel %s", data.Arrival.Baggage)
	}

	status := data.liveStatus()
	if details {
		status = data.liveDetails()
	}
	if status == "" {
		status = data.flightStatus()
	}
	b.WriteString(status)

	dep := data.Departure.Delay > 0
	arr := data.Arrival.Delay > 0
	either, both := dep || arr, dep && arr
	if either {
		b.WriteString("  It is currently delayed by")
		if dep {
			fmt.Fprintf(b, " %.0fm departing", data.Departure.Delay)
		}
		if both {
			fmt.Fprint(b, " and")
		}
		if arr {
			fmt.Fprintf(b, " %.0fm on arrival", data.Arrival.Delay)
		}
		b.WriteString(".")
	}
	return b.String()
}

func (data *Data) liveStatus() string {
	if data.Live.Updated == "" {
		return ""
	}
	if data.Live.Ground {
		return " is on the ground."
	}
	return " is airborne."
}

func (data *Data) liveDetails() string {
	if data.Live.Updated == "" || data.Live.Ground {
		return ""
	}
	lat := math.Abs(data.Live.Latitude)
	latN := "N"
	if data.Live.Latitude < 0 {
		latN = "S"
	}
	lon := math.Abs(data.Live.Longitude)
	lonE := "E"
	if data.Live.Longitude < 0 {
		lonE = "W"
	}
	code := " Aircraft"
	if data.Aircraft.Reg != "" {
		code += " " + data.Aircraft.Reg
	} else if data.Aircraft.Code != "" {
		code += " " + data.Aircraft.Code
	}
	return fmt.Sprintf(". %s is flying at %.0f km/h, altitude %.2f km, bearing %.0f° at %.4f%s, %.4f%s.",
		code, data.Live.Speed, data.Live.Altitude / 1000, data.Live.Direction, lat, latN, lon, lonE)
}

func (data *Data) flightStatus() string {
	switch data.FlightStatus {
	case "active":
		return " is airborne."
	case "landed":
		return " has landed!"
	case "cancelled":
		return " has been cancelled :-("
	case "diverted":
		return " has been diverted :-("
	default:
		return fmt.Sprintf(" is %s.", data.FlightStatus)
	}
}

func (data *Data) Updated(prev *Data) bool {
	return prev.String() != data.String()
}

func (data *Data) Landed(prev *Data) bool {
	if data.FlightStatus == "landed" || data.FlightStatus == "cancelled" {
		return true
	}
	// Frustratingly the API does not always update the FlightStatus.
	// Sooooo... heuristics time.
	if prev == nil {
		// If we don't have a previous instance to compare with, we can't do
		// anything.
		return false
	}
	if prev.Live.Updated != "" && !prev.Live.Ground &&
		(data.Live.Updated == "" ||
			(data.Live.Updated != "" && data.Live.Ground )) {
		// Live tracking transitioned from present and not-ground to
		// either not present or present and ground between two updates
		return true
	}
	return false
}

func QueryAll(flights []string) map[string]*Data {
	if *aviationStackAPIKey == "" {
		logging.Warn("flights.QueryAll: API key not set")
		return nil
	}
	results := make(map[string]*Data)
	for _, flight := range flights {
		data, err := QueryOne(flight)
		if err != nil {
			logging.Error("flights.QueryAll: %v", err)
		} else {
			results[flight] = data
		}
		// Don't send more than 1QPS to the API
		time.Sleep(time.Second)
	}
	return results
}

func QueryOne(flight string) (*Data, error) {
	if *aviationStackAPIKey == "" {
		return nil, fmt.Errorf("querying for flight %q: API key not set", flight)
	}
	data, err := doQuery(flight, "iata")
	if data == nil || err != nil {
		data, err = doQuery(flight, "icao")
	}
	if err != nil {
		return nil, fmt.Errorf("querying for flight %q: %w", flight, err)
	}
	if data == nil {
		return nil, fmt.Errorf("querying for flight %q: no results", flight)
	}

	return data, nil
}

func doQuery(flight, queryType string) (*Data, error) {
	u, _ := url.Parse(apiBase)
	q := u.Query()
	q.Set("access_key", *aviationStackAPIKey)
	q.Set("flight_" + queryType, flight)
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}
	if len(apiResp.Data) == 0 {
		return nil, nil
	}
	logging.Info("json body:\n\n%s\n\n", body)

	// Copy so apiResp etc can be discarded.
	data := apiResp.Data[0]
	return &data, nil
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
