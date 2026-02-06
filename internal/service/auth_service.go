// Package service contains core business logic.
package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"github.com/projuktisheba/pse-api-v1/pkg/identity"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication workflows for both employees and customers.
type AuthService struct {
	employees *repository.EmployeeRepository
	customers *repository.CustomerRepository
	tokens    *repository.RedisTokenRepo
	secret    string
}

// NewAuthService constructs an AuthService.
func NewAuthService(e *repository.EmployeeRepository, u *repository.CustomerRepository, t *repository.RedisTokenRepo, secret string) *AuthService {
	return &AuthService{
		employees: e,
		customers: u,
		tokens:    t,
		secret:    secret,
	}
}

// ==================== EMPLOYEE (ADMIN) AUTH ====================

// EmployeeRegister creates a new employee with a hashed password.
func (s *AuthService) EmployeeRegister(ctx context.Context, email, password, name, mobile, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.employees.Create(ctx, &model.Employee{
		Name:       name,
		Email:      email,
		Password:   string(hash),
		Mobile:     mobile,
		Role:       role,
		IsActive:   true,
		IsVerified: false,
	})
}

// EmployeeLogin authenticates an employee and returns access + refresh tokens.
func (s *AuthService) EmployeeLogin(ctx context.Context, email, password string) (*model.Employee, string, string, error) {
	employee, err := s.employees.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}

	if !employee.IsActive {
		return nil, "", "", errors.New("account is inactive")
	}

	if bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password)) != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	// Generate tokens with "employee" type prefix
	access, err := s.generateJWT(employee.ID, "employee", employee.Role, 15*time.Minute)
	if err != nil {
		return nil, "", "", err
	}

	refresh, err := utils.GenerateRandomToken()
	if err != nil {
		return nil, "", "", err
	}

	// Store refresh token with employee prefix
	tokenKey := "employee:" + strconv.Itoa(employee.ID)
	if err := s.tokens.Save(ctx, refresh, tokenKey, 7*24*time.Hour); err != nil {
		return nil, "", "", err
	}

	// Sanitize output
	employee.Password = ""
	return employee, access, refresh, nil
}

// EmployeeRefresh validates a refresh token and issues a new access token for employee.
func (s *AuthService) EmployeeRefresh(ctx context.Context, token string) (string, error) {
	data, err := s.tokens.Get(ctx, token)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Parse employee ID from stored data
	var uid int
	if _, err := utils.Sscanf(data, "employee:%d", &uid); err != nil {
		return "", errors.New("invalid token type")
	}

	// Fetch employee to get role
	employee, err := s.employees.FindByID(ctx, uid)
	if err != nil {
		return "", errors.New("customer not found")
	}

	access, err := s.generateJWT(uid, "employee", employee.Role, 15*time.Minute)
	if err != nil {
		return "", err
	}

	return access, nil
}

// ==================== CUSTOMER (CUSTOMER) AUTH ====================

// CustomerRegister creates a new customer customer with a hashed password.
func (s *AuthService) CustomerRegister(ctx context.Context, req *model.CustomerCreateRequest) (*model.Customer, error) {
	// Check if email already exists
	if req.Email != "" {
		exists, err := s.customers.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already exists")
		}
	}

	// Check if phone already exists
	exists, err := s.customers.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("phone number already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	customer := &model.Customer{
		Name:        model.ToNullString(req.Name),
		FName:       model.ToNullString(req.FName),
		LName:       model.ToNullString(req.LName),
		Phone:       req.Phone,
		Email:       model.ToNullString(req.Email),
		Password:    string(hash),
		Image:       "def.png",
		IsActive:    true,
		AppLanguage: "en",
	}

	if err := s.customers.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// CustomerLogin authenticates a customer customer and returns access + refresh tokens.
func (s *AuthService) CustomerLogin(ctx context.Context, emailOrPhone, password string) (*model.Customer, string, string, error) {
	var customer *model.Customer
	var err error

	// Find by email if isEmail is true
	isEmail := identity.IsEmail(emailOrPhone)

	if isEmail {
		customer, err = s.customers.GetByEmail(ctx, emailOrPhone)
		if err != nil {
			return nil, "", "", errors.New("invalid credentials")
		}
	} else {
		customer, err = s.customers.GetByPhone(ctx, emailOrPhone)
		if err != nil {
			return nil, "", "", errors.New("invalid credentials")
		}
	}

	if !customer.IsActive {
		return nil, "", "", errors.New("account is inactive")
	}

	if customer.IsTempBlocked {
		return nil, "", "", errors.New("account is temporarily blocked")
	}

	if bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password)) != nil {
		// Increment login attempts
		_ = s.customers.IncrementLoginHitCount(ctx, customer.ID, customer.LoginHitCount >= 4)
		return nil, "", "", errors.New("invalid credentials")
	}

	// Reset login attempts on successful login
	_ = s.customers.ResetLoginHitCount(ctx, customer.ID)

	// Generate tokens with "customer" type prefix
	access, err := s.generateCustomerJWT(customer.ID, 15*time.Minute)
	if err != nil {
		return nil, "", "", err
	}

	refresh, err := utils.GenerateRandomToken()
	if err != nil {
		return nil, "", "", err
	}

	// Store refresh token with customer prefix
	tokenKey := "customer:" + strconv.FormatInt(customer.ID, 10)
	if err := s.tokens.Save(ctx, refresh, tokenKey, 7*24*time.Hour); err != nil {
		return nil, "", "", err
	}

	// Sanitize output
	customer.Password = ""
	return customer, access, refresh, nil
}

// CustomerRefresh validates a refresh token and issues a new access token for customer customer.
func (s *AuthService) CustomerRefresh(ctx context.Context, token string) (string, error) {
	data, err := s.tokens.Get(ctx, token)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Parse customer ID from stored data
	var uid int64
	if _, err := utils.Sscanf(data, "customer:%d", &uid); err != nil {
		return "", errors.New("invalid token type")
	}

	access, err := s.generateCustomerJWT(uid, 15*time.Minute)
	if err != nil {
		return "", err
	}

	return access, nil
}

// ==================== TOKEN GENERATION ====================

// generateJWT creates a signed JWT access token for employees.
func (s *AuthService) generateJWT(uid int, customerType, role string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  uid,
		"type": customerType,
		"role": role,
		"exp":  time.Now().Add(ttl).Unix(),
		"iat":  time.Now().Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.secret))
}

// generateCustomerJWT creates a signed JWT access token for customer customers.
func (s *AuthService) generateCustomerJWT(uid int64, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  uid,
		"type": "customer",
		"exp":  time.Now().Add(ttl).Unix(),
		"iat":  time.Now().Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.secret))
}
