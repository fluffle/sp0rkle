package db

import (
	"reflect"

	"github.com/fluffle/goirc/logging"
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

func (b *Both) Next(k Key) (int, error) {
	return b.BoltC.Next(k)
}

func (b *Both) Mongo() *mgo.Collection {
	return b.MongoC.Mongo()
}
