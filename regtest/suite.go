package regtest

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/fluffle/goirc/client"
)

var (
	mu          sync.Mutex
	globalHarness *Harness
	globalBot     *BotProcess
)

// Start connects the harness IRC client to the server, forks the bot,
// joins the test channel, and self-validates.
func Start() (*Harness, error) {
	mu.Lock()
	defer mu.Unlock()

	if globalHarness != nil || globalBot != nil {
		return nil, fmt.Errorf("regtest: already started; call Stop() first")
	}

	server := os.Getenv("REGTEST_SERVER")
	if server == "" {
		return nil, fmt.Errorf("REGTEST_SERVER environment variable not set")
	}

	channel := generateChannel()

	cfg := client.NewConfig("regtest-bot", "regtest", "regtest harness")
	cfg.Server = server
	conn := client.Client(cfg)

	if err := conn.Connect(); err != nil {
		return nil, fmt.Errorf("harness connect: %w", err)
	}

	h := &Harness{Conn: conn, channel: channel}
	h.SetBotNick("sp0rklf-test")

	botPath, err := findBotBinary()
	if err != nil {
		conn.Quit("regtest: find bot binary failed")
		conn.Close()
		return nil, fmt.Errorf("find bot binary: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "sp0rkle-test-*")
	if err != nil {
		conn.Quit("regtest: temp dir failed")
		conn.Close()
		return nil, fmt.Errorf("temp dir: %w", err)
	}

	bot := &BotProcess{
		path:    botPath,
		args:    []string{"--servers=" + server, "--channels=" + channel, "--boltdb=" + tmpDir + "/sp0rkle.boltdb", "--backup_dir=", "--nick=sp0rklf-test"},
		tempDir: tmpDir,
	}

	if err := bot.Start(); err != nil {
		bot.Stop()
		conn.Quit("regtest: bot start failed")
		conn.Close()
		return nil, fmt.Errorf("bot start: %w", err)
	}

	h.globalBot = bot
	h.globalTmpDir = tmpDir

	globalHarness = h
	globalBot = bot

	conn.Join(channel)

	if err := waitForBotJoin(h, "sp0rklf-test", 5*time.Second); err != nil {
		cleanup(h, bot, tmpDir)
		return nil, fmt.Errorf("wait for bot join: %w", err)
	}

	if err := h.selfValidate(); err != nil {
		cleanup(h, bot, tmpDir)
		return nil, fmt.Errorf("self-validate: %w", err)
	}

	return h, nil
}

// Stop disconnects the harness IRC client and kills the bot process.
func Stop() error {
	mu.Lock()
	defer mu.Unlock()

	var retErr error

	if globalBot != nil {
		if err := globalBot.Stop(); err != nil {
			retErr = fmt.Errorf("bot stop: %w", err)
		}
		globalBot = nil
	}

	if globalHarness != nil {
		h := globalHarness
		globalHarness = nil

		h.Quit("regtest stop")
		if err := h.Conn.Close(); err != nil {
			if retErr != nil {
				retErr = fmt.Errorf("%w; close: %w", retErr, err)
			} else {
				retErr = fmt.Errorf("close: %w", err)
			}
		}
	}

	return retErr
}

// cleanup tears down all resources created by Start().
func cleanup(h *Harness, bot *BotProcess, tmpDir string) {
	if bot != nil {
		bot.Stop()
	}
	if h != nil {
		h.Quit("regtest: cleanup")
		h.Conn.Close()
	}
	if tmpDir != "" {
		os.RemoveAll(tmpDir)
	}
}

func generateChannel() string {
	b := make([]byte, 7)
	if _, err := rand.Read(b); err != nil {
		panic("regtest: rand.Read failed: " + err.Error())
	}
	return "#spt-" + hex.EncodeToString(b)[:7]
}

func findBotBinary() (string, error) {
	binary := os.Getenv("REGTEST_BINARY")
	if binary != "" {
		if _, err := os.Stat(binary); err != nil {
			return "", fmt.Errorf("REGTEST_BINARY %q not found: %w", binary, err)
		}
		return binary, nil
	}
	cmd := exec.Command("go", "build", "-o", "sp0rkle-test", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build bot binary: %w", err)
	}
	return "./sp0rkle-test", nil
}

func waitForBotJoin(h *Harness, botNick string, timeout time.Duration) error {
	resultCh := make(chan *client.Line, 1)
	var remover client.Remover
	remover = h.HandleFunc(client.JOIN, func(conn *client.Conn, line *client.Line) {
		if line.Nick == botNick {
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
	case <-resultCh:
		return nil
	case <-timer.C:
		return fmt.Errorf("bot did not join %q within %v", h.channel, timeout)
	}
}

func (h *Harness) selfValidate() error {
	resultCh := make(chan *client.Line, 1)

	myNick := h.Me().Nick
	if myNick == "" {
		return fmt.Errorf("self-validate: Me().Nick is empty")
	}

	var remover client.Remover
	remover = h.HandleFunc(client.PRIVMSG, func(conn *client.Conn, line *client.Line) {
		if line.Nick == myNick {
			remover.Remove()
			select {
			case resultCh <- line:
			default:
			}
		}
	})
	defer remover.Remove()

	h.Privmsg(myNick, "hello")

	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()

	select {
	case <-resultCh:
		return nil
	case <-timer.C:
		return fmt.Errorf("self-validate: no response to own nick within 2s")
	}
}
