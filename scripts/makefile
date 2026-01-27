# =========================
# Makefile Usage
# =========================

# Build all binaries (Linux, macOS, Windows, ARM64):
#   $ make all
#   Output: ./bin/app-linux, ./bin/app-mac, ./bin/app-windows.exe, ./bin/app-linux-arm64

# Build a specific platform:
#   $ make build-linux       # Linux x86_64
#   $ make build-mac         # macOS x86_64
#   $ make build-windows     # Windows x86_64
#   $ make build-linux-arm64 # Linux ARM64

# Clean binaries:
#   $ make clean             # Removes all ./bin/app-* files

# Notes:
#   - Override arch: make build-linux GOARCH=arm64
#   - Binaries are placed in ./bin/
#   - Cross-compilation works on one machine


# =========================
# Go REST API Build Makefile
# =========================

# Project paths
CMD_DIR := ./cmd/api
BIN_DIR := ./bin

# Default Go environment
GO := go
GOOS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH ?= amd64

# Binary names
BINARY_LINUX := $(BIN_DIR)/app-linux
BINARY_MAC := $(BIN_DIR)/app-mac
BINARY_WINDOWS := $(BIN_DIR)/app-windows.exe

# =========================
# Create bin folder if missing
# =========================
$(shell mkdir -p $(BIN_DIR))

# =========================
# Build Targets
# =========================

.PHONY: all
all: build-linux build-mac build-windows
	@echo "All builds complete!"

# Linux x86_64
.PHONY: build-linux
build-linux:
	@echo "🔹 Building Linux binary..."
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY_LINUX) $(CMD_DIR)
	@echo "Linux build: $(BINARY_LINUX)"

# macOS x86_64
.PHONY: build-mac
build-mac:
	@echo "🔹 Building macOS binary..."
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BINARY_MAC) $(CMD_DIR)
	@echo "macOS build: $(BINARY_MAC)"

# Windows x86_64
.PHONY: build-windows
build-windows:
	@echo "🔹 Building Windows binary..."
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY_WINDOWS) $(CMD_DIR)
	@echo "Windows build: $(BINARY_WINDOWS)"

# Optional: ARM builds (Linux)
.PHONY: build-linux-arm64
build-linux-arm64:
	@echo "🔹 Building Linux ARM64 binary..."
	GOOS=linux GOARCH=arm64 $(GO) build -o $(BIN_DIR)/app-linux-arm64 $(CMD_DIR)
	@echo "Linux ARM64 build: $(BIN_DIR)/app-linux-arm64"

# Clean binaries
.PHONY: clean
clean:
	@echo "Cleaning binaries..."
	rm -f $(BIN_DIR)/app-*
	@echo "Clean complete"
