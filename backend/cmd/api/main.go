package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/tanydotai/tanyai/backend/internal/config"
	"github.com/tanydotai/tanyai/backend/internal/db"
	"github.com/tanydotai/tanyai/backend/internal/server"
)

func main() {
	logFile, err := setupLogger()
	if err != nil {
		log.Fatalf("setup logger: %v", err)
	}
	if logFile != nil {
		defer logFile.Close()
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database, err := db.Open(dbCtx, cfg.PostgresURL, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	srv, err := server.New(database, cfg)
	if err != nil {
		log.Fatalf("init server: %v", err)
	}
	if err := srv.Run(ctx); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func setupLogger() (*os.File, error) {
	logMode := os.Getenv("LOG_MODE")
	appEnv := os.Getenv("APP_ENV")
	
	// Production atau LOG_MODE=stdout: log hanya ke stdout
	if logMode == "stdout" || appEnv == "production" {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
		slog.SetDefault(slog.New(handler))
		return nil, nil
	}

	// Development: log ke file dan stdout
	if err := os.MkdirAll("logs", 0o755); err != nil {
		return nil, err
	}

	path := filepath.Join("logs", "api.log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(io.MultiWriter(os.Stdout, file), &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	slog.SetDefault(slog.New(handler))

	return file, nil
}
