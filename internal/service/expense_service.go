package service

import (
	"context"
	"errors"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

type ExpenseService struct {
	repo *repository.ExpenseRepository
}

func NewExpenseService(repo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) CreateCategory(ctx context.Context, c *model.ExpenseCategory) error {
	if c.Name == "" {
		return errors.New("category name is required")
	}
	return s.repo.CreateCategory(ctx, c)
}

func (s *ExpenseService) GetCategories(ctx context.Context, branchID *int64) ([]model.ExpenseCategory, error) {
	return s.repo.GetCategories(ctx, branchID)
}

func (s *ExpenseService) CreateExpense(ctx context.Context, e *model.Expense) error {
	if e.Amount <= 0 {
		return errors.New("expense amount must be greater than zero")
	}
	if e.CategoryID <= 0 {
		return errors.New("expense category is required")
	}
	if e.BranchID <= 0 {
		return errors.New("branch id is required")
	}
	return s.repo.CreateExpense(ctx, e)
}

func (s *ExpenseService) GetExpenses(ctx context.Context, branchID int64) ([]model.Expense, error) {
	return s.repo.GetExpenses(ctx, branchID)
}
