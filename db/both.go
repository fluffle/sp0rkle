package db

import (
	"reflect"
	"sort"
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
func (b *Both) Migrated() bool {
	return b.Checker.Migrated()
}

// This function rigourously tested for all of 5 minutes
// at http://play.golang.org/p/IwZQ17Bpjt ;-)
func dupe(in interface{}) interface{} {
	vv := reflect.ValueOf(in)
	vt := reflect.TypeOf(in)
	if vv.Kind() == reflect.Ptr {
		vt = vv.Elem().Type()
		return reflect.New(vt).Interface()
	}
	return reflect.New(vt).Elem().Interface()
}

func (b *Both) Get(key Key, value interface{}) error {
	var mErr, bErr error
	other := dupe(value)
	if b.Migrated() {
		mErr = b.MongoC.Get(key, other)
		bErr = b.BoltC.Get(key, value)
	} else {
		mErr = b.MongoC.Get(key, value)
		bErr = b.BoltC.Get(key, other)
	}
	if mErr != bErr {
		logging.Warn("Get() errors differ: %v != %v", mErr, bErr)
	}
	if !reflect.DeepEqual(value, other) {
		logging.Warn("Get() mismatch for %s.", key)
		if b.Migrated() {
			logging.Debug("Mongo: %#v", other)
			logging.Debug("Bolt: %#v", value)
		} else {
			logging.Debug("Mongo: %#v", value)
			logging.Debug("Bolt: %#v", other)
		}
	}
	if b.Migrated() {
		return bErr
	}
	return mErr
}

func (b *Both) Match(key, re string, value interface{}) error {
	var mErr, bErr error
	other := dupe(value)
	if b.Migrated() {
		mErr = b.MongoC.Match(key, re, other)
		bErr = b.BoltC.Match(key, re, value)
	} else {
		mErr = b.MongoC.Match(key, re, value)
		bErr = b.BoltC.Match(key, re, other)
	}
	if mErr != bErr {
		logging.Warn("Match() errors differ: %v != %v", mErr, bErr)
	}
	vdiff, vok := value.(Diffable)
	odiff, ook := other.(Diffable)
	if ook && vok {
		vstr := vdiff.Strings()
		ostr := odiff.Strings()
		sort.Strings(vstr)
		sort.Strings(ostr)
		unified, err := diff.Unified(vstr, ostr)
		if err != nil {
			logging.Debug("Match() Diff: %v\n%s", err, strings.Join(unified, "\n"))
		}
	} else if !reflect.DeepEqual(value, other) {
		logging.Warn("Match() mismatch for %s.", key)
		if b.Migrated() {
			logging.Debug("Mongo: %#v", other)
			logging.Debug("Bolt: %#v", value)
		} else {
			logging.Debug("Mongo: %#v", value)
			logging.Debug("Bolt: %#v", other)
		}
	}
	if b.Migrated() {
		return bErr
	}
	return mErr
}

func (b *Both) All(key Key, value interface{}) error {
	var mErr, bErr error
	other := dupe(value)
	if b.Migrated() {
		mErr = b.MongoC.All(key, other)
		bErr = b.BoltC.All(key, value)
	} else {
		mErr = b.MongoC.All(key, value)
		bErr = b.BoltC.All(key, other)
	}
	if mErr != bErr {
		logging.Warn("All() errors differ: %v != %v", mErr, bErr)
	}
	if !reflect.DeepEqual(value, other) {
		logging.Warn("All() mismatch for %s.", key)
		if b.Migrated() {
			logging.Debug("Mongo: %#v", other)
			logging.Debug("Bolt: %#v", value)
		} else {
			logging.Debug("Mongo: %#v", value)
			logging.Debug("Bolt: %#v", other)
		}
	}
	if b.Migrated() {
		return bErr
	}
	return mErr
}

func (b *Both) Put(value interface{}) error {
	mErr := b.MongoC.Put(value)
	bErr := b.BoltC.Put(value)
	if mErr != bErr {
		logging.Warn("Put() errors differ: %v != %v", mErr, bErr)
	}
	if b.Migrated() {
		return bErr
	}
	return mErr
}

func (b *Both) BatchPut(value interface{}) error {
	return b.BoltC.BatchPut(value)
}

func (b *Both) Del(value interface{}) error {
	mErr := b.MongoC.Del(value)
	bErr := b.BoltC.Del(value)
	if mErr != bErr {
		logging.Warn("Del() errors differ: %v != %v", mErr, bErr)
	}
	if b.Migrated() {
		return bErr
	}
	return mErr
}

func (b *Both) Next(k Key, set ...uint64) (int, error) {
	return b.BoltC.Next(k, set...)
}

func (b *Both) Mongo() *mgo.Collection {
	return b.MongoC.Mongo()
}
