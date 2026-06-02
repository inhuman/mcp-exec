package exectool

import (
	"fmt"
	"strings"
	"time"

	"github.com/inhuman/mcp-exec/internal/config"
	"github.com/inhuman/mcp-exec/internal/isolator"
)

// ExecRequest is the public input schema of the exec tool.
type ExecRequest struct {
	Code     string `json:"code" jsonschema:"Python source to execute (required)"`
	TimeoutS int    `json:"timeout_s,omitempty" jsonschema:"wall-clock timeout in seconds (default 30, clamped to [1,300])"`
	Stdin    string `json:"stdin,omitempty" jsonschema:"data piped to the program's stdin (max 1 MiB)"`
}

// ExecResult is the public output schema of the exec tool. The field set is a
// versioned contract (constitution III) and is identical across all transports.
type ExecResult struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	DurationMS int64  `json:"duration_ms"`
	Truncated  bool   `json:"truncated"`
	TimedOut   bool   `json:"timed_out"`
}

// validate checks the input and maps it to an isolator.Request, applying the
// configured stdin limit and clamping timeout_s into [1, MaxTimeoutS]. A
// returned error means an invalid tool call (execution must not start).
func validate(in ExecRequest, cfg config.Config) (isolator.Request, error) {
	if strings.TrimSpace(in.Code) == "" {
		return isolator.Request{}, fmt.Errorf("code must not be empty")
	}
	if len(in.Stdin) > cfg.MaxStdinBytes {
		return isolator.Request{}, fmt.Errorf("stdin exceeds limit of %d bytes", cfg.MaxStdinBytes)
	}

	t := in.TimeoutS
	if t == 0 {
		t = cfg.DefaultTimeoutS
	}
	if t < 1 {
		t = 1
	}
	if t > cfg.MaxTimeoutS {
		t = cfg.MaxTimeoutS
	}

	return isolator.Request{
		Code:           in.Code,
		Stdin:          in.Stdin,
		Timeout:        time.Duration(t) * time.Second,
		MaxOutputBytes: cfg.MaxOutputBytes,
	}, nil
}
