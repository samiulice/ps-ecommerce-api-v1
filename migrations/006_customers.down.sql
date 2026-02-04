-- Migration: Drop customers table


DROP INDEX IF EXISTS idx_customers_referral_code;
DROP INDEX IF EXISTS idx_customers_phone;
DROP TABLE IF EXISTS customers;
