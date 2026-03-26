#!/bin/bash
# =========================
# Run all UP migrations in order
# =========================

# Exit on any error
set -e

# Ask for DB URL
echo "Enter your PostgreSQL DB URL:"
read -r DB_URL

# Optional: validate input
if [ -z "$DB_URL" ]; then
    echo "❌ DB URL cannot be empty!"
    exit 1
fi

# Migrations folder
MIGRATIONS_DIR="./"

echo "🔹 Running UP migrations from $MIGRATIONS_DIR..."

# Loop through all *.up.sql files sorted by prefix
for file in $(find "$MIGRATIONS_DIR" -name "*_*.up.sql" | sort); do
    echo "➡ Applying $file"
    psql "$DB_URL" -f "$file"
done

echo "✅ All UP migrations applied!"