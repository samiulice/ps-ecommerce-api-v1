-- Migration: Create purchase order, purchase, and purchase return tables (PostgreSQL)

CREATE TABLE IF NOT EXISTS purchase_orders (
    id            BIGSERIAL PRIMARY KEY,
    order_date    DATE NOT NULL,
    due_date      DATE,
    prefix_code   VARCHAR(255),
    count_id      VARCHAR(255),
    order_code    VARCHAR(255) NOT NULL,
    order_status  VARCHAR(255) NOT NULL,
    party_id      BIGINT NOT NULL,
    state_id      BIGINT,
    note          TEXT,
    round_off     NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    grand_total   NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    paid_amount   NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    currency_id   BIGINT,
    exchange_rate NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    created_by    BIGINT,
    updated_by    BIGINT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT purchase_orders_order_code_unique UNIQUE (order_code),
    CONSTRAINT purchase_orders_party_id_fkey FOREIGN KEY (party_id) REFERENCES suppliers(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_purchase_orders_order_date ON purchase_orders(order_date DESC);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_party_id ON purchase_orders(party_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_status ON purchase_orders(order_status);

CREATE TABLE IF NOT EXISTS purchases (
    id                               BIGSERIAL PRIMARY KEY,
    purchase_date                    DATE NOT NULL,
    prefix_code                      VARCHAR(255),
    count_id                         VARCHAR(255),
    purchase_code                    VARCHAR(255) NOT NULL,
    reference_no                     VARCHAR(255),
    purchase_order_id                BIGINT,
    party_id                         BIGINT NOT NULL,
    state_id                         BIGINT,
    carrier_id                       BIGINT,
    note                             TEXT,
    shipping_charge                  NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    is_shipping_charge_distributed   BOOLEAN NOT NULL DEFAULT FALSE,
    round_off                        NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    grand_total                      NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    change_return                    INTEGER,
    paid_amount                      NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    currency_id                      BIGINT,
    exchange_rate                    NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    created_by                       BIGINT,
    updated_by                       BIGINT,
    created_at                       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT purchases_purchase_code_unique UNIQUE (purchase_code),
    CONSTRAINT purchases_party_id_fkey FOREIGN KEY (party_id) REFERENCES suppliers(id) ON DELETE RESTRICT,
    CONSTRAINT purchases_purchase_order_id_fkey FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_purchases_purchase_date ON purchases(purchase_date DESC);
CREATE INDEX IF NOT EXISTS idx_purchases_party_id ON purchases(party_id);
CREATE INDEX IF NOT EXISTS idx_purchases_purchase_order_id ON purchases(purchase_order_id);

CREATE TABLE IF NOT EXISTS purchase_return (
    id            BIGSERIAL PRIMARY KEY,
    return_date   DATE NOT NULL,
    prefix_code   VARCHAR(255),
    count_id      VARCHAR(255),
    return_code   VARCHAR(255) NOT NULL,
    reference_no  VARCHAR(255),
    party_id      BIGINT NOT NULL,
    state_id      BIGINT,
    note          TEXT,
    round_off     NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    grand_total   NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    paid_amount   NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    currency_id   BIGINT,
    exchange_rate NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    created_by    BIGINT,
    updated_by    BIGINT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT purchase_return_return_code_unique UNIQUE (return_code),
    CONSTRAINT purchase_return_party_id_fkey FOREIGN KEY (party_id) REFERENCES suppliers(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_purchase_return_date ON purchase_return(return_date DESC);
CREATE INDEX IF NOT EXISTS idx_purchase_return_party_id ON purchase_return(party_id);
