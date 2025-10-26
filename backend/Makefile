ifneq (,$(wildcard .env))
include .env
export $(shell sed -n 's/^\([^#=]*\)=.*/\1/p' .env)
endif

.PHONY: migrate seed test build dev

migrate:
	@go run ./cmd/migrate

seed:
	@go run ./cmd/seed

test:
	@go test ./...

build:
	@go build ./...

dev:
	@./dev.sh
