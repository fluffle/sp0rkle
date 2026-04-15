package conf

import (
	"fmt"
	"strings"

	"github.com/fluffle/sp0rkle/db"
)

const (
	COLLECTION = "conf"
	// Conf namespace for per-nick timezones
	zoneNs = "timezones"
)

var bolt db.C

func Bolt(ns string) *namespace {
	bolt.Init(db.Bolt.Keyed(), COLLECTION, nil)
	return &namespace{ns: ns, Collection: &bolt}
}

func Ns(ns string) Namespace {
	return Bolt(ns)
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
	Value   any
}

func (e Entry) K() db.Key {
	return db.K{db.S{"ns", e.Ns}, db.S{"key", e.Key}}
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
