package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mcp-try/internal/config"
	"mcp-try/internal/logger"
	"mcp-try/internal/mcp"
	"mcp-try/internal/metrics"
	"mcp-try/internal/middleware"
	"mcp-try/internal/storage"
)

type Server struct {
	cfg   *config.Config
	srv   *http.Server
	store storage.Store
}

func New(cfg *config.Config) (*Server, error) {
	var store storage.Store
	var err error

	if cfg.Database.Path != "" {
		store, err = storage.NewSQLiteStore(cfg.Database.Path)
		if err != nil {
			logger.Warn(context.Background(), "failed to init SQLite, using memory store", "error", err)
			store = storage.NewMemoryStore()
		}
	} else {
		store = storage.NewMemoryStore()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", healthHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ready", readyHandler)
	mux.Handle("/metrics", metrics.Handler())

	mcpHandler := mcp.NewHandler(store)
	mux.Handle("/mcp", mcpHandler)

	handler := middleware.Chain(
		mux,
		middleware.Recovery,
		middleware.RequestID,
		middleware.Logging,
		middleware.SecurityHeaders,
		middleware.CORS(cfg.Security),
		middleware.APIKeyAuth(cfg.Security),
		middleware.MaxBodySize(cfg.Security.MaxRequestSize),
		middleware.RateLimit(cfg.Security.RateLimitRPS),
		middleware.Timeout(cfg.Server.WriteTimeout),
	)

	srv := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{cfg: cfg, srv: srv, store: store}, nil
}

func (s *Server) Start() error {
	ctx := context.Background()

	errCh := make(chan error, 1)
	go func() {
		logger.Info(ctx, "server starting", "addr", s.cfg.Server.Port)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-quit:
		logger.Info(ctx, "shutdown signal received", "signal", sig.String())
	}

	return s.Shutdown()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.ShutdownTimeout)
	defer cancel()

	logger.Info(ctx, "shutting down server")

	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, "server shutdown error", "error", err)
		return err
	}

	if s.store != nil {
		s.store.Close()
	}

	logger.Info(ctx, "server stopped gracefully")
	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok","service":"mcp-s2xd"}`)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"ready":true}`)
}
