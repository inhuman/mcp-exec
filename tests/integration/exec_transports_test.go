package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"

	"github.com/inhuman/mcp-exec/internal/config"
	"github.com/inhuman/mcp-exec/internal/exectool"
	"github.com/inhuman/mcp-exec/internal/isolator"
	"github.com/inhuman/mcp-exec/internal/server"
)

func parityCfg() config.Config {
	return config.Config{DefaultTimeoutS: 30, MaxTimeoutS: 300, MaxOutputBytes: 1 << 20, MaxStdinBytes: 1 << 20}
}

// newServer builds a server backed by a deterministic Noop isolator so results
// are identical across transports regardless of timing.
func newServer() *mcp.Server {
	iso := isolator.NewNoop()
	iso.Result = isolator.Result{Stdout: "parity\n", ExitCode: 0}
	return server.New(parityCfg(), iso, zap.NewNop())
}

func callExec(t *testing.T, cs *mcp.ClientSession) exectool.ExecResult {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "exec",
		Arguments: map[string]any{"code": "print('parity')"},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("unexpected tool error: %+v", res.Content)
	}
	raw, err := json.Marshal(res.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var out exectool.ExecResult
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	out.DurationMS = 0 // timing differs between calls; not part of parity
	return out
}

func connect(t *testing.T, transport mcp.Transport) *mcp.ClientSession {
	t.Helper()
	client := mcp.NewClient(&mcp.Implementation{Name: "parity-test", Version: "v0"}, nil)
	cs, err := client.Connect(context.Background(), transport, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func TestTransportParity(t *testing.T) {
	ctx := context.Background()

	// stdio (in-memory pipe).
	srvStdio := newServer()
	clientT, serverT := mcp.NewInMemoryTransports()
	go func() { _ = srvStdio.Run(ctx, serverT) }()
	stdioRes := callExec(t, connect(t, clientT))

	// streamable HTTP.
	httpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return newServer() }, nil)
	httpTS := httptest.NewServer(httpHandler)
	t.Cleanup(func() { httpTS.CloseClientConnections(); httpTS.Close() })
	httpRes := callExec(t, connect(t, &mcp.StreamableClientTransport{Endpoint: httpTS.URL}))

	// SSE.
	sseHandler := mcp.NewSSEHandler(func(*http.Request) *mcp.Server { return newServer() }, nil)
	sseTS := httptest.NewServer(sseHandler)
	t.Cleanup(func() { sseTS.CloseClientConnections(); sseTS.Close() })
	sseRes := callExec(t, connect(t, &mcp.SSEClientTransport{Endpoint: sseTS.URL}))

	if !reflect.DeepEqual(stdioRes, httpRes) {
		t.Errorf("stdio vs http mismatch:\n stdio=%+v\n http =%+v", stdioRes, httpRes)
	}
	if !reflect.DeepEqual(stdioRes, sseRes) {
		t.Errorf("stdio vs sse mismatch:\n stdio=%+v\n sse  =%+v", stdioRes, sseRes)
	}
	if stdioRes.Stdout != "parity\n" {
		t.Errorf("unexpected stdout: %q", stdioRes.Stdout)
	}
}
