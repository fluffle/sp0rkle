package factdriver

import (
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/lib/factoids"
	"github.com/fluffle/sp0rkle/lib/util"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"labix.org/v2/mgo/bson"
	"strings"
)

const driverName string = "factoids"

type factoidDriver struct {
	*factoids.FactoidCollection

	// Keep a reference to the last factoid looked up around
	// for use with 'edit that' and 'delete that' commands.
	// Do this on a per-channel basis to avoid (too much) confusion.
	lastseen map[string]bson.ObjectId

	// A list of text processing plugins to apply to factoid values
	plugins []base.Plugin
}

func FactoidDriver(db *db.Database) *factoidDriver {
	fc := factoids.Collection(db)
	return &factoidDriver{
		FactoidCollection: fc,
		lastseen:          make(map[string]bson.ObjectId),
		plugins:           make([]base.Plugin, 0),
	}
}

func (fd *factoidDriver) Name() string {
	return driverName
}

func (fd *factoidDriver) AddPlugin(p base.Plugin) {
	fd.plugins = append(fd.plugins, p)
}

func (fd *factoidDriver) ApplyPlugins(val string, line *base.Line) string {
	for _, p := range fd.plugins {
		val = p.Apply(val, line)
	}
	return val
}

func (fd *factoidDriver) Lastseen(ch string, id ...bson.ObjectId) bson.ObjectId {
	if len(id) > 0 {
		old, ok := fd.lastseen[ch]
		fd.lastseen[ch] = id[0]
		if ok && old != "" {
			return old
		}
	} else if lastseen, ok := fd.lastseen[ch]; ok {
		return lastseen
	}
	return ""
}

// Does some standard processing on s to make it key-like.
func ToKey(s string, prefixes bool) string {
	// Lowercase and strip leading/trailing spaces and (some) punctuation
	s = strings.ToLower(strings.Trim(s, "!?., "))
	// Remove IRC formatting and colours
	s = util.RemoveColours(util.RemoveFormatting(s))
	if prefixes {
		// Remove commonly-written "prefixes" (see lib/util/prefix.rl)
		// We don't always want this, so guard it with a boolean.
		s = util.RemovePrefixes(s)
	}
	return s
}
