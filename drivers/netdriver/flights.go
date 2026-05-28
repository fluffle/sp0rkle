package netdriver

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"
	"github.com/fluffle/goirc/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util/datetime"
	"github.com/fluffle/sp0rkle/util/flights"
)

const (
	pollFreq    = 1 * time.Hour
	maxStalkAge = 24 * time.Hour
)

type subscriber struct {
	target    string
	network   string
}

type flightInfo struct {
	flight string
	data   *flights.Data
	start  time.Time
	last   time.Time
	subs   map[subscriber]bool
}

type flightUpdate struct {
	updated bool
	landed bool
}

func (fi *flightInfo) old() bool {
	return time.Now().Sub(fi.start) > 24*time.Hour
}

func (fi *flightInfo) update(data *flights.Data) flightUpdate {
	last := fi.data
	fi.data = data
	fi.last = time.Now()
	return flightUpdate{
		updated: data.Updated(last),
		landed: data.Landed(last),
	}
}

func (fi *flightInfo) sendMsgs(ctxs map[string]*bot.Context) {
	msg := fmt.Sprintf("%s Tracked since %s.", fi.data.String(), datetime.Format(fi.start))
	for sub := range fi.subs {
		ctx, ok := ctxs[sub.network]
		if !ok {
			// No longer connected to this network, drop subscriber
			delete(fi.subs, sub)
			continue
		}
		ctx.Privmsg(sub.target, msg)
	}
}

type flightPoller struct {
	*sync.Mutex
	tracking map[string]*flightInfo
	aliases map[string]string
}

func (fp *flightPoller) alias(flight, alias string) bool {
	fp.Lock()
	defer fp.Unlock()
	flight = strings.ToUpper(flight)
	if flight == strings.ToUpper(alias) {
		return false
	}
	if _, ok := fp.tracking[flight]; !ok {
		return false
	}
	fp.aliases[alias] = flight
	return true
}

func (fp *flightPoller) get(flight string) *flightInfo {
	fp.Lock()
	defer fp.Unlock()
	return fp.getLocked(flight)
}

func (fp *flightPoller) getLocked(flight string) *flightInfo {
	upper := strings.ToUpper(flight)
	if info, ok := fp.tracking[upper]; ok {
		return info
	}
	if aliased, ok := fp.aliases[flight]; ok {
		return fp.getLocked(aliased)
	}
	return nil
}

func (fp *flightPoller) stalk(flight string, data *flights.Data, sub subscriber) bool {
	fp.Lock()
	defer fp.Unlock()
	info := fp.getLocked(flight)
	if info != nil {
		last := info.data
		if data.Landed(last) {
			// stop tracking if we're tracking and the flight landed
			fp.cleanupLocked(info.flight)
			return false
		}
		// permit re-tracking to force-update status
		info.data = data
		// permit re-tracking to reset expiry timer
		info.start = time.Now()
		// Add subscriber
		info.subs[sub] = true
		return true
	}
	// If we aren't tracking already, at least check the flight status
	// before tracking; don't track landed flights.
	if data.Landed(nil)	{
		return false
	}
	flight = strings.ToUpper(flight)
	// Not tracking yet, add new info.
	fp.tracking[flight] = &flightInfo{
		flight: flight,
		data:   data,
		start:  time.Now(),
		last:   time.Now(),
		subs:   map[subscriber]bool{sub: true},
	}
	return true
}

func (fp *flightPoller) unstalk(flight string, sub subscriber) bool {
	fp.Lock()
	defer fp.Unlock()
	info := fp.getLocked(flight)
	if info == nil {
		return false
	}
	delete(info.subs, sub)
	if len(info.subs) == 0 {
		fp.cleanupLocked(info.flight)
	}
	return true
}

func (fp *flightPoller) cleanupLocked(flight string) {
	delete(fp.tracking, flight)
	for alias, num := range fp.aliases {
		if flight == num {
			delete(fp.aliases, alias)
		}
	}
}

func (fp *flightPoller) flights() []string {
	fp.Lock()
	defer fp.Unlock()
	return slices.Collect(maps.Keys(fp.tracking))
}

func (fp *flightPoller) Poll(ctxs []*bot.Context) {
	if len(ctxs) == 0 {
		return
	}
	toQuery := fp.flights()
	if len(toQuery) == 0 {
		return
	}
	results := flights.QueryAll(toQuery)
	if len(results) == 0 {
		logging.Warn("No query responses from API.")
		return
	}

	fp.Lock()
	defer fp.Unlock()
	for _, flight := range toQuery {
		data := results[flight]
		info, ok := fp.tracking[flight]
		if !ok {
			// flight deleted while we were querying
			continue
		}
		upd := info.update(data)
		if upd.updated {
			info.sendMsgs(networkMap(ctxs))
		}
		if upd.landed {
			// stop stalking tracked flights when they land
			fp.cleanupLocked(flight)
		}
	}
}

func networkMap(ctxs []*bot.Context) map[string]*bot.Context {
	m := map[string]*bot.Context{}
	for _, ctx := range ctxs {
		m[ctx.Network()] = ctx
	}
	return m
}

func (fp *flightPoller) Start() {}

func (fp *flightPoller) Stop()  {}

func (fp *flightPoller) Tick() time.Duration {
	return pollFreq
}

func (d *Driver) stalk(ctx *bot.Context) {
	flight := strings.ToUpper(strings.TrimSpace(ctx.Text()))
	if flight == "" {
		ctx.ReplyN("Which flight do you want to stalk?")
		return
	}
	data, err := flights.QueryOne(flight)
	if err != nil {
		ctx.ReplyN("Not stalking: %v", err)
		return
	}
	sub := subscriber{
		target: ctx.Target(),
		network: ctx.Network(),
	}
	if d.fp.stalk(flight, data, sub) {
		for _, code := range []string{data.Flight.IATA, data.Flight.ICAO} {
			d.fp.alias(flight, code)
		}
		ctx.ReplyN("Stalking flight %s: %s", flight, data.String())
	} else {
		ctx.ReplyN("Not stalking flight %s: %s", flight, data.String())
	}
}

func (d *Driver) unstalk(ctx *bot.Context) {
	flight := strings.TrimSpace(ctx.Text())
	if flight == "" {
		ctx.ReplyN("Which flight do you want to stop stalking?")
		return
	}
	sub := subscriber{
		target:    ctx.Target(),
		network:   ctx.Network(),
	}
	if d.fp.unstalk(flight, sub) {
		ctx.ReplyN("Stopped stalking %s here.", flight)
	} else {
		ctx.ReplyN("I wasn't stalking %s here.", flight)
	}
}

func (d *Driver) status(ctx *bot.Context) {
	flight := strings.TrimSpace(ctx.Text())
	if flight == "" {
		ctx.ReplyN("Which flight do you want status for?")
		return
	}
	if info := d.fp.get(flight); info != nil {
		ctx.ReplyN("%s Last updated at %s.", info.data.Details(), datetime.Format(info.last))
	} else {
		ctx.ReplyN("I'm not currently stalking %s.", flight)
	}
}

func (d *Driver) alias(ctx *bot.Context) {
	fields := strings.Fields(ctx.Text())
	if len(fields) < 2 {
		ctx.ReplyN("I need a flight number and an alias?")
		return
	}
	flight := strings.ToUpper(fields[0])
	alias := fields[1]
	if d.fp.alias(flight, alias) {
		ctx.ReplyN("Added alias %q for %s.", alias, flight)
	} else {
		ctx.ReplyN("I'm not currently stalking %s.", flight)
	}
}
