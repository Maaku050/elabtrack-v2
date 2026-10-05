.PHONY: help run tidy fmt vet test build migrate-up migrate-down migrate-create seed

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

migrate-create: ## Create a new migration: make migrate-create name=create_products
	go run ./cmd/api --migrate-create=$(name)

seed: ## Run database seeders
	go run ./cmd/api --seed
