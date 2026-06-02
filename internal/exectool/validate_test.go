package exectool

import (
	"testing"
	"time"

	"github.com/inhuman/mcp-exec/internal/config"
)

func testCfg() config.Config {
	return config.Config{
		DefaultTimeoutS: 30,
		MaxTimeoutS:     300,
		MaxOutputBytes:  1 << 20,
		MaxStdinBytes:   1 << 20,
		Python:          "python3",
	}
}

func TestValidate_EmptyCode(t *testing.T) {
	if _, err := validate(ExecRequest{Code: "   "}, testCfg()); err == nil {
		t.Fatal("expected error for empty code")
	}
}

func TestValidate_StdinTooLarge(t *testing.T) {
	cfg := testCfg()
	cfg.MaxStdinBytes = 4
	if _, err := validate(ExecRequest{Code: "x", Stdin: "toolong"}, cfg); err == nil {
		t.Fatal("expected error for oversized stdin")
	}
}

func TestValidate_TimeoutClamp(t *testing.T) {
	cfg := testCfg()
	cases := []struct {
		in   int
		want time.Duration
	}{
		{0, 30 * time.Second},   // default
		{-5, 1 * time.Second},   // below min
		{10, 10 * time.Second},  // in range
		{900, 300 * time.Second}, // above max
	}
	for _, c := range cases {
		req, err := validate(ExecRequest{Code: "x", TimeoutS: c.in}, cfg)
		if err != nil {
			t.Fatalf("unexpected error for timeout %d: %v", c.in, err)
		}
		if req.Timeout != c.want {
			t.Errorf("timeout %d → got %v, want %v", c.in, req.Timeout, c.want)
		}
	}
}

func TestValidate_MapsFields(t *testing.T) {
	req, err := validate(ExecRequest{Code: "print(1)", Stdin: "data"}, testCfg())
	if err != nil {
		t.Fatal(err)
	}
	if req.Code != "print(1)" || req.Stdin != "data" || req.MaxOutputBytes != 1<<20 {
		t.Errorf("unexpected mapped request: %+v", req)
	}
}
