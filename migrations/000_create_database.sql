-- =================================================
-- PostgreSQL Database & User Creation Script
-- =================================================
-- How to run:
--   psql -U postgres -f 000_create_database.sql
--
-- If you see:
--   FATAL:  Peer authentication failed for user "postgres"
--
-- Run instead:
--   sudo -i -u postgres
--   psql -f 000_create_database.sql
-- =================================================

-- =========================
-- CONFIG VARIABLES
-- =========================
-- Change these in ONE place only

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'super_shop_dev_user') THEN
        CREATE ROLE super_shop_dev_user
            LOGIN
            PASSWORD 'QmaDNHGpVtdD8sCv40MIvZFono48XZrW'
            NOSUPERUSER
            NOCREATEDB
            NOCREATEROLE;
    END IF;
END
$$;

-- =========================
-- TERMINATE ACTIVE CONNECTIONS
-- =========================
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = 'super_shop_dev_db'
  AND pid <> pg_backend_pid();

-- =========================
-- DROP & RECREATE DATABASE
-- =========================
DROP DATABASE IF EXISTS super_shop_dev_db;

CREATE DATABASE super_shop_dev_db
    WITH
    OWNER = super_shop_dev_user
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TEMPLATE = template0;

-- =========================
-- GRANT PRIVILEGES
-- =========================
GRANT ALL PRIVILEGES ON DATABASE super_shop_dev_db TO super_shop_dev_user;

-- =========================
-- CONNECT TO NEW DB
-- =========================
\connect super_shop_dev_db;

-- =========================
-- SCHEMA + EXTENSIONS (OPTIONAL BUT RECOMMENDED)
-- =========================

-- Lock down public schema
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO super_shop_dev_user;

-- Required extensions for auth system
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================
-- DEFAULT PRIVILEGES
-- =========================
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO super_shop_dev_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE, SELECT ON SEQUENCES TO super_shop_dev_user;

-- =========================
-- VERIFY (OPTIONAL)
-- =========================
-- psql "postgres://super_shop_dev_user:QmaDNHGpVtdD8sCv40MIvZFono48XZrW@localhost:5432/super_shop_dev_db?sslmode=disable"
-- \l     -- list databases
-- \du    -- list roles/users
-- \dt    -- list tables
