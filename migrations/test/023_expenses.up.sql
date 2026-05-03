CREATE TABLE IF NOT EXISTS expense_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    branch_id BIGINT NULL, -- Null if it's a global category, otherwise specific to a branch
    total_amount NUMERIC(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS expenses (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES expense_categories(id) ON DELETE CASCADE,
    branch_id BIGINT NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    expense_date TIMESTAMP NOT NULL,
    description TEXT,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
