package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type AttributeRepo struct {
	db *pgxpool.Pool
}

func NewAttributeRepo(db *pgxpool.Pool) *AttributeRepo {
	return &AttributeRepo{db: db}
}

// Create inserts a new attribute
func (r *AttributeRepo) Create(ctx context.Context, u *model.Attribute) error {
	query := `INSERT INTO attributes (name) 
	          VALUES ($1) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, u.Name).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	// Reuse your existing unique violation helper if available, or check manually
	if err != nil && isUniqueViolation(err) {
		return fmt.Errorf("attribute name '%s' already exists", u.Name)
	}
	return err
}

// Update modifies an existing attribute
func (r *AttributeRepo) Update(ctx context.Context, u *model.Attribute) error {
	query := `UPDATE attributes SET name=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2`

	_, err := r.db.Exec(ctx, query, u.Name, u.ID)

	if err != nil && isUniqueViolation(err) {
		return fmt.Errorf("attribute name '%s' already exists", u.Name)
	}
	return err
}

// Delete removes a attribute
func (r *AttributeRepo) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM attributes WHERE id = $1", id)
	if err == nil && tag.RowsAffected() == 0 {
		return errors.New("attribute not found")
	}
	return err
}

// GetByID fetches a single attribute
func (r *AttributeRepo) GetByID(ctx context.Context, id int) (*model.Attribute, error) {
	query := `SELECT id, name, created_at, updated_at FROM attributes WHERE id = $1`

	var u model.Attribute
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.New("attribute not found")
	}
	return &u, err
}

// GetAll retrieves all attributes
func (r *AttributeRepo) GetAll(ctx context.Context) ([]*model.Attribute, error) {
	query := `SELECT id, name, created_at, updated_at FROM attributes ORDER BY id ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attributes []*model.Attribute
	for rows.Next() {
		var u model.Attribute
		if err := rows.Scan(&u.ID, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		attributes = append(attributes, &u)
	}
	return attributes, nil
}
