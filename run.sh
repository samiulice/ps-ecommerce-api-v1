#!/bin/bash
# =========================
# Run Go REST API Server
# =========================

set -e

# 1️⃣ Load .env safely
if [ -f .env ]; then
    echo "🔹 Loading .env variables..."
    set -o allexport
    source .env
    set +o allexport
else
    echo "⚠️  .env file not found. Please create one."
    exit 1
fi

# 2️⃣ Optional: Run migrations
# echo "🔹 Running migrations..."
# ./migrate_up.sh

# 3️⃣ Build and run server
echo "🔹 Building Go server..."
go build -o ./bin/app ./cmd/api

echo "🔹 Starting server..."
./bin/app
