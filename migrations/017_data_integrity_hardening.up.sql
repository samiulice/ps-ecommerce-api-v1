-- Migration: Harden integrity rules and lookup indexes for purchases, suppliers, customers, and social links.

-- -----------------------------
-- purchase_orders constraints
-- -----------------------------
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'purchase_orders_due_date_after_order_chk'
    ) THEN
        ALTER TABLE purchase_orders
            ADD CONSTRAINT purchase_orders_due_date_after_order_chk
            CHECK (due_date IS NULL OR due_date >= order_date)
            NOT VALID;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'purchase_orders_totals_non_negative_chk'
    ) THEN
        ALTER TABLE purchase_orders
            ADD CONSTRAINT purchase_orders_totals_non_negative_chk
            CHECK (
                round_off >= 0
                AND grand_total >= 0
                AND paid_amount >= 0
                AND exchange_rate >= 0
                AND paid_amount <= grand_total
            )
            NOT VALID;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_purchase_orders_due_date
    ON purchase_orders (due_date);

-- -----------------------------
-- purchases constraints
-- -----------------------------
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'purchases_totals_non_negative_chk'
    ) THEN
        ALTER TABLE purchases
            ADD CONSTRAINT purchases_totals_non_negative_chk
            CHECK (
                shipping_charge >= 0
                AND round_off >= 0
                AND grand_total >= 0
                AND paid_amount >= 0
                AND exchange_rate >= 0
                AND paid_amount <= grand_total
                AND (change_return IS NULL OR change_return >= 0)
            )
            NOT VALID;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_purchases_reference_no
    ON purchases (reference_no)
    WHERE reference_no IS NOT NULL;

-- -----------------------------
-- purchase_return constraints
-- -----------------------------
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'purchase_return_totals_non_negative_chk'
    ) THEN
        ALTER TABLE purchase_return
            ADD CONSTRAINT purchase_return_totals_non_negative_chk
            CHECK (
                round_off >= 0
                AND grand_total >= 0
                AND paid_amount >= 0
                AND exchange_rate >= 0
                AND paid_amount <= grand_total
            )
            NOT VALID;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_purchase_return_reference_no
    ON purchase_return (reference_no)
    WHERE reference_no IS NOT NULL;

-- -----------------------------
-- suppliers constraints/indexes
-- -----------------------------
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'suppliers_non_negative_values_chk'
    ) THEN
        ALTER TABLE suppliers
            ADD CONSTRAINT suppliers_non_negative_values_chk
            CHECK (
                credit_limit >= 0
                AND outstanding_balance >= 0
                AND (lead_time_days IS NULL OR lead_time_days >= 0)
            )
            NOT VALID;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'suppliers_required_trimmed_fields_chk'
    ) THEN
        ALTER TABLE suppliers
            ADD CONSTRAINT suppliers_required_trimmed_fields_chk
            CHECK (
                btrim(supplier_code) <> ''
                AND btrim(name) <> ''
                AND btrim(phone) <> ''
            )
            NOT VALID;
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS idx_suppliers_email_ci_unique
    ON suppliers (LOWER(email))
    WHERE email IS NOT NULL;

-- -----------------------------
-- customers constraints/indexes
-- -----------------------------
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'customers_required_trimmed_phone_chk'
    ) THEN
        ALTER TABLE customers
            ADD CONSTRAINT customers_required_trimmed_phone_chk
            CHECK (btrim(phone) <> '')
            NOT VALID;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'customers_non_negative_values_chk'
    ) THEN
        ALTER TABLE customers
            ADD CONSTRAINT customers_non_negative_values_chk
            CHECK (
                (wallet_balance IS NULL OR wallet_balance >= 0)
                AND (loyalty_point IS NULL OR loyalty_point >= 0)
                AND login_hit_count >= 0
            )
            NOT VALID;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_customers_is_active_created_at
    ON customers (is_active, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_customers_is_wholesaler
    ON customers (is_wholesaler);

CREATE UNIQUE INDEX IF NOT EXISTS idx_customers_email_ci_unique
    ON customers (LOWER(email))
    WHERE email IS NOT NULL;

-- -----------------------------
-- social_links constraints/indexes
-- -----------------------------
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'social_links_required_trimmed_fields_chk'
    ) THEN
        ALTER TABLE social_links
            ADD CONSTRAINT social_links_required_trimmed_fields_chk
            CHECK (
                btrim(name) <> ''
                AND btrim(icon) <> ''
                AND btrim(link) <> ''
            )
            NOT VALID;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'social_links_display_order_non_negative_chk'
    ) THEN
        ALTER TABLE social_links
            ADD CONSTRAINT social_links_display_order_non_negative_chk
            CHECK (display_order >= 0)
            NOT VALID;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_social_links_name
    ON social_links (name);
