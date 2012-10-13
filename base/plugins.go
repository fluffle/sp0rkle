package base

import (
	"sync"
)

type Plugin interface {
	Apply(string, *Line) string
}

type pluginSet struct {
	sync.RWMutex
	set []Plugin
}

func NewPluginSet() *pluginSet {
	return &pluginSet{set: make([]Plugin, 0, 10)}
}

func (ps *pluginSet) Add(p Plugin) {
	ps.Lock()
	defer ps.Unlock()
	ps.set = append(ps.set, p)
}

func (ps *pluginSet) Apply(in string, l *Line) string {
	ps.RLock()
	defer ps.RUnlock()
	for _, p := range ps.set {
		in = p.Apply(in, l)
	}
	return in
}
