-- DOWN MIGRATION: 024_accounting_and_finance.down.sql

-- 1. Drop View
DROP VIEW IF EXISTS view_daily_financial_report;

-- 2. Drop Triggers
DROP TRIGGER IF EXISTS trigger_log_order_transaction ON orders;
DROP TRIGGER IF EXISTS trigger_log_pos_sale_transaction ON pos_sales;
DROP TRIGGER IF EXISTS trigger_log_purchase_transaction ON purchases;
DROP TRIGGER IF EXISTS trigger_log_purchase_return_transaction ON purchase_return;
DROP TRIGGER IF EXISTS trigger_log_operational_expense ON expenses;
DROP TRIGGER IF EXISTS trigger_log_delivery_expense ON order_deliveries;

-- 3. Drop Trigger Functions
DROP FUNCTION IF EXISTS log_order_transaction();
DROP FUNCTION IF EXISTS log_pos_sale_transaction();
DROP FUNCTION IF EXISTS log_purchase_transaction();
DROP FUNCTION IF EXISTS log_purchase_return_transaction();
DROP FUNCTION IF EXISTS log_operational_expense();
DROP FUNCTION IF EXISTS log_delivery_expense();

-- 4. Drop Table and Type
DROP TABLE IF EXISTS account_transactions;
DROP TYPE IF EXISTS transaction_type;