package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/bridge"
	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var programLevel = new(slog.LevelVar)
	switch strings.ToLower(cfg.System.LogLevel) {
	case "debug":
		programLevel.Set(slog.LevelDebug)
	case "warn":
		programLevel.Set(slog.LevelWarn)
	case "error":
		programLevel.Set(slog.LevelError)
	default:
		programLevel.Set(slog.LevelInfo)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: programLevel,
	}))
	slog.SetDefault(logger)

	b, err := bridge.NewBridge(cfg)
	if err != nil {
		slog.Error("Failed to initialize bridge", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- b.Start(ctx)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		slog.Info("Received termination signal, shutting down...")
		cancel()
	case err := <-errChan:
		slog.Error("Bridge stopped unexpectedly", "error", err)
		os.Exit(1)
	}

	// Give some time for cleanup if needed
	time.Sleep(100 * time.Millisecond)
	slog.Info("Shutdown complete")
}
