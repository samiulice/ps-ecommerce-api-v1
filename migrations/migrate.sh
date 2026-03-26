#!/bin/bash
# =========================
# PostgreSQL Migration Tool with Numeric Confirmation
# =========================

set -e

MIGRATIONS_DIR="./"

# Step 1: Ask for DB URL
echo "Enter your PostgreSQL DB URL:"
read -r DB_URL

if [ -z "$DB_URL" ]; then
    echo "❌ DB URL cannot be empty!"
    exit 1
fi

# Step 2: Test DB connection
echo "🔍 Checking database connection..."
if ! psql "$DB_URL" -c "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ Failed to connect to database!"
    exit 1
fi
echo "✅ Database connection successful!"

# Step 3: Create migration_history table if not exists
psql "$DB_URL" -c "
CREATE TABLE IF NOT EXISTS migration_history (
    id SERIAL PRIMARY KEY,
    filename TEXT UNIQUE NOT NULL,
    direction TEXT NOT NULL,
    applied_at TIMESTAMP DEFAULT now()
);
" > /dev/null

# Step 4: Ask global confirmation once (1 = yes, 2 = no)
echo "Enable global auto-confirm for all migrations?"
echo "1) Yes"
echo "2) No"
read -r GLOBAL_CONFIRM

if [[ "$GLOBAL_CONFIRM" != "1" && "$GLOBAL_CONFIRM" != "2" ]]; then
    echo "❌ Invalid selection, exiting."
    exit 1
fi

# Step 5: Main menu loop
while true; do
    echo ""
    echo "====== Migration Menu ======"
    echo "1. Run UP migrations"
    echo "2. Run DOWN migrations"
    echo "0. Exit"
    echo "============================"
    echo -n "Select an option: "
    read -r OPTION

    case $OPTION in
        1)
            echo "🔹 Running UP migrations..."
            for file in $(find "$MIGRATIONS_DIR" -name "*_*.up.sql" | sort); do
                # Skip if already applied
                APPLIED=$(psql "$DB_URL" -tAc "SELECT 1 FROM migration_history WHERE filename='$(basename "$file")' AND direction='UP';")
                if [ "$APPLIED" = "1" ]; then
                    echo "⏭ Skipping $file (already applied)"
                    continue
                fi

                # Confirm if global = no (2)
                if [ "$GLOBAL_CONFIRM" = "2" ]; then
                    echo "⚠️ Apply $file?"
                    echo "1) Yes"
                    echo "2) No"
                    read -r CONFIRM
                    if [ "$CONFIRM" != "1" ]; then
                        echo "⏭ Skipped $file"
                        continue
                    fi
                fi

                # Apply migration
                echo "➡ Applying $file"
                psql "$DB_URL" -f "$file"

                # Record in migration_history
                psql "$DB_URL" -c "INSERT INTO migration_history(filename,direction) VALUES('$(basename "$file")','UP');"
            done
            echo "✅ All UP migrations done!"
            ;;

        2)
            echo "🔹 Running DOWN migrations..."
            for file in $(find "$MIGRATIONS_DIR" -name "*_*.down.sql" | sort -r); do
                # Skip if UP migration was never applied
                APPLIED=$(psql "$DB_URL" -tAc "SELECT 1 FROM migration_history WHERE filename='$(basename "$file")' AND direction='UP';")
                if [ "$APPLIED" != "1" ]; then
                    echo "⏭ Skipping $file (not applied before)"
                    continue
                fi

                # Confirm if global = no (2)
                if [ "$GLOBAL_CONFIRM" = "2" ]; then
                    echo "⚠️ Revert $file?"
                    echo "1) Yes"
                    echo "2) No"
                    read -r CONFIRM
                    if [ "$CONFIRM" != "1" ]; then
                        echo "⏭ Skipped $file"
                        continue
                    fi
                fi

                # Apply DOWN migration
                echo "➡ Reverting $file"
                psql "$DB_URL" -f "$file"

                # Remove from migration_history
                psql "$DB_URL" -c "DELETE FROM migration_history WHERE filename='$(basename "$file")' AND direction='UP';"
            done
            echo "✅ All DOWN migrations done!"
            ;;

        0)
            echo "👋 Exiting..."
            exit 0
            ;;

        *)
            echo "❌ Invalid option! Please choose 0, 1, or 2."
            ;;
    esac
done