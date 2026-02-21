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

type BrandRepo struct {
	db *pgxpool.Pool
}

func NewBrandRepo(db *pgxpool.Pool) *BrandRepo {
	return &BrandRepo{db: db}
}

// ---------------------------------------------------------------------
// LEVEL 1: BRAND CRUD
// ---------------------------------------------------------------------

func (r *BrandRepo) Create(ctx context.Context, b *model.Brand) error {
	query := `INSERT INTO brands (name, priority, thumbnail, is_active) 
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, b.Name, b.Priority, b.Thumbnail, b.IsActive).
		Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
	if isUniqueViolation(err) {
		return fmt.Errorf("brand name '%s' already exists", b.Name)
	}
	return err
}

func (r *BrandRepo) Update(ctx context.Context, b *model.Brand) error {
	query := `UPDATE brands SET name=$1, priority=$2, thumbnail=$3, is_active=$4, updated_at=CURRENT_TIMESTAMP
	          WHERE id=$5 RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, query, b.Name, b.Priority, b.Thumbnail, b.IsActive, b.ID).
		Scan(&b.CreatedAt, &b.UpdatedAt)
	if isUniqueViolation(err) {
		return fmt.Errorf("brand name '%s' already exists", b.Name)
	}
	return err
}

func (r *BrandRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM brands WHERE id = $1", id)
	if err == nil && tag.RowsAffected() == 0 {
		return errors.New("brand not found")
	}
	return err
}

// GetByID is needed to preserve the old LogoURL during updates if no new image is sent
func (r *BrandRepo) GetByID(ctx context.Context, id int64) (*model.Brand, error) {
	query := `SELECT id, name, priority, thumbnail, is_active, created_at, updated_at FROM brands WHERE id=$1`
	var b model.Brand
	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.Name, &b.Priority, &b.Thumbnail, &b.IsActive, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.New("brand not found")
	}
	return &b, err
}
func (r *BrandRepo) GetBrands(ctx context.Context, status string) ([]*model.Brand, error) {

	baseQuery := `
		SELECT id, name, priority, thumbnail, is_active, created_at, updated_at
		FROM brands
	`

	var conditions []string
	var args []any
	argPos := 1

	// status filter
	if status == "active" {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, true)
		argPos++
	} else if status == "inactive" {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, false)
		argPos++
	}

	// build final query
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY priority ASC"

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []*model.Brand

	for rows.Next() {
		var brand model.Brand
		err := rows.Scan(
			&brand.ID,
			&brand.Name,
			&brand.Priority,
			&brand.Thumbnail,
			&brand.IsActive,
			&brand.CreatedAt,
			&brand.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		brands = append(brands, &brand)
	}

	return brands, nil
}