-- =========================
-- EXTENSIONS
-- =========================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================
-- EMPLOYEES TABLE
-- =========================
CREATE TABLE IF NOT EXISTS employees (
    id            BIGSERIAL PRIMARY KEY,
    uuid          UUID NOT NULL DEFAULT uuid_generate_v4(),
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    name          TEXT NOT NULL DEFAULT '',
    mobile        VARCHAR(100) NOT NULL DEFAULT '',
    role          VARCHAR(100) NOT NULL DEFAULT '',
    branch_id INTEGER NOT NULL DEFAULT 1;
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    is_verified   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_employees_email ON employees(email);
CREATE INDEX IF NOT EXISTS idx_employees_uuid ON employees(uuid);

-- =========================
-- REFRESH TOKENS TABLE (OPTIONAL if not using Redis)
-- =========================
-- If you prefer Postgres storage instead of Redis
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          BIGSERIAL PRIMARY KEY,
    employee_id     BIGINT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at  TIMESTAMPTZ,

    CONSTRAINT uq_employee_active_token UNIQUE (employee_id, revoked_at)
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_employee_id ON refresh_tokens(employee_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);

-- INSERT INTO public.employees
-- (id, "uuid", email, password_hash, "name", mobile, "role", is_active, is_verified, created_at, updated_at)
-- VALUES(1, 'b5046bdb-7b0a-41da-9c5b-9b300b09991a'::uuid, 'noorsupermart@gmail.com', '$2a$10$Dcwf7EbwRiUDfKzuc1i8Lu7POM0BkPEhJiApCD7ZQLLhYHlt8sM3W', '', '', '', true, true, '2026-02-02 05:43:38.458', '2026-02-02 05:43:38.458');