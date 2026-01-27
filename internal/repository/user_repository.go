// Package repository implements data access logic.
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/projuktisheba/pse-api-v1/internal/model"
)

// UserRepository is a pgx-based implementation of UserRepository.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepo creates a new UserRepository.
func NewUserRepo(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, uuid, name, email, password_hash, mobile, role, is_active, is_verified, created_at, updated_at FROM users WHERE email = $1`,
		email,
	)

	u := &models.User{}
	err := row.Scan(&u.ID, &u.UUID, &u.Name, &u.Email, &u.Password, &u.Mobile, &u.Role, &u.IsActive, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users(name, email, password_hash, mobile, role, is_active, is_verified) VALUES($1, $2, $3, $4, $5, $6, $7)`,
		u.Name, u.Email, u.Password, u.Mobile, u.Role, u.IsActive, u.IsVerified,
	)
	return err
}
