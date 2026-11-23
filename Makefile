# Single entry point for every build, test, quality and utility command of timefy.
.DEFAULT_GOAL := help
.PHONY: help dev run run-example build build-cli build-all examples test test-ci tidy deps-upgrade deps-clean-cache lint race bench cover tree fmt fmt-check check clean

LOG_DIR  := logs
BIN_NAME := timefy
BIN_DIR  := bin

# Detect OS for binary extension
ifeq ($(OS),Windows_NT)
	BIN_EXT := .exe
else
	BIN_EXT :=
endif

BIN_OUT := $(BIN_DIR)/$(BIN_NAME)$(BIN_EXT)

# ==============================================================================
# Help
# Lists every documented target (default goal).
help: ## Show this help message
# Print the heading.
	@echo "Available targets:"
# Extract "target: ## description" lines, sort them and print them in colour.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ==============================================================================
# Development
# Runs the development entry point (main/main.go) for quick reference.
dev: ## Run the development entry point (main/main.go)
# Run the development entry point.
	go run ./main/main.go

# ==============================================================================
# Running the main application
# Executes cmd/timefy/main.go, useful for development and quick testing.
run: ## Run the CLI (cmd/timefy)
# Run the CLI from source.
	go run ./cmd/timefy

# Runs one sample program from examples/ (Usage: make run-example NAME=audit).
run-example: ## Run an example (Usage: make run-example NAME=audit)
# Require NAME and ensure the example exists.
	@if [ -z "$(NAME)" ] || [ ! -d "examples/$(NAME)" ]; then echo "Usage: make run-example NAME=<audit|billing|global_team|reporting|scheduler>"; exit 1; fi
# Run the selected example.
	go run "./examples/$(NAME)"

# ==============================================================================
# Building the application
# Compiles cmd/timefy into a binary at bin/timefy (or bin/timefy.exe on Windows).
build: ## Build the CLI binary into bin/
# Ensure the output directory exists.
	@mkdir -p $(BIN_DIR)
# Compile the CLI binary.
	go build -o $(BIN_OUT) ./cmd/timefy

# Same as build, explicitly for the CLI binary.
build-cli: ## Build the CLI binary (alias of build)
# Ensure the output directory exists.
	@mkdir -p $(BIN_DIR)
# Compile the CLI binary.
	go build -o $(BIN_OUT) ./cmd/timefy

# Compiles every package (the same check CI runs), discarding the output.
build-all: ## Compile every package, as CI does
# Verbose build of all packages.
	go build -v ./...

# Compiles every sample program under examples/ to prove they still build.
examples: ## Compile all programs under examples/
# Build all example packages (output is discarded).
	go build ./examples/...

# ==============================================================================
# Module support and testing
# Runs tests across all packages in the project, showing code coverage.
test: ## Run all tests with coverage
# Run uncached tests with coverage.
	go test -count=1 -cover ./...

# Runs the CI coverage script (race detector, atomic coverage) into logs/coverage.txt.
test-ci: ## Run the CI coverage script (sh/go_deps.sh)
# Execute the per-package race and coverage script.
	bash sh/go_deps.sh

# Cleans up the module by removing unused dependencies.
tidy: ## Tidy go.mod and go.sum
# Drop unused dependencies and add missing ones.
	go mod tidy

# ==============================================================================
# Upgrading dependencies
# Updates all direct dependencies to their latest minor/patch versions,
# then re-tidies the module.
deps-upgrade: ## Upgrade dependencies to latest minor/patch
# Upgrade dependencies, including those of tests.
	go get -u -t -d -v ./...
# Re-tidy the module after the upgrade.
	go mod tidy

# Removes all items from the Go module cache.
deps-clean-cache: ## Clear the Go module cache
# Wipe the module cache.
	go clean -modcache

# ==============================================================================
# Quality
# Runs go vet across all packages.
lint: ## Run go vet on all packages
# Static checks with go vet.
	go vet ./...

# Runs all tests with the race detector enabled.
race: ## Run tests with the race detector
# Race-enabled, uncached tests.
	go test -race -count=1 ./...

# Runs all benchmarks with memory allocation reporting.
bench: ## Run benchmarks with allocation stats
# Run every benchmark.
	go test -bench=. -benchmem ./...

# Generates an HTML coverage report at logs/coverage.html.
cover: ## Generate an HTML coverage report in logs/
# Ensure the log directory exists.
	@mkdir -p $(LOG_DIR)
# Collect the coverage profile.
	go test -coverprofile=./$(LOG_DIR)/coverage.out ./...
# Render the profile as HTML.
	go tool cover -html=./$(LOG_DIR)/coverage.out -o ./$(LOG_DIR)/coverage.html

# Runs formatting check, vet and tests in one go.
check: fmt-check lint test ## Run fmt-check, lint and test
# Reached only when all prerequisites passed.
	@echo "All checks passed."

# ==============================================================================
# Formatting
# Formats all Go source files in-place using gofmt.
fmt: ## Format all Go files
# Apply gofmt to every package.
	go fmt ./...

# Checks whether all Go source files are properly formatted.
# Exits with a non-zero status if any files need formatting.
fmt-check: ## Fail if any Go file is not gofmt-formatted
# Collect unformatted files; 1) list them, 2) fail if any exist.
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not formatted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

# ==============================================================================
# Utilities
# Creates a text file representing the project's directory structure.
tree: ## Write the directory tree to logs/ and print it
# Ensure the log directory exists.
	@mkdir -p $(LOG_DIR)
# Write the tree, excluding tooling and build folders.
	tree -I ".gradle|.idea|build|logs|.vscode|.git|.github|vendor" > ./$(LOG_DIR)/tree_source_oss.txt
# Print the generated tree.
	cat ./$(LOG_DIR)/tree_source_oss.txt

# ==============================================================================
# Clean
# Removes the build directory and log directory (generated files only).
clean: ## Remove bin/ and logs/
# Delete compiled binaries.
	rm -rf $(BIN_DIR)
# Delete logs and coverage output.
	rm -rf $(LOG_DIR)
