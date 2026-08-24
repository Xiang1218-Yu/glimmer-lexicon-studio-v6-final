package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"glimmer-lexicon-studio/internal/api"
	"glimmer-lexicon-studio/internal/core"
	"glimmer-lexicon-studio/internal/ops"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger = logger.With("service", "glimmer-lexicon-studio")
	engine := core.NewEngine()
	// Assemble the module registry after the server begins listening so the
	// HTTP layer can answer readiness probes with a stable "not ready"
	// response during initialization instead of an empty module list.
	go func() {
		engine.Init()
		logger.Info("engine initialized", "status", "ready", "modules", len(engine.Modules()))
	}()
	config := ops.LoadConfig()
	server := api.New(engine, logger)
	httpServer := &http.Server{
		Addr: config.Address, Handler: server.Handler(),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second,
		WriteTimeout: 20 * time.Second, IdleTimeout: 45 * time.Second,
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go func() {
		<-ctx.Done()
		shutdown, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		_ = httpServer.Shutdown(shutdown)
	}()
	logger.Info("glimmer-lexicon-studio listening", "addr", httpServer.Addr, "status", "initializing")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
