.PHONY: help build test test-unit test-integration run docker-up docker-down clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Backend targets
build-backend: ## Build the backend application
	cd backend && go build -o bin/ecom-backend ./cmd/main.go

test: test-unit ## Run all tests

test-unit: ## Run unit tests (domain + application layers)
	cd backend && go test -v -short ./domain/... ./application/...

test-integration: ## Run integration tests (requires PostgreSQL)
	cd backend && bash test_setup.sh && go test -v ./infrastructure/... ./api/...

test-coverage: ## Run tests with coverage report
	cd backend && go test -cover ./...

run-backend: ## Run the backend server locally
	cd backend && go run cmd/main.go

# Frontend targets
install-frontend: ## Install frontend dependencies
	cd frontend && npm install

build-frontend: ## Build the frontend for production
	cd frontend && npm run build

run-frontend: ## Run the frontend development server
	cd frontend && npm run dev

# Docker targets
docker-up: ## Start all services with Docker Compose
	docker-compose up --build

docker-down: ## Stop all Docker services
	docker-compose down

docker-clean: ## Stop services and remove volumes
	docker-compose down -v

# Database targets
db-setup: ## Set up test database
	cd backend && bash test_setup.sh

# Utility targets
clean: ## Clean build artifacts
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules

lint-backend: ## Lint backend code
	cd backend && go vet ./...
	cd backend && gofmt -s -w .

lint-frontend: ## Lint frontend code
	cd frontend && npm run lint

format: ## Format all code
	cd backend && gofmt -s -w .
	cd frontend && npm run format 2>/dev/null || echo "No formatter configured"

# Development helpers
dev: ## Run backend and frontend in parallel (requires tmux or separate terminals)
	@echo "Run these in separate terminals:"
	@echo "  Terminal 1: make run-backend"
	@echo "  Terminal 2: make run-frontend"

check: test-unit lint-backend ## Run quick checks before commit
	@echo "All checks passed."

.DEFAULT_GOAL := help
