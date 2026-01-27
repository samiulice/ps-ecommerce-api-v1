CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    category_id BIGINT REFERENCES categories(id),
    sub_category_id BIGINT DEFAULT 0,
    sub_sub_category_id BIGINT DEFAULT 0,
    brand_id BIGINT DEFAULT 0,

    sku TEXT UNIQUE,
    unit VARCHAR(50) DEFAULT '',
    search_tags TEXT  DEFAULT '',
    thumbnail TEXT DEFAULT '',
    additional_thumbnails TEXT, -- comma-separated

    unit_price NUMERIC(12,2) NOT NULL DEFAULT 0,
    min_order_qty NUMERIC(12,2)  NOT NULL DEFAULT 1,
    current_stock_qty NUMERIC(12,2)  NOT NULL DEFAULT 0,
    stock_alert_qty NUMERIC(12,2)  NOT NULL DEFAULT 0,

    discount_type VARCHAR(20), -- percent | flat
    discount_amount NUMERIC(12,2) DEFAULT 0,

    tax_amount NUMERIC(12,2) DEFAULT 0,
    tax_calculation VARCHAR(20), -- inclusive | exclusive

    shipping_cost NUMERIC(12,2) DEFAULT 0,
    shipping_cost_type VARCHAR(20), -- static | per_qty

    has_variation BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);

CREATE TABLE IF NOT EXISTS product_variations (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES products(id) ON DELETE CASCADE,

    name TEXT NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    sku TEXT UNIQUE,
    stock INT NOT NULL DEFAULT 0,
    thumbnail TEXT,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_variations_product_id ON product_variations(product_id);
