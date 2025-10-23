package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/tanydotai/tanyai/backend/internal/config"
	"github.com/tanydotai/tanyai/backend/internal/db"
	"github.com/tanydotai/tanyai/backend/internal/server"
)

func main() {
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
