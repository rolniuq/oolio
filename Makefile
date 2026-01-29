.PHONY: help build run test lint clean migrate-up migrate-down sqlc-compile
.PHONY: docker-up docker-down docker-logs docker-dev docker-services docker-migrate docker-clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development Commands
build: ## Build the application
	go build -o bin/oolio ./cmd/main.go

run: ## Run the application
	go run ./cmd/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

# Database Commands
sqlc-compile: ## Generate SQLC code
	sqlc generate

migrate-up: ## Run database migrations up
	migrate -path migrations -database "postgresql://oolio:oolio_password@localhost:5432/oolio_db?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "postgresql://oolio:oolio_password@localhost:5432/oolio_db?sslmode=disable" down

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	migrate create -ext sql -dir migrations -seq $(NAME)

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

dev-setup: ## Setup development environment
	go mod download
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker Commands - Production
docker-up: ## Start all services (API + Web + DB + Redis)
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## Show logs for all services
	docker-compose logs -f

docker-build: ## Build Docker images
	docker-compose build

docker-restart: ## Restart all services
	docker-compose restart

# Docker Commands - Development
docker-dev: ## Start in development mode with hot reload
	docker-compose --profile dev up -d

docker-dev-down: ## Stop development services
	docker-compose --profile dev down

docker-dev-logs: ## Show development logs
	docker-compose --profile dev logs -f

# Docker Commands - Services Only
docker-services: ## Start only DB and Redis (for local development)
	docker-compose up -d db redis

docker-services-down: ## Stop DB and Redis
	docker-compose stop db redis

# Docker Commands - Migration
docker-migrate: ## Run database migrations in Docker
	docker-compose --profile migration up migrate

# Docker Commands - Cleanup
docker-clean: ## Remove all containers, volumes, and images
	docker-compose down -v
	docker system prune -f

docker-clean-all: ## Remove everything including images
	docker-compose down -v --rmi all
	docker system prune -af

# Docker Commands - Individual Services
docker-api: ## Start only API service
	docker-compose up -d api

docker-web: ## Start only Web service
	docker-compose up -d web

docker-db: ## Start only Database service
	docker-compose up -d db

docker-redis: ## Start only Redis service
	docker-compose up -d redis
