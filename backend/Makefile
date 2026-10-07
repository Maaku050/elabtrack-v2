.PHONY: help run tidy fmt vet test build migrate-up migrate-down migrate-status migrate-create migrate-adopt-legacy seed sessions-cleanup

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

run: ## Run the API server
	go run ./cmd/api

tidy: ## Tidy module dependencies
	go mod tidy

fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

test: ## Run all tests
	go test ./...

build: ## Build the API binary into bin/
	go build -o bin/api ./cmd/api

migrate-up: ## Apply pending migrations
	go run ./cmd/api --migrate-up

migrate-down: ## Roll back the latest migration
	go run ./cmd/api --migrate-down

migrate-status: ## Show migration state without schema mutation
	go run ./cmd/api --migrate-status

migrate-create: ## Create a new migration: make migrate-create name=create_products
	go run ./cmd/api --migrate-create=$(name)

seed: ## Explicitly run development-only seeds (denied in test/production)
	go run ./cmd/api --seed

sessions-cleanup: ## Explicitly delete one bounded batch of terminal sessions older than seven days
	go run ./cmd/api --sessions-cleanup

migrate-adopt-legacy: ## Attest known development history and add checksums (read docs first)
	go run ./cmd/api --migrate-adopt-legacy
