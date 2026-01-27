#!/bin/bash
# =========================
# Build script for Linux
# =========================
# Make it executable: chmod +x scripts/build_linux.sh
# Run: ./scripts/build_linux.sh

set -e

echo "🔹 Building API for Linux..."

# Output folder
mkdir -p ../bin

# Build binary
GOOS=linux GOARCH=amd64 go build -o ../bin/app-linux ../cmd/api

echo "Build complete: ../bin/app-linux"
