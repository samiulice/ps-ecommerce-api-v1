-- Migration: Drop Delivery Management tables
-- Down Migration

DROP TABLE IF EXISTS withdraw_requests CASCADE;
DROP TABLE IF EXISTS delivery_wallets CASCADE;
DROP TABLE IF EXISTS order_deliveries CASCADE;
DROP TABLE IF EXISTS delivery_men CASCADE;
DROP TABLE IF EXISTS delivery_methods CASCADE;

DROP FUNCTION IF EXISTS update_delivery_tables_updated_at() CASCADE;