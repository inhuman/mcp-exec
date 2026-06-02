package main

import (
	"context"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/inhuman/mcp-exec/internal/config"
	"github.com/inhuman/mcp-exec/internal/isolator"
	"github.com/inhuman/mcp-exec/internal/server"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	iso := isolator.NewProcess(cfg.Python)
	srv := server.New(cfg, iso, log)

	if err := server.Run(ctx, cfg, srv, log); err != nil {
		log.Fatal("server stopped", zap.Error(err))
	}
}
