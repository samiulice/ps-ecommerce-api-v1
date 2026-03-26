-- Migration: Create suppliers table (PostgreSQL)

CREATE TABLE IF NOT EXISTS suppliers (
    id                  BIGSERIAL PRIMARY KEY,
    supplier_code       VARCHAR(50) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    company_name        VARCHAR(255),
    contact_person      VARCHAR(255),
    phone               VARCHAR(25) NOT NULL,
    email               VARCHAR(255),
    website             VARCHAR(255),
    tax_id              VARCHAR(100),
    trade_license_no    VARCHAR(100),
    payment_terms       VARCHAR(120),
    credit_limit        NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    outstanding_balance NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    lead_time_days      INTEGER,
    rating              NUMERIC(3, 2),
    street_address      VARCHAR(250),
    country             VARCHAR(50),
    city                VARCHAR(50),
    zip                 VARCHAR(20),
    notes               TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT suppliers_supplier_code_unique UNIQUE (supplier_code),
    CONSTRAINT suppliers_phone_unique UNIQUE (phone),
    CONSTRAINT suppliers_email_unique UNIQUE (email),
    CONSTRAINT suppliers_rating_range CHECK (rating IS NULL OR (rating >= 0 AND rating <= 5))
);

CREATE INDEX IF NOT EXISTS idx_suppliers_name ON suppliers(name);
CREATE INDEX IF NOT EXISTS idx_suppliers_is_active ON suppliers(is_active);
CREATE INDEX IF NOT EXISTS idx_suppliers_country_city ON suppliers(country, city);
CREATE INDEX IF NOT EXISTS idx_suppliers_created_at ON suppliers(created_at DESC);
