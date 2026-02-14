// Package service implements business logic for customers.
package service

import (
	"context"
	"database/sql"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// Common validation errors.
var (
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrInvalidPhone       = errors.New("phone number is required")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("password must be at least 6 characters")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPhoneAlreadyExists = errors.New("phone number already exists")
	ErrCustomerInactive   = errors.New("customer account is inactive")
	ErrCustomerBlocked    = errors.New("customer account is temporarily blocked")
)

// CustomerService handles business logic for customer operations.
type CustomerService struct {
	repo *repository.CustomerRepository
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// Create validates and creates a new customer (Multipart).
func (s *CustomerService) Create(ctx context.Context, c *model.Customer, file multipart.File, header *multipart.FileHeader) (*model.Customer, error) {
	// 1. Basic Validation
	if c.Phone == "" {
		return nil, ErrInvalidPhone
	}
	// Password is required for creation
	if len(c.Password) != 0 && len(c.Password) < 6 {
		return nil, ErrInvalidPassword
	}

	// 2. Check Duplicates
	if c.Email.Valid && c.Email.String != "" {
		exists, err := s.repo.ExistsByEmail(ctx, c.Email.String)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailAlreadyExists
		}
	}
	exists, err := s.repo.ExistsByPhone(ctx, c.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPhoneAlreadyExists
	}

	// 3. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	c.Password = string(hashedPassword)

	// 4. Handle Image URL generation
	if file != nil {
		// Use Phone or Name as the unique identifier for the filename
		ext := strings.ToLower(filepath.Ext(header.Filename))
		c.Image = utils.GetCustomerImageURL(c.Phone, ext)
	} else {
		c.Image = "def.png"
	}

	// Set Defaults
	c.AppLanguage = "en"

	// 5. Insert into Database
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	// 6. Save Image File (if database insert succeeded)
	if file != nil {
		defer file.Close()
		_, err := utils.SaveMultipartImage(file, header, utils.GetCustomerFolderPath(""), c.Phone)
		if err != nil {
			// Optional: Log error, but typically we return it so the controller knows
			return nil, err
		}
	}

	return c, nil
}

// Update validates and updates an existing customer (Multipart).
func (s *CustomerService) Update(ctx context.Context, c *model.Customer, file multipart.File, header *multipart.FileHeader) (*model.Customer, error) {
	if c.ID <= 0 {
		return nil, errors.New("customer ID is required")
	}

	// 1. Fetch Existing Customer
	existing, err := s.GetByID(ctx, c.ID)
	if err != nil {
		return nil, err
	}

	// 2. Validate Uniqueness (if Phone/Email changed)
	if c.Phone != existing.Phone {
		exists, err := s.repo.ExistsByPhone(ctx, c.Phone)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrPhoneAlreadyExists
		}
	}
	// Note: Add Email uniqueness check here if needed

	// 3. Handle Password
	// If new password is provided, hash it. Otherwise, keep the old one.
	if c.Password != "" {
		if len(c.Password) < 6 {
			return nil, ErrInvalidPassword
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		c.Password = string(hashedPassword)
	} else {
		c.Password = existing.Password
	}

	// 4. Handle Image URL
	if file != nil {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		c.Image = utils.GetCustomerImageURL(c.Phone, ext)
	} else {
		c.Image = existing.Image // Keep existing image if no new file uploaded
	}

	// 5. Update Database
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}

	// 6. Save New Image & Cleanup Old
	if file != nil {
		defer file.Close()
		_, err := utils.SaveMultipartImage(file, header, utils.GetCustomerFolderPath(""), c.Phone)
		if err != nil {
			return nil, err
		}

		// Delete old image if it wasn't the default and it's different
		if existing.Image != "def.png" && existing.Image != c.Image {
			utils.DeleteFile(utils.GetCustomerFolderPath(filepath.Base(existing.Image)))
		}
	}

	return c, nil
}

// ListCustomers retrieves customers with filters
func (s *CustomerService) ListCustomers(ctx context.Context, filter model.CustomerFilter) ([]*model.CustomerResponse, int64, error) {
	customers, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []*model.CustomerResponse
	for _, customer := range customers {
		responses = append(responses, customer.ToResponse())
	}

	return responses, total, nil
}

// GetByID retrieves a customer by ID.
func (s *CustomerService) GetByID(ctx context.Context, id int64) (*model.Customer, error) {
	if id <= 0 {
		return nil, ErrCustomerNotFound
	}

	customer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}

	return customer, nil
}

// UpdateAccountStatus updates an existing customer account status.
func (s *CustomerService) UpdateAccountStatus(ctx context.Context, accountStatus bool, customerId int64) error {
	// Save updates
	return s.repo.UpdateAccountStatus(ctx, accountStatus, customerId)
}

// Delete removes a customer by ID.
func (s *CustomerService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrCustomerNotFound
	}

	return s.repo.Delete(ctx, id)
}

// GetByEmail retrieves a customer by email.
func (s *CustomerService) GetByEmail(ctx context.Context, email string) (*model.Customer, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	customer, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}

	return customer, nil
}

// GetByPhone retrieves a customer by phone number.
func (s *CustomerService) GetByPhone(ctx context.Context, phone string) (*model.Customer, error) {
	if phone == "" {
		return nil, ErrInvalidPhone
	}

	customer, err := s.repo.GetByPhone(ctx, phone)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}

	return customer, nil
}

// VerifyEmail marks a customer's email as verified.
func (s *CustomerService) VerifyEmail(ctx context.Context, id int64) error {
	return s.repo.UpdateEmailVerification(ctx, id)
}

// VerifyPhone marks a customer's phone as verified.
func (s *CustomerService) VerifyPhone(ctx context.Context, id int64) error {
	return s.repo.UpdatePhoneVerification(ctx, id)
}

// ChangePassword updates the customer's password.
func (s *CustomerService) ChangePassword(ctx context.Context, id int64, newPassword string) error {
	if len(newPassword) < 6 {
		return ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, id, string(hashedPassword))
}

// validateCreateRequest validates customer creation request.
func (s *CustomerService) validateCreateRequest(req *model.CustomerCreateRequest) error {
	if req.Phone == "" {
		return ErrInvalidPhone
	}
	if len(req.Password) < 6 {
		return ErrInvalidPassword
	}
	if req.Email != "" && !isValidEmail(req.Email) {
		return ErrInvalidEmail
	}
	return nil
}

// Helper functions

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
func isValidEmail(email string) bool {
	// Basic email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
