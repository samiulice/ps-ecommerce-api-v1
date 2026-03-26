package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type SupplierRepo struct {
	db *pgxpool.Pool
}

func NewSupplierRepo(db *pgxpool.Pool) *SupplierRepo {
	return &SupplierRepo{db: db}
}

func (r *SupplierRepo) Create(ctx context.Context, s *model.Supplier) error {
	query := `
        INSERT INTO suppliers (
            supplier_code, name, company_name, contact_person, phone, email,
            website, tax_id, trade_license_no, payment_terms,
            credit_limit, outstanding_balance, lead_time_days, rating,
            street_address, country, city, zip, notes, is_active
        ) VALUES (
            $1, $2, $3, $4, $5, $6,
            $7, $8, $9, $10,
            $11, $12, $13, $14,
            $15, $16, $17, $18, $19, $20
        ) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		s.SupplierCode, s.Name, s.CompanyName, s.ContactPerson, s.Phone, s.Email,
		s.Website, s.TaxID, s.TradeLicenseNo, s.PaymentTerms,
		s.CreditLimit, s.OutstandingBalance, s.LeadTimeDays, s.Rating,
		s.StreetAddress, s.Country, s.City, s.Zip, s.Notes, s.IsActive,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)

	if mappedErr := mapSupplierUniqueViolation(err); mappedErr != nil {
		return mappedErr
	}
	return err
}

func (r *SupplierRepo) FindByID(ctx context.Context, id int64) (*model.Supplier, error) {
	query := `
        SELECT
            id, supplier_code, name, company_name, contact_person, phone, email,
            website, tax_id, trade_license_no, payment_terms,
            credit_limit, outstanding_balance, lead_time_days, rating,
            street_address, country, city, zip, notes, is_active,
            created_at, updated_at
        FROM suppliers
        WHERE id = $1`

	s := &model.Supplier{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.SupplierCode, &s.Name, &s.CompanyName, &s.ContactPerson, &s.Phone, &s.Email,
		&s.Website, &s.TaxID, &s.TradeLicenseNo, &s.PaymentTerms,
		&s.CreditLimit, &s.OutstandingBalance, &s.LeadTimeDays, &s.Rating,
		&s.StreetAddress, &s.Country, &s.City, &s.Zip, &s.Notes, &s.IsActive,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("supplier with id %d not found", id)
	}
	return s, err
}

func (r *SupplierRepo) Update(ctx context.Context, s *model.Supplier) error {
	query := `
        UPDATE suppliers SET
            supplier_code = $1,
            name = $2,
            company_name = $3,
            contact_person = $4,
            phone = $5,
            email = $6,
            website = $7,
            tax_id = $8,
            trade_license_no = $9,
            payment_terms = $10,
            credit_limit = $11,
            outstanding_balance = $12,
            lead_time_days = $13,
            rating = $14,
            street_address = $15,
            country = $16,
            city = $17,
            zip = $18,
            notes = $19,
            is_active = $20,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $21
        RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		s.SupplierCode, s.Name, s.CompanyName, s.ContactPerson, s.Phone, s.Email,
		s.Website, s.TaxID, s.TradeLicenseNo, s.PaymentTerms,
		s.CreditLimit, s.OutstandingBalance, s.LeadTimeDays, s.Rating,
		s.StreetAddress, s.Country, s.City, s.Zip, s.Notes, s.IsActive,
		s.ID,
	).Scan(&s.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("supplier with id %d not found", s.ID)
	}
	if mappedErr := mapSupplierUniqueViolation(err); mappedErr != nil {
		return mappedErr
	}
	return err
}

func (r *SupplierRepo) UpdateAccountStatus(ctx context.Context, status bool, supplierID int64) error {
	query := `
        UPDATE suppliers SET
            is_active = $1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $2`

	tag, err := r.db.Exec(ctx, query, status, supplierID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("supplier with id %d not found", supplierID)
	}
	return nil
}

func (r *SupplierRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM suppliers WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("supplier with id %d not found", id)
	}
	return nil
}

func (r *SupplierRepo) ExistsByCode(ctx context.Context, supplierCode string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM suppliers WHERE supplier_code = $1)`
	err := r.db.QueryRow(ctx, query, supplierCode).Scan(&exists)
	return exists, err
}

func (r *SupplierRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM suppliers WHERE email = $1)`
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *SupplierRepo) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM suppliers WHERE phone = $1)`
	err := r.db.QueryRow(ctx, query, phone).Scan(&exists)
	return exists, err
}

func (r *SupplierRepo) List(ctx context.Context, filter model.SupplierFilter) ([]model.Supplier, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		conditions = append(conditions, fmt.Sprintf(
			`(supplier_code ILIKE $%d OR name ILIKE $%d OR company_name ILIKE $%d OR contact_person ILIKE $%d OR phone ILIKE $%d OR email ILIKE $%d)`,
			argIndex, argIndex, argIndex, argIndex, argIndex, argIndex,
		))
		args = append(args, searchPattern)
		argIndex++
	}

	if filter.CheckAccountStatus {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, filter.IsActive)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM suppliers %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count suppliers: %w", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (filter.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
        SELECT
            id, supplier_code, name, company_name, contact_person, phone, email,
            website, tax_id, trade_license_no, payment_terms,
            credit_limit, outstanding_balance, lead_time_days, rating,
            street_address, country, city, zip, notes, is_active,
            created_at, updated_at
        FROM suppliers
        %s
        ORDER BY id DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query suppliers: %w", err)
	}
	defer rows.Close()

	var suppliers []model.Supplier
	for rows.Next() {
		var s model.Supplier
		err := rows.Scan(
			&s.ID, &s.SupplierCode, &s.Name, &s.CompanyName, &s.ContactPerson, &s.Phone, &s.Email,
			&s.Website, &s.TaxID, &s.TradeLicenseNo, &s.PaymentTerms,
			&s.CreditLimit, &s.OutstandingBalance, &s.LeadTimeDays, &s.Rating,
			&s.StreetAddress, &s.Country, &s.City, &s.Zip, &s.Notes, &s.IsActive,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan supplier: %w", err)
		}
		suppliers = append(suppliers, s)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating supplier rows: %w", err)
	}

	return suppliers, total, nil
}

func mapSupplierUniqueViolation(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil
	}

	switch pgErr.ConstraintName {
	case "suppliers_supplier_code_unique":
		return fmt.Errorf("supplier_code already exists")
	case "suppliers_phone_unique":
		return fmt.Errorf("phone already exists")
	case "suppliers_email_unique":
		return fmt.Errorf("email already exists")
	default:
		return fmt.Errorf("duplicate supplier record")
	}
}
