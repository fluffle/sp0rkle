package main

import (
	"fmt"
	"github.com/fluffle/goirc/client"
	"lib/util"
	"os"
	"rand"
	"strings"
	"strconv"
	"time"
)

type FactoidPlugin func(string, *client.Line) string

type PluginManager interface {
	AddPlugin(p FactoidPlugin)
	ApplyPlugins(string, *client.Line) string
}

type PluginProvider interface {
	RegisterPlugins(PluginManager)
}

func (fd *factoidDriver) AddPlugin(p FactoidPlugin) {
	fd.plugins = append(fd.plugins, p)
}

func (fd *factoidDriver) ApplyPlugins(val string, line *client.Line) string {
	for _, p := range fd.plugins {
		val = p(val, line)
	}
	return val
}

func (fd *factoidDriver) RegisterPlugins(pm PluginManager) {
	// pm == fd in this case, but meh.
	pm.AddPlugin(FactoidPlugin(plug_identifiers))
	pm.AddPlugin(FactoidPlugin(plug_rand))
}

// Replicate perlfu's $<stuff> identifiers
func plug_identifiers(val string, line *client.Line) string {
	ts := time.LocalTime()
	return id_replacer(val, line, ts)
}

// Split this out so we can inject a deterministic time for testing.
func id_replacer(val string, line *client.Line, ts *time.Time) string {
	val = strings.Replace(val, "$nick", line.Nick, -1)
	val = strings.Replace(val, "$chan", line.Args[0], -1)
	val = strings.Replace(val, "$username", line.Ident, -1)
	val = strings.Replace(val, "$user", line.Ident, -1)
	val = strings.Replace(val, "$host", line.Host, -1)
	val = strings.Replace(val, "$date", ts.Format(time.ANSIC), -1)
	val = strings.Replace(val, "$time", ts.Format("15:04:05"), -1)
	return val
}

// Replicate the "rand" plugin's behaviour.
var myrand *rand.Rand = util.NewRand(time.Nanoseconds())

func plug_rand(val string, line *client.Line) string {
	return rand_replacer(val, myrand)
}

// Split this out so we can inject a deterministic rand.Rand for testing.
// It's at times like this I miss easy number -> string conversion
// and first-class regex constructs. Doing without is fun!
func rand_replacer(val string, r *rand.Rand) string {
	for {
		var lo, hi float32
		var err os.Error
		format := "%.0f"
		// Work out the indices of the plugin start and end.
		ps := strings.Index(val, "<plugin=rand ")
		if ps == -1 {
			break
		}
		pe := strings.Index(val[ps:], ">")
		if pe == -1 {
			// WTF!?
			break
		}
		pe += ps
		// Mid is where the plugin args start.
		mid := ps + 13
		// If there's a space before the plugin ends, we also have a format.
		sp := strings.Index(val[mid:pe], " ")
		if sp != -1 {
			sp += mid
			format = strings.TrimSpace(val[sp:pe])
		} else {
			sp = pe
		}
		// If there's a dash before the space or the plugin ends, we have a
		// range lo-hi, rather than just 0-hi.
		if dash := strings.Index(val[mid:sp], "-"); dash != -1 {
			dash += mid
			if lo, err = strconv.Atof32(val[mid:dash]); err != nil {
				lo = 0
			}
			if hi, err = strconv.Atof32(val[dash+1 : sp]); err != nil {
				hi = 0
			}
		} else {
			lo = 0
			if hi, err = strconv.Atof32(val[mid:sp]); err != nil {
				hi = 0
			}
		}
		rnd := r.Float32()*(hi-lo) + lo
		val = val[:ps] + fmt.Sprintf(format, rnd) + val[pe+1:]
	}
	return val
}
