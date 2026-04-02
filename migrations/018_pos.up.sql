-- =============================================
-- POS System Tables
-- =============================================

-- POS Sales
CREATE TABLE IF NOT EXISTS pos_sales (
    id BIGSERIAL PRIMARY KEY,
    reference_no VARCHAR(50) UNIQUE NOT NULL,
    customer_id BIGINT REFERENCES customers(id) ON DELETE SET NULL, -- optional
    branch_id BIGINT REFERENCES branches(id) ON DELETE SET NULL,
    sale_type VARCHAR(50) DEFAULT 'retail' NOT NULL,
    subtotal DECIMAL(12, 2) NOT NULL DEFAULT 0,
    discount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    tax DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total DECIMAL(12, 2) NOT NULL DEFAULT 0,
    amount_paid DECIMAL(12, 2) NOT NULL DEFAULT 0,
    payment_method VARCHAR(50) NOT NULL DEFAULT 'cash', -- cash, card, mobile_banking
    payment_status payment_status NOT NULL DEFAULT 'paid', -- typically paid for POS
    sale_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sale_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- POS Sale Items
CREATE TABLE IF NOT EXISTS pos_sale_items (
    id BIGSERIAL PRIMARY KEY,
    pos_sale_id BIGINT NOT NULL REFERENCES pos_sales(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_variation_id BIGINT REFERENCES product_variations(id) ON DELETE SET NULL, 
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(12, 2) NOT NULL,
    subtotal DECIMAL(12, 2) NOT NULL,
    tax_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total DECIMAL(12, 2) NOT NULL
);

-- Update products stock trigger for POS
CREATE OR REPLACE FUNCTION pos_sale_item_insert_trigger()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE products
    SET current_stock_qty = current_stock_qty - NEW.quantity
    WHERE id = NEW.product_id;
    
    IF NEW.product_variation_id IS NOT NULL THEN
        UPDATE product_variations
        SET stock_qty = stock_qty - NEW.quantity
        WHERE id = NEW.product_variation_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_pos_sale_item_insert
AFTER INSERT ON pos_sale_items
FOR EACH ROW
EXECUTE FUNCTION pos_sale_item_insert_trigger();