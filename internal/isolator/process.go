package isolator

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Process runs each request as a python subprocess inside the locked container:
// an ephemeral tmpdir as cwd, its own process-group (so timeouts kill all
// children), and a shared cap on combined stdout+stderr. Network/privilege
// isolation is enforced by how the container/pod is launched (see quickstart),
// not inside this process.
type Process struct {
	python string
}

func NewProcess(python string) *Process {
	if python == "" {
		python = "python3"
	}
	return &Process{python: python}
}

func (p *Process) Run(ctx context.Context, req Request) (Result, error) {
	tmpdir, err := os.MkdirTemp("", "mcpexec-")
	if err != nil {
		return Result{}, err
	}
	defer os.RemoveAll(tmpdir)

	script := filepath.Join(tmpdir, "main.py")
	if err := os.WriteFile(script, []byte(req.Code), 0o600); err != nil {
		return Result{}, err
	}

	runCtx, cancel := context.WithTimeout(ctx, req.Timeout)
	defer cancel()

	// Run the script file so process stdin stays available for the program.
	cmd := exec.CommandContext(runCtx, p.python, script)
	cmd.Dir = tmpdir
	cmd.Stdin = strings.NewReader(req.Stdin)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// On timeout kill the whole process-group, not just the leader, so no child
	// outlives the deadline. The process may already be gone (ESRCH) — best-effort.
	cmd.Cancel = func() error {
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return nil
	}

	cap := &outputCapper{remaining: req.MaxOutputBytes}
	cmd.Stdout = cap.writer(&cap.stdout)
	cmd.Stderr = cap.writer(&cap.stderr)

	start := time.Now()
	runErr := cmd.Run()
	dur := time.Since(start)

	timedOut := errors.Is(runCtx.Err(), context.DeadlineExceeded)

	exitCode := 0
	switch {
	case timedOut:
		exitCode = -1
	case runErr != nil:
		var ee *exec.ExitError
		if errors.As(runErr, &ee) {
			exitCode = ee.ExitCode()
		} else {
			return Result{}, runErr
		}
	}

	return Result{
		Stdout:    safeUTF8(cap.stdout.String()),
		Stderr:    safeUTF8(cap.stderr.String()),
		ExitCode:  exitCode,
		Duration:  dur,
		Truncated: cap.truncated,
		TimedOut:  timedOut,
	}, nil
}

func safeUTF8(s string) string { return strings.ToValidUTF8(s, "�") }

// outputCapper caps the combined size of stdout and stderr. Both writers share
// one counter under a mutex, since the two streams are written concurrently.
type outputCapper struct {
	mu        sync.Mutex
	remaining int
	truncated bool
	stdout    bytes.Buffer
	stderr    bytes.Buffer
}

func (c *outputCapper) writer(b *bytes.Buffer) *capWriter { return &capWriter{c: c, b: b} }

type capWriter struct {
	c *outputCapper
	b *bytes.Buffer
}

func (w *capWriter) Write(p []byte) (int, error) {
	w.c.mu.Lock()
	defer w.c.mu.Unlock()
	if w.c.remaining <= 0 {
		w.c.truncated = true
		return len(p), nil
	}
	if len(p) > w.c.remaining {
		w.b.Write(p[:w.c.remaining])
		w.c.remaining = 0
		w.c.truncated = true
		return len(p), nil
	}
	w.b.Write(p)
	w.c.remaining -= len(p)
	return len(p), nil
}
