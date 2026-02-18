-- =============================================
-- 1. Units Table
-- =============================================
CREATE TABLE IF NOT EXISTS units (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    symbol VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 2. ATTRIBUTES
-- =============================================
CREATE TABLE IF NOT EXISTS attributes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 3. Products Table
-- =============================================
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',

    category_id BIGINT REFERENCES categories(id),
    sub_category_id BIGINT DEFAULT NULL,
    sub_sub_category_id BIGINT DEFAULT NULL,
    brand_id BIGINT DEFAULT NULL,

    sku TEXT UNIQUE,
    status SMALLINT NOT NULL DEFAULT 1,

    unit_id INT REFERENCES units(id),

    tags TEXT DEFAULT '',
    thumbnail TEXT DEFAULT '',
    gallery_images TEXT[],

    unit_price NUMERIC(12,2) NOT NULL DEFAULT 0,
    purchase_price NUMERIC(12,2) DEFAULT 0,

    min_order_qty NUMERIC(12,2) NOT NULL DEFAULT 1,
    current_stock_qty NUMERIC(12,2) NOT NULL DEFAULT 0,
    stock_alert_qty NUMERIC(12,2) NOT NULL DEFAULT 0,
    total_sold NUMERIC(12,2) NOT NULL DEFAULT 0,

    discount_type VARCHAR(20) DEFAULT 'percentage',
    discount_amount NUMERIC(12,2) DEFAULT 0,

    tax_amount NUMERIC(12,2) DEFAULT 0,
    tax_type VARCHAR(20) DEFAULT 'exclusive',

    shipping_cost NUMERIC(12,2) DEFAULT 0,
    shipping_type VARCHAR(20) DEFAULT 'static',

    has_variation BOOLEAN NOT NULL DEFAULT FALSE,
    variation_attributes JSONB NOT NULL DEFAULT '{}'::jsonb,

    total_reviews BIGINT DEFAULT 0,
    avg_rating NUMERIC(3,2) DEFAULT 0,

    five_star_count BIGINT NOT NULL DEFAULT 0,
    four_star_count BIGINT NOT NULL DEFAULT 0,
    three_star_count BIGINT NOT NULL DEFAULT 0,
    two_star_count BIGINT NOT NULL DEFAULT 0,
    one_star_count BIGINT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =============================
-- PRODUCT INDEXES (Optimized)
-- =============================

CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_sub_category ON products(sub_category_id);
CREATE INDEX IF NOT EXISTS idx_products_sub_subcategory ON products(sub_sub_category_id);
CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand_id);

-- IMPORTANT: FK index (Postgres does NOT auto create)
CREATE INDEX IF NOT EXISTS idx_products_unit_id ON products(unit_id);

CREATE INDEX IF NOT EXISTS idx_products_stock ON products(current_stock_qty);

-- Sorting indexes
CREATE INDEX IF NOT EXISTS idx_products_price ON products(unit_price);
CREATE INDEX IF NOT EXISTS idx_products_total_sold ON products(total_sold);
CREATE INDEX IF NOT EXISTS idx_products_avg_rating ON products(avg_rating);
CREATE INDEX IF NOT EXISTS idx_products_total_reviews ON products(total_reviews);

CREATE INDEX IF NOT EXISTS idx_products_five_star_count ON products(five_star_count);
CREATE INDEX IF NOT EXISTS idx_products_four_star_count ON products(four_star_count);
CREATE INDEX IF NOT EXISTS idx_products_three_star_count ON products(three_star_count);
CREATE INDEX IF NOT EXISTS idx_products_two_star_count ON products(two_star_count);
CREATE INDEX IF NOT EXISTS idx_products_one_star_count ON products(one_star_count);

-- Full text search
CREATE INDEX IF NOT EXISTS idx_products_tags
ON products
USING GIN (to_tsvector('simple', tags));

-- JSONB filter optimization
CREATE INDEX IF NOT EXISTS idx_products_variation_attributes
ON products
USING GIN (variation_attributes);


-- =============================================
-- 4. PRODUCT VARIATIONS
-- =============================================
CREATE TABLE IF NOT EXISTS product_variations (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES products(id) ON DELETE CASCADE,

    variation_attributes JSONB NOT NULL DEFAULT '{}'::jsonb,

    sku TEXT NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    stock_qty NUMERIC(12,2) NOT NULL DEFAULT 0,

    thumbnail TEXT DEFAULT '',

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(sku, product_id)
);

-- Variation Indexes
CREATE INDEX IF NOT EXISTS idx_variations_product_id
ON product_variations(product_id);

CREATE INDEX IF NOT EXISTS idx_variations_attributes
ON product_variations
USING GIN (variation_attributes);


-- =============================================
-- 5. Product Reviews
-- =============================================
CREATE TABLE IF NOT EXISTS product_reviews (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customers(id),

    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    title VARCHAR(255),
    comment TEXT,

    review_images TEXT[],

    is_verified_purchase BOOLEAN DEFAULT FALSE,
    status SMALLINT DEFAULT 1,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(product_id, customer_id)
);

-- Review Indexes (Optimized)
CREATE INDEX IF NOT EXISTS idx_reviews_product ON product_reviews(product_id);

-- Important for review listing page
CREATE INDEX IF NOT EXISTS idx_reviews_product_created
ON product_reviews(product_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_reviews_status ON product_reviews(status);

CREATE INDEX IF NOT EXISTS idx_reviews_rating ON product_reviews(rating);

CREATE INDEX IF NOT EXISTS idx_reviews_customer ON product_reviews(customer_id);
