package regtest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/goirc/logging"
)

type testLogger struct {
	t *testing.T
}

func (tl testLogger) Debug(fmt string, args ...any) {
	tl.t.Logf(fmt, args...)
}

func (tl testLogger) Info(fmt string, args ...any) {
	tl.t.Logf(fmt, args...)
}

func (tl testLogger) Warn(fmt string, args ...any) {
	tl.t.Logf(fmt, args...)
}

func (tl testLogger) Error(fmt string, args ...any) {
	tl.t.Logf(fmt, args...)
}

func EnableClientLogging(t *testing.T) {
	logging.SetLogger(testLogger{t})
	t.Cleanup(DisableClientLogging)
}

func DisableClientLogging() {
	logging.SetLogger(nil)
}

var (
	mu      sync.Mutex
	running bool
)

// Start connects the harness IRC client to the server, forks the bot,
// joins the test channel, and self-validates.
func Start(ctx context.Context) (*Harness, error) {
	mu.Lock()
	defer mu.Unlock()

	if running {
		return nil, fmt.Errorf("regtest: already started; call Stop() first")
	}

	botBinary, err := getBinaryPath("REGTEST_BOT")
	if err != nil {
		return nil, err
	}
	ircdBinary, err := getBinaryPath("REGTEST_IRCD")
	if err != nil {
		return nil, err
	}

	tmpDir, err := os.MkdirTemp("", "sp0rkle-test-*")
	if err != nil {
		return nil, fmt.Errorf("temp dir: %w", err)
	}

	h := &Harness{
		Channel: generateChannel(),
		BotNick: generateNick(),
		tempDir: tmpDir,
	}

	// First we must spin up the ergo ircd
	localAddr, err := freeLocalAddr()
	if err != nil {
		h.cleanup()
		return nil, fmt.Errorf("find free port: %w", err)
	}
	conf, err := templateConfig(tmpDir, localAddr)
	if err != nil {
		h.cleanup()
		return nil, fmt.Errorf("template ircd conf: %w", err)
	}

	h.ircd = Exec(ircdBinary, "run", "--conf=" + conf)
	// Only way to be sure ergo is booted and running is to
	// watch what it prints to stderr and wait for a log line...
	watcher := NewWriteWatcher(os.Stdout,
		fmt.Sprintf("now listening on %s,", localAddr))
	h.ircd.Stderr = watcher

	if err := h.ircd.Start(ctx); err != nil {
		h.cleanup()
		return nil, fmt.Errorf("ircd start: %w", err)
	}
	select {
		case <-watcher.Found:
		case <-time.After(5*time.Second):
		h.cleanup()
		return nil, fmt.Errorf("harness did not see ergo listening in 5s")
	}

	cfg := client.NewConfig("test0rkle", "regtest", "kicking sp0rkle in the test0rkles")
	cfg.Server = localAddr
	// disable flood protections cos they mess with timings
	cfg.Flood = true
	conn := client.Client(cfg)
	h.Conn = conn
	connected := make(chan struct{})
	h.HandleFunc(client.CONNECTED, func(c *client.Conn, l *client.Line) {
		c.Join(h.Channel)
		close(connected)
	})

	if err := conn.ConnectContext(ctx); err != nil {
		h.cleanup()
		return nil, fmt.Errorf("harness connect: %w", err)
	}
	// wait here until our client has definitely connected, so we are guaranteed
	// to be in the channel by the time the bot being tested is.
	select {
	case <-connected:
	case <-time.After(5 * time.Second):
		h.cleanup()
		return nil, fmt.Errorf("harness did not receive connected event in 5s")
	}

	if err := h.selfValidate(); err != nil {
		h.cleanup()
		return nil, fmt.Errorf("self-validate: %w", err)
	}

	h.sp0rkle = Exec(
		botBinary,
		"--servers=" + localAddr,
		"--channels=" + h.Channel,
		"--boltdb=" + tmpDir + "/sp0rkle.boltdb",
		"--backup_dir=" + tmpDir,
		"--nick=" + h.BotNick,
	)

	if err := h.sp0rkle.Start(ctx); err != nil {
		h.cleanup()
		return nil, fmt.Errorf("bot start: %w", err)
	}

	if err := h.waitForBotJoin(); err != nil {
		h.cleanup()
		return nil, fmt.Errorf("wait for bot join: %w", err)
	}

	running = true
	return h, nil
}

