package netdriver

import (
	"encoding/json"
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"net/url"
	"strings"
	"time"
)

const udUrl = "http://api.urbandictionary.com/v0/define?term=%s"

// TODO(fluffle): Put this in util and clean up the various copies of it.
const TimeFormat = "15:04:05, Monday 2 January 2006"

// These shamelessly stolen from StalkR:
// https://github.com/StalkR/misc/blob/master/urbandictionary/urbandictionary.go
type udResult struct {
	Type       string   `json:"result_type"`
	HasRelated bool     `json:"has_related_words"`
	Pages      int      `json:"pages,omitempty"`
	Total      int      `json:"total,omitempty"`
	Sounds     []string `json:"sounds,omitempty"`
	List       []udDef  `json:"list"`
}

type udDef struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
	Example    string `json:"example"`
	Author     string `json:"author"`
	Id         int    `json:"defid"`
	Url        string `json:"permalink"`
	Vote       string `json:"current_vote"`
	Upvotes    int    `json:"thumbs_up"`
	Downvotes  int    `json:"thumbs_down"`
	Term       string `json:"term,omitempty"`
	Type       string `json:"type,omitempty"`
}

type udCacheEntry struct {
	result *udResult
	stamp  time.Time
}

type udCache map[string]udCacheEntry

func (udc udCache) prune() {
	for k, v := range udc {
		if time.Since(v.stamp) > 24*time.Hour {
			delete(udc, k)
		}
	}
}

func (udc udCache) fetch(term string) (entry udCacheEntry, ok bool, err error) {
	udc.prune()
	entry, ok = udc[term]
	if ok { return }
	entry.result = &udResult{}
	data, err := get(fmt.Sprintf(udUrl, url.QueryEscape(term)))
	if err != nil { return }
	if err = json.Unmarshal(data, entry.result); err != nil {
		logging.Debug("JSON: %s", data)
		return
	}
	// Abuse Pages and Total for our own ends here
	entry.result.Pages, entry.result.Total = -1, len(entry.result.List)
	entry.stamp = time.Now()
	udc[term] = entry
	return
}

var cache = udCache{}

func urbanDictionary(ctx *bot.Context) {
	entry, ok, err := cache.fetch(strings.ToLower(ctx.Text()))
	if err != nil {
		ctx.ReplyN("ud request failed: %#v", err)
		return
	}
	cached, r := "", entry.result
	if ok {
		cached = fmt.Sprintf(", result cached at %s",
			entry.stamp.Format(TimeFormat))
	}
	if r.Total == 0 || r.Type == "no_results" {
		ctx.ReplyN("%s isn't defined yet%s.", ctx.Text(), cached)
		return
	}
	// Cycle through all the definitions on repeated calls for the same term
	r.Pages = (r.Pages + 1) % r.Total
	def := r.List[r.Pages]
	ctx.Reply("[%d/%d] %s (%d up, %d down%s)", r.Pages+1, r.Total,
		strings.Replace(def.Definition, "\r\n", " ", -1),
		def.Upvotes, def.Downvotes, cached)
}
