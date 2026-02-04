// Package service implements business logic for customers.
package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Common validation errors.
var (
	ErrCustomerNotFound       = errors.New("customer not found")
	ErrInvalidPhone       = errors.New("phone number is required")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("password must be at least 6 characters")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPhoneAlreadyExists = errors.New("phone number already exists")
	ErrCustomerInactive       = errors.New("customer account is inactive")
	ErrCustomerBlocked        = errors.New("customer account is temporarily blocked")
)

// CustomerService handles business logic for customer operations.
type CustomerService struct {
	repo *repository.CustomerRepository
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// Create validates and creates a new customer.
func (s *CustomerService) Create(ctx context.Context, req *model.CustomerCreateRequest) (*model.Customer, error) {
	// Validate required fields
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check if email already exists
	if req.Email != "" {
		exists, err := s.repo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailAlreadyExists
		}
	}

	// Check if phone already exists
	exists, err := s.repo.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPhoneAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Build customer entity
	customer := &model.Customer{
		Name:        toNullString(req.Name),
		FName:       toNullString(req.FName),
		LName:       toNullString(req.LName),
		Phone:       req.Phone,
		Email:       toNullString(req.Email),
		Password:    string(hashedPassword),
		Image:       "def.png",
		IsActive:    true,
		AppLanguage: "en",
	}

	// Create customer in database
	if err := s.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
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

// Update validates and updates an existing customer.
func (s *CustomerService) Update(ctx context.Context, id int64, req *model.CustomerUpdateRequest) (*model.Customer, error) {
	// Fetch existing customer
	customer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		customer.Name = toNullString(req.Name)
	}
	if req.FName != "" {
		customer.FName = toNullString(req.FName)
	}
	if req.LName != "" {
		customer.LName = toNullString(req.LName)
	}
	if req.Phone != "" {
		// Check if new phone conflicts with another customer
		if req.Phone != customer.Phone {
			exists, err := s.repo.ExistsByPhone(ctx, req.Phone)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrPhoneAlreadyExists
			}
		}
		customer.Phone = req.Phone
	}
	if req.Image != "" {
		customer.Image = req.Image
	}
	if req.StreetAddress != "" {
		customer.StreetAddress = toNullString(req.StreetAddress)
	}
	if req.Country != "" {
		customer.Country = toNullString(req.Country)
	}
	if req.City != "" {
		customer.City = toNullString(req.City)
	}
	if req.Zip != "" {
		customer.Zip = toNullString(req.Zip)
	}
	if req.HouseNo != "" {
		customer.HouseNo = toNullString(req.HouseNo)
	}
	if req.ApartmentNo != "" {
		customer.ApartmentNo = toNullString(req.ApartmentNo)
	}
	if req.AppLanguage != "" {
		customer.AppLanguage = req.AppLanguage
	}

	// Save updates
	if err := s.repo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
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
