-- Add retail and wholesale unit columns referencing the units table
ALTER TABLE products ADD COLUMN retail_unit_id INT REFERENCES units(id);
ALTER TABLE products ADD COLUMN wholesale_unit_id INT REFERENCES units(id);

-- Migrate existing unit_id values to both new unit columns
UPDATE products SET retail_unit_id = unit_id, wholesale_unit_id = unit_id;

-- Add performance indexes for foreign key lookups
CREATE INDEX IF NOT EXISTS idx_products_retail_unit_id ON products(retail_unit_id);
CREATE INDEX IF NOT EXISTS idx_products_wholesale_unit_id ON products(wholesale_unit_id);
