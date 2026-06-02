package isolator

import (
	"context"
	"time"
)

// Request is a single sandboxed execution request.
type Request struct {
	Code           string
	Stdin          string
	Timeout        time.Duration
	MaxOutputBytes int
}

// Result is the outcome of a single execution. It is never persisted.
type Result struct {
	Stdout    string
	Stderr    string
	ExitCode  int
	Duration  time.Duration
	Truncated bool
	TimedOut  bool
}

// Isolator runs caller code in an isolated environment and returns its result.
// An error is returned only for infrastructure failures (e.g. interpreter not
// found); a non-zero exit of the executed code is a normal Result, not an error.
type Isolator interface {
	Run(ctx context.Context, req Request) (Result, error)
}
