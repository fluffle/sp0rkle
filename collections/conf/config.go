package conf

import (
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/db"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"sync"
)

const COLLECTION string = "conf"

type namespace string

type Namespace interface {
	All() []Entry
	String(key string, value ...string) string
	Int(key string, value ...int) int
	Float(key string, value ...float64) float64
	Value(key string, value ...interface{}) interface{}
	Delete(key string)
}

type Entry struct {
	Ns	    namespace
	Key     string
	Value   interface{}
}

func (e Entry) String() string {
	return fmt.Sprintf("%s<%s: %v>", e.Ns, e.Key, e.Value)
}

func (e Entry) Id() bson.M {
	return bson.M{"ns": e.Ns, "key": e.Key}
}

var conf *mgo.Collection
var lock sync.Mutex

func Ns(ns string) namespace {
	lock.Lock()
	defer lock.Unlock()
	if conf == nil {
		conf = db.Init().C(COLLECTION)
		err := conf.EnsureIndex(mgo.Index{Key: []string{"ns", "key"}, Unique: true})
		if err != nil {
			logging.Error("Couldn't create index on sp0rkle.conf: %s", err)
		}
	}
	return namespace(ns)
}

func (ns namespace) key(key string) bson.M {
	return bson.M{"ns": ns, "key": key}
}

func (ns namespace) set(key string, value interface{}) {
	e := Entry{Ns: ns, Key: key, Value: value}
	if _, err := conf.Upsert(e.Id(), e); err != nil {
		logging.Error("Couldn't upsert config entry '%s': %s", e, err)
	}
}

func (ns namespace) get(key string) interface{} {
	var e Entry
	if err := conf.Find(ns.key(key)).One(&e); err == nil {
		return e.Value
	}
	return nil
}

func (ns namespace) All() []Entry {
	var e []Entry
	if err := conf.Find(bson.M{"ns": ns}).All(&e); err == nil {
		return e
	}
	return []Entry{}
}

func (ns namespace) String(key string, value ...string) string {
	if len(value) > 0 {
		ns.set(key, value[0])
		return value[0]
	}
	if val, ok := ns.get(key).(string); ok {
		return val
	}
	return ""
}

func (ns namespace) Int(key string, value ...int) int {
	if len(value) > 0 {
		ns.set(key, value[0])
		return value[0]
	}
	if val, ok := ns.get(key).(int); ok {
		return val
	}
	return 0
}

func (ns namespace) Float(key string, value ...float64) float64 {
	if len(value) > 0 {
		ns.set(key, value[0])
		return value[0]
	}
	if val, ok := ns.get(key).(float64); ok {
		return val
	}
	return 0
}
	
func (ns namespace) Value(key string, value ...interface{}) interface{} {
	if len(value) > 0 {
		ns.set(key, value[0])
		return value[0]
	}
	return ns.get(key)
}

func (ns namespace) Delete(key string) {
	conf.Remove(ns.key(key))
}
