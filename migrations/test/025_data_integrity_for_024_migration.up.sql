BEGIN;

-- 1. Backfill Paid Orders (Income)
INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount, created_at)
SELECT 
    DATE(created_at), 
    'sale_income', 
    id, 
    order_number, 
    'Historical Income from Order: ' || order_number, 
    total, 
    0, 
    created_at
FROM orders 
WHERE payment_status = 'paid'
AND NOT EXISTS (
    SELECT 1 FROM account_transactions 
    WHERE reference_id = orders.id AND type = 'sale_income'
);

-- 2. Backfill Refunded Orders (Loss/Refund)
INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount, created_at)
SELECT 
    DATE(updated_at), 
    'sale_refund', 
    id, 
    order_number, 
    'Historical Refund for Order: ' || order_number, 
    0, 
    total, 
    updated_at
FROM orders 
WHERE payment_status = 'refunded'
AND NOT EXISTS (
    SELECT 1 FROM account_transactions 
    WHERE reference_id = orders.id AND type = 'sale_refund'
);

-- 3. Backfill POS Sales (Income)
INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount, created_at)
SELECT 
    DATE(sale_date), 
    'sale_income', 
    id, 
    reference_no, 
    'Historical POS Sale: ' || reference_no, 
    amount_paid, 
    0, 
    created_at
FROM pos_sales
WHERE amount_paid > 0
AND NOT EXISTS (
    SELECT 1 FROM account_transactions 
    WHERE reference_id = pos_sales.id AND type = 'sale_income' AND description LIKE 'Historical POS%'
);

-- 4. Backfill Purchase Payments (Expenses)
INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount, created_at)
SELECT 
    DATE(purchase_date), 
    'purchase_expense', 
    id, 
    purchase_code, 
    'Historical Purchase Payment: ' || purchase_code, 
    0, 
    paid_amount, 
    created_at
FROM purchases
WHERE paid_amount > 0
AND NOT EXISTS (
    SELECT 1 FROM account_transactions 
    WHERE reference_id = purchases.id AND type = 'purchase_expense'
);

-- 5. Backfill Purchase Returns (Income Recovery)
INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount, created_at)
SELECT 
    DATE(return_date), 
    'purchase_refund', 
    id, 
    return_code, 
    'Historical Purchase Return: ' || return_code, 
    paid_amount, 
    0, 
    created_at
FROM purchase_return
WHERE paid_amount > 0
AND NOT EXISTS (
    SELECT 1 FROM account_transactions 
    WHERE reference_id = purchase_return.id AND type = 'purchase_refund'
);

-- 6. Backfill General Expenses (Operational Expense) - FIXED
INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount, created_at)
SELECT 
    DATE(expense_date), 
    'operational_expense', 
    id, 
    'EXP-' || id::text, -- Added cast to text
    description, 
    0, 
    amount, 
    created_at
FROM expenses
WHERE NOT EXISTS ( -- Added missing WHERE
    SELECT 1 FROM account_transactions 
    WHERE reference_id = expenses.id AND type = 'operational_expense' -- Removed quotes from type
);

-- 7. Backfill Delivery Earnings (Operational Expense)
INSERT INTO account_transactions (transaction_date, type, reference_id, reference_number, description, credit_amount, debit_amount, created_at)
SELECT 
    DATE(updated_at), 
    'operational_expense', 
    id, 
    'DEL-' || order_id::text, -- Added cast to text
    'Historical Delivery Earning for Order #' || order_id::text, 
    0, 
    delivery_man_earning, 
    updated_at
FROM order_deliveries
WHERE delivery_status = 'delivered' AND delivery_man_earning > 0
AND NOT EXISTS (
    SELECT 1 FROM account_transactions 
    WHERE reference_id = order_deliveries.id AND description LIKE 'Historical Delivery%'
);

COMMIT;