package conf

import (
	"github.com/fluffle/goirc/logging"
	"github.com/fluffle/sp0rkle/db"
)

type Namespace interface {
	String(key string, value ...string) string
	Int(key string, value ...int) int
	Float(key string, value ...float64) float64
	Delete(key string)
}

type namespace struct {
	collection db.KeyedCollection[*Entry]
	ns         string
}

func (ns *namespace) K() db.Key {
	return db.K{db.S{"ns", ns.ns}}
}

func (ns *namespace) CollectionName() string { return "conf" }

var _ db.Keyer = (*namespace)(nil)

func (ns *namespace) set(key string, value any) {
	e := &Entry{Ns: ns.ns, Key: key, Value: value}
	if err := ns.collection.Put(e); err != nil {
		logging.Error("Couldn't set config entry %q: %v", e, err)
	}
}

func (ns *namespace) get(key string) any {
	e := &Entry{Ns: ns.ns, Key: key}
	e, err := ns.collection.Get(e.K())
	if err != nil {
		logging.Error("Couldn't get config entry %q: %v", e, err)
		return nil
	}
	if e == nil {
		// key not found
		return nil
	}
	return e.Value
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

func (ns *namespace) Delete(key string) {
	ns.collection.Del(&Entry{Ns: ns.ns, Key: key})
}
