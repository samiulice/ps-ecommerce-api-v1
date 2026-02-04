-- Migration: Drop orders and order_items tables
-- Down Migration

-- Drop triggers first
DROP TRIGGER IF EXISTS trigger_orders_updated_at ON orders;

-- Drop function
DROP FUNCTION IF EXISTS update_orders_updated_at();
DROP FUNCTION IF EXISTS generate_order_number();

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

-- Drop enum types
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS order_status;
