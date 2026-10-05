# Root Makefile - convenience wrappers around backend/frontend tooling.
# Linux development; equivalent commands are documented in README.md.

.PHONY: help dev backend frontend migrate-up migrate-down migrate-status migrate-create compose-migrate-up seed sessions-cleanup test lint fmt vet install

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

dev: ## Run backend and frontend together (requires concurrently installed or two terminals)
	@echo "Open two terminals: 'make backend' and 'make frontend'"

backend: ## Run Go backend
	cd backend && go run ./cmd/api

frontend: ## Run Vite frontend
	cd frontend && npm run dev

migrate-up: ## Apply database migrations
	cd backend && go run ./cmd/api --migrate-up

migrate-down: ## Roll back the latest migration
	cd backend && go run ./cmd/api --migrate-down

migrate-status: ## Show migration state without schema mutation
	cd backend && go run ./cmd/api --migrate-status

migrate-create: ## Create a new migration: make migrate-create name=describe_change
	cd backend && go run ./cmd/api --migrate-create=$(name)

seed: ## Explicitly run development-only seeds (denied in test/production)
	cd backend && go run ./cmd/api --seed

sessions-cleanup: ## Explicitly delete one bounded batch of terminal sessions older than seven days
	cd backend && go run ./cmd/api --sessions-cleanup

test: ## Run all tests
	cd backend && go test ./...
	cd frontend && npm run test:run

lint: ## Lint backend and frontend
	cd backend && go vet ./...
	cd frontend && npm run lint

fmt: ## Format Go code
	cd backend && go fmt ./...

vet: ## Run go vet
	cd backend && go vet ./...

install: ## Install backend and frontend dependencies
	cd backend && go mod download
	cd frontend && npm install

compose-migrate-up: ## Explicitly migrate the local Compose database
	docker compose --profile full run --rm --build backend --migrate-up