// Stop disconnects the harness IRC client and kills the bot process.
func (h *Harness) Stop() error {
	mu.Lock()
	defer mu.Unlock()
	err := h.cleanup()
	if err == nil {
		running = false
	}
	return err
}


// cleanup tears down all resources created by Start(), in reverse
// order to their creation.
func (h *Harness) cleanup() error {
	var errs []string

	// Close down bot process if it is running.
	if h.sp0rkle != nil {
		// Don't return straight away, cos we have to clean up properly.
		if err := h.sp0rkle.Stop(); err != nil {
			errs = append(errs, err.Error())
		}
		h.sp0rkle = nil
	}

	// Safely disconnect from IRC if connected.
	if h.Conn != nil {
		if h.Connected() {
			disconnected := make(chan struct{})
			h.HandleFunc(client.DISCONNECTED, func(c *client.Conn, l *client.Line) {
				close(disconnected)
			})
			h.Quit("regtest: cleanup")
			select {
			case <-disconnected:
			case <-time.After(5*time.Second):
			}
			if err := h.Conn.Close(); err != nil {
				errs = append(errs, err.Error())
			}
		}
		h.Conn = nil
	}

	// Close down ircd process if it is running.
	if h.ircd != nil {
		// Don't return straight away, cos we have to clean up properly.
		if err := h.ircd.Stop(); err != nil {
			errs = append(errs, err.Error())
		}
		h.ircd = nil
	}

	// Clean up temp dir if we have one.
	if h.tempDir != "" {
		if err := os.RemoveAll(h.tempDir); err != nil {
			errs = append(errs, err.Error())
		}
		h.tempDir = ""
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("Encountered %d errors during cleanup:\n\t%s",
		len(errs), strings.Join(errs, "\n\t"))
}

func getBinaryPath(env string) (string, error) {
	path := os.Getenv(env)
	if path == "" {
		return "", fmt.Errorf("%s: environment variable not set", env)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("%s: could not get abspath for %q: %w", env, path, err)
	}
	stat, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("%s: abspath %q not found: %w", env, abs, err)
	}
	if !stat.Mode().IsRegular() || (stat.Mode().Perm() & 0500) != 0500 {
		return "", fmt.Errorf("%s: abspath %q not regular o+rx file", env, abs)
	}
	return abs, nil
}

func randSuffix() string {
	b := make([]byte, 7)
	if _, err := rand.Read(b); err != nil {
		panic("regtest: rand.Read failed: " + err.Error())
	}
	return hex.EncodeToString(b)[:7]
}

func generateChannel() string {
	return "#spt-" + randSuffix()
}

func generateNick() string {
	return "sp0rkle^" + randSuffix()
}

func (h *Harness) waitForBotJoin() error {
	match := func(line *client.Line) bool {
		return line.Nick == h.BotNick
	}
	if _, err := h.ExpectEvent(client.JOIN, PatternFunc(match)).Wait(); err != nil {
		return fmt.Errorf("bot did not join %q: %v", h.Channel, err)
	}
	return nil
}

func (h *Harness) selfValidate() error {
	myNick := h.Me().Nick
	if myNick == "" {
		return fmt.Errorf("self-validate: Me().Nick is empty")
	}

	w := h.Expect(exactPattern{nick: myNick, text: "hello"})
	h.Privmsg(myNick, "hello")
	if _, err := w.Wait(); err != nil {
		return fmt.Errorf("self-validate PRIVMSG %s: %v", myNick, err)
	}
	return nil
}
