#!/bin/bash
# =========================
# Run all DOWN migrations in reverse order
# =========================

# Exit on any error
set -e

# Database URL (change if needed)
DB_URL="postgres://super_shop_dev_user:QmaDNHGpVtdD8sCv40MIvZFono48XZrW@localhost:5432/super_shop_dev_db?sslmode=disable"

# Migrations folder
MIGRATIONS_DIR="./"

echo "🔹 Running DOWN migrations from $MIGRATIONS_DIR..."

# Loop through all *.down.sql files in reverse sorted order
for file in $(ls $MIGRATIONS_DIR/*_*.down.sql | sort -r); do
    echo "➡ Reverting $file"
    psql "$DB_URL" -f "$file"
done

echo "All DOWN migrations applied!"
