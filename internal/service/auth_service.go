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
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication workflows.
type AuthService struct {
	users  *repository.UserRepository
	tokens *repository.RedisTokenRepo
	secret string
}

// NewAuthService constructs an AuthService.
func NewAuthService(u *repository.UserRepository, t *repository.RedisTokenRepo, secret string) *AuthService {
	return &AuthService{
		users:  u,
		tokens: t,
		secret: secret,
	}
}

// Register creates a new user with a hashed password.
func (s *AuthService) Register(ctx context.Context, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.users.Create(ctx, &model.User{
		Email:    email,
		Password: string(hash),
	})
}

// Login authenticates a user and returns access + refresh tokens.
func (s *AuthService) Login(ctx context.Context, email, password string) (*model.User, string, string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	access, err := s.generateJWT(user.ID, 15*time.Minute)
	if err != nil {
		return nil, "", "", err
	}

	refresh, err := utils.GenerateRandomToken()
	if err != nil {
		return nil, "", "", err
	}

	if err := s.tokens.Save(ctx, refresh, user.ID, 7*24*time.Hour); err != nil {
		return nil, "", "", err
	}
	//sanitize output
	user.Password = ""
	// TODO: Remove hard-coded branch_id
	user.BranchID = 1
	return user, access, refresh, nil
}

// Refresh validates a refresh token and issues a new access token.
func (s *AuthService) Refresh(ctx context.Context, token string) (string, error) {
	userIDStr, err := s.tokens.Get(ctx, token)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	uid, _ := strconv.Atoi(userIDStr)
	access, err := s.generateJWT(uid, 15*time.Minute)
	if err != nil {
		return "", err
	}

	return access, nil
}

// generateJWT creates a signed JWT access token.
func (s *AuthService) generateJWT(uid int, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": uid,
		"exp": time.Now().Add(ttl).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.secret))
}
