.PHONY: help run build test test-coverage clean docker-build docker-run lint fmt

# Variables
APP_NAME := fizzbuzz-api
GO := go
GOFLAGS := -v
DOCKER_IMAGE := $(APP_NAME):latest
PORT ?= 8080

# Help target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*##"; printf "\033[36m\033[0m"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

run: ## Run the application locally
	$(GO) run cmd/server/main.go

dev: ## Run with hot reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

build: ## Build the application binary
	$(GO) build $(GOFLAGS) -o bin/$(APP_NAME) cmd/server/main.go

clean: ## Clean build artifacts
	rm -rf bin/ dist/ tmp/

##@ Testing

test: ## Run unit tests
	$(GO) test $(GOFLAGS) ./...

test-coverage: ## Run tests with coverage report
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration: ## Run integration tests
	$(GO) test $(GOFLAGS) -tags=integration ./test/integration/...

benchmark: ## Run benchmarks
	$(GO) test -bench=. -benchmem ./...

##@ Code Quality

lint: ## Run linter (requires golangci-lint)
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

fmt: ## Format code
	$(GO) fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	$(GO) vet ./...

##@ Dependencies

deps: ## Download dependencies
	$(GO) mod download

deps-update: ## Update dependencies
	$(GO) get -u ./...
	$(GO) mod tidy

deps-clean: ## Clean module cache
	$(GO) clean -modcache

##@ Docker

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run application in Docker
	docker run -p $(PORT):$(PORT) --rm $(DOCKER_IMAGE)

docker-push: ## Push Docker image to registry
	docker tag $(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)

##@ Database

migrate-up: ## Run database migrations up
	@echo "No migrations yet"

migrate-down: ## Run database migrations down
	@echo "No migrations yet"

##@ Deployment

deploy-dev: ## Deploy to development environment
	@echo "Deploying to development..."

deploy-prod: ## Deploy to production environment
	@echo "Deploying to production..."

##@ Utilities

check: lint vet test ## Run all checks (lint, vet, test)

install-tools: ## Install development tools
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/vektra/mockery/v2@latest

.DEFAULT_GOAL := help