package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type BranchRepo struct {
	db *pgxpool.Pool
}

func NewBranchRepo(db *pgxpool.Pool) *BranchRepo {
	return &BranchRepo{db: db}
}

func (r *BranchRepo) Create(ctx context.Context, b *model.Branch) error {
	query := `INSERT INTO branches 
		(name, country, city, address, mobile, telephone, email, latitude, longitude) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		b.Name, b.Country, b.City, b.Address, b.Mobile, b.Telephone, b.Email, b.Latitude, b.Longitude).
		Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)

	if isUniqueViolation(err) {
		return fmt.Errorf("branch name '%s' already exists", b.Name)
	}
	return err
}

func (r *BranchRepo) Update(ctx context.Context, b *model.Branch) error {
	query := `UPDATE branches SET 
		name=$1, country=$2, city=$3, address=$4, mobile=$5, telephone=$6, email=$7, latitude=$8, longitude=$9, updated_at=CURRENT_TIMESTAMP
		WHERE id=$10 RETURNING created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		b.Name, b.Country, b.City, b.Address, b.Mobile, b.Telephone, b.Email, b.Latitude, b.Longitude, b.ID).
		Scan(&b.CreatedAt, &b.UpdatedAt)

	if isUniqueViolation(err) {
		return fmt.Errorf("branch name '%s' already exists", b.Name)
	}
	return err
}

func (r *BranchRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM branches WHERE id = $1", id)
	if err == nil && tag.RowsAffected() == 0 {
		return errors.New("branch not found")
	}
	return err
}

func (r *BranchRepo) GetByID(ctx context.Context, id int64) (*model.Branch, error) {
	query := `SELECT id, name, country, city, address, mobile, telephone, email, latitude, longitude, created_at, updated_at 
	          FROM branches WHERE id=$1`
	var b model.Branch
	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.Name, &b.Country, &b.City, &b.Address, &b.Mobile, &b.Telephone, &b.Email, &b.Latitude, &b.Longitude, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.New("branch not found")
	}
	return &b, err
}

func (r *BranchRepo) GetBranches(ctx context.Context) ([]*model.Branch, error) {
	query := `SELECT id, name, country, city, address, mobile, telephone, email, latitude, longitude, created_at, updated_at 
	          FROM branches ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []*model.Branch

	for rows.Next() {
		var b model.Branch
		err := rows.Scan(
			&b.ID, &b.Name, &b.Country, &b.City, &b.Address, &b.Mobile, &b.Telephone, &b.Email, &b.Latitude, &b.Longitude, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		branches = append(branches, &b)
	}

	return branches, nil
}
