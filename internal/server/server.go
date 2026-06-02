package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"

	"github.com/inhuman/mcp-exec/internal/config"
	"github.com/inhuman/mcp-exec/internal/exectool"
	"github.com/inhuman/mcp-exec/internal/isolator"
)

// Version is the server/tool contract version. Overridable at build time via
// -ldflags "-X github.com/inhuman/mcp-exec/internal/server.Version=...".
var Version = "v0.1.0"

// New builds an MCP server exposing the single exec tool. The same server is
// reused across every transport, so the tool set is identical everywhere.
func New(cfg config.Config, iso isolator.Isolator, log *zap.Logger) *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{Name: "mcp-exec", Version: Version}, nil)
	h := exectool.New(cfg, iso, log)
	mcp.AddTool(srv, &mcp.Tool{Name: "exec", Description: exectool.Description(cfg)}, h.Handle)
	return srv
}

// Run serves the given MCP server over the transport selected by config.
func Run(ctx context.Context, cfg config.Config, srv *mcp.Server, log *zap.Logger) error {
	switch cfg.Transport {
	case "stdio":
		log.Info("serving over stdio")
		return srv.Run(ctx, &mcp.StdioTransport{})
	case "http":
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return srv }, nil)
		mux := http.NewServeMux()
		mux.Handle("/mcp", handler)
		log.Info("serving streamable HTTP", zap.String("addr", cfg.Addr), zap.String("path", "/mcp"), zap.Bool("auth", cfg.AuthToken != ""))
		return listenAndServe(ctx, cfg.Addr, withAuth(cfg.AuthToken, mux))
	case "sse":
		handler := mcp.NewSSEHandler(func(*http.Request) *mcp.Server { return srv }, nil)
		log.Info("serving SSE", zap.String("addr", cfg.Addr), zap.Bool("auth", cfg.AuthToken != ""))
		return listenAndServe(ctx, cfg.Addr, withAuth(cfg.AuthToken, handler))
	default:
		return fmt.Errorf("unknown transport %q (want stdio|http|sse)", cfg.Transport)
	}
}

func listenAndServe(ctx context.Context, addr string, handler http.Handler) error {
	httpSrv := &http.Server{Addr: addr, Handler: handler}
	go func() {
		<-ctx.Done()
		_ = httpSrv.Close() // best-effort shutdown on context cancellation
	}()
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
