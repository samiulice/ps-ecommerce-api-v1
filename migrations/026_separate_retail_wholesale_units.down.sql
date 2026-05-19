-- Drop indexes
DROP INDEX IF EXISTS idx_products_retail_unit_id;
DROP INDEX IF EXISTS idx_products_wholesale_unit_id;

-- Drop retail and wholesale unit columns
ALTER TABLE products DROP COLUMN IF EXISTS retail_unit_id;
ALTER TABLE products DROP COLUMN IF EXISTS wholesale_unit_id;
