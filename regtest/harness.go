package regtest

import (
	"fmt"
	"time"

	"github.com/fluffle/goirc/client"
)

// Harness manages the harness IRC client used to drive tests.
type Harness struct {
	*client.Conn
	botNick    string
	channel    string
	globalBot  *BotProcess
	globalTmpDir string
}

// BotNick returns the nick of the bot being tested.
func (h *Harness) BotNick() string {
	return h.botNick
}

// Channel returns the test channel name.
func (h *Harness) Channel() string {
	return h.channel
}

// SetBotNick sets the nick of the bot being tested.
func (h *Harness) SetBotNick(nick string) {
	h.botNick = nick
}

// SendAndExpect sends a PRIVMSG to the test channel and waits for a response
// matching the pattern. Returns the matched Line or an error on timeout.
func (h *Harness) SendAndExpect(msg string, p Pattern, timeout time.Duration) (*client.Line, error) {
	if h.Conn == nil {
		return nil, fmt.Errorf("send and expect: harness Conn is nil")
	}
	if h.channel == "" {
		return nil, fmt.Errorf("send and expect: channel is empty")
	}
	h.Privmsg(h.channel, msg)
	return h.Expect(p, timeout)
}

// Expect waits for the next IRC line matching the pattern from any source.
// Returns the matched Line or an error on timeout.
func (h *Harness) Expect(p Pattern, timeout time.Duration) (*client.Line, error) {
	if h.Conn == nil {
		return nil, fmt.Errorf("expect: harness Conn is nil")
	}
	if p == nil {
		return nil, fmt.Errorf("expect: pattern is nil")
	}

	resultCh := make(chan *client.Line, 1)

	var remover client.Remover
	remover = h.HandleFunc(client.PRIVMSG, func(conn *client.Conn, line *client.Line) {
		if p.Match(line) {
			remover.Remove()
			select {
			case resultCh <- line:
			default:
			}
		}
	})
	defer remover.Remove()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case line := <-resultCh:
		return line, nil
	case <-timer.C:
		return nil, fmt.Errorf("expect: timeout after %v", timeout)
	}
}
