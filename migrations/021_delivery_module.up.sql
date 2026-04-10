-- Migration: Create Delivery Management tables
-- Up Migration

-- Delivery Methods Table
CREATE TABLE IF NOT EXISTS delivery_methods (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    base_cost DECIMAL(10, 2) NOT NULL DEFAULT 0,
    estimated_days VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Delivery Men Table (Links to employees table for auth)
CREATE TABLE IF NOT EXISTS delivery_men (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT UNIQUE REFERENCES employees(id) ON DELETE CASCADE,       
    identity_type VARCHAR(50),
    identity_number VARCHAR(100),
    identity_image VARCHAR(255),
    vehicle_type VARCHAR(50),
    vehicle_number VARCHAR(100),
    bank_name VARCHAR(100),
    account_no VARCHAR(100),
    account_holder_name VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_online BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Order Deliveries Table
CREATE TABLE IF NOT EXISTS order_deliveries (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT UNIQUE NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    delivery_man_id BIGINT REFERENCES delivery_men(id) ON DELETE SET NULL,
    delivery_status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, accepted, out_for_delivery, delivered, failed, cancelled
    delivery_fee_collected DECIMAL(10, 2) NOT NULL DEFAULT 0,
    delivery_man_earning DECIMAL(10, 2) NOT NULL DEFAULT 0,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Delivery Wallets Table
CREATE TABLE IF NOT EXISTS delivery_wallets (
    id BIGSERIAL PRIMARY KEY,
    delivery_man_id BIGINT UNIQUE NOT NULL REFERENCES delivery_men(id) ON DELETE CASCADE,
    total_earned DECIMAL(12, 2) NOT NULL DEFAULT 0,
    total_withdrawn DECIMAL(12, 2) NOT NULL DEFAULT 0,
    current_balance DECIMAL(12, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Withdraw Requests Table
CREATE TABLE IF NOT EXISTS withdraw_requests (
    id BIGSERIAL PRIMARY KEY,
    delivery_man_id BIGINT NOT NULL REFERENCES delivery_men(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, approved, rejected
    admin_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_delivery_tables_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_delivery_methods_updated_at BEFORE UPDATE ON delivery_methods FOR EACH ROW EXECUTE FUNCTION update_delivery_tables_updated_at();
CREATE TRIGGER trigger_delivery_men_updated_at BEFORE UPDATE ON delivery_men FOR EACH ROW EXECUTE FUNCTION update_delivery_tables_updated_at();
CREATE TRIGGER trigger_order_deliveries_updated_at BEFORE UPDATE ON order_deliveries FOR EACH ROW EXECUTE FUNCTION update_delivery_tables_updated_at();
CREATE TRIGGER trigger_delivery_wallets_updated_at BEFORE UPDATE ON delivery_wallets FOR EACH ROW EXECUTE FUNCTION update_delivery_tables_updated_at();
CREATE TRIGGER trigger_withdraw_requests_updated_at BEFORE UPDATE ON withdraw_requests FOR EACH ROW EXECUTE FUNCTION update_delivery_tables_updated_at();
-- Insert Delivery Module Permissions
INSERT INTO permissions (key, display_name, module, description)
VALUES 
    ('delivery.manage', 'Manage Delivery Settings', 'delivery', 'Can configure delivery methods, driver onboarding, and handle approvals'),
    ('delivery.assign', 'Assign Deliveries', 'delivery', 'Can assign orders to delivery men')
ON CONFLICT (key) DO NOTHING;

-- Insert Delivery Man Role
INSERT INTO roles (name, slug, description, is_active)
VALUES 
    ('Delivery Man', 'delivery_man', 'Platform delivery rider executing orders on the road', TRUE)
ON CONFLICT (slug) DO NOTHING;

-- Assign Delivery Permissions to Admins
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.slug IN ('chairman', 'manager') AND p.key IN ('delivery.manage', 'delivery.assign')
ON CONFLICT DO NOTHING;
