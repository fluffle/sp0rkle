package netdriver

import (
	"encoding/xml"
	"fmt"
	"github.com/fluffle/sp0rkle/bot"
	"net/url"
	"strings"
)

const udUrl = "http://www.urbandictionary.com/define.php?term=%s"

func urbanDictionary(ctx *bot.Context) {
	uri := fmt.Sprintf(udUrl, url.QueryEscape(ctx.Text()))
	d, err := get(uri)
	if err != nil {
		ctx.ReplyN("GET failed: %v", err)
		return
	}
	// Parsing HTML with encoding/xml is a bit meh.
	// First skip to <div class="definition">, then ...
	for {
		tok, err := d.Token()
		if err != nil {
			ctx.ReplyN("HTML parse error: %v", err)
			break
		}
		se, ok := tok.(xml.StartElement)
		if ok &&
			se.Name.Local == "div" &&
			len(se.Attr) == 1 &&
			se.Attr[0].Value == "definition" {
			break
		}
	}
	// ... the next token should be the start of the definition.
	// At this point we assemble a slice of strings containing the
	// text inside this div from the xml.CharData tokens, bolding
	// everything inside <a> tags as other definable names.
	defn := []string{}
LOOP:
	for {
		t, err := d.Token()
		if err != nil {
			ctx.ReplyN("HTML parse error: %v", err)
			break
		}
		switch tok := t.(type) {
		case xml.CharData:
			// These can contain newlines, so turn them to spaces.
			s := strings.Replace(string(tok), "\n", " ", -1)
			defn = append(defn, s)
		case xml.StartElement:
			switch tok.Name.Local {
			case "a":
				// <a>
				defn = append(defn, "\x02")
			case "br":
				defn = append(defn, " ")
			}
		case xml.EndElement:
			switch tok.Name.Local {
			case "a":
				// </a>
				defn = append(defn, "\x02")
			case "div":
				// </div>; bail out
				break LOOP
			}
		}
	}
	str, out := strings.Join(defn, ""), ""
	for len(str) > 350 {
		idx := strings.Index(str[250:], ". ")
		if idx == -1 {
			idx = strings.Index(str[250:], " ")
		} else {
			idx++
		}
		idx += 251 // after ". " or " " or 250th character
		out, str = str[:idx], str[idx:]
		ctx.ReplyN("%s ...", out)
	}
	ctx.ReplyN("%s", str)
}
