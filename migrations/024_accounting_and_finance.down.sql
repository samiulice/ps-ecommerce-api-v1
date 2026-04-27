-- DOWN MIGRATION: 024_accounting_and_finance.down.sql

DROP VIEW IF EXISTS view_daily_financial_report;

DROP TRIGGER IF EXISTS trigger_log_operational_expense ON expenses;
DROP FUNCTION IF EXISTS log_operational_expense();

DROP TRIGGER IF EXISTS trigger_log_purchase_transaction ON purchases;
DROP FUNCTION IF EXISTS log_purchase_transaction();

DROP TRIGGER IF EXISTS trigger_log_order_transaction ON orders;
DROP FUNCTION IF EXISTS log_order_transaction();

DROP TABLE IF EXISTS account_transactions;
DROP TYPE IF EXISTS transaction_type;
