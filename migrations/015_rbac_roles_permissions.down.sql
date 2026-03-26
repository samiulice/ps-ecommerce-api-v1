ALTER TABLE employees DROP CONSTRAINT IF EXISTS employees_role_id_fkey;
DROP INDEX IF EXISTS idx_employees_role_id;
ALTER TABLE employees DROP COLUMN IF EXISTS role_id;

DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;