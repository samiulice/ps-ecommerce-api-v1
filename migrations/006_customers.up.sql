-- Migration: Create customers table (PostgreSQL)
-- Converted from MySQL schema

CREATE TABLE IF NOT EXISTS customers (
    id                      BIGSERIAL PRIMARY KEY,
    name                    VARCHAR(80),
    f_name                  VARCHAR(255),
    l_name                  VARCHAR(255),
    phone                   VARCHAR(25) NOT NULL,
    image                   VARCHAR(30) NOT NULL DEFAULT 'def.png',
    email                   VARCHAR(255),
    email_verified_at       TIMESTAMPTZ,
    password                VARCHAR(80) NOT NULL,
    remember_token          VARCHAR(100),
    created_at              TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    street_address          VARCHAR(250),
    country                 VARCHAR(50),
    city                    VARCHAR(50),
    zip                     VARCHAR(20),
    house_no                VARCHAR(50),
    apartment_no            VARCHAR(50),
    cm_firebase_token       VARCHAR(191),
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    payment_card_last_four  VARCHAR(191),
    payment_card_brand      VARCHAR(191),
    payment_card_fawry_token TEXT,
    login_medium            VARCHAR(191),
    social_id               VARCHAR(191),
    is_phone_verified       BOOLEAN NOT NULL DEFAULT FALSE,
    temporary_token         VARCHAR(191),
    is_email_verified       BOOLEAN NOT NULL DEFAULT FALSE,
    wallet_balance          NUMERIC(8, 2),
    loyalty_point           NUMERIC(18, 4) DEFAULT 0.0000,
    login_hit_count         SMALLINT NOT NULL DEFAULT 0,
    is_temp_blocked         BOOLEAN NOT NULL DEFAULT FALSE,
    temp_block_time         TIMESTAMPTZ,
    referral_code           VARCHAR(255),
    referred_by             INTEGER,
    app_language            VARCHAR(191) NOT NULL DEFAULT 'en',

    CONSTRAINT customers_email_unique UNIQUE (email)
);

-- Index for phone lookups (commonly queried)
CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);

-- Index for referral code lookups
CREATE INDEX IF NOT EXISTS idx_customers_referral_code ON customers(referral_code);

