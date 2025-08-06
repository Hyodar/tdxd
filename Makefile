# TDXS Makefile

# Build variables
BINARY_NAME := tdxs
CMD_PATH := ./cmd/$(BINARY_NAME)
BUILD_DIR := build
INSTALL_PREFIX := /usr/local
INSTALL_BIN := $(INSTALL_PREFIX)/bin

# Version info
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")

# Build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT)"
GO_BUILD := go build $(LDFLAGS)

# Constellation variables
CONSTELLATION_REPO := https://github.com/edgelesssys/constellation.git
CONSTELLATION_VERSION := v2.23.1
CONSTELLATION_INTERNAL := internal/constellation

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: all build install clean test sync-constellation help

## Default target
all: sync-constellation build

## Build the tdxs binary
build: sync-constellation
	@echo "$(CYAN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## Install the tdxs binary to system
install: build
	@echo "$(CYAN)Installing $(BINARY_NAME) to $(INSTALL_BIN)...$(NC)"
	@sudo mkdir -p $(INSTALL_BIN)
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_BIN)/
	@sudo chmod 755 $(INSTALL_BIN)/$(BINARY_NAME)
	@echo "$(GREEN)Installation complete$(NC)"

## Uninstall the tdxs binary
uninstall:
	@echo "$(YELLOW)Uninstalling $(BINARY_NAME) from $(INSTALL_BIN)...$(NC)"
	@sudo rm -f $(INSTALL_BIN)/$(BINARY_NAME)
	@echo "$(GREEN)Uninstall complete$(NC)"

## Sync Constellation internal packages
sync-constellation:
	@echo "$(CYAN)Syncing Constellation internal packages ($(CONSTELLATION_VERSION))...$(NC)"
	@if [ ! -d "$(CONSTELLATION_INTERNAL)" ]; then \
		echo "$(YELLOW)Cloning Constellation repository to temp directory...$(NC)"; \
		rm -rf ./tmp/constellation $(CONSTELLATION_INTERNAL); \
		mkdir -p ./tmp; \
		git clone --depth 1 --branch $(CONSTELLATION_VERSION) \
			$(CONSTELLATION_REPO) ./tmp/constellation; \
		echo "$(CYAN)Extracting internal directory...$(NC)"; \
		mkdir -p $(CONSTELLATION_INTERNAL); \
		cp -r ./tmp/constellation/internal/* $(CONSTELLATION_INTERNAL)/; \
		echo "$(CYAN)Cleaning up temp directory...$(NC)"; \
		rm -rf ./tmp/constellation; \
		echo "$(GREEN)Constellation internal packages extracted$(NC)"; \
		echo "$(CYAN)Updating package imports...$(NC)"; \
		find $(CONSTELLATION_INTERNAL) -type f -name "*.go" ! -name "*.pb.go" -exec sed -i 's|github.com/edgelesssys/constellation/v2/internal|github.com/Hyodar/tdxs/internal/constellation|g' {} +; \
		echo "$(GREEN)Package imports updated$(NC)"; \
		go mod tidy; \
	else \
		echo "$(YELLOW)Constellation internal packages already exist. Skipping...$(NC)"; \
	fi

## Clean build artifacts and dependencies
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)Clean complete$(NC)"

## Clean temporary directories
clean-temp:
	@echo "$(YELLOW)Cleaning temporary directories...$(NC)"
	@rm -rf ./tmp
	@echo "$(GREEN)Temp directories cleaned$(NC)"

## Deep clean including Constellation packages
clean-all: clean clean-temp
	@echo "$(RED)Removing Constellation internal packages...$(NC)"
	@rm -rf $(CONSTELLATION_INTERNAL)
	@echo "$(GREEN)Deep clean complete$(NC)"

## Run tests
test:
	@echo "$(CYAN)Running tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)Tests complete$(NC)"

## Run tests with coverage
test-coverage:
	@echo "$(CYAN)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## Format code
fmt:
	@echo "$(CYAN)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Format complete$(NC)"

## Lint code
lint:
	@echo "$(CYAN)Linting code...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not installed, skipping...$(NC)"; \
	fi
	@echo "$(GREEN)Lint complete$(NC)"

## Generate dependencies
deps:
	@echo "$(CYAN)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

## Run the daemon locally
run: build
	@echo "$(CYAN)Running $(BINARY_NAME)...$(NC)"
	$(BUILD_DIR)/$(BINARY_NAME) start --config config.toml --verbose

## Show version information
version:
	@echo "$(CYAN)Version Information:$(NC)"
	@echo "  Version: $(VERSION)"
	@echo "  Commit:  $(COMMIT)"
	@echo "  Date:    $(BUILD_DATE)"

## Display help message
help:
	@echo "$(CYAN)TDXS Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)Usage:$(NC)"
	@echo "  make [target]"
	@echo ""
	@echo "$(YELLOW)Targets:$(NC)"
	@awk '/^##/ { \
		getline target; \
		gsub(/^[^:]*:/, "", target); \
		gsub(/^## /, "", $$0); \
		printf "  $(CYAN)%-20s$(NC) %s\n", target, $$0 \
	}' $(MAKEFILE_LIST) | grep -v 'MAKEFILE_LIST'
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make                 # Sync constellation and build"
	@echo "  make build          # Build the binary"
	@echo "  make install        # Build and install to system"
	@echo "  make test           # Run tests"
	@echo "  make clean-all      # Remove all artifacts including deps"