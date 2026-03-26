DROP INDEX IF EXISTS idx_social_links_name;
ALTER TABLE social_links DROP CONSTRAINT IF EXISTS social_links_display_order_non_negative_chk;
ALTER TABLE social_links DROP CONSTRAINT IF EXISTS social_links_required_trimmed_fields_chk;

DROP INDEX IF EXISTS idx_customers_email_ci_unique;
DROP INDEX IF EXISTS idx_customers_is_wholesaler;
DROP INDEX IF EXISTS idx_customers_is_active_created_at;
ALTER TABLE customers DROP CONSTRAINT IF EXISTS customers_non_negative_values_chk;
ALTER TABLE customers DROP CONSTRAINT IF EXISTS customers_required_trimmed_phone_chk;

DROP INDEX IF EXISTS idx_suppliers_email_ci_unique;
ALTER TABLE suppliers DROP CONSTRAINT IF EXISTS suppliers_required_trimmed_fields_chk;
ALTER TABLE suppliers DROP CONSTRAINT IF EXISTS suppliers_non_negative_values_chk;

DROP INDEX IF EXISTS idx_purchase_return_reference_no;
ALTER TABLE purchase_return DROP CONSTRAINT IF EXISTS purchase_return_totals_non_negative_chk;

DROP INDEX IF EXISTS idx_purchases_reference_no;
ALTER TABLE purchases DROP CONSTRAINT IF EXISTS purchases_totals_non_negative_chk;

DROP INDEX IF EXISTS idx_purchase_orders_due_date;
ALTER TABLE purchase_orders DROP CONSTRAINT IF EXISTS purchase_orders_totals_non_negative_chk;
ALTER TABLE purchase_orders DROP CONSTRAINT IF EXISTS purchase_orders_due_date_after_order_chk;
