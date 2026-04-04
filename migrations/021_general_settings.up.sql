CREATE TABLE IF NOT EXISTS general_settings (
    id SERIAL PRIMARY KEY,
    
    company_name VARCHAR(255) DEFAULT 'ProjuktiSheba',
    company_logo VARCHAR(255) DEFAULT '',
    currency_symbol VARCHAR(10) DEFAULT '৳',
    currency_code VARCHAR(10) DEFAULT 'BDT',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Ensure Singleton Row
CREATE UNIQUE INDEX only_one_general_setting ON general_settings ((1));

-- Seed Initial Row
INSERT INTO general_settings (id, company_name) VALUES (1, 'ProjuktiSheba') ON CONFLICT DO NOTHING;