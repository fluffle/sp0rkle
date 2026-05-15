package regtest

import (
	"context"
	"testing"
	"time"
)

func TestStart_RequiresPath(t *testing.T) {
	p := Exec("")
	err := p.Start(context.Background())
	if err == nil {
		t.Error("Start with empty path should error")
	}
}

func TestIsRunning_NilCmd(t *testing.T) {
	p := Exec("")
	if p.IsRunning() {
		t.Error("IsRunning should return false for nil cmd")
	}
}

func TestStopNilCmd(t *testing.T) {
	p := Exec("")
	err := p.Stop()
	if err != nil {
		t.Errorf("Stop with nil cmd should not error, got: %v", err)
	}
}

func TestStart_NonexistentPath(t *testing.T) {
	p := Exec("/nonexistent")
	err1 := p.Start(context.Background())
	// First call fails (binary doesn't exist), but cmd is reset to nil.
	if err1 == nil {
		t.Fatal("Start err = nil, want error")
	}

	// Second call should also fail (same reason), not panic.
	err2 := p.Start(context.Background())
	if err2 == nil {
		t.Fatal("Duplicate Start err = nil, want error")
	}

	// Stop should be safe even after failed Start.
	err := p.Stop()
	if err != nil {
		t.Errorf("Stop after failed Start should not error, got: %v", err)
	}
}

func TestStartStop_ProcessExitsNaturally(t *testing.T) {
	p := Exec("/bin/echo", "hello")
	err1 := p.Start(context.Background())
	// First call should succeed
	if err1 != nil {
		t.Fatalf("Start err = %v, want nil", err1)
	}

	// Second call should fail because command is running / has run, not panic.
	err2 := p.Start(context.Background())
	if err2 == nil {
		t.Fatal("Duplicate Start err = nil, want error")
	}

	// Wait a bit.
	<-time.After(10 * time.Millisecond)

	// Process should not be running.
	if got := p.IsRunning(); got != false {
		t.Errorf("IsRunning = %t after 1s, want false", got)
	}

	// Stop should succeed.
	err := p.Stop()
	if err != nil {
		t.Errorf("Stop err = %v, want nil", err)
	}
}

func TestStop_KillsProcess(t *testing.T) {
	p := Exec("/bin/sleep", "60")
	err1 := p.Start(context.Background())
	// First call should succeed
	if err1 != nil {
		t.Fatalf("Start err = %v, want nil", err1)
	}

	// Wait a bit.
	<-time.After(100 * time.Millisecond)

	// Process should be running.
	if got := p.IsRunning(); got != true {
		t.Errorf("IsRunning = %t after 1s, want true", got)
	}

	// Stop should succeed.
	err := p.Stop()
	if err != nil {
		t.Errorf("Stop err = %v, want nil", err)
	}
}
