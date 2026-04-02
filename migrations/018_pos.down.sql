-- 1. Drop Trigger and Function
DROP TRIGGER IF EXISTS trg_pos_sale_item_insert ON pos_sale_items;
DROP FUNCTION IF EXISTS pos_sale_item_insert_trigger();

-- 2. Drop Tables
DROP TABLE IF EXISTS pos_sale_items;
DROP TABLE IF EXISTS pos_sales;

-- 3. Optionally drop the Enum (uncomment if you want to remove it entirely)
DROP TYPE IF EXISTS payment_status;