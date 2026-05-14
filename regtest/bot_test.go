package regtest

import (
	"os/exec"
	"testing"
)

func TestBotProcessFields(t *testing.T) {
	// Compile-time check: these operations must compile.
	// If they don't, the BotProcess struct fields are wrong.
	var p BotProcess
	_ = p.cmd
	_ = p.path
	_ = p.args
	_ = p.tempDir

	// Verify the struct can be used as expected.
	p.cmd = &exec.Cmd{}
	p.path = "/usr/bin/sp0rkle"
	p.args = []string{"--test"}
	p.tempDir = "/tmp/test"

	if p.path != "/usr/bin/sp0rkle" {
		t.Errorf("path field = %q, want %q", p.path, "/usr/bin/sp0rkle")
	}
	if p.tempDir != "/tmp/test" {
		t.Errorf("tempDir field = %q, want %q", p.tempDir, "/tmp/test")
	}
}

func TestBotProcessStartRequiresPath(t *testing.T) {
	p := &BotProcess{}
	err := p.Start()
	if err == nil {
		t.Error("Start with empty path should error")
	}
}

func TestBotProcessIsRunning(t *testing.T) {
	p := &BotProcess{}
	if p.IsRunning() {
		t.Error("IsRunning should return false for nil cmd")
	}
}

func TestBotProcessStopNilCmd(t *testing.T) {
	p := &BotProcess{}
	err := p.Stop()
	if err != nil {
		t.Errorf("Stop with nil cmd should not error, got: %v", err)
	}
}

func TestBotProcessStartIdempotent(t *testing.T) {
	p := &BotProcess{path: "/nonexistent"}
	err1 := p.Start()
	// First call fails (binary doesn't exist), but cmd is reset to nil.
	if err1 == nil {
		t.Fatal("expected Start to fail with nonexistent binary")
	}

	// Second call should also fail (same reason), not panic.
	err2 := p.Start()
	if err2 == nil {
		t.Fatal("expected second Start to fail with nonexistent binary")
	}

	// Stop should be safe even after failed Start.
	err := p.Stop()
	if err != nil {
		t.Errorf("Stop after failed Start should not error, got: %v", err)
	}
}

func TestBotProcessStopIgnoresESRCH(t *testing.T) {
	// Test that Stop doesn't error if the process already exited.
	p := &BotProcess{path: "/nonexistent"}
	p.Start() // This fails but the process never starts.
	// The process never existed, so Signal(0) would fail.
	// Stop should handle this gracefully.
	err := p.Stop()
	if err != nil {
		t.Errorf("Stop should not error for process that never existed, got: %v", err)
	}
}
