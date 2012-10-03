package bot

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/sp0rkle/base"
	"strings"
	"sync"
)

// Mostly gratuitously stolen from net/http ;-)
type cmdFn func(*Sp0rkle, *base.Line)

type CommandFunc struct {
	fn cmdFn
	help string
}

func (cf CommandFunc) Execute(bot *Sp0rkle, line *base.Line) {
	cf.fn(bot, line)
}

func (cf CommandFunc) Help() string {
	return cf.help
}

type Command interface {
	Execute(*Sp0rkle, *base.Line)
	Help() string
}

type CommandSet struct {
	sync.RWMutex
	set map[string]Command
}

func NewCommandSet() *CommandSet {
	return &CommandSet{set: make(map[string]Command)}
}

var commands = NewCommandSet()

func Cmd(cmd Command, prefix string) {
	if cmd == nil || prefix == "" {
		logging.Error("Can't handle prefix '%s' with supplied order.", prefix)
		return
	}
	commands.Lock()
	defer commands.Unlock()
	if _, ok := commands.set[prefix]; ok {
		logging.Error("Prefix '%s' already registered.", prefix)
		return
	}
	commands.set[prefix] = cmd
}

func CmdFunc(fn cmdFn, prefix, help string) {
	Cmd(&CommandFunc{fn, help}, prefix)
}

func commandMatch(txt string) Command {
	commands.RLock()
	defer commands.RUnlock()

	var final Command
	prefixlen := 0
	for prefix, cmd := range commands.set {
		if !strings.HasPrefix(txt, prefix) {
			continue
		}
		if final == nil || len(prefix) > prefixlen {
			prefixlen = len(prefix)
			final = cmd
		}
	}
	return final
}
