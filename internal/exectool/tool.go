package exectool

import (
	"context"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"

	"github.com/inhuman/mcp-exec/internal/config"
	"github.com/inhuman/mcp-exec/internal/isolator"
)

// Handler implements the exec tool. Runs are serialized within the instance
// (constitution VIII): only one execution proceeds at a time.
type Handler struct {
	iso isolator.Isolator
	cfg config.Config
	log *zap.Logger
	mu  sync.Mutex
}

func New(cfg config.Config, iso isolator.Isolator, log *zap.Logger) *Handler {
	return &Handler{iso: iso, cfg: cfg, log: log}
}

// Handle is the MCP tool handler. Invalid input returns an error (tool-call
// error); a non-zero exit or a timeout is a successful result with the
// corresponding fields set.
func (h *Handler) Handle(ctx context.Context, _ *mcp.CallToolRequest, in ExecRequest) (*mcp.CallToolResult, ExecResult, error) {
	req, err := validate(in, h.cfg)
	if err != nil {
		return nil, ExecResult{}, err
	}

	h.mu.Lock()
	res, runErr := h.iso.Run(ctx, req)
	h.mu.Unlock()
	if runErr != nil {
		h.log.Error("exec failed", zap.Error(runErr))
		return nil, ExecResult{}, runErr
	}

	out := ExecResult{
		Stdout:     res.Stdout,
		Stderr:     res.Stderr,
		ExitCode:   res.ExitCode,
		DurationMS: res.Duration.Milliseconds(),
		Truncated:  res.Truncated,
		TimedOut:   res.TimedOut,
	}

	// Metadata only: never the code, stdin or full output (constitution IX).
	h.log.Info("exec done",
		zap.Int("exit_code", out.ExitCode),
		zap.Int64("duration_ms", out.DurationMS),
		zap.Int("output_bytes", len(out.Stdout)+len(out.Stderr)),
		zap.Bool("truncated", out.Truncated),
		zap.Bool("timed_out", out.TimedOut),
	)

	return nil, out, nil
}
