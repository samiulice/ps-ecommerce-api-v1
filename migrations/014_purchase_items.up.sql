CREATE TABLE IF NOT EXISTS purchase_items (
    id          BIGSERIAL PRIMARY KEY,
    purchase_id BIGINT NOT NULL,
    item_type   VARCHAR(20) NOT NULL,
    product_id  BIGINT,
    item_name   VARCHAR(255) NOT NULL,
    quantity    NUMERIC(20, 4) NOT NULL DEFAULT 0,
    unit_price  NUMERIC(20, 4) NOT NULL DEFAULT 0,
    total_price NUMERIC(20, 4) NOT NULL DEFAULT 0,
    note        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT purchase_items_purchase_id_fkey
        FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE,
    CONSTRAINT purchase_items_product_id_fkey
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL,
    CONSTRAINT purchase_items_item_type_check
        CHECK (item_type IN ('product', 'material'))
);

CREATE INDEX IF NOT EXISTS idx_purchase_items_purchase_id
    ON purchase_items(purchase_id);

CREATE INDEX IF NOT EXISTS idx_purchase_items_product_id
    ON purchase_items(product_id);
