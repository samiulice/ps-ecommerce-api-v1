CREATE TABLE IF NOT EXISTS roles (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(120) NOT NULL,
    slug        VARCHAR(120) NOT NULL UNIQUE,
    description TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS permissions (
    id           BIGSERIAL PRIMARY KEY,
    key          VARCHAR(160) NOT NULL UNIQUE,
    display_name VARCHAR(160) NOT NULL,
    module       VARCHAR(120) NOT NULL,
    description  TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

ALTER TABLE employees
    ADD COLUMN IF NOT EXISTS role_id BIGINT;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'employees' AND constraint_name = 'employees_role_id_fkey'
    ) THEN
        ALTER TABLE employees
            ADD CONSTRAINT employees_role_id_fkey
            FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_employees_role_id ON employees(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);

INSERT INTO permissions (key, display_name, module, description)
VALUES
    ('role.view', 'View Roles', 'roles', 'Can view role list and details'),
    ('role.create', 'Create Roles', 'roles', 'Can create new roles'),
    ('role.edit', 'Edit Roles', 'roles', 'Can update existing roles'),
    ('role.delete', 'Delete Roles', 'roles', 'Can delete roles'),
    ('user.view', 'View Employees', 'employees', 'Can view employee list and details'),
    ('user.create', 'Create Employees', 'employees', 'Can create employees'),
    ('user.edit', 'Edit Employees', 'employees', 'Can update employees and assign roles'),
    ('user.delete', 'Delete Employees', 'employees', 'Can delete employees'),
    ('product.view', 'View Products', 'products', 'Can view products'),
    ('product.create', 'Create Products', 'products', 'Can create products'),
    ('product.edit', 'Edit Products', 'products', 'Can update products'),
    ('product.delete', 'Delete Products', 'products', 'Can delete products'),
    ('category.view', 'View Categories', 'categories', 'Can view categories'),
    ('category.create', 'Create Categories', 'categories', 'Can create categories'),
    ('category.edit', 'Edit Categories', 'categories', 'Can update categories'),
    ('category.delete', 'Delete Categories', 'categories', 'Can delete categories'),
    ('brand.view', 'View Brands', 'brands', 'Can view brands'),
    ('brand.create', 'Create Brands', 'brands', 'Can create brands'),
    ('brand.edit', 'Edit Brands', 'brands', 'Can update brands'),
    ('brand.delete', 'Delete Brands', 'brands', 'Can delete brands'),
    ('supplier.view', 'View Suppliers', 'suppliers', 'Can view suppliers'),
    ('supplier.create', 'Create Suppliers', 'suppliers', 'Can create suppliers'),
    ('supplier.edit', 'Edit Suppliers', 'suppliers', 'Can update suppliers'),
    ('supplier.delete', 'Delete Suppliers', 'suppliers', 'Can delete suppliers'),
    ('purchase.view', 'View Purchases', 'purchases', 'Can view purchases'),
    ('purchase.create', 'Create Purchases', 'purchases', 'Can create purchases'),
    ('purchase.edit', 'Edit Purchases', 'purchases', 'Can update purchases'),
    ('purchase.delete', 'Delete Purchases', 'purchases', 'Can delete purchases'),
    ('branch.view', 'View Branches', 'branches', 'Can view branches'),
    ('branch.create', 'Create Branches', 'branches', 'Can create branches'),
    ('branch.edit', 'Edit Branches', 'branches', 'Can update branches'),
    ('branch.delete', 'Delete Branches', 'branches', 'Can delete branches'),
    ('unit.view', 'View Units', 'units', 'Can view units'),
    ('unit.create', 'Create Units', 'units', 'Can create units'),
    ('unit.edit', 'Edit Units', 'units', 'Can update units'),
    ('unit.delete', 'Delete Units', 'units', 'Can delete units'),
    ('attribute.view', 'View Attributes', 'attributes', 'Can view attributes'),
    ('attribute.create', 'Create Attributes', 'attributes', 'Can create attributes'),
    ('attribute.edit', 'Edit Attributes', 'attributes', 'Can update attributes'),
    ('attribute.delete', 'Delete Attributes', 'attributes', 'Can delete attributes'),
    ('settings.view', 'View Settings', 'settings', 'Can view site settings'),
    ('settings.edit', 'Edit Settings', 'settings', 'Can update site settings'),
    ('order.view', 'View Orders', 'orders', 'Can view orders'),
    ('order.edit', 'Edit Orders', 'orders', 'Can update order status and payment state'),
    ('order.delete', 'Delete Orders', 'orders', 'Can delete orders')
ON CONFLICT (key) DO NOTHING;

INSERT INTO roles (name, slug, description, is_active)
VALUES
    ('Chairman', 'chairman', 'Full access role', TRUE),
    ('Manager', 'manager', 'Operational manager role', TRUE),
    ('Staff', 'staff', 'Limited employee role', TRUE)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.slug = 'chairman'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.key IN (
    'user.view', 'user.create', 'user.edit',
    'role.view',
    'product.view', 'product.create', 'product.edit',
    'category.view', 'category.create', 'category.edit',
    'brand.view', 'brand.create', 'brand.edit',
    'supplier.view', 'supplier.create', 'supplier.edit',
    'purchase.view', 'purchase.create', 'purchase.edit',
    'branch.view', 'unit.view', 'attribute.view', 'settings.view', 'settings.edit', 'order.view', 'order.edit'
)
WHERE r.slug = 'manager'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.key IN (
    'product.view', 'category.view', 'brand.view',
    'supplier.view', 'purchase.view', 'purchase.create',
    'order.view', 'unit.view', 'attribute.view'
)
WHERE r.slug = 'staff'
ON CONFLICT DO NOTHING;

UPDATE employees e
SET role_id = r.id,
    role = r.slug
FROM roles r
WHERE e.role_id IS NULL
  AND (
      (LOWER(TRIM(e.role)) = 'chairman' AND r.slug = 'chairman') OR
      (LOWER(TRIM(e.role)) = 'manager' AND r.slug = 'manager') OR
      (LOWER(TRIM(e.role)) = 'staff' AND r.slug = 'staff')
  );

UPDATE employees e
SET role_id = r.id,
    role = r.slug
FROM roles r
WHERE e.role_id IS NULL
  AND r.slug = 'chairman';