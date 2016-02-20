package db

import (
	"strings"
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// RSEP is the ascii record separator non-printable character.
	RSEP = "\x1e"
	// USEP is the ascii unit separator non-printable character.
	USEP = "\x1f"
)

type Database interface {
	Init(db string) error
	Close()
	C(name string) Collection
}

type Collection interface {
	Get(Key, interface{}) error
	All(Key, interface{}) error
	Put(Key, interface{}) error
	Del(Key) error
	// So we don't have to do everything at once.
	Mongo() *mgo.Collection
}

type C struct {
	Collection
	sync.Once
}

func (c *C) Init(db Database, name string, f func(Collection)) {
	c.Do(func() {
		c.Collection = db.C(name)
		if f != nil {
			f(c)
		}
	})
}

type Key interface {
	M() bson.M
	B() []byte
}

// Basically bson.D but only string->string.
type Elem struct {
	Name, Value string
}
type K []Elem

// This is one-way, loses ordering.
func (k K) M() bson.M {
	m := bson.M{}
	for _, v := range k {
		m[v.Name] = v.Value
	}
	return m
}

// Ordered version of the above, reversible.
func (k K) D() bson.D {
	d := make(bson.D, 0, len(k))
	for _, v := range k {
		d = append(d, bson.DocElem{v.Name, v.Value})
	}
	return d
}

// This is reversible and suitable for a BoltDB key.
func (k K) B() []byte {
	items := make([]string, 0, len(k))
	for _, v := range k {
		items = append(items, v.Name+USEP+string(v.Value))
	}
	return []byte(strings.Join(items, RSEP))
}
