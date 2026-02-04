// Package repository implements data access logic for customers.
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

// CustomerRepository handles database operations for customers.
type CustomerRepository struct {
	db *pgxpool.Pool
}

// NewCustomerRepo creates a new CustomerRepository.
func NewCustomerRepo(db *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// isCustomerUniqueViolation checks if the error is a Postgres unique constraint violation.
func isCustomerUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}

// Create inserts a new customer into the database.
func (r *CustomerRepository) Create(ctx context.Context, u *model.Customer) error {
	query := `
		INSERT INTO customers (
			name, f_name, l_name, phone, image, email, password,
			street_address, country, city, zip, house_no, apartment_no,
			is_active, is_phone_verified, is_email_verified, app_language, referral_code
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		u.Name, u.FName, u.LName, u.Phone, u.Image, u.Email, u.Password,
		u.StreetAddress, u.Country, u.City, u.Zip, u.HouseNo, u.ApartmentNo,
		u.IsActive, u.IsPhoneVerified, u.IsEmailVerified, u.AppLanguage, u.ReferralCode,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if isCustomerUniqueViolation(err) {
		return fmt.Errorf("customer with email '%s' already exists", u.Email.String)
	}
	return err
}

// FindByID retrieves a customer by ID.
func (r *CustomerRepository) FindByID(ctx context.Context, id int64) (*model.Customer, error) {
	query := `
		SELECT 
			id, name, f_name, l_name, phone, image, email, email_verified_at,
			password, remember_token, created_at, updated_at,
			street_address, country, city, zip, house_no, apartment_no,
			cm_firebase_token, is_active, payment_card_last_four, payment_card_brand,
			payment_card_fawry_token, login_medium, social_id, is_phone_verified,
			temporary_token, is_email_verified, wallet_balance, loyalty_point,
			login_hit_count, is_temp_blocked, temp_block_time, referral_code,
			referred_by, app_language
		FROM customers WHERE id = $1`

	u := &model.Customer{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.FName, &u.LName, &u.Phone, &u.Image, &u.Email, &u.EmailVerifiedAt,
		&u.Password, &u.RememberToken, &u.CreatedAt, &u.UpdatedAt,
		&u.StreetAddress, &u.Country, &u.City, &u.Zip, &u.HouseNo, &u.ApartmentNo,
		&u.CMFirebaseToken, &u.IsActive, &u.PaymentCardLastFour, &u.PaymentCardBrand,
		&u.PaymentCardFawryToken, &u.LoginMedium, &u.SocialID, &u.IsPhoneVerified,
		&u.TemporaryToken, &u.IsEmailVerified, &u.WalletBalance, &u.LoyaltyPoint,
		&u.LoginHitCount, &u.IsTempBlocked, &u.TempBlockTime, &u.ReferralCode,
		&u.ReferredBy, &u.AppLanguage,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("customer with id %d not found", id)
	}
	return u, err
}

// GetByEmail retrieves a customer by email.
func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*model.Customer, error) {
	query := `
		SELECT 
			id, name, f_name, l_name, phone, image, email, email_verified_at,
			password, remember_token, created_at, updated_at,
			street_address, country, city, zip, house_no, apartment_no,
			cm_firebase_token, is_active, payment_card_last_four, payment_card_brand,
			payment_card_fawry_token, login_medium, social_id, is_phone_verified,
			temporary_token, is_email_verified, wallet_balance, loyalty_point,
			login_hit_count, is_temp_blocked, temp_block_time, referral_code,
			referred_by, app_language
		FROM customers WHERE email = $1`

	u := &model.Customer{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.FName, &u.LName, &u.Phone, &u.Image, &u.Email, &u.EmailVerifiedAt,
		&u.Password, &u.RememberToken, &u.CreatedAt, &u.UpdatedAt,
		&u.StreetAddress, &u.Country, &u.City, &u.Zip, &u.HouseNo, &u.ApartmentNo,
		&u.CMFirebaseToken, &u.IsActive, &u.PaymentCardLastFour, &u.PaymentCardBrand,
		&u.PaymentCardFawryToken, &u.LoginMedium, &u.SocialID, &u.IsPhoneVerified,
		&u.TemporaryToken, &u.IsEmailVerified, &u.WalletBalance, &u.LoyaltyPoint,
		&u.LoginHitCount, &u.IsTempBlocked, &u.TempBlockTime, &u.ReferralCode,
		&u.ReferredBy, &u.AppLanguage,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("customer with email '%s' not found", email)
	}
	return u, err
}

// GetByPhone retrieves a customer by phone number.
func (r *CustomerRepository) GetByPhone(ctx context.Context, phone string) (*model.Customer, error) {
	query := `
		SELECT 
			id, name, f_name, l_name, phone, image, email, email_verified_at,
			password, remember_token, created_at, updated_at,
			street_address, country, city, zip, house_no, apartment_no,
			cm_firebase_token, is_active, payment_card_last_four, payment_card_brand,
			payment_card_fawry_token, login_medium, social_id, is_phone_verified,
			temporary_token, is_email_verified, wallet_balance, loyalty_point,
			login_hit_count, is_temp_blocked, temp_block_time, referral_code,
			referred_by, app_language
		FROM customers WHERE phone = $1`

	u := &model.Customer{}
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&u.ID, &u.Name, &u.FName, &u.LName, &u.Phone, &u.Image, &u.Email, &u.EmailVerifiedAt,
		&u.Password, &u.RememberToken, &u.CreatedAt, &u.UpdatedAt,
		&u.StreetAddress, &u.Country, &u.City, &u.Zip, &u.HouseNo, &u.ApartmentNo,
		&u.CMFirebaseToken, &u.IsActive, &u.PaymentCardLastFour, &u.PaymentCardBrand,
		&u.PaymentCardFawryToken, &u.LoginMedium, &u.SocialID, &u.IsPhoneVerified,
		&u.TemporaryToken, &u.IsEmailVerified, &u.WalletBalance, &u.LoyaltyPoint,
		&u.LoginHitCount, &u.IsTempBlocked, &u.TempBlockTime, &u.ReferralCode,
		&u.ReferredBy, &u.AppLanguage,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("customer with phone '%s' not found", phone)
	}
	return u, err
}

// Update updates an existing customer in the database.
func (r *CustomerRepository) Update(ctx context.Context, u *model.Customer) error {
	query := `
		UPDATE customers SET
			name = $1, f_name = $2, l_name = $3, phone = $4, image = $5,
			email = $6, street_address = $7, country = $8, city = $9,
			zip = $10, house_no = $11, apartment_no = $12, cm_firebase_token = $13,
			is_active = $14, is_phone_verified = $15, is_email_verified = $16,
			wallet_balance = $17, loyalty_point = $18, app_language = $19,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		u.Name, u.FName, u.LName, u.Phone, u.Image,
		u.Email, u.StreetAddress, u.Country, u.City,
		u.Zip, u.HouseNo, u.ApartmentNo, u.CMFirebaseToken,
		u.IsActive, u.IsPhoneVerified, u.IsEmailVerified,
		u.WalletBalance, u.LoyaltyPoint, u.AppLanguage,
		u.ID,
	).Scan(&u.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("customer with id %d not found", u.ID)
	}
	if isCustomerUniqueViolation(err) {
		return fmt.Errorf("customer with email '%s' already exists", u.Email.String)
	}
	return err
}

