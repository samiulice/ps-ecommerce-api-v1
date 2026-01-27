-- =========================
-- DROP TRIGGERS
-- =========================
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;

-- =========================
-- DROP FUNCTIONS
-- =========================
DROP FUNCTION IF EXISTS set_updated_at();

-- =========================
-- DROP TABLES
-- =========================
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
