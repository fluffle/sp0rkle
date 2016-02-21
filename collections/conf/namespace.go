package conf

import (
	boltdb "github.com/boltdb/bolt"
	"github.com/fluffle/goirc/logging"
	"github.com/fluffle/sp0rkle/db"
	"gopkg.in/mgo.v2"
)

type Namespace interface {
	All() []Entry
	String(key string, value ...string) string
	Int(key string, value ...int) int
	Float(key string, value ...float64) float64
	Value(key string, value ...interface{}) interface{}
	Delete(key string)
}

type namespace struct {
	db.Collection
	ns string
}

func (ns *namespace) K(key ...string) db.K {
	if len(key) > 0 {
		return db.K{{"ns", ns.ns}, {"key", key[0]}}
	}
	return db.K{{"ns", ns.ns}}
}

func (ns *namespace) set(key string, value interface{}) {
	e := Entry{Ns: ns.ns, Key: key, Value: value}
	if err := ns.Put(ns.K(key), &e); err != nil {
		logging.Error("Couldn't set config entry %q: %v", e, err)
	}
}

func (ns *namespace) get(key string) interface{} {
	var e Entry
	if err := ns.Get(ns.K(key), &e); err != nil && err != mgo.ErrNotFound && err != boltdb.ErrTxNotWritable {
		logging.Error("Couldn't get config entry for ns=%q key=%q: %v", ns.ns, key, err)
		return nil
	}
	return e.Value
}

func (ns *namespace) All() []Entry {
	var e []Entry
	if err := ns.Collection.All(ns.K(), &e); err == nil {
		return e
	}
	return []Entry{}
}

func (ns *namespace) String(key string, value ...string) string {
	if len(value) > 0 {
		ns.set(key, value[0])
		return value[0]
	}
	if val, ok := ns.get(key).(string); ok {
		return val
	}
	return ""
}

func (ns *namespace) Int(key string, value ...int) int {
	if len(value) > 0 {
		ns.set(key, value[0])
		return value[0]
	}
	if val, ok := ns.get(key).(int); ok {
		return val
	}
	return 0
}

func (ns *namespace) Float(key string, value ...float64) float64 {
	if len(value) > 0 {
		ns.set(key, value[0])
		return value[0]
	}
	if val, ok := ns.get(key).(float64); ok {
		return val
	}
	return 0
}

func (ns *namespace) Value(key string, value ...interface{}) interface{} {
	if len(value) > 0 {
		ns.set(key, value[0])
		return value[0]
	}
	return ns.get(key)
}

func (ns *namespace) Delete(key string) {
	ns.Del(ns.K(key))
}
