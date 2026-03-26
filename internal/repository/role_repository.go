package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepo(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) ListPermissions(ctx context.Context) ([]model.Permission, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, key, display_name, module, COALESCE(description, ''), created_at, updated_at
		FROM permissions
		ORDER BY module, key
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var permission model.Permission
		if err := rows.Scan(&permission.ID, &permission.Key, &permission.DisplayName, &permission.Module, &permission.Description, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, rows.Err()
}

func (r *RoleRepository) List(ctx context.Context) ([]model.Role, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, slug, COALESCE(description, ''), is_active, created_at, updated_at
		FROM roles
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]model.Role, 0)
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Slug, &role.Description, &role.IsActive, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		permissions, err := r.permissionsByRoleID(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		role.Permissions = permissions
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *RoleRepository) FindByID(ctx context.Context, id int64) (*model.Role, error) {
	var role model.Role
	err := r.db.QueryRow(ctx, `
		SELECT id, name, slug, COALESCE(description, ''), is_active, created_at, updated_at
		FROM roles
		WHERE id = $1
	`, id).Scan(&role.ID, &role.Name, &role.Slug, &role.Description, &role.IsActive, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	permissions, err := r.permissionsByRoleID(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions
	return &role, nil
}

func (r *RoleRepository) FindBySlug(ctx context.Context, slug string) (*model.Role, error) {
	var role model.Role
	err := r.db.QueryRow(ctx, `
		SELECT id, name, slug, COALESCE(description, ''), is_active, created_at, updated_at
		FROM roles
		WHERE LOWER(slug) = LOWER($1)
	`, slug).Scan(&role.ID, &role.Name, &role.Slug, &role.Description, &role.IsActive, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	permissions, err := r.permissionsByRoleID(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions
	return &role, nil
}

func (r *RoleRepository) Create(ctx context.Context, role *model.Role, permissionKeys []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO roles (name, slug, description, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, role.Name, role.Slug, nullableString(role.Description), role.IsActive).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return err
	}

	if err := r.replaceRolePermissions(ctx, tx, role.ID, permissionKeys); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *RoleRepository) Update(ctx context.Context, role *model.Role, permissionKeys []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, `
		UPDATE roles
		SET name = $1,
			slug = $2,
			description = $3,
			is_active = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`, role.Name, role.Slug, nullableString(role.Description), role.IsActive, role.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("role with id %d not found", role.ID)
	}

	if _, err := tx.Exec(ctx, `UPDATE employees SET role = $1 WHERE role_id = $2`, role.Slug, role.ID); err != nil {
		return err
	}

	if err := r.replaceRolePermissions(ctx, tx, role.ID, permissionKeys); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	var assignedCount int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM employees WHERE role_id = $1`, id).Scan(&assignedCount); err != nil {
		return err
	}
	if assignedCount > 0 {
		return fmt.Errorf("role %d is assigned to %d employee(s)", id, assignedCount)
	}

	result, err := r.db.Exec(ctx, `DELETE FROM roles WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("role with id %d not found", id)
	}
	return nil
}

func (r *RoleRepository) permissionsByRoleID(ctx context.Context, roleID int64) ([]model.Permission, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.key, p.display_name, p.module, COALESCE(p.description, ''), p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
		ORDER BY p.module, p.key
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]model.Permission, 0)
	for rows.Next() {
		var permission model.Permission
		if err := rows.Scan(&permission.ID, &permission.Key, &permission.DisplayName, &permission.Module, &permission.Description, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, rows.Err()
}

func (r *RoleRepository) replaceRolePermissions(ctx context.Context, tx pgx.Tx, roleID int64, permissionKeys []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, roleID); err != nil {
		return err
	}

	keys := sanitizePermissionKeys(permissionKeys)
	if len(keys) == 0 {
		return nil
	}

	query := `
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT $1, p.id
		FROM permissions p
		WHERE p.key = ANY($2)
	`
	_, err := tx.Exec(ctx, query, roleID, keys)
	return err
}

func sanitizePermissionKeys(keys []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, key)
	}
	return result
}

func nullableString(v string) any {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return strings.TrimSpace(v)
}
