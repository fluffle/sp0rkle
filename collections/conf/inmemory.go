package conf

import (
	"sync"
)

type inMem struct {
	sync.Mutex
	ns   string
	data map[string]interface{}
}

func InMem(ns string) Namespace {
	return &inMem{ns: ns, data: make(map[string]interface{})}
}

func (ns *inMem) All() Entries {
	ns.Lock()
	defer ns.Unlock()
	e := make(Entries, 0, len(ns.data))
	for k, v := range ns.data {
		e = append(e, Entry{ns.ns, k, v})
	}
	return e
}

func (ns *inMem) String(key string, value ...string) string {
	ns.Lock()
	defer ns.Unlock()
	if len(value) > 0 {
		ns.data[key] = value[0]
		return value[0]
	}
	if val, ok := ns.data[key].(string); ok {
		return val
	}
	return ""
}

func (ns *inMem) Int(key string, value ...int) int {
	ns.Lock()
	defer ns.Unlock()
	if len(value) > 0 {
		ns.data[key] = value[0]
		return value[0]
	}
	if val, ok := ns.data[key].(int); ok {
		return val
	}
	return 0
}

func (ns *inMem) Float(key string, value ...float64) float64 {
	ns.Lock()
	defer ns.Unlock()
	if len(value) > 0 {
		ns.data[key] = value[0]
		return value[0]
	}
	if val, ok := ns.data[key].(float64); ok {
		return val
	}
	return 0
}

func (ns *inMem) Value(key string, value ...interface{}) interface{} {
	ns.Lock()
	defer ns.Unlock()
	if len(value) > 0 {
		ns.data[key] = value[0]
		return value[0]
	}
	if val, ok := ns.data[key]; ok {
		return val
	}
	return nil
}

func (ns *inMem) Delete(key string) {
	ns.Lock()
	defer ns.Unlock()
	delete(ns.data, key)
}
