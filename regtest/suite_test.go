package regtest

import (
	"os"
	"testing"
	"time"

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

func TestStartRequiresServer(t *testing.T) {
	old := os.Getenv("REGTEST_SERVER")
	os.Unsetenv("REGTEST_SERVER")
	defer os.Setenv("REGTEST_SERVER", old)

	_, err := Start()
	if err == nil {
		t.Error("Start without REGTEST_SERVER should error")
	}
}

func TestWaitForBotJoinTimeout(t *testing.T) {
	cfg := client.NewConfig("testbot", "test", "test bot")
	conn := client.Client(cfg)
	h := &Harness{Conn: conn, channel: "#test"}
	err := waitForBotJoin(h, "nonexistent-bot", 50*time.Millisecond)
	if err == nil {
		t.Error("waitForBotJoin with no bot should error")
	}
}
