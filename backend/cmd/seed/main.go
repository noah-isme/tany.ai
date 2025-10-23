package main

import (
	"context"
	"log"
	"time"

	"github.com/tanydotai/tanyai/backend/internal/config"
	"github.com/tanydotai/tanyai/backend/internal/db"
	"github.com/tanydotai/tanyai/backend/internal/seed"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	database, err := db.Open(ctx, cfg.PostgresURL, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	if err := seed.Seed(ctx, database); err != nil {
		log.Fatalf("seed database: %v", err)
	}

	log.Println("database seeded with sample data")
}
