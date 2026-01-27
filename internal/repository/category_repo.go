package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

// isUniqueViolation checks if the error is a Postgres unique constraint violation (Code 23505)
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}

// ---------------------------------------------------------------------
// EFFICIENT TREE RETRIEVAL
// ---------------------------------------------------------------------

// ListFullTree fetches Categories -> Sub -> SubSub in a SINGLE database query.
func (r *CategoryRepo) ListFullTree(ctx context.Context, onlyActive bool) ([]model.Category, error) {
	// We construct a JSON object in SQL to avoid N+1 query problems.
	query := `
		SELECT 
			c.id, c.name, c.logo_url, c.priority, c.is_active, c.created_at, c.updated_at,
			COALESCE((
				SELECT json_agg(sub ORDER BY sub.priority ASC)
				FROM (
					SELECT 
						s.id, s.category_id, s.name, s.priority, s.is_active, s.created_at, s.updated_at,
						COALESCE((
							SELECT json_agg(ss ORDER BY ss.priority ASC)
							FROM sub_sub_categories ss
							WHERE ss.sub_category_id = s.id AND ($1 = FALSE OR ss.is_active = TRUE)
						), '[]'::json) AS sub_sub_categories
					FROM sub_categories s
					WHERE s.category_id = c.id AND ($1 = FALSE OR s.is_active = TRUE)
				) sub
			), '[]'::json) AS sub_categories
		FROM categories c
		WHERE ($1 = FALSE OR c.is_active = TRUE)
		ORDER BY c.priority ASC, c.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, onlyActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		// pgx automatically unmarshals the JSON column into the Struct slice
		err := rows.Scan(
			&c.ID, &c.Name, &c.LogoURL, &c.Priority, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
			&c.SubCategories,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// ---------------------------------------------------------------------
// LEVEL 1: CATEGORY CRUD
// ---------------------------------------------------------------------

func (r *CategoryRepo) Create(ctx context.Context, c *model.Category) error {
	query := `INSERT INTO categories (name, logo_url, priority, is_active) 
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, c.Name, c.LogoURL, c.Priority, c.IsActive).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if isUniqueViolation(err) {
		return fmt.Errorf("category name '%s' already exists", c.Name)
	}
	return err
}

func (r *CategoryRepo) Update(ctx context.Context, c *model.Category) error {
	query := `UPDATE categories SET name=$1, logo_url=$2, priority=$3, is_active=$4, updated_at=NOW() 
	          WHERE id=$5 RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, query, c.Name, c.LogoURL, c.Priority, c.IsActive, c.ID).
		Scan(&c.CreatedAt, &c.UpdatedAt)
	if isUniqueViolation(err) {
		return fmt.Errorf("category name '%s' already exists", c.Name)
	}
	return err
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err == nil && tag.RowsAffected() == 0 {
		return errors.New("category not found")
	}
	return err
}

// GetByID is needed to preserve the old LogoURL during updates if no new image is sent
func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	query := `SELECT id, name, logo_url, priority, is_active, created_at, updated_at FROM categories WHERE id=$1`
	var c model.Category
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.LogoURL, &c.Priority, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.New("category not found")
	}
	return &c, err
}

// ---------------------------------------------------------------------
// LEVEL 2: SUB-CATEGORY CRUD
// ---------------------------------------------------------------------

func (r *CategoryRepo) CreateSub(ctx context.Context, s *model.SubCategory) error {
	query := `INSERT INTO sub_categories (category_id, name, priority, is_active) 
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, s.CategoryID, s.Name, s.Priority, s.IsActive).
		Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	if isUniqueViolation(err) {
		return fmt.Errorf("sub-category name '%s' already exists", s.Name)
	}
	return err
}

func (r *CategoryRepo) UpdateSub(ctx context.Context, s *model.SubCategory) error {
	query := `UPDATE sub_categories SET name=$1, priority=$2, is_active=$3, updated_at=NOW() 
	          WHERE id=$4 RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, query, s.Name, s.Priority, s.IsActive, s.ID).
		Scan(&s.CreatedAt, &s.UpdatedAt)
	if isUniqueViolation(err) {
		return fmt.Errorf("sub-category name '%s' already exists", s.Name)
	}
	return err
}

func (r *CategoryRepo) DeleteSub(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM sub_categories WHERE id = $1", id)
	if err == nil && tag.RowsAffected() == 0 {
		return errors.New("sub-category not found")
	}
	return err
}

func (r *CategoryRepo) GetSubByID(ctx context.Context, id int64) (*model.SubCategory, error) {
	query := `SELECT id, category_id, name, priority, is_active, created_at, updated_at 
	          FROM sub_categories WHERE id = $1`
	var s model.SubCategory
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.CategoryID, &s.Name, &s.Priority, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.New("sub-category not found")
	}
	return &s, err
}

// ---------------------------------------------------------------------
// LEVEL 3: SUB-SUB-CATEGORY CRUD
// ---------------------------------------------------------------------

func (r *CategoryRepo) CreateSubSub(ctx context.Context, ss *model.SubSubCategory) error {
	query := `INSERT INTO sub_sub_categories (sub_category_id, name, priority, is_active) 
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, ss.SubCategoryID, ss.Name, ss.Priority, ss.IsActive).
		Scan(&ss.ID, &ss.CreatedAt, &ss.UpdatedAt)
	if isUniqueViolation(err) {
		return fmt.Errorf("sub-sub-category name '%s' already exists", ss.Name)
	}
	return err
}

func (r *CategoryRepo) UpdateSubSub(ctx context.Context, ss *model.SubSubCategory) error {
	query := `UPDATE sub_sub_categories SET name=$1, priority=$2, is_active=$3, updated_at=NOW() 
	          WHERE id=$4 RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, query, ss.Name, ss.Priority, ss.IsActive, ss.ID).
		Scan(&ss.CreatedAt, &ss.UpdatedAt)
	if isUniqueViolation(err) {
		return fmt.Errorf("sub-sub-category name '%s' already exists", ss.Name)
	}
	return err
}

func (r *CategoryRepo) DeleteSubSub(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM sub_sub_categories WHERE id = $1", id)
	if err == nil && tag.RowsAffected() == 0 {
		return errors.New("sub-sub-category not found")
	}
	return err
}
func (r *CategoryRepo) GetSubSubByID(ctx context.Context, id int64) (*model.SubSubCategory, error) {
	query := `SELECT id, sub_category_id, name, priority, is_active, created_at, updated_at 
	          FROM sub_sub_categories WHERE id = $1`
	var ss model.SubSubCategory
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ss.ID, &ss.SubCategoryID, &ss.Name, &ss.Priority, &ss.IsActive, &ss.CreatedAt, &ss.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.New("sub-sub-category not found")
	}
	return &ss, err
}