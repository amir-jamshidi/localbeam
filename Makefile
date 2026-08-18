BINARY_NAME=localbeam
VERSION=1.0.0
BUILD_DIR=build
CMD_PATH=./cmd/localbeam

.PHONY: all build build-all run install clean help

all: build

## build: Build binary for current platform
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY_NAME) $(CMD_PATH)
	@echo "✅ Done: ./$(BINARY_NAME)"

## build-all: Build for Linux, macOS, Windows
build-all:
	@mkdir -p $(BUILD_DIR)
	@echo "🔨 Building for all platforms..."
	GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64    $(CMD_PATH)
	GOOS=linux   GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64    $(CMD_PATH)
	GOOS=darwin  GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64   $(CMD_PATH)
	GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64   $(CMD_PATH)
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	@echo "✅ Binaries in ./$(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)/

## run: Build and run with default settings
run: build
	./$(BINARY_NAME)

## run-port: Run on custom port (usage: make run-port PORT=9000)
run-port: build
	./$(BINARY_NAME) --port $(PORT)

## install: Install to /usr/local/bin
install: build
	@echo "📦 Installing to /usr/local/bin/$(BINARY_NAME)..."
	@cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Installed! Run: $(BINARY_NAME)"

## init-config: Create default config file
init-config: build
	./$(BINARY_NAME) --init-config

## clean: Remove build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)

## help: Show this help
help:
	@echo "LocalBeam v$(VERSION) — Makefile commands:"
	@grep -E '^## [a-z]' Makefile | sed 's/## /  make /'
