// Package repository implements data access logic for employees (admin users).
package repository

import (
	"context"
	"errors"
	"fmt"

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
		INSERT INTO employees (name, email, password_hash, mobile, role, is_active, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uuid, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		e.Name, e.Email, e.Password, e.Mobile, e.Role, e.IsActive, e.IsVerified,
	).Scan(&e.ID, &e.UUID, &e.CreatedAt, &e.UpdatedAt)

	return err
}

// GetByEmail retrieves an employee by email.
func (r *EmployeeRepository) GetByEmail(ctx context.Context, email string) (*model.Employee, error) {
	query := `
		SELECT id, uuid, branch_id, name, email, password_hash, mobile, role, is_active, is_verified, created_at, updated_at
		FROM employees WHERE email = $1`

	e := &model.Employee{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&e.ID, &e.UUID, &e.BranchID, &e.Name, &e.Email, &e.Password,
		&e.Mobile, &e.Role, &e.IsActive, &e.IsVerified, &e.CreatedAt, &e.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("employee with email '%s' not found", email)
	}
	return e, err
}

// FindByID retrieves an employee by ID.
func (r *EmployeeRepository) FindByID(ctx context.Context, id int) (*model.Employee, error) {
	query := `
		SELECT id, uuid, branch_id, name, email, password_hash, mobile, role, is_active, is_verified, created_at, updated_at
		FROM employees WHERE id = $1`

	e := &model.Employee{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.UUID, &e.BranchID, &e.Name, &e.Email, &e.Password,
		&e.Mobile, &e.Role, &e.IsActive, &e.IsVerified, &e.CreatedAt, &e.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("employee with id %d not found", id)
	}
	return e, err
}
