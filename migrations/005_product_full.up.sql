-- Drop existing products table if exists
DROP TABLE IF EXISTS product_variations CASCADE;
DROP TABLE IF EXISTS products CASCADE;

-- Products table (converted from MySQL)
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    added_by VARCHAR(191) DEFAULT NULL,
    user_id BIGINT DEFAULT NULL,
    name VARCHAR(80) DEFAULT NULL,
    slug VARCHAR(120) DEFAULT NULL,
    product_type VARCHAR(20) NOT NULL DEFAULT 'physical',
    category_ids VARCHAR(80) DEFAULT NULL,
    category_id VARCHAR(191) DEFAULT NULL,
    sub_category_id VARCHAR(191) DEFAULT NULL,
    sub_sub_category_id VARCHAR(191) DEFAULT NULL,
    brand_id BIGINT DEFAULT NULL,
    unit VARCHAR(191) DEFAULT NULL,
    min_qty INT NOT NULL DEFAULT 1,
    refundable BOOLEAN NOT NULL DEFAULT TRUE,
    digital_product_type VARCHAR(30) DEFAULT NULL,
    digital_file_ready VARCHAR(191) DEFAULT NULL,
    digital_file_ready_storage_type VARCHAR(10) DEFAULT 'public',
    images TEXT DEFAULT NULL,
    color_image TEXT NOT NULL DEFAULT '',
    thumbnail VARCHAR(255) DEFAULT NULL,
    thumbnail_storage_type VARCHAR(10) DEFAULT 'public',
    preview_file VARCHAR(255) DEFAULT NULL,
    preview_file_storage_type VARCHAR(255) DEFAULT 'public',
    featured VARCHAR(255) DEFAULT NULL,
    flash_deal VARCHAR(255) DEFAULT NULL,
    video_provider VARCHAR(30) DEFAULT NULL,
    video_url VARCHAR(150) DEFAULT NULL,
    colors VARCHAR(150) DEFAULT NULL,
    variant_product BOOLEAN NOT NULL DEFAULT FALSE,
    attributes VARCHAR(255) DEFAULT NULL,
    choice_options TEXT DEFAULT NULL,
    variation TEXT DEFAULT NULL,
    digital_product_file_types TEXT DEFAULT NULL,
    digital_product_extensions TEXT DEFAULT NULL,
    published BOOLEAN NOT NULL DEFAULT FALSE,
    unit_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    purchase_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    tax VARCHAR(191) NOT NULL DEFAULT '0.00',
    tax_type VARCHAR(80) DEFAULT NULL,
    tax_model VARCHAR(20) NOT NULL DEFAULT 'exclude',
    discount VARCHAR(191) NOT NULL DEFAULT '0.00',
    discount_type VARCHAR(80) DEFAULT NULL,
    current_stock INT DEFAULT NULL,
    minimum_order_qty INT NOT NULL DEFAULT 1,
    details TEXT DEFAULT NULL,
    free_shipping BOOLEAN NOT NULL DEFAULT FALSE,
    attachment VARCHAR(191) DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    status BOOLEAN NOT NULL DEFAULT TRUE,
    featured_status BOOLEAN NOT NULL DEFAULT TRUE,
    meta_title VARCHAR(191) DEFAULT NULL,
    meta_description VARCHAR(191) DEFAULT NULL,
    meta_image VARCHAR(191) DEFAULT NULL,
    request_status BOOLEAN NOT NULL DEFAULT FALSE,
    denied_note TEXT DEFAULT NULL,
    shipping_cost NUMERIC(8,2) DEFAULT NULL,
    multiply_qty BOOLEAN DEFAULT NULL,
    temp_shipping_cost NUMERIC(8,2) DEFAULT NULL,
    is_shipping_cost_updated BOOLEAN DEFAULT NULL,
    code VARCHAR(191) DEFAULT NULL
);

-- Indexes for better query performance
CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_brand_id ON products(brand_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_published ON products(published);
CREATE INDEX idx_products_user_id ON products(user_id);
