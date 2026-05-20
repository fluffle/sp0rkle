package regtest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type WriteWatcher struct {
	Found chan struct{}
	mu *sync.Mutex
	buf *bytes.Buffer
	dest io.Writer
	needle []byte
}

func NewWriteWatcher(dest io.Writer, needle string) *WriteWatcher {
	return &WriteWatcher{
		mu: &sync.Mutex{},
		buf: bytes.NewBuffer(make([]byte, 0, 8192)),
		dest: dest,
		needle: []byte(needle),
		Found: make(chan struct{}),
	}
}

func (w *WriteWatcher) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buf == nil {
		return w.dest.Write(p)
	}
	w.buf.Write(p)
	if bytes.Contains(w.buf.Bytes(), w.needle) {
		w.buf = nil
		close(w.Found)
	}
	return w.dest.Write(p)
}

// Process manages a single process.
type Process struct {
	*sync.Mutex
	Stdout  io.Writer
	Stderr  io.Writer
	cmd     *exec.Cmd
	path    string
	args    []string
	cancel  context.CancelFunc
	exit    chan error
}

func Exec(path string, args ...string) *Process {
	return &Process{
		Mutex: &sync.Mutex{},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		path: path,
		args: args,
		exit: make(chan error, 1),
	}
}

// Start forks the binary and optionally waits for a particular string.
func (p *Process) Start(ctx context.Context) error {
	p.Lock()
	defer p.Unlock()
	if p.path == "" {
		return fmt.Errorf("process: empty binary path")
	}

	if p.cmd != nil {
		return fmt.Errorf("%s: already started", p.path)
	}
	cctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	p.cmd = exec.CommandContext(cctx, p.path, p.args...)
	p.cmd.Stdout = p.Stdout
	p.cmd.Stderr = p.Stderr

	if err := p.cmd.Start(); err != nil {
		p.cmd = nil
		cancel()
		return fmt.Errorf("%s: start: %w", p.path, err)
	}
	go p.wait()

	return nil
}

func (p *Process) wait() {
	p.exit <- p.cmd.Wait()
	close(p.exit)
	p.Lock()
	p.cmd = nil
	p.Unlock()
}

func (p *Process) kill() error {
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Signal(os.Interrupt)
	}
	select {
	case waitErr := <-p.exit:
		// Already exited, cancel just cleans up the context.
		p.cancel()
		return waitErr
	case <-time.After(1*time.Second):
		// No exit after 1s, cancel should send SIGKILL.
		p.cancel()
	}
	return <-p.exit
}

// Stop kills the process and waits for it to exit.
func (p *Process) Stop() error {
	p.Lock()
	defer p.Unlock()
	if p.cmd == nil {
		return nil
	}

	waitErr := p.kill()
	if waitErr == nil {
		return nil
	}
	// an error from us killing the process is fine.
	if ee, ok := errors.AsType[*exec.ExitError](waitErr); ok {
		if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
			sig := ws.Signal()
			if sig == syscall.SIGINT || sig == syscall.SIGKILL {
				return nil
			}
		}
	}
	return fmt.Errorf("%s: exited with error: %w", p.path, waitErr)
}

// IsRunning returns true if the process is still running.
func (p *Process) IsRunning() bool {
	p.Lock()
	defer p.Unlock()
	return p.cmd != nil
}
