package main

import (
	"context"
	"os"

	"mcp-try/internal/config"
	"mcp-try/internal/httpserver"
	"mcp-try/internal/logger"
	"mcp-try/internal/tools"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Error(context.Background(), "config load failed", "error", err)
		os.Exit(1)
	}

	logger.Init(cfg.Logging.Level, cfg.Logging.Format)

	tools.RegisterDefaults()
	tools.RegisterProvenanceTools()
	tools.RegisterStatisticalTools()

	srv, err := httpserver.New(cfg)
	if err != nil {
		logger.Error(context.Background(), "server init failed", "error", err)
		os.Exit(1)
	}

	if err := srv.Start(); err != nil {
		logger.Error(context.Background(), "server failed", "error", err)
		os.Exit(1)
	}
}
