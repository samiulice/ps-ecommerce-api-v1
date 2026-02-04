-- Migration: Create orders and order_items tables
-- Up Migration

-- Order Status Enum Type
DO $$ BEGIN
    CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'returned');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Payment Status Enum Type
DO $$ BEGIN
    CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed', 'refunded', 'partially_refunded');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Orders Table
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id BIGINT REFERENCES customers(id) ON DELETE SET NULL,
    customer_name VARCHAR(100) NOT NULL,
    customer_mobile VARCHAR(20) NOT NULL,
    customer_email VARCHAR(100),
    customer_area VARCHAR(255),
    customer_city VARCHAR(100),
    payment_method VARCHAR(50) NOT NULL DEFAULT 'COD',
    payment_status payment_status NOT NULL DEFAULT 'pending',
    order_status order_status NOT NULL DEFAULT 'pending',
    subtotal DECIMAL(12, 2) NOT NULL DEFAULT 0,
    shipping_cost DECIMAL(10, 2) NOT NULL DEFAULT 0,
    discount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    tax DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total DECIMAL(12, 2) NOT NULL DEFAULT 0,
    order_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    delivered_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancelled_reason TEXT
);

-- Order Items Table
CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT REFERENCES products(id) ON DELETE SET NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total_price DECIMAL(12, 2) NOT NULL DEFAULT 0,
    discount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    tax DECIMAL(10, 2) NOT NULL DEFAULT 0,
    variation_info JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for orders table
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_status ON orders(order_status);
CREATE INDEX IF NOT EXISTS idx_orders_payment_status ON orders(payment_status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number);
CREATE INDEX IF NOT EXISTS idx_orders_customer_mobile ON orders(customer_mobile);

-- Indexes for order_items table
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);

-- Trigger function to update updated_at
CREATE OR REPLACE FUNCTION update_orders_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to orders table
DROP TRIGGER IF EXISTS trigger_orders_updated_at ON orders;
CREATE TRIGGER trigger_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_orders_updated_at();

-- Function to generate order number
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS TEXT AS $$
DECLARE
    prefix TEXT := 'ORD';
    date_part TEXT;
    seq_num INT;
    order_num TEXT;
BEGIN
    date_part := TO_CHAR(NOW(), 'YYYYMMDD');
    
    -- Get count of orders today + 1
    SELECT COUNT(*) + 1 INTO seq_num
    FROM orders
    WHERE DATE(created_at) = CURRENT_DATE;
    
    order_num := prefix || '-' || date_part || '-' || LPAD(seq_num::TEXT, 4, '0');
    
    RETURN order_num;
END;
$$ LANGUAGE plpgsql;

-- Comments
COMMENT ON TABLE orders IS 'Stores customer orders';
COMMENT ON TABLE order_items IS 'Stores individual items within an order';
COMMENT ON COLUMN orders.order_number IS 'Unique order identifier for display (e.g., ORD-20260205-0001)';
COMMENT ON COLUMN orders.payment_method IS 'Payment method: COD, ONLINE, CARD, etc.';
