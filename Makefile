# Go Data API Microservices Makefile

# Variables
GO_VERSION := 1.19
APP_NAME := go-data-api-microservices
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S%z)

# Go related variables
GO_BUILD_FLAGS := -ldflags='-w -s -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}'
GO_TEST_FLAGS := -race -coverprofile=coverage.out -covermode=atomic

# Docker related variables
DOCKER_REGISTRY := docker.io
DOCKER_IMAGE := ${APP_NAME}
DOCKER_TAG := ${VERSION}

# Services
SERVICES := api-gateway data-service auth-service analytics-service

# Default target
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Development commands
.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod tidy

.PHONY: build
build: deps ## Build all services
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		go build $(GO_BUILD_FLAGS) -o bin/$$service cmd/$$service/main.go; \
	done

.PHONY: build-api-gateway
build-api-gateway: deps ## Build API Gateway service
	go build $(GO_BUILD_FLAGS) -o bin/api-gateway cmd/api-gateway/main.go

.PHONY: build-data-service
build-data-service: deps ## Build Data service
	go build $(GO_BUILD_FLAGS) -o bin/data-service cmd/data-service/main.go

.PHONY: build-auth-service
build-auth-service: deps ## Build Auth service
	go build $(GO_BUILD_FLAGS) -o bin/auth-service cmd/auth-service/main.go

.PHONY: build-analytics-service
build-analytics-service: deps ## Build Analytics service
	go build $(GO_BUILD_FLAGS) -o bin/analytics-service cmd/analytics-service/main.go

# Run commands
.PHONY: run
run: build ## Run all services (in background)
	@echo "Starting all services..."
	@./bin/api-gateway & \
	./bin/data-service & \
	./bin/auth-service & \
	./bin/analytics-service &
	@echo "All services started in background"

.PHONY: run-api-gateway
run-api-gateway: build-api-gateway ## Run API Gateway service
	./bin/api-gateway

.PHONY: run-data-service
run-data-service: build-data-service ## Run Data service
	./bin/data-service

.PHONY: run-auth-service
run-auth-service: build-auth-service ## Run Auth service
	./bin/auth-service

.PHONY: run-analytics-service
run-analytics-service: build-analytics-service ## Run Analytics service
	./bin/analytics-service

# Testing commands
.PHONY: test
test: ## Run all tests
	go test $(GO_TEST_FLAGS) ./...

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: benchmark
benchmark: ## Run benchmarks
	go test -bench=. -benchmem ./...

# Code quality commands
.PHONY: lint
lint: ## Run linter
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

# Docker commands
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) .
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest .

.PHONY: docker-push
docker-push: docker-build ## Push Docker image
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest

.PHONY: docker-compose-up
docker-compose-up: ## Start services with docker-compose
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## Stop services with docker-compose
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

# Database commands
.PHONY: db-migrate-up
db-migrate-up: ## Run database migrations up
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/data_api?sslmode=disable" up

.PHONY: db-migrate-down
db-migrate-down: ## Run database migrations down
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/data_api?sslmode=disable" down

# Clean commands
.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

.PHONY: clean-docker
clean-docker: ## Clean Docker images
	docker rmi $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) || true
	docker rmi $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest || true

# Development setup
.PHONY: setup
setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env from .env.example"; fi
	@mkdir -p bin
	@echo "Development environment setup complete"

.PHONY: all
all: clean deps fmt vet lint test build ## Run all checks and build
