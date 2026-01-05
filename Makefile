RED				:= \033[0;31m
GREEN			:= \033[0;32m
YELLOW			:= \033[1;33m
BLUE			:= \033[0;34m
PURPLE			:= \033[0;35m
CYAN			:= \033[0;36m
NC				:= \033[0m # No Color

# Docker settings
DOCKER_IMAGE    := s21-api-client-builder
DOCKER_WORK_DIR := /app
DOCKER_CONTAINER := s21-api-client

# Go settings
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOGET := $(GOCMD) get
GOFMT := gofmt

# Build directories
BUILD_DIR := build
CMD_DIR := cmd

.PHONY: all
all: build

## build: Build the application
.PHONY: build
build:
	@echo -e "$(CYAN)Building application...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -buildvcs=false -v -o $(BUILD_DIR)/client ./$(CMD_DIR)/client
	@echo -e "$(GREEN)Build complete!$(NC)"

## test: Run all tests
.PHONY: test
test: test-mock
	@echo -e "$(GREEN)All tests passed!$(NC)"

## test-real: Run tests with real credentials (requires S21_LOGIN and S21_PASSWORD env vars)
.PHONY: test-real
test-real:
	@echo -e "$(CYAN)Running tests with real credentials...$(NC)"
	@if [ -z "$(S21_LOGIN)" ] || [ -z "$(S21_PASSWORD)" ]; then \
		echo -e "$(RED)Error: S21_LOGIN and S21_PASSWORD environment variables must be set$(NC)"; \
		exit 1; \
	fi
	$(GOTEST) -v -tags=integration ./tests/integration/...

## test-mock: Run tests with mock API
.PHONY: test-mock
test-mock:
	@echo -e "$(CYAN)Running tests with mock API...$(NC)"
	$(GOTEST) -v -tags=mock ./tests/unit/...

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo -e "$(YELLOW)Cleaning...$(NC)"
	rm -rf $(BUILD_DIR)
	@echo -e "$(GREEN)Clean complete!$(NC)"

## fmt: Format Go code
.PHONY: fmt
fmt:
	@echo -e "$(CYAN)Formatting code...$(NC)"
	$(GOFMT) -s -w .

## lint: Run linter
.PHONY: lint
lint:
	@echo -e "$(CYAN)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo -e "$(YELLOW)golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin$(NC)"; \
	fi

## deps: Download dependencies
.PHONY: deps
deps:
	@echo -e "$(CYAN)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy

## mod-verify: Verify dependencies
.PHONY: mod-verify
mod-verify:
	@echo -e "$(CYAN)Verifying dependencies...$(NC)"
	$(GOMOD) verify

## Docker targets

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	@echo -e "$(CYAN)Building Docker image...$(NC)"
	docker compose build
	@echo -e "$(GREEN)Docker image built successfully!$(NC)"

## docker-up: Start Docker container
.PHONY: docker-up
docker-up:
	@echo -e "$(CYAN)Starting Docker container...$(NC)"
	docker compose up -d
	@echo -e "$(GREEN)Docker container started!$(NC)"

## docker-down: Stop Docker container
.PHONY: docker-down
docker-down:
	@echo -e "$(YELLOW)Stopping Docker container...$(NC)"
	docker compose down
	@echo -e "$(GREEN)Docker container stopped!$(NC)"

## docker-shell: Open shell in Docker container
.PHONY: docker-shell
docker-shell:
	@echo -e "$(BLUE)Opening shell in Docker container...$(NC)"
	docker compose exec $(DOCKER_CONTAINER) sh

## docker-build-project: Build project in Docker
.PHONY: docker-build-project
docker-build-project:
	@echo -e "$(CYAN)Building project in Docker...$(NC)"
	docker compose exec $(DOCKER_CONTAINER) make build

## docker-clean: Clean project in Docker
.PHONY: docker-clean
docker-clean:
	@echo -e "$(CYAN)Cleaning project in Docker...$(NC)"
	docker compose exec $(DOCKER_CONTAINER) make clean

## docker-test: Run tests in Docker
.PHONY: docker-test
docker-test:
	@echo -e "$(CYAN)Running tests in Docker...$(NC)"
	docker compose exec $(DOCKER_CONTAINER) make test

## docker-test-mock: Run mock tests in Docker
.PHONY: docker-test-mock
docker-test-mock:
	@echo -e "$(CYAN)Running mock tests in Docker...$(NC)"
	docker compose exec $(DOCKER_CONTAINER) make test-mock

## docker-fmt: Format code in Docker
.PHONY: docker-fmt
docker-fmt:
	@echo -e "$(CYAN)Formatting code in Docker...$(NC)"
	docker compose exec $(DOCKER_CONTAINER) make fmt

## docker-lint: Run linter in Docker
.PHONY: docker-lint
docker-lint:
	@echo -e "$(CYAN)Running linter in Docker...$(NC)"
	docker compose exec $(DOCKER_CONTAINER) make lint

## docker-deps: Download dependencies in Docker
.PHONY: docker-deps
docker-deps:
	@echo -e "$(CYAN)Downloading dependencies in Docker...$(NC)"
	docker compose exec $(DOCKER_CONTAINER) make deps

## docker-all: Build and test in Docker
.PHONY: docker-all
docker-all: docker-up docker-build-project docker-test docker-down
	@echo -e "$(GREEN)All Docker operations completed!$(NC)"

## help: Show this help message
.PHONY: help
help:
	@echo -e "$(CYAN)Usage: make [target]$(NC)"
	@echo ""
	@echo -e "$(BLUE)Available targets:$(NC)"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /' | column -t -s ':'
