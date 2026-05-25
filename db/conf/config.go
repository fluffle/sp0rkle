package conf

import (
	"fmt"

	"github.com/fluffle/sp0rkle/db"
)

const (
	COLLECTION = "conf"
)

func New(c db.KeyedCollection[*Entry]) *Registry {
	return &Registry{collection: c}
}

type Registry struct {
	collection db.KeyedCollection[*Entry]
}

func (r *Registry) Ns(name string) Namespace {
	return &namespace{collection: r.collection, ns: name}
}

type Entry struct {
	Ns, Key string
	Value   any
}

func (e *Entry) K() db.Key {
	return db.K{db.S{"ns", e.Ns}, db.S{"key", e.Key}}
}

func (e *Entry) CollectionName() string { return "conf" }

var _ db.Keyer = (*Entry)(nil)

func (e Entry) String() string {
	return fmt.Sprintf("%s<%s: %v>", e.Ns, e.Key, e.Value)
}


