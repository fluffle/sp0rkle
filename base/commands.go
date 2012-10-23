package base

import (
	"github.com/fluffle/golog/logging"
	"strings"
	"sync"
)

type Handler interface {
	Execute(*Line)
}

type HandlerFunc func(*Line)

func (hf HandlerFunc) Execute(line *Line) {
	hf(line)
}

type Command interface {
	Handler
	Help() string
}

type command struct {
	fn HandlerFunc
	help string
}

func NewCommand(fn HandlerFunc, help string) Command {
	return &command{fn, help}
}

func (c *command) Execute(line *Line) {
	c.fn(line)
}

func (c *command) Help() string {
	return c.help
}

// CommandSet mostly gratuitously stolen from net/http ;-)
type commandSet struct {
	sync.RWMutex
	set map[string]Command
}

func NewCommandSet() *commandSet {
	return &commandSet{set: make(map[string]Command)}
}

func (cs *commandSet) Add(cmd Command, prefix string) {
	if cmd == nil || prefix == "" {
		logging.Error("Can't handle prefix '%s' with command.", prefix)
		return
	}
	cs.Lock()
	defer cs.Unlock()
	if _, ok := cs.set[prefix]; ok {
		logging.Error("Prefix '%s' already registered.", prefix)
		return
	}
	cs.set[prefix] = cmd
}

func (cs *commandSet) Match(txt string) (final Command, prefixlen int) {
	cs.RLock()
	defer cs.RUnlock()

	for prefix, cmd := range cs.set {
		if !strings.HasPrefix(txt, prefix) {
			continue
		}
		if final == nil || len(prefix) > prefixlen {
			prefixlen = len(prefix)
			final = cmd
		}
	}
	return
}
