package netdriver

import (
	"encoding/xml"
	"github.com/fluffle/sp0rkle/bot"
	"net/http"
)

func get(req string) (*xml.Decoder, error) {
	res, err := http.Get(req)
	if err != nil {
		return nil, err
	}
	d := xml.NewDecoder(res.Body)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity
	return d, nil
}

func Init() {
	bot.CommandFunc(urbanDictionary, "ud", "ud <term>  -- "+
		"Look up <term> on UrbanDictionary.")
}
