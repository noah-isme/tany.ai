SHELL := /bin/bash

.PHONY: dev api web migrate seed backend-test backend-build frontend-lint frontend-test frontend-build

dev:
	@echo "Starting tany.ai API and web app..."
	@bash -c 'trap "kill 0" EXIT; \
		(cd backend && set -a && [ -f .env ] && . .env && set +a && go run ./cmd/api) & \
		(cd frontend && npm run dev)'

api:
	@cd backend && set -a && [ -f .env ] && . .env && set +a && go run ./cmd/api

web:
	@cd frontend && npm run dev

migrate:
	@cd backend && make migrate

seed:
	@cd backend && make seed

backend-test:
	@cd backend && go test ./...

backend-build:
	@cd backend && go build ./...

frontend-lint:
	@cd frontend && npm run lint

frontend-test:
	@cd frontend && npm test

frontend-build:
	@cd frontend && npm run build
