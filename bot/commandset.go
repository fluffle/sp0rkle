package bot

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/util"
	"strings"
	"sync"
)

type HandlerFunc func(*Context)

func (hf HandlerFunc) Handle(conn *client.Conn, line *client.Line) {
	if ctx := context(conn, line); ctx != nil {
		hf(ctx)
	}
}

type Runner interface {
	Run(context *Context)
	Help() string
}

type command struct {
	fn   HandlerFunc
	help string
}

func (c *command) Run(ctx *Context) {
	c.fn(ctx)
}

func (c *command) Help() string {
	return c.help
}

type CommandSet interface {
	client.Handler
	Add(command Runner, prefix string)
}

type commandSet struct {
	sync.RWMutex
	set map[string]Runner
}

func newCommandSet() *commandSet {
	cs := &commandSet{set: make(map[string]Runner)}
	// commandSet implements Runner to provide help for itself.
	cs.Add(cs, "help")
	return cs
}

func (cs *commandSet) Add(r Runner, prefix string) {
	if r == nil || prefix == "" {
		logging.Error("Prefix or runner empty when adding command.", prefix)
		return
	}
	cs.Lock()
	defer cs.Unlock()
	if _, ok := cs.set[prefix]; ok {
		logging.Error("Prefix '%s' already registered.", prefix)
		return
	}
	cs.set[prefix] = r
}

// commandSet.match() mostly gratuitously stolen from net/http ;-)
func (cs *commandSet) match(txt string) (final Runner, prefixlen int) {
	cs.RLock()
	defer cs.RUnlock()

	for prefix, r := range cs.set {
		if !strings.HasPrefix(txt, prefix) {
			continue
		}
		if final == nil || len(prefix) > prefixlen {
			prefixlen = len(prefix)
			final = r
		}
	}
	return
}

// Implement client.Handler so commandSet can Handle things directly.
func (cs *commandSet) Handle(conn *client.Conn, line *client.Line) {
	// This is a dirty hack to treat factoid additions as a special
	// case, since they may begin with command string prefixes.
	ctx := context(conn, line)
	if ctx == nil || util.IsFactoidAddition(line.Text()) {
		return
	}
	if r, ln := cs.match(ctx.Text()); ctx.Addressed && r != nil {
		// Cut command off, trim and compress spaces.
		ctx.Args[1] = strings.Join(strings.Fields(ctx.Args[1][ln:]), " ")
		r.Run(ctx)
	}
}

func (cs *commandSet) Run(ctx *Context) {
	if r, _ := cs.match(ctx.Text()); r != nil {
		ctx.ReplyN("%s", r.Help())
	} else if len(ctx.Text()) == 0 {
		ctx.ReplyN("https://github.com/fluffle/sp0rkle/wiki " +
			"-- pull requests welcome ;-)")
	} else {
		ctx.ReplyN("Unrecognised command '%s'.", ctx.Text())
	}
}

func (cs *commandSet) Help() string {
	return "If you have to ask, you're beyond help."
}
