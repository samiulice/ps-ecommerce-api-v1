// Package repository implements data access logic for employees (admin users).
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

// EmployeeRepository handles database operations for employees.
type EmployeeRepository struct {
	db *pgxpool.Pool
}

// NewEmployeeRepo creates a new EmployeeRepository.
func NewEmployeeRepo(db *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// Create inserts a new employee into the database.
func (r *EmployeeRepository) Create(ctx context.Context, e *model.Employee) error {
	query := `
		INSERT INTO employees (name, email, password_hash, mobile, role_id, role, branch_id, is_active, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, uuid, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		e.Name, e.Email, e.Password, e.Mobile, nullableRoleID(e.RoleID), e.Role, e.BranchID, e.IsActive, e.IsVerified,
	).Scan(&e.ID, &e.UUID, &e.CreatedAt, &e.UpdatedAt)

	return err
}

// GetByEmail retrieves an employee by email.
func (r *EmployeeRepository) GetByEmail(ctx context.Context, email string) (*model.Employee, error) {
	e, err := r.findOne(ctx, "e.email = $1", email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("employee with email '%s' not found", email)
	}
	return e, err
}

// FindByID retrieves an employee by ID.
func (r *EmployeeRepository) FindByID(ctx context.Context, id int) (*model.Employee, error) {
	e, err := r.findOne(ctx, "e.id = $1", id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("employee with id %d not found", id)
	}
	return e, err
}

func (r *EmployeeRepository) List(ctx context.Context) ([]model.Employee, error) {
	query := employeeBaseQuery("") + ` ORDER BY e.id DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []model.Employee
	for rows.Next() {
		e, scanErr := scanEmployee(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		e.Password = ""
		employees = append(employees, *e)
	}
	return employees, rows.Err()
}

func (r *EmployeeRepository) UpdateAdmin(ctx context.Context, e *model.Employee) error {
	query := `
		UPDATE employees
		SET name = $1,
			email = $2,
			mobile = $3,
			branch_id = $4,
			role_id = $5,
			role = $6,
			is_active = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8`

	result, err := r.db.Exec(ctx, query,
		e.Name,
		e.Email,
		e.Mobile,
		e.BranchID,
		nullableRoleID(e.RoleID),
		e.Role,
		e.IsActive,
		e.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("employee with id %d not found", e.ID)
	}
	return nil
}

func (r *EmployeeRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx, `DELETE FROM employees WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("employee with id %d not found", id)
	}
	return nil
}

func (r *EmployeeRepository) ExistsByEmailExcludingID(ctx context.Context, email string, excludeID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM employees WHERE LOWER(email) = LOWER($1) AND id <> $2)`,
		email, excludeID,
	).Scan(&exists)
	return exists, err
}

func (r *EmployeeRepository) findOne(ctx context.Context, condition string, arg any) (*model.Employee, error) {
	query := employeeBaseQuery("WHERE " + condition)
	row := r.db.QueryRow(ctx, query, arg)
	return scanEmployee(row)
}

func employeeBaseQuery(where string) string {
	return `
		SELECT
			e.id,
			e.uuid,
			e.branch_id,
			e.name,
			e.email,
			e.password_hash,
			e.mobile,
			COALESCE(e.role_id, 0) AS role_id,
			COALESCE(r.slug, e.role, '') AS role_slug,
			COALESCE(r.name, '') AS role_name,
			COALESCE(ARRAY_AGG(DISTINCT p.key) FILTER (WHERE p.key IS NOT NULL), '{}') AS permissions,
			e.is_active,
			e.is_verified,
			e.created_at,
			e.updated_at
		FROM employees e
		LEFT JOIN roles r ON r.id = e.role_id
		LEFT JOIN role_permissions rp ON rp.role_id = r.id
		LEFT JOIN permissions p ON p.id = rp.permission_id
		` + where + `
		GROUP BY e.id, r.id`
}

type employeeScanner interface {
	Scan(dest ...any) error
}

func scanEmployee(scanner employeeScanner) (*model.Employee, error) {
	e := &model.Employee{}
	permissions := make([]string, 0)
	err := scanner.Scan(
		&e.ID,
		&e.UUID,
		&e.BranchID,
		&e.Name,
		&e.Email,
		&e.Password,
		&e.Mobile,
		&e.RoleID,
		&e.Role,
		&e.RoleName,
		&permissions,
		&e.IsActive,
		&e.IsVerified,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	filtered := make([]string, 0, len(permissions))
	for _, key := range permissions {
		if strings.TrimSpace(key) != "" {
			filtered = append(filtered, key)
		}
	}
	e.Permissions = filtered
	return e, nil
}

func nullableRoleID(roleID int64) any {
	if roleID <= 0 {
		return nil
	}
	return roleID
}
