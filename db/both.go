package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fluffle/goirc/logging"
	"github.com/fluffle/sp0rkle/util/diff"
	"gopkg.in/mgo.v2"
)

// TODO(fluffle): There is a lot of copypasta in here.

// Both implements Collection by writing to both
// and comparing reads. If Migrated() returns true
// reads return BoltDB data/errors, otherwise MongoDB.
type Both struct {
	Checker M
	MongoC  C
	BoltC   C
}

func (b *Both) Debug(on bool) {
	b.MongoC.Debug(on)
	b.BoltC.Debug(on)
}

// Having the Checker field be an M is really helpful
// but it means Both can't use Go's embedding to
// automatically delegate the Migrated method.
func (b *Both) Check() MigrationState {
	return b.Checker.Check()
}

func dupeR(vt reflect.Type, vv reflect.Value) reflect.Value {
	switch vv.Kind() {
	case reflect.Ptr:
		duped := dupeR(vv.Elem().Type(), vv.Elem())
		ptr := reflect.New(vv.Elem().Type())
		ptr.Elem().Set(duped)
		return ptr
	case reflect.Slice:
		return reflect.MakeSlice(vt, 0, vv.Cap())
	default:
		return reflect.New(vt).Elem()
	}
}

// This function rigourously tested for all of 15 minutes
// at https://play.golang.org/p/IrEWIxm_PEH ;-)
func dupe(in interface{}) interface{} {
	return dupeR(reflect.TypeOf(in), reflect.ValueOf(in)).Interface()
}

func (b *Both) compareErr(method string, mErr, bErr error) error {
	if mErr != bErr {
		logging.Warn("%s() errors differ: %v != %v", method, mErr, bErr)
	}
	if b.Check() <= MONGO_PRIMARY {
		return mErr
	}
	return bErr
}

func (b *Both) compare(method, key string, mValue, bValue interface{}, mErr, bErr error) error {
	// Mongo returns ErrNotFound, Bolt returns nil, nil.
	if errors.Is(mErr, mgo.ErrNotFound) && bErr == nil && bValue == nil {
		return nil
	}
	// If we can diff, compare by diffing, otherwise just do a DeepEqual.
	unified, err := diff.SortDiff(mValue, bValue)
	if err == diff.ErrDiff {
		logging.Debug("%s() Diff for key %s (-mongo, +bolt): %v\n%s",
			method, key, err, strings.Join(unified, "\n"))
	} else if err == diff.ErrNotDiffable && !reflect.DeepEqual(mValue, bValue) {
		logging.Warn("%s() mismatch for key %s.", method, key)
		logging.Debug("Mongo: %#v", mValue)
		logging.Debug("Bolt: %#v", bValue)
	}
	return b.compareErr(method, mErr, bErr)
}

func (b *Both) Get(key Key, value interface{}) error {
	other := dupe(value)
	switch b.Check() {
	case MONGO_ONLY:
		return b.MongoC.Get(key, value)
	case MONGO_PRIMARY:
		return b.compare("Get", key.String(), value, other,
			b.MongoC.Get(key, value), b.BoltC.Get(key, other))
	case BOLT_PRIMARY:
		return b.compare("Get", key.String(), other, value,
			b.MongoC.Get(key, other), b.BoltC.Get(key, value))
	case BOLT_ONLY:
		return b.BoltC.Get(key, value)
	}
	return ErrInvalidState
}

func (b *Both) Match(key, re string, value interface{}) error {
	other := dupe(value)
	switch b.Check() {
	case MONGO_ONLY:
		return b.MongoC.Match(key, re, value)
	case MONGO_PRIMARY:
		return b.compare("Match", key, value, other,
			b.MongoC.Match(key, re, value), b.BoltC.Match(key, re, other))
	case BOLT_PRIMARY:
		return b.compare("Match", key, other, value,
			b.MongoC.Match(key, re, other), b.BoltC.Match(key, re, value))
	case BOLT_ONLY:
		return b.BoltC.Match(key, re, value)
	}
	return ErrInvalidState
}

func (b *Both) All(key Key, value interface{}) error {
	other := dupe(value)
	switch b.Check() {
	case MONGO_ONLY:
		return b.MongoC.All(key, value)
	case MONGO_PRIMARY:
		return b.compare("All", key.String(), value, other,
			b.MongoC.All(key, value), b.BoltC.All(key, other))
	case BOLT_PRIMARY:
		return b.compare("All", key.String(), other, value,
			b.MongoC.All(key, other), b.BoltC.All(key, value))
	case BOLT_ONLY:
		return b.BoltC.All(key, value)
	}
	return ErrInvalidState
}

func (b *Both) Put(value interface{}) error {
	switch b.Check() {
	case MONGO_ONLY:
		return b.MongoC.Put(value)
	case MONGO_PRIMARY, BOLT_PRIMARY:
		return b.compareErr("Put", b.MongoC.Put(value), b.BoltC.Put(value))
	case BOLT_ONLY:
		return b.BoltC.Put(value)
	}
	return ErrInvalidState
}

func (b *Both) BatchPut(value interface{}) error {
	switch b.Check() {
	case MONGO_ONLY:
		// BatchPut is a bolt thing, fail before migration
		return fmt.Errorf("unable to BatchPut in MONGO_ONLY migration state\n\n%#v\n", value)
	case MONGO_PRIMARY, BOLT_PRIMARY, BOLT_ONLY:
		return b.BoltC.BatchPut(value)
	}
	return ErrInvalidState
}

func (b *Both) Del(value interface{}) error {
	switch b.Check() {
	case MONGO_ONLY:
		return b.MongoC.Del(value)
	case MONGO_PRIMARY, BOLT_PRIMARY:
		return b.compareErr("Del", b.MongoC.Del(value), b.BoltC.Del(value))
	case BOLT_ONLY:
		return b.BoltC.Del(value)
	}
	return ErrInvalidState
}

func (b *Both) Next(key Key, set ...int) (int, error) {
	switch b.Check() {
	case MONGO_ONLY:
		// Next is a bolt think, fail before migration
		return 0, fmt.Errorf("unable to Next(%s, %v) in MONGO_ONLY migration state", key, set)
	case MONGO_PRIMARY, BOLT_PRIMARY, BOLT_ONLY:
		return b.BoltC.Next(key, set...)
	}
	return 0, ErrInvalidState
}

func (b *Both) Mongo() *mgo.Collection {
	return b.MongoC.Mongo()
}
