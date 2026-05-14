package regtest

import (
	"testing"
	"time"

	"github.com/fluffle/goirc/client"
)

func TestHarnessBotNick(t *testing.T) {
	cfg := client.NewConfig("testbot", "test", "test bot")
	conn := client.Client(cfg)
	h := &Harness{Conn: conn, botNick: "testbot"}
	if h.BotNick() != "testbot" {
		t.Errorf("BotNick() = %q, want %q", h.BotNick(), "testbot")
	}
}

func TestHarnessChannel(t *testing.T) {
	cfg := client.NewConfig("testbot", "test", "test bot")
	conn := client.Client(cfg)
	h := &Harness{Conn: conn, botNick: "testbot", channel: "#test"}
	if h.Channel() != "#test" {
		t.Errorf("Channel() = %q, want %q", h.Channel(), "#test")
	}
}

func TestHarnessSetBotNick(t *testing.T) {
	cfg := client.NewConfig("testbot", "test", "test bot")
	conn := client.Client(cfg)
	h := &Harness{Conn: conn, botNick: "testbot"}
	h.SetBotNick("newnick")
	if h.botNick != "newnick" {
		t.Errorf("SetBotNick() did not update botNick")
	}
}

func TestHarnessExpectTimeout(t *testing.T) {
	cfg := client.NewConfig("testbot", "test", "test bot")
	conn := client.Client(cfg)
	h := &Harness{Conn: conn, botNick: "testbot", channel: "#test"}
	_, err := h.Expect(Exact("world"), 10*time.Millisecond)
	if err == nil {
		t.Error("Expect with no server should error")
	}
}
