package exectool

import (
	"context"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/inhuman/mcp-exec/internal/isolator"
)

// TestHandle_LogPrivacy verifies SC-008 / constitution IX: the code, stdin and
// full output never reach the logs; only metadata does.
func TestHandle_LogPrivacy(t *testing.T) {
	const (
		secretCode   = "SECRET_CODE_DO_NOT_LOG"
		secretStdin  = "SECRET_STDIN_DO_NOT_LOG"
		secretOutput = "SECRET_OUTPUT_DO_NOT_LOG"
	)

	core, logs := observer.New(zap.InfoLevel)
	iso := isolator.NewNoop()
	iso.Result = isolator.Result{Stdout: secretOutput, ExitCode: 0}

	h := New(testCfg(), iso, zap.New(core))
	_, out, err := h.Handle(context.Background(), nil, ExecRequest{Code: secretCode, Stdin: secretStdin})
	if err != nil {
		t.Fatal(err)
	}
	if out.Stdout != secretOutput {
		t.Fatalf("output not returned to caller: %q", out.Stdout)
	}

	for _, e := range logs.All() {
		blob := e.Message
		for k, v := range e.ContextMap() {
			blob += " " + k + "="
			if s, ok := v.(string); ok {
				blob += s
			}
		}
		for _, secret := range []string{secretCode, secretStdin, secretOutput} {
			if strings.Contains(blob, secret) {
				t.Errorf("log leaked sensitive data %q in entry: %s", secret, blob)
			}
		}
	}

	// Metadata must be present.
	if logs.FilterField(zap.Int("exit_code", 0)).Len() == 0 {
		t.Error("expected exit_code metadata in logs")
	}
}

func TestHandle_InvalidInputIsToolError(t *testing.T) {
	h := New(testCfg(), isolator.NewNoop(), zap.NewNop())
	if _, _, err := h.Handle(context.Background(), nil, ExecRequest{Code: ""}); err == nil {
		t.Fatal("expected tool error for empty code")
	}
}
