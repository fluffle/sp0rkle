package regtest

import (
	"fmt"
	"time"

	"github.com/fluffle/goirc/client"
)

// Harness manages the harness IRC client used to drive tests.
type Harness struct {
	*client.Conn
	sp0rkle *Process
	ircd    *Process
	BotNick string
	Channel string
	tempDir string
}

// Command is a helper function to send explicit commands to the test bot.
func (h *Harness) Command(msg string) {
	h.Privmsg(h.Channel, h.BotNick + ": " + msg)
}

// CommandAndExpect is a Wrapper around SendAndExpect that calls Command first.
func (h *Harness) CommandAndExpect(msg string, p Pattern) (*client.Line, error) {
	return h.SendAndExpect(h.BotNick + ": " + msg, p)
}

// SendAndExpect sends a PRIVMSG to the test channel and waits for a response
// matching the pattern. Returns the matched Line or an error on timeout.
func (h *Harness) SendAndExpect(msg string, p Pattern) (*client.Line, error) {
	if h.Conn == nil {
		return nil, fmt.Errorf("send and expect: harness Conn is nil")
	}
	if h.Channel == "" {
		return nil, fmt.Errorf("send and expect: channel is empty")
	}
	h.Privmsg(h.Channel, msg)
	return h.Expect(p)
}

// ExpectFunc forwards a PatternFunc to Expect.
func (h *Harness) ExpectFunc(pf PatternFunc) (*client.Line, error) {
	return h.Expect(pf)
}

func (h *Harness) Expect(p Pattern) (*client.Line, error) {
	return h.ExpectEvent(client.PRIVMSG, p)
}

// ExpectEvent waits for the next IRC line matching the pattern from any source.
// Returns the matched Line or an error on timeout.
func (h *Harness) ExpectEvent(event string, p Pattern) (*client.Line, error) {
	if h.Conn == nil {
		return nil, fmt.Errorf("expect: harness Conn is nil")
	}
	if p == nil {
		return nil, fmt.Errorf("expect: pattern is nil")
	}

	resultCh := make(chan *client.Line, 1)

	var remover client.Remover
	remover = h.HandleFunc(event, func(conn *client.Conn, line *client.Line) {
		if p.Match(line) {
			remover.Remove()
			select {
			case resultCh <- line:
			default:
			}
		}
	})
	defer remover.Remove()

	// 2 seconds should be more than enough; making this configurable makes the
	// API quite unwieldy and this is an upper limit for failing tests anyway.
	timeout := 2*time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case line := <-resultCh:
		return line, nil
	case <-timer.C:
		return nil, fmt.Errorf("expect: timeout after %s", timeout)
	}
}
