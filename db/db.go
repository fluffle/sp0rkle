package db

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/fluffle/sp0rkle/util/datetime"
	"github.com/fluffle/sp0rkle/util/bson"
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

// collection[T] is the internal generic collection interface.
// External users should only reference KeyedCollection[T] or IndexedCollection[T].
type collection[T any] interface {
	Get(key Key) (T, error)
	All(key Key) ([]T, error)
	Match(pred func(T) bool) ([]T, error)
	Put(value T) error
	BatchPut(values []T) error
	Del(value T) error
	Next(key Key, set ...int) (int, error)
	Debug(show bool)
}

// KeyedCollection[T] and IndexedCollection[T] are type aliases that express
// the constraint on T. They are the only public entry points for external code.
type KeyedCollection[T Keyer] = collection[T]
type IndexedCollection[T Indexer] = collection[T]

type ValueValidator[T any] interface {
	// ValidateValues receives all stored values and returns the subset
	// that should be deleted. Called on a nil receiver, must not dereference it.
	ValidateValues(values []T) ([]T, error)
}

type Database interface {
	Live() bool
	Close() error
}

type Elem interface {
	Pair() (string, any)
	Bytes() []byte
	String() string
}

// String key element.
type S struct {
	Name, Value string
}

func (e S) Pair() (string, any) {
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

func (e I) Pair() (string, any) {
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

func (e T) Pair() (string, any) {
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

func (e ID) Pair() (string, any) {
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

// Timestamp key element.
// BSON stores timestamps with millisecond precision, so this element
// automatically truncates time.Time to milliseconds before encoding,
// ensuring index keys remain consistent after a BSON roundtrip.
type TS struct {
	Name  string
	Value time.Time
}

func (e TS) Pair() (string, any) {
	return e.Name, uint64(e.Value.UnixMilli())
}

func (e TS) Bytes() []byte {
	return I{Name: e.Name, Value: uint64(e.Value.UnixMilli())}.Bytes()
}

func (e TS) String() string {
	return fmt.Sprintf("%s: %s", e.Name, datetime.Format(e.Value))
}

type Key interface {
	String() string
	// BoltDB repr
	B() ([][]byte, []byte)
	// Valid checks that no element name contains USEP.
	Valid() error
}

type K []Elem

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

// Valid checks that no element name contains USEP.
// This prevents key collisions where different name/value pairs produce
// identical byte keys (e.g. S{"a", "\x1fb"} collides with S{"a\x1f", "b"}).
// Values are not checked, because a USEP-free name guarantees the first
// USEP byte is always the separator, eliminating ambiguity, and the binary
// serializations of integers and object IDs may contain USEP.
func (k K) Valid() error {
	for i, e := range k {
		name, _ := e.Pair()
		if strings.IndexByte(name, USEP) >= 0 {
			return fmt.Errorf("db.K[%d] (%T) name contains USEP: %q", i, e, name)
		}
	}
	return nil
}
