package db

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/fluffle/golog/logging"
	"go.etcd.io/bbolt"
	"gopkg.in/mgo.v2/bson"
)

type rowScanner interface {
	fmt.Stringer
	scan([]byte) error
}

type allScanner struct {
	sp *slicePtr
}

func (allScanner) String() string        { return "All()" }
func (s allScanner) scan(v []byte) error { return bson.Unmarshal(v, s.sp.ponyElem()) }

type matchScanner struct {
	re    string
	rx    *regexp.Regexp
	field string
	sp    *slicePtr
}

func (s matchScanner) String() string { return fmt.Sprintf("Match(%q, /%s/)", s.field, s.re) }

func (s matchScanner) scan(v []byte) error {
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

type indexScanner struct {
	sp   *slicePtr
	vals *bbolt.Bucket
	// When scanning over indexes, we might encounter multiple pointers to the
	// same value. Returning duplicates in this case would be unhelpful.
	seen map[string]bool
}

func (indexScanner) String() string { return "All()" }

func (s indexScanner) scan(v []byte) error {
	if s.seen[string(v)] {
		return nil
	}
	s.seen[string(v)] = true
	data := s.vals.Get(v)
	if !isBson(data) {
		logging.Warn("%s: encountered dangling pointer %q", s, v)
		return nil
	}
	return bson.Unmarshal(suffix(data), s.sp.ponyElem())
}

func scanTx(b *bbolt.Bucket, scanner rowScanner) error {
	cs := []*bbolt.Cursor{b.Cursor()}
	var c *bbolt.Cursor

	for len(cs) > 0 {
		c, cs = cs[0], cs[1:]
		for k, v := c.First(); k != nil; k, v = c.Next() {
			switch {
			case v == nil:
				// Flatten the nested buckets under key.
				if nest := b.Bucket(k); nest != nil {
					cs = append(cs, nest.Cursor())
				}
			case isPointer(v):
				// To future me, if this bites me in the ass: sorry.
				// indexScanner transparently handles pointer resolution.
				if err := scanner.scan(v); err != nil {
					return fmt.Errorf("scan/unmarshal pointer: %w", err)
				}
			case isBson(v):
				if err := scanner.scan(suffix(v)); err != nil {
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
