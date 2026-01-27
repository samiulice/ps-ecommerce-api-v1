#!/bin/bash
# =========================
# Build script for macOS
# =========================
# Make it executable: chmod +x scripts/build_mac.sh
# Run: ./scripts/build_mac.sh

set -e

echo "🔹 Building API for macOS..."

# Output folder
mkdir -p ../bin

# Build binary
GOOS=darwin GOARCH=amd64 go build -o ../bin/app-mac ../cmd/api

echo "Build complete: ../bin/app-mac"
