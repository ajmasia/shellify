# Shellify Makefile

# Variables
BINARY_NAME := sfy
CMD_DIR := ./cmd/sfy
BIN_DIR := bin
GUI_DIR := gui
VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOFMT := gofmt
GOMOD := $(GOCMD) mod

# Ldflags for version injection
LDFLAGS := -ldflags "-s -w \
	-X github.com/ajmasia/shellify/internal/interfaces/cli.Version=$(VERSION) \
	-X github.com/ajmasia/shellify/internal/interfaces/cli.Commit=$(COMMIT) \
	-X github.com/ajmasia/shellify/internal/interfaces/cli.BuildDate=$(BUILD_DATE)"

.PHONY: all build build-with-gui run test test-coverage lint fmt tidy clean install help \
	gui-install gui-dev gui-build gui-lint gui-check

## all: Build the binary (default target)
all: build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)"

## run: Build and run the binary with optional ARGS
run: build
	@$(BIN_DIR)/$(BINARY_NAME) $(ARGS)

## test: Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## test-race: Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	$(GOTEST) -v -race ./...

## lint: Run golangci-lint
lint:
	@test -f "$$(go env GOPATH)/bin/golangci-lint" || (echo "Installing golangci-lint v1.57.2..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2)
	@"$$(go env GOPATH)/bin/golangci-lint" run ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@test -f "$$(go env GOPATH)/bin/goimports" || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	@"$$(go env GOPATH)/bin/goimports" -w -local github.com/ajmasia/shellify .

## tidy: Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html

## install: Install binary to ~/.local/bin
install: build
	@echo "Installing $(BINARY_NAME) to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	@cp $(BIN_DIR)/$(BINARY_NAME) ~/.local/bin/$(BINARY_NAME)
	@echo "Installed: ~/.local/bin/$(BINARY_NAME)"

## version: Show version
version:
	@echo "$(VERSION)"

## release-snapshot: Build snapshot release locally
release-snapshot:
	@test -f "$$(go env GOPATH)/bin/goreleaser" || (echo "Installing goreleaser..." && go install github.com/goreleaser/goreleaser@latest)
	@"$$(go env GOPATH)/bin/goreleaser" build --snapshot --clean

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

# GUI targets
## gui-install: Install GUI dependencies
gui-install:
	@echo "Installing GUI dependencies..."
	cd $(GUI_DIR) && npm ci

## gui-dev: Start GUI development server
gui-dev:
	@echo "Starting GUI development server..."
	cd $(GUI_DIR) && npm run dev

## gui-build: Build GUI for production
gui-build:
	@echo "Building GUI..."
	cd $(GUI_DIR) && npm run build

## gui-lint: Run GUI linter
gui-lint:
	@echo "Running GUI linter..."
	cd $(GUI_DIR) && npm run lint

## gui-check: Run GUI lint and typecheck
gui-check:
	@echo "Running GUI checks..."
	cd $(GUI_DIR) && npm run lint && npm run typecheck

## build-with-gui: Build binary with embedded GUI
build-with-gui: gui-build
	@echo "Building $(BINARY_NAME) v$(VERSION) with embedded GUI..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -tags embed_gui $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME) (with GUI)"

## dev: Start development (API server with GUI dev mode)
dev:
	@echo "Start API server with: make run ARGS='server'"
	@echo "Start GUI dev server with: make gui-dev"
