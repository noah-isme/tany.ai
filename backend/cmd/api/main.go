package main

import (
	"log"

	"github.com/tanydotai/tanyai/backend/internal/server"
)

func main() {
	srv := server.New()
	if err := srv.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
