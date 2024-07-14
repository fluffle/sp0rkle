package db

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// RSEP is the ascii record separator non-printable character.
	RSEP = '\x1e'
	// USEP is the ascii unit separator non-printable character.
	USEP = '\x1f'
	// TRUE and FALSE are used in constructing keys from booleans.
	TRUE  = '\xff'
	FALSE = '\x00'
)

type Database interface {
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
	Next(Key, ...int) (int, error)
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
	Value uint64
}

func (e I) Pair() (string, interface{}) {
	return e.Name, e.Value
}

func (e I) Bytes() []byte {
	v := make([]byte, 8)
	// Big endian is lexographically sortable, handy for indexes.
	binary.BigEndian.PutUint64(v, e.Value)
	b := bytes.NewBuffer(make([]byte, 0, len(e.Name)+9))
	b.WriteString(e.Name)
	b.WriteByte(USEP)
	b.Write(v)
	return b.Bytes()
}

func (e I) String() string {
	return fmt.Sprintf("%s: %d", e.Name, e.Value)
}

// Boolean key element.
type T struct {
	Name  string
	Value bool
}

func (e T) Pair() (string, interface{}) {
	return e.Name, e.Value
}

func (e T) Bytes() []byte {
	b := bytes.NewBuffer(make([]byte, 0, len(e.Name)+2))
	b.WriteString(e.Name)
	b.WriteByte(USEP)
	if e.Value {
		b.WriteByte(TRUE)
	} else {
		b.WriteByte(FALSE)
	}
	return b.Bytes()
}

func (e T) String() string {
	return fmt.Sprintf("%s: %t", e.Name, e.Value)
}

// ObjectId key element, because aaargh of course casting it to a string
// fucks it up, even though bson.ObjectId is just a string type.
type ID struct {
	Value bson.ObjectId
}

func (e ID) Pair() (string, interface{}) {
	return "_id", e.Value
}

func (e ID) Bytes() []byte {
	b := bytes.NewBuffer(make([]byte, 0, len(e.Value)+4))
	b.WriteString("_id")
	b.WriteByte(USEP)
	b.WriteString(string(e.Value))
	return b.Bytes()
}

func (e ID) String() string {
	return fmt.Sprintf("_id: %s", e.Value)
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
	items := make([][]byte, len(k))
	for i, e := range k {
		items[i] = e.Bytes()
	}
	return items[:len(items)-1], items[len(items)-1]
}

func (k K) String() string {
	s := make([]string, len(k))
	for i, e := range k {
		s[i] = e.String()
	}
	return "K<" + strings.Join(s, ", ") + ">"
}
