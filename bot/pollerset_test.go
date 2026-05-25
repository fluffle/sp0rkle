package bot

import (
	"testing"
	"time"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
)

type dummyPoller struct{}

func (d *dummyPoller) Poll(ctxs []*Context) {}
func (d *dummyPoller) Start()                {}
func (d *dummyPoller) Stop()                 {}
func (d *dummyPoller) Tick() time.Duration   { return 10 * time.Second }

func TestNilChannelPanic(t *testing.T) {
	logging.InitFromFlags()
	ps := newPollerSet(newRewriteSet())
	ps.Add(&dummyPoller{}) // Add poller before any connections - startOne returns nil
	ps.Handle(nil, &client.Line{Cmd: client.DISCONNECTED}) // Should not panic
}