// Delete removes a customer from the database by ID.
func (r *CustomerRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM customers WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer with id %d not found", id)
	}
	return nil
}

// UpdatePassword updates the customer's password.
func (r *CustomerRepository) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	query := `UPDATE customers SET password = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	tag, err := r.db.Exec(ctx, query, hashedPassword, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer with id %d not found", id)
	}
	return nil
}

// UpdateEmailVerification marks email as verified.
func (r *CustomerRepository) UpdateEmailVerification(ctx context.Context, id int64) error {
	query := `UPDATE customers SET is_email_verified = TRUE, email_verified_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer with id %d not found", id)
	}
	return nil
}

// UpdatePhoneVerification marks phone as verified.
func (r *CustomerRepository) UpdatePhoneVerification(ctx context.Context, id int64) error {
	query := `UPDATE customers SET is_phone_verified = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer with id %d not found", id)
	}
	return nil
}

// IncrementLoginHitCount increments login attempts and optionally blocks the customer.
func (r *CustomerRepository) IncrementLoginHitCount(ctx context.Context, id int64, block bool) error {
	var query string
	if block {
		query = `UPDATE customers SET login_hit_count = login_hit_count + 1, is_temp_blocked = TRUE, temp_block_time = CURRENT_TIMESTAMP WHERE id = $1`
	} else {
		query = `UPDATE customers SET login_hit_count = login_hit_count + 1 WHERE id = $1`
	}
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// ResetLoginHitCount resets login attempts after successful login.
func (r *CustomerRepository) ResetLoginHitCount(ctx context.Context, id int64) error {
	query := `UPDATE customers SET login_hit_count = 0, is_temp_blocked = FALSE, temp_block_time = NULL WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// ExistsByEmail checks if a customer exists with the given email.
func (r *CustomerRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE email = $1)`
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

// ExistsByPhone checks if a customer exists with the given phone.
func (r *CustomerRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE phone = $1)`
	err := r.db.QueryRow(ctx, query, phone).Scan(&exists)
	return exists, err
}
