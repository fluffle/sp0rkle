package regtest

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// BotProcess manages a single sp0rkle bot process.
type BotProcess struct {
	cmd     *exec.Cmd
	path    string
	args    []string
	tempDir string
}

// Start forks the sp0rkle binary.
func (bp *BotProcess) Start() error {
	if bp.path == "" {
		return fmt.Errorf("bot process: empty binary path")
	}

	if bp.cmd != nil {
		return fmt.Errorf("bot process: already started")
	}

	bp.cmd = exec.Command(bp.path, bp.args...)
	bp.cmd.Stdout = os.Stdout
	bp.cmd.Stderr = os.Stderr

	if err := bp.cmd.Start(); err != nil {
		bp.cmd = nil
		return fmt.Errorf("bot process: start: %w", err)
	}

	return nil
}

// Stop kills the bot process and waits for it to exit.
func (bp *BotProcess) Stop() error {
	if bp.cmd == nil {
		return nil
	}

	bp.cmd.Process.Signal(os.Interrupt)
	waitErr := bp.cmd.Wait()

	// Clean up tempDir regardless of wait outcome.
	var cleanupErr error
	if bp.tempDir != "" {
		cleanupErr = os.RemoveAll(bp.tempDir)
	}

	bp.cmd = nil

	// Ignore "no such process" - process may have exited on its own.
	if waitErr != nil && !errors.Is(waitErr, syscall.ESRCH) {
		if cleanupErr != nil {
			return fmt.Errorf("bot process: wait: %w, cleanup: %w", waitErr, cleanupErr)
		}
		return fmt.Errorf("bot process: wait: %w", waitErr)
	}

	if cleanupErr != nil {
		return fmt.Errorf("bot process: cleanup: %w", cleanupErr)
	}

	return nil
}

// IsRunning returns true if the bot process is still running.
func (bp *BotProcess) IsRunning() bool {
	if bp.cmd == nil {
		return false
	}
	err := bp.cmd.Process.Signal(syscall.Signal(0))
	return err == nil
}
