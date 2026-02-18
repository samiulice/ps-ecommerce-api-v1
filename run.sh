#!/bin/bash
set -e

echo "🚀 Running in DEV mode..."

if [ -f .env ]; then
    set -o allexport
    source .env
    set +o allexport
else
    echo ".env file not found!"
    exit 1
fi

go run ./cmd/api
