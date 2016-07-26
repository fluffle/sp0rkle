package conf

import (
	"fmt"
	"strings"

	"github.com/fluffle/goirc/logging"
	"github.com/fluffle/sp0rkle/db"
	"gopkg.in/mgo.v2"
)

const (
	COLLECTION = "conf"
	// Conf namespace for per-nick timezones
	zoneNs = "timezones"
)

var mongo db.C

func mongoIndexes(c db.Collection) {
	err := c.Mongo().EnsureIndex(mgo.Index{Key: []string{"ns", "key"}, Unique: true})
	if err != nil {
		logging.Error("Couldn't create index on sp0rkle.conf: %s", err)
	}
}

func Mongo(ns string) *namespace {
	mongo.Init(db.Mongo, COLLECTION, mongoIndexes)
	return &namespace{ns: ns, Collection: &mongo}
}

var bolt db.C

func Bolt(ns string) *namespace {
	bolt.Init(db.Bolt, COLLECTION, nil)
	return &namespace{ns: ns, Collection: &bolt}
}

var checker db.M

func Ns(ns string) *both {
	checker.Init(migrator{}, COLLECTION)
	return &both{bolt: Bolt(ns), mongo: Mongo(ns), Checker: checker}
}

// Lazy, I shouldn't really do this ;-)
func Zone(nick string, tz ...string) string {
	if len(tz) > 0 && tz[0] == "" {
		Ns(zoneNs).Delete(strings.ToLower(nick))
		return ""
	}
	return Ns(zoneNs).String(strings.ToLower(nick), tz...)
}

type Entry struct {
	Ns, Key string
	Value   interface{}
}

func (e Entry) K() db.Key {
	return db.K{{"ns", e.Ns}, {"key", e.Key}}
}

var _ db.Keyer = (*Entry)(nil)

func (e Entry) String() string {
	return fmt.Sprintf("%s<%s: %v>", e.Ns, e.Key, e.Value)
}

// To make migration easier.
type Entries []Entry

func (es Entries) Strings() []string {
	s := make([]string, len(es))
	for i, e := range es {
		s[i] = e.String()
	}
	return s
}
