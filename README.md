# PS-Ecommerce API — Project Setup & Documentation

## 1️⃣ Project Structure

Following Go best practices and recommended layout:

```
pse-api-v1/
├───assets
│   ├───private
│   │   ├───images
│   │   └───pdf
│   └───public
│       ├───fonts
│       └───images
│           ├───brands
│           └───categories
├───bin
├───cmd
│   └───api
├───docs
├───internal
│   ├───config
│   ├───database
│   ├───handler
│   ├───middleware
│   ├───model
│   ├───repository
│   ├───routes
│   ├───server
│   └───service
├───migrations
├───pkg
│   └───utils
├───scripts
└───tests
    └───integration
```

## 2️⃣ Configuration

* `.env` holds environment variables.
* `internal/config/config.go` defines structured configs:

```go
type JWTConfig struct { ... }     // Access & refresh tokens
type DBConfig struct { ... }      // Postgres URL
type ServerConfig struct { ... }  // HTTP server settings
type AppConfig struct { ... }     // App metadata & feature flags
```

* `.env.example` provides a template for required variables.

## 3️⃣ Database Setup

* PostgreSQL database and user created via `000_create_database.sql`:

```sql
CREATE ROLE super_shop_dev_user LOGIN PASSWORD 'your_password';
CREATE DATABASE super_shop_dev_db OWNER super_shop_dev_user;
GRANT ALL PRIVILEGES ON DATABASE super_shop_dev_db TO super_shop_dev_user;
```

* Tables and triggers created via migration scripts:

  * `001_init.up.sql` — sets up `users` table, triggers, indexes, and extensions
  * `002_user.up.sql` — additional user-related tables (if any)

* DOWN scripts (`*.down.sql`) allow rolling back migrations safely (tables, triggers, indexes), **excluding extensions**.

## 4️⃣ Migrations

Scripts to run migrations in order:

### migrate_up.sh

```bash
#!/bin/bash
# Run all UP migrations in order
set -e
DB_URL="postgres://super_shop_dev_user:password@localhost:5432/super_shop_dev_db?sslmode=disable"
MIGRATIONS_DIR="./migrations"
for file in $(ls $MIGRATIONS_DIR/*_*.up.sql | sort); do
    echo "Applying $file"
    psql "$DB_URL" -f "$file"
done
echo "All UP migrations applied!"
```

### migrate_down.sh

```bash
#!/bin/bash
# Run all DOWN migrations in reverse order
set -e
DB_URL="postgres://super_shop_dev_user:password@localhost:5432/super_shop_dev_db?sslmode=disable"
MIGRATIONS_DIR="./migrations"
for file in $(ls $MIGRATIONS_DIR/*_*.down.sql | sort -r); do
    echo "Reverting $file"
    psql "$DB_URL" -f "$file"
done
echo "All DOWN migrations applied!"
```

* These scripts **automatically detect migration files by prefix** (001, 002, …) and run them in order.

## 5️⃣ Server Setup

* `run.sh` script to build and run the Go server:

```bash
#!/bin/bash
set -e
# Load .env
if [ -f .env ]; then
    set -o allexport
    source .env
    set +o allexport
else
    echo ".env file not found!"
    exit 1
fi

# Build and start server
go build -o ./bin/api ./cmd/api
./bin/api
```

* Loads configuration from `.env`
* Builds server binary in `./bin/api`
* Starts HTTP server with Chi router, middleware, and routes defined in `internal/routes/`

## 6️⃣ Makefile Commands

* `make all` — build binaries for all platforms (Linux, macOS, Windows, ARM64)
* `make build-linux` — build only Linux binary
* `make build-mac` — build only macOS binary
* `make build-windows` — build only Windows binary
* `make build-linux-arm64` — build Linux ARM64 binary
* `make clean` — remove all generated binaries from `./bin/`

## 7️⃣ Notes

* Postgres DB URL example for connection:

```text
postgres://super_shop_dev_user:password@localhost:5432/super_shop_dev_db?sslmode=disable
```

* Migration scripts are **idempotent** — safe to run multiple times.
* Server is fully configurable via `.env`.
* Redis integration can be added later, but the system works with Postgres-only for refresh tokens.
