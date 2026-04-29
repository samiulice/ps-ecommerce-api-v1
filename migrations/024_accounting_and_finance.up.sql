-- UP MIGRATION: 024_accounting_and_finance.up.sql
-- Create transactions ledger and financial calculation views

-- 1. Create Transaction Ledger Table
DO $$ BEGIN
    CREATE TYPE transaction_type AS ENUM (
        'sale_income', 
        'purchase_expense', 
        'operational_expense', 
        'sale_refund', 
        'purchase_refund'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

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

-- ---------------------------------------------------------
-- 2. Trigger Functions for Automation
-- ---------------------------------------------------------

-- A) Orders Trigger (E-commerce Sales & Refunds)
CREATE OR REPLACE FUNCTION log_order_transaction()
RETURNS TRIGGER AS $$
BEGIN
    -- Log Income when an order is paid
    IF NEW.payment_status = 'paid' AND (OLD.payment_status IS NULL OR OLD.payment_status != 'paid') THEN
        INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount)
        VALUES (CURRENT_DATE, 'sale_income', NEW.id, NEW.order_number, 'Income from Order: ' || NEW.order_number, NEW.total, 0);
    
    -- Log Refund if payment status moves to refunded
    ELSIF NEW.payment_status = 'refunded' AND OLD.payment_status != 'refunded' THEN
        INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount)
        VALUES (CURRENT_DATE, 'sale_refund', NEW.id, NEW.order_number, 'Refund for Order: ' || NEW.order_number, 0, NEW.total);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_log_order_transaction ON orders;
CREATE TRIGGER trigger_log_order_transaction
    AFTER UPDATE OF payment_status ON orders
    FOR EACH ROW
    EXECUTE FUNCTION log_order_transaction();


-- B) POS Sales Trigger (Immediate Income)
CREATE OR REPLACE FUNCTION log_pos_sale_transaction()
RETURNS TRIGGER AS $$
BEGIN
    -- POS sales are typically paid instantly
    INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount)
    VALUES (DATE(NEW.sale_date), 'sale_income', NEW.id, NEW.reference_no, 'POS Sale: ' || NEW.reference_no, NEW.amount_paid, 0);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_log_pos_sale_transaction ON pos_sales;
CREATE TRIGGER trigger_log_pos_sale_transaction
    AFTER INSERT ON pos_sales
    FOR EACH ROW
    EXECUTE FUNCTION log_pos_sale_transaction();


-- C) Purchase Trigger (Inventory Expense with Delta Tracking)
CREATE OR REPLACE FUNCTION log_purchase_transaction()
RETURNS TRIGGER AS $$
DECLARE
    paid_delta NUMERIC(15,2);
BEGIN
    paid_delta := NEW.paid_amount - COALESCE(OLD.paid_amount, 0);
    
    IF paid_delta > 0 THEN
        INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount)
        VALUES (CURRENT_DATE, 'purchase_expense', NEW.id, NEW.purchase_code, 'Payment for Purchase: ' || NEW.purchase_code, 0, paid_delta);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_log_purchase_transaction ON purchases;
CREATE TRIGGER trigger_log_purchase_transaction
    AFTER INSERT OR UPDATE OF paid_amount ON purchases
    FOR EACH ROW
    EXECUTE FUNCTION log_purchase_transaction();


-- D) Purchase Return Trigger (Recovered Money)
CREATE OR REPLACE FUNCTION log_purchase_return_transaction()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.paid_amount > 0 THEN
        INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount)
        VALUES (DATE(NEW.return_date), 'purchase_refund', NEW.id, NEW.return_code, 'Refund from Purchase Return: ' || NEW.return_code, NEW.paid_amount, 0);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_log_purchase_return_transaction ON purchase_return;
CREATE TRIGGER trigger_log_purchase_return_transaction
    AFTER INSERT ON purchase_return
    FOR EACH ROW
    EXECUTE FUNCTION log_purchase_return_transaction();


-- E) General Expenses Trigger
CREATE OR REPLACE FUNCTION log_operational_expense()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount)
    VALUES (DATE(NEW.expense_date), 'operational_expense', NEW.id, 'EXP-'||NEW.id, NEW.description, 0, NEW.amount);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_log_operational_expense ON expenses;
CREATE TRIGGER trigger_log_operational_expense
    AFTER INSERT ON expenses
    FOR EACH ROW
    EXECUTE FUNCTION log_operational_expense();


-- F) Delivery Earnings Trigger (Automatic Operational Expense)
CREATE OR REPLACE FUNCTION log_delivery_expense()
RETURNS TRIGGER AS $$
BEGIN
    -- When delivery is complete, the earning for the delivery man is a cost to the company
    IF NEW.delivery_status = 'delivered' AND OLD.delivery_status != 'delivered' AND NEW.delivery_man_earning > 0 THEN
        INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount)
        VALUES (CURRENT_DATE, 'operational_expense', NEW.id, 'DEL-'||NEW.order_id, 'Delivery Man Earning for Order #' || NEW.order_id, 0, NEW.delivery_man_earning);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_log_delivery_expense ON order_deliveries;
CREATE TRIGGER trigger_log_delivery_expense
    AFTER UPDATE OF delivery_status ON order_deliveries
    FOR EACH ROW
    EXECUTE FUNCTION log_delivery_expense();


-- ---------------------------------------------------------
-- 3. Reporting View
-- ---------------------------------------------------------

CREATE OR REPLACE VIEW view_daily_financial_report AS
SELECT 
    transaction_date,
    -- Income Categories
    SUM(CASE WHEN type = 'sale_income' THEN credit_amount ELSE 0 END) AS gross_sales_income,
    SUM(CASE WHEN type = 'purchase_refund' THEN credit_amount ELSE 0 END) AS purchase_refunds_received,
    
    -- Expense Categories
    SUM(CASE WHEN type = 'sale_refund' THEN debit_amount ELSE 0 END) AS sales_refunds_paid,
    SUM(CASE WHEN type = 'purchase_expense' THEN debit_amount ELSE 0 END) AS inventory_purchase_costs,
    SUM(CASE WHEN type = 'operational_expense' THEN debit_amount ELSE 0 END) AS operational_expenses,
    
    -- Totals
    (SUM(credit_amount)) AS total_cash_in,
    (SUM(debit_amount)) AS total_cash_out,
    
    -- Net Profit = (All Income) - (All Costs/Refunds)
    (SUM(credit_amount) - SUM(debit_amount)) AS net_cash_flow
FROM account_transactions
GROUP BY transaction_date;