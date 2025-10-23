package main

import (
	"log"

	"github.com/tanydotai/tanyai/backend/internal/config"
	"github.com/tanydotai/tanyai/backend/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := db.RunMigrations("migrations", cfg.PostgresURL); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	log.Println("database migrations applied")
}
