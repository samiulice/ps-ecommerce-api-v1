package service

import (
	"context"
	"errors"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

type BranchService struct {
	repo *repository.BranchRepo
}

func NewBranchService(repo *repository.BranchRepo) *BranchService {
	return &BranchService{repo: repo}
}

func (s *BranchService) Create(ctx context.Context, b *model.Branch) error {
	if b.Name == "" {
		return errors.New("branch name is required")
	}
	if b.City == "" {
		return errors.New("city is required")
	}
	if b.Address == "" {
		return errors.New("address is required")
	}
	return s.repo.Create(ctx, b)
}

func (s *BranchService) Update(ctx context.Context, b *model.Branch) error {
	if b.ID == 0 {
		return errors.New("id is required")
	}
	if b.Name == "" {
		return errors.New("branch name is required")
	}
	return s.repo.Update(ctx, b)
}

func (s *BranchService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *BranchService) GetByID(ctx context.Context, id int64) (*model.Branch, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BranchService) GetBranches(ctx context.Context) ([]*model.Branch, error) {
	return s.repo.GetBranches(ctx)
}