package db

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/fluffle/golog/logging"
	"go.etcd.io/bbolt"
	"gopkg.in/mgo.v2/bson"
)

type rowScanner interface {
	fmt.Stringer
	scan([]byte, []byte) error
}

type allScanner struct {
	sp *slicePtr
}

func (allScanner) String() string           { return "All()" }
func (s allScanner) scan(_, v []byte) error { return bson.Unmarshal(v, s.sp.ponyElem()) }

type matchScanner struct {
	re    string
	rx    *regexp.Regexp
	field string
	sp    *slicePtr
}

func (s matchScanner) String() string { return fmt.Sprintf("Match(%q, /%s/)", s.field, s.re) }

func (s matchScanner) scan(_, v []byte) error {
	ev := s.sp.newElem()
	if err := bson.Unmarshal(v, ev.Addr().Interface()); err != nil {
		return err
	}
	cev := ev
	for cev.Kind() == reflect.Ptr {
		cev = cev.Elem()
	}
	if s.rx.MatchString(cev.FieldByName(s.field).String()) {
		s.sp.appendElem(ev)
	}
	return nil
}

type badPointerReason string

const (
	nonBsonData badPointerReason = "pointer points to non-BSON data"
	obsoleteKey badPointerReason = "key is obsolete for referenced data"
)

type badPointerErr struct {
	k, v   []byte
	data   any
	reason badPointerReason
}

func (err *badPointerErr) Error() string {
	return fmt.Sprintf("key %q -> pointer %q -> %#v: %s",
		err.k, err.v, err.data, err.reason)
}

func (err *badPointerErr) Is(other error) bool {
	_, ok := other.(*badPointerErr)
	return ok
}

type indexScanner struct {
	sp   *slicePtr
	vals *bbolt.Bucket
	// When scanning over indexes, we might encounter multiple pointers to the
	// same value. Returning duplicates in this case would be unhelpful.
	seen map[string]bool
}

func (indexScanner) String() string { return "All()" }

func (s indexScanner) scan(k, v []byte) error {
	if s.seen[string(v)] {
		return nil
	}
	s.seen[string(v)] = true
	data := s.vals.Get(v)
	if !isBson(data) {
		return &badPointerErr{k: k, v: v, data: data, reason: nonBsonData}
	}
	return bson.Unmarshal(suffix(data), s.sp.ponyElem())
}

type fsckIndex struct {
	et   reflect.Type
	vals *bbolt.Bucket
}

func (fsckIndex) String() string { return "fsckIndex()" }

func (s fsckIndex) scan(k, v []byte) error {
	data := s.vals.Get(v)
	if !isBson(data) {
		return &badPointerErr{k: k, v: v, data: data, reason: nonBsonData}
	}
	elem := reflect.New(s.et).Interface()
	if err := bson.Unmarshal(suffix(data), elem); err != nil {
		return err
	}
	idx, ok := elem.(Indexer)
	if !ok {
		logging.Error("fsckIndex: unmarshaled value %#v for key %q pointer %q is not an indexer", elem, k, v)
		// keep scanning
		return nil
	}
	found := false
	for _, key := range idx.Indexes() {
		_, last := key.B()
		found = found || bytes.Equal(last, k)
	}
	if !found {
		return &badPointerErr{k: k, v: v, data: elem, reason: obsoleteKey}
	}
	return nil
}

type fsckValues struct {
	et   reflect.Type
	idxs *bbolt.Bucket
}

func (fsckValues) String() string { return "fsckValues()" }

func (s fsckValues) scan(k, v []byte) error {
	elem := reflect.New(s.et).Interface()
	if err := bson.Unmarshal(v, elem); err != nil {
		return err
	}
	idx, ok := elem.(Indexer)
	if !ok {
		logging.Error("fsckValues: unmarshaled value %#v for key %q pointer %q is not an indexer", elem, k, v)
		// keep scanning
		return nil
	}
	ptr := toPointer(idx)
	if !bytes.Equal(ptr, k) {
		logging.Error("fsckValues: key %q derived from value %#v does not match actual key %q", ptr, idx, k)
		// todo: fix?
	}
INDEXES:
	for _, key := range idx.Indexes() {
		elems, last := key.B()
		b := s.idxs
		for _, elem := range elems {
			nest := b.Bucket(elem)
			if nest != nil {
				b = nest
				continue
			}
			logging.Error("fsckValues: index bucket %q for value %#v does not exist, creating", elem, idx)
			var err error
			nest, err = b.CreateBucket(elem)
			if err != nil {
				logging.Error("fsckValues: creating index bucket %q: %v", elem, err)
				// keep fixing
				continue INDEXES
			}
			b = nest
		}
		// b now contains the correct bucket for the final index element
		idxptr := b.Get(last)
		if idxptr == nil || !bytes.Equal(idxptr, ptr) {
			logging.Error("fsckValues: final key element %q for value %#v incorrect (%q != %q), fixing", last, idx, idxptr, ptr)
			if err := b.Put(last, ptr); err != nil {
				logging.Error("fsckValues: writing index pointer %q: %v", last, err)
			}
		}
	}
	return nil
}

func scanTx(b *bbolt.Bucket, scanner rowScanner) error {
	cs := []*bbolt.Cursor{b.Cursor()}
	var c *bbolt.Cursor
	writable := b.Writable()

	for len(cs) > 0 {
		c, cs = cs[0], cs[1:]
		for k, v := c.First(); k != nil; k, v = c.Next() {
			switch {
			case v == nil:
				// Flatten the nested buckets under key.
				if nest := c.Bucket().Bucket(k); nest != nil {
					cs = append(cs, nest.Cursor())
				} else {
					logging.Error("nested bucket %q returned nil", k)
				}
			case isPointer(v):
				// To future me, if this bites me in the ass: sorry.
				// indexScanner transparently handles pointer resolution.
				err := scanner.scan(k, v)
				if err == nil {
					continue
				}
				if !errors.Is(err, &badPointerErr{}) {
					return fmt.Errorf("scan/unmarshal pointer: %w", err)
				}
				logging.Debug("encountered bad pointer: %v", err)
				if !writable {
					continue
				}
				logging.Warn("deleting key %q", k)
				if delErr := c.Delete(); delErr != nil {
					return fmt.Errorf("delete key %q: %w", k, delErr)
				}
			case isBson(v):
				if err := scanner.scan(k, suffix(v)); err != nil {
					return fmt.Errorf("scan/unmarshal value: %w", err)
				}
			default:
				// Reasonably sure we shouldn't hit this condition.
				logging.Warn("%s: unexpected data k=%q v=%q", scanner, k, v)
			}
		}
	}
	return nil
}
