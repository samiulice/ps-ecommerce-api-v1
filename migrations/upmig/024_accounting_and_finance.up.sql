-- UP MIGRATION: 024_accounting_and_finance.up.sql
-- Create transactions ledger and financial calculation views

-- 1. Create Transaction Ledger Table
CREATE TYPE transaction_type AS ENUM ('sale_income', 'purchase_expense', 'operational_expense', 'sale_refund', 'purchase_refund');

CREATE TABLE IF NOT EXISTS account_transactions (
    id BIGSERIAL PRIMARY KEY,
    transaction_date DATE NOT NULL DEFAULT CURRENT_DATE,
    type transaction_type NOT NULL,
    reference_id BIGINT,                  -- ID of the order, purchase, or expense
    reference_number VARCHAR(255),        -- order_number, purchase_code, etc.
    description TEXT,
    credit_amount NUMERIC(15, 2) DEFAULT 0.00, -- Money in (Income)
    debit_amount NUMERIC(15, 2) DEFAULT 0.00,  -- Money out (Expense/Purchases)
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_acc_transactions_date ON account_transactions(transaction_date DESC);
CREATE INDEX IF NOT EXISTS idx_acc_transactions_type ON account_transactions(type);

-- 2. Create Triggers to auto-fill the ledger so application code doesn't become complicated

-- A) Trigger for Orders (Income)
CREATE OR REPLACE FUNCTION log_order_transaction()
RETURNS TRIGGER AS $$
BEGIN
    -- Log when an order is paid
    IF NEW.payment_status = 'paid' AND (OLD.payment_status IS NULL OR OLD.payment_status != 'paid') THEN
        INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount)
        VALUES (DATE(NEW.created_at), 'sale_income', NEW.id, NEW.order_number, 'Income from Order ' || NEW.order_number, NEW.total, 0);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_log_order_transaction
    AFTER UPDATE OF payment_status ON orders
    FOR EACH ROW
    EXECUTE FUNCTION log_order_transaction();


-- B) Trigger for Purchases (Expense - Inventory)
CREATE OR REPLACE FUNCTION log_purchase_transaction()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.paid_amount > 0 THEN
        -- Only insert if it's a new or the paid amount changed (simplified logic)
        -- For robust systems, you update or calculate delta. Here we'll stick to new inserts for simplicity
        INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount)
        VALUES (NEW.purchase_date, 'purchase_expense', NEW.id, NEW.purchase_code, 'Payment for Purchase ' || NEW.purchase_code, 0, NEW.paid_amount);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_log_purchase_transaction
    AFTER INSERT ON purchases
    FOR EACH ROW
    EXECUTE FUNCTION log_purchase_transaction();


-- C) Trigger for Expenses (Operational)
CREATE OR REPLACE FUNCTION log_operational_expense()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount)
    VALUES (DATE(NEW.expense_date), 'operational_expense', NEW.id, 'EXP-'||NEW.id, NEW.description, 0, NEW.amount);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_log_operational_expense
    AFTER INSERT ON expenses
    FOR EACH ROW
    EXECUTE FUNCTION log_operational_expense();


-- 3. Create a Reporting View for Performance
-- Fetch Daily Income, Expense, Purchase and Net Profit instantly without joining multiple huge tables
CREATE OR REPLACE VIEW view_daily_financial_report AS
SELECT 
    transaction_date,
    SUM(CASE WHEN type IN ('sale_income') THEN credit_amount ELSE 0 END) AS total_sales_income,
    SUM(CASE WHEN type IN ('sale_refund') THEN debit_amount ELSE 0 END) AS total_sales_refunds,
    SUM(CASE WHEN type IN ('purchase_expense') THEN debit_amount ELSE 0 END) AS total_purchases,
    SUM(CASE WHEN type IN ('operational_expense') THEN debit_amount ELSE 0 END) AS total_expenses,
    -- Net Profit = (All Income) - (All Expenses & Purchases)
    (
        SUM(CASE WHEN type IN ('sale_income', 'purchase_refund') THEN credit_amount ELSE 0 END) 
        - 
        SUM(CASE WHEN type IN ('purchase_expense', 'operational_expense', 'sale_refund') THEN debit_amount ELSE 0 END)
    ) AS net_profit
FROM account_transactions
GROUP BY transaction_date
ORDER BY transaction_date DESC;
