package integration

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/inhuman/mcp-exec/internal/isolator"
)

func requirePython(t *testing.T) string {
	t.Helper()
	p, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not available")
	}
	return p
}

func run(t *testing.T, code, stdin string, timeout time.Duration, maxOut int) isolator.Result {
	t.Helper()
	iso := isolator.NewProcess(requirePython(t))
	res, err := iso.Run(context.Background(), isolator.Request{
		Code:           code,
		Stdin:          stdin,
		Timeout:        timeout,
		MaxOutputBytes: maxOut,
	})
	if err != nil {
		t.Fatalf("isolator error: %v", err)
	}
	return res
}

func TestProcess_BasicStdout(t *testing.T) {
	res := run(t, "print(2 + 2)", "", 10*time.Second, 1<<20)
	if res.Stdout != "4\n" || res.ExitCode != 0 || res.TimedOut {
		t.Errorf("unexpected result: %+v", res)
	}
	if res.Duration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestProcess_NonZeroExit(t *testing.T) {
	res := run(t, "import sys; sys.exit(3)", "", 10*time.Second, 1<<20)
	if res.ExitCode != 3 || res.TimedOut {
		t.Errorf("expected exit 3, got %+v", res)
	}
}

func TestProcess_Stdin(t *testing.T) {
	res := run(t, "import sys; print(sys.stdin.read().strip().upper())", "hello", 10*time.Second, 1<<20)
	if res.Stdout != "HELLO\n" {
		t.Errorf("stdin not piped: %q", res.Stdout)
	}
}

func TestProcess_Timeout(t *testing.T) {
	start := time.Now()
	res := run(t, "while True:\n    pass", "", 1*time.Second, 1<<20)
	if !res.TimedOut {
		t.Errorf("expected timed_out=true, got %+v", res)
	}
	if res.ExitCode != -1 {
		t.Errorf("expected exit_code=-1 on timeout, got %d", res.ExitCode)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("timeout took too long: %v", elapsed)
	}
}

func TestProcess_OutputTruncated(t *testing.T) {
	res := run(t, "print('x' * 100000)", "", 10*time.Second, 1024)
	if !res.Truncated {
		t.Error("expected truncated=true")
	}
	if len(res.Stdout)+len(res.Stderr) > 1024 {
		t.Errorf("output exceeded cap: %d bytes", len(res.Stdout)+len(res.Stderr))
	}
}

func TestProcess_EphemeralTmpdir(t *testing.T) {
	// First run writes a file into its cwd; second run must not see it.
	r1 := run(t, "open('marker.txt', 'w').write('x')", "", 10*time.Second, 1<<20)
	if r1.ExitCode != 0 {
		t.Fatalf("write run failed: %+v", r1)
	}
	r2 := run(t, "import os; print(os.path.exists('marker.txt'))", "", 10*time.Second, 1<<20)
	if r2.Stdout != "False\n" {
		t.Errorf("tmpdir not ephemeral, second run saw marker: %q", r2.Stdout)
	}
}
