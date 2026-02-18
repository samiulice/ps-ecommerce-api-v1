package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type UnitRepo struct {
	db *pgxpool.Pool
}

func NewUnitRepo(db *pgxpool.Pool) *UnitRepo {
	return &UnitRepo{db: db}
}

// Create inserts a new unit
func (r *UnitRepo) Create(ctx context.Context, u *model.Unit) error {
	query := `INSERT INTO units (name, symbol) 
	          VALUES ($1, $2) RETURNING id, created_at, updated_at`
	
	err := r.db.QueryRow(ctx, query, u.Name, u.Symbol).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	
	// Reuse your existing unique violation helper if available, or check manually
	if err != nil && isUniqueViolation(err) {
		return fmt.Errorf("unit name '%s' already exists", u.Name)
	}
	return err
}

// Update modifies an existing unit
func (r *UnitRepo) Update(ctx context.Context, u *model.Unit) error {
	query := `UPDATE units SET name=$1, symbol=$2, updated_at=CURRENT_TIMESTAMP WHERE id=$3`
	
	_, err := r.db.Exec(ctx, query, u.Name, u.Symbol, u.ID)
	
	if err != nil && isUniqueViolation(err) {
		return fmt.Errorf("unit name '%s' already exists", u.Name)
	}
	return err
}

// Delete removes a unit
func (r *UnitRepo) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM units WHERE id = $1", id)
	if err == nil && tag.RowsAffected() == 0 {
		return errors.New("unit not found")
	}
	return err
}

// GetByID fetches a single unit
func (r *UnitRepo) GetByID(ctx context.Context, id int) (*model.Unit, error) {
	query := `SELECT id, name, symbol, created_at, updated_at FROM units WHERE id = $1`
	
	var u model.Unit
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Symbol, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.New("unit not found")
	}
	return &u, err
}

// GetAll retrieves all units
func (r *UnitRepo) GetAll(ctx context.Context) ([]*model.Unit, error) {
	query := `SELECT id, name, symbol, created_at, updated_at FROM units ORDER BY id ASC`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*model.Unit
	for rows.Next() {
		var u model.Unit
		if err := rows.Scan(&u.ID, &u.Name, &u.Symbol, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		units = append(units, &u)
	}
	return units, nil
}