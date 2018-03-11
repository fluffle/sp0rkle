package db

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/fluffle/sp0rkle/util"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// RSEP is the ascii record separator non-printable character.
	RSEP = '\x1e'
	// USEP is the ascii unit separator non-printable character.
	USEP = '\x1f'
)

type Database interface {
	Close()
	C(name string) Collection
}

type Collection interface {
	Get(Key, interface{}) error
	// GetPR(Key, interface{}) error ?
	Match(string, string, interface{}) error
	All(Key, interface{}) error
	Put(interface{}) error
	BatchPut(interface{}) error
	Del(interface{}) error
	Next(Key, ...uint64) (int, error)
	// Turn on debugging for this collection.
	Debug(bool)
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

// A value that is stored directly at Key in BoltDB.
// The method is not called Key because conf.Entry has
// a field named Key which references data in mongo
// but still needs to implement this interface.
// Naming is hard, but this is probably fine because
// they will most likely be returning a db.K anyway.
type Keyer interface {
	K() Key
}

// A value that is stored directly at K{{"_id", ObjectId}}
// with pointers for each Key in Indexes in BoltDB.
type Indexer interface {
	Id() bson.ObjectId
	Indexes() []Key
}

type Elem interface {
	Pair() (string, interface{})
	Bytes() []byte
	String() string
}

// String key element.
type S struct {
	Name, Value string
}

func (e S) Pair() (string, interface{}) {
	return e.Name, e.Value
}

func (e S) Bytes() []byte {
	b := bytes.NewBuffer(make([]byte, 0, len(e.Name)+len(e.Value)+1))
	b.WriteString(e.Name)
	b.WriteByte(USEP)
	b.WriteString(e.Value)
	return b.Bytes()
}

func (e S) String() string {
	return e.Name + ": " + e.Value
}

// Integer key element.
type I struct {
	Name  string
	Value int
}

func (e I) Pair() (string, interface{}) {
	return e.Name, e.Value
}

func (e I) Bytes() []byte {
	v := util.EncodeVarint(uint64(e.Value))
	b := bytes.NewBuffer(make([]byte, 0, len(e.Name)+len(v)+1))
	b.WriteString(e.Name)
	b.WriteByte(USEP)
	b.Write(v)
	return b.Bytes()
}

func (e I) String() string {
	return fmt.Sprintf("%s: %d", e.Name, e.Value)
}

type Key interface {
	String() string
	// MongoDB repr
	M() bson.M
	// BoltDB repr
	B() ([][]byte, []byte)
}

type K []Elem

// This is one-way, loses ordering.
func (k K) M() bson.M {
	m := bson.M{}
	for _, e := range k {
		n, v := e.Pair()
		m[n] = v
	}
	return m
}

// Successive key elements create nested BoltDB buckets.
// The final key element is used as the bucket key.
func (k K) B() ([][]byte, []byte) {
	if len(k) == 0 {
		return nil, nil
	}
	items := make([][]byte, 0, len(k))
	for _, e := range k {
		items = append(items, e.Bytes())
	}
	return items[:len(items)-1], items[len(items)-1]
}

func (k K) String() string {
	s := make([]string, 0, len(k))
	for _, e := range k {
		s = append(s, e.String())
	}
	return "K<" + strings.Join(s, ", ") + ">"
}
