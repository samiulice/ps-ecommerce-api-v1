#!/bin/bash
# =========================
# Run all UP migrations in order
# =========================

# Exit on any error
set -e

# Database URL (change if needed)
DB_URL="postgres://super_shop_dev_user:QmaDNHGpVtdD8sCv40MIvZFono48XZrW@localhost:5432/super_shop_dev_db?sslmode=disable"

# Migrations folder
MIGRATIONS_DIR="./"

echo "🔹 Running UP migrations from $MIGRATIONS_DIR..."

# Loop through all *.up.sql files sorted by prefix
for file in $(ls $MIGRATIONS_DIR/*_*.up.sql | sort); do
    echo "➡ Applying $file"
    psql "$DB_URL" -f "$file"
done

echo "All UP migrations applied!"
