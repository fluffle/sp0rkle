package regtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/fluffle/goirc/client"
)

// A Waiter waits for a line from the server.
type Waiter struct {
	// Can be modified from the default of 2s if desired.
	Timeout time.Duration
	remover client.Remover
	seen    chan *client.Line
	result  chan *client.Line
}

func (w *Waiter) Wait() (*client.Line, error) {
	defer w.remover.Remove()
	timer := time.NewTimer(w.Timeout)
	defer timer.Stop()

	lastSeen := ""
	for {
		select {
		case line := <-w.seen:
			if line != nil {
				lastSeen = fmt.Sprintf("; last line seen:\n%s", line.Raw)
			}
		case line := <-w.result:
			return line, nil
		case <-timer.C:
			return nil, fmt.Errorf("expect: timeout after %s%s", w.Timeout, lastSeen)
		}
	}
}

func (w *Waiter) MustWait(t *testing.T) {
	if _, err := w.Wait(); err != nil {
		t.Fatalf("MustWait: %v", err)
	}
}

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
	w := h.Expect(p)
	h.Privmsg(h.Channel, msg)
	return w.Wait()
}

// ExpectFunc forwards a PatternFunc to Expect.
func (h *Harness) ExpectFunc(pf PatternFunc) *Waiter {
	return h.Expect(pf)
}

func (h *Harness) Expect(p Pattern) *Waiter {
	return h.ExpectEvent(client.PRIVMSG, p)
}

// ExpectEvent creates a handler that matches lines for the given event
// against the provided pattern. It returns a Waiter that can be used
// to wait for a matching line or a timeout.
// These steps are separated so that handlers can be added _before_ any
// communication with the server that may cause a response line is sent.
// See SendAndExpect above for an example.
func (h *Harness) ExpectEvent(event string, p Pattern) *Waiter {
	if h.Conn == nil {
		panic("expect: harness Conn is nil")
	}
	if p == nil {
		panic("expect: pattern is nil")
	}

	seenCh := make(chan *client.Line, 1)
	resultCh := make(chan *client.Line, 1)

	var remover client.Remover
	remover = h.HandleFunc(event, func(conn *client.Conn, line *client.Line) {
		select {
		case seenCh <- line:
		default:
		}
		if p.Match(line) {
			remover.Remove()
			select {
			case resultCh <- line:
			default:
			}
			close(seenCh)
			close(resultCh)
		}
	})

	return &Waiter{
		Timeout: 2*time.Second,
		remover: remover,
		result:  resultCh,
		seen:    seenCh,
	}
}

func (h *Harness) MustRenick(t *testing.T, to string) string {
	from := h.Me().Nick
	w := h.ExpectEvent(client.NICK, cmdPattern{who: from, what: to})
	h.Nick(to)
	if _, err := w.Wait(); err != nil {
		t.Fatalf("expected nick change, got: %v", err)
	}
	return from
}

func (h *Harness) MustPart(t *testing.T) {
	w := h.ExpectEvent(client.PART, cmdPattern{who: h.Me().Nick, what: h.Channel})
	h.Part(h.Channel)
	if _, err := w.Wait(); err != nil {
		t.Fatalf("expected PART, got: %v", err)
	}
}

func (h *Harness) MustJoin(t *testing.T) {
	w := h.ExpectEvent(client.JOIN, cmdPattern{who: h.Me().Nick, what: h.Channel})
	h.Join(h.Channel)
	if _, err := w.Wait(); err != nil {
		t.Fatalf("expected JOIN, got: %v", err)
	}
}
