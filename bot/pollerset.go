package bot

import (
	"sync"
	"time"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
)

type Poller interface {
	Poll([]*Context)
	Start()
	Stop()
	Tick() time.Duration
}

type PollerSet interface {
	Add(Poller)
	client.Handler
}

type pollerSet struct {
	sync.RWMutex
	// Pollers are started once when we first connect to a server
	// and are stopped when we disconnect from the last server.
	set map[Poller]chan struct{}
	// Each time Poll() is called, it is passed a list of Contexts
	// that represent the set of servers currently connected to.
	conns map[*client.Conn]*Context
}

func newPollerSet() *pollerSet {
	return &pollerSet{
		set:   make(map[Poller]chan struct{}),
		conns: make(map[*client.Conn]*Context),
	}
}

func (ps *pollerSet) Add(p Poller) {
	ps.Lock()
	defer ps.Unlock()
	ps.set[p] = ps.startOne(p)
	logging.Debug("Add: # conns: %d, # pollers: %d", len(ps.conns), len(ps.set))
}

// pollerSet handles both CONNECTED and DISCONNECTED events
func (ps *pollerSet) Handle(conn *client.Conn, line *client.Line) {
	ps.Lock()
	defer ps.Unlock()
	switch line.Cmd {
	case client.CONNECTED:
		ps.conns[conn] = context(conn, line)
		logging.Debug("Conn: # conns: %d, # pollers: %d", len(ps.conns), len(ps.set))
		if len(ps.conns) == 1 {
			for p := range ps.set {
				ps.set[p] = ps.startOne(p)
			}
		}
	case client.DISCONNECTED:
		delete(ps.conns, conn)
		logging.Debug("Disc: # conns: %d, # pollers: %d", len(ps.conns), len(ps.set))
		if len(ps.conns) == 0 {
			for p, quit := range ps.set {
				close(quit)
				ps.set[p] = nil
			}
		}
	}
}

func (ps *pollerSet) startOne(p Poller) chan struct{} {
	if len(ps.conns) == 0 {
		return nil
	}
	logging.Debug("Starting poller %#v at %s intervals.", p, p.Tick())
	tick := time.NewTicker(p.Tick())
	quit := make(chan struct{})
	go func() {
		p.Start()
		for {
			select {
			case <-tick.C:
				p.Poll(ps.contexts())
			case <-quit:
				tick.Stop()
				p.Stop()
				return
			}
		}
	}()
	return quit
}

func (ps *pollerSet) contexts() []*Context {
	ps.RLock()
	defer ps.RUnlock()
	ctxs := make([]*Context, 0, len(ps.conns))
	for _, ctx := range ps.conns {
		ctxs = append(ctxs, ctx)
	}
	return ctxs
}
