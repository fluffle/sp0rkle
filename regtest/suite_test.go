package regtest

import (
	"context"
	"os"
	"testing"

	"github.com/fluffle/goirc/client"
)

func TestChannelGeneration(t *testing.T) {
	ch := generateChannel()
	if len(ch) != 12 { // #spt-XXXXXXX
		t.Errorf("generateChannel() length = %d, want 12", len(ch))
	}
	if ch[0] != '#' {
		t.Errorf("generateChannel() = %q, should start with #", ch)
	}
	if ch[:5] != "#spt-" {
		t.Errorf("generateChannel() = %q, should start with #spt-", ch)
	}
}

func TestChannelUniqueness(t *testing.T) {
	channels := make(map[string]bool)
	for i := 0; i < 100; i++ {
		ch := generateChannel()
		if channels[ch] {
			t.Errorf("generateChannel() produced duplicate: %q", ch)
		}
		channels[ch] = true
	}
}

func TestStartRequiresBotBinary(t *testing.T) {
	old := os.Getenv("REGTEST_BOT")
	os.Unsetenv("REGTEST_BOT")
	defer os.Setenv("REGTEST_BOT", old)

	_, err := Start(context.Background())
	if err == nil {
		t.Error("Start without REGTEST_BOT should error")
	}
}

func TestStartRequiresIrcdBinary(t *testing.T) {
	old := os.Getenv("REGTEST_IRCD")
	os.Unsetenv("REGTEST_IRCD")
	defer os.Setenv("REGTEST_IRCD", old)

	_, err := Start(context.Background())
	if err == nil {
		t.Error("Start without REGTEST_IRCD should error")
	}
}

func TestWaitForBotJoinTimeout(t *testing.T) {
	cfg := client.NewConfig("testbot", "test", "test bot")
	conn := client.Client(cfg)
	h := &Harness{Conn: conn, Channel: "#test"}
	err := h.waitForBotJoin()
	if err == nil {
		t.Error("waitForBotJoin with no bot should error")
	}
}
