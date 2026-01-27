package service

import (
	"context"
	"errors"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepo
}

func NewCategoryService(repo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetTree returns the full hierarchy
func (s *CategoryService) GetTree(ctx context.Context, onlyActive bool) ([]model.Category, error) {
	return s.repo.ListFullTree(ctx, onlyActive)
}

// --- Level 1 ---
func (s *CategoryService) Create(ctx context.Context, c *model.Category) error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	return s.repo.Create(ctx, c)
}
func (s *CategoryService) Update(ctx context.Context, c *model.Category) error {
	if c.ID == 0 {
		return errors.New("id is required")
	}
	return s.repo.Update(ctx, c)
}
func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
func (s *CategoryService) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	return s.repo.GetByID(ctx, id)
}

// --- Level 2 ---
func (s *CategoryService) CreateSub(ctx context.Context, sc *model.SubCategory) error {
	if sc.CategoryID == 0 {
		return errors.New("parent category_id is required")
	}
	if sc.Name == "" {
		return errors.New("name is required")
	}
	return s.repo.CreateSub(ctx, sc)
}
func (s *CategoryService) UpdateSub(ctx context.Context, sc *model.SubCategory) error {
	if sc.ID == 0 {
		return errors.New("id is required")
	}
	return s.repo.UpdateSub(ctx, sc)
}
func (s *CategoryService) DeleteSub(ctx context.Context, id int64) error {
	return s.repo.DeleteSub(ctx, id)
}
func (s *CategoryService) GetSubByID(ctx context.Context, id int64) (*model.SubCategory, error) {
	return s.repo.GetSubByID(ctx, id)
}

// --- Level 3 ---
func (s *CategoryService) CreateSubSub(ctx context.Context, ssc *model.SubSubCategory) error {
	if ssc.SubCategoryID == 0 {
		return errors.New("parent sub_category_id is required")
	}
	if ssc.Name == "" {
		return errors.New("name is required")
	}
	return s.repo.CreateSubSub(ctx, ssc)
}
func (s *CategoryService) UpdateSubSub(ctx context.Context, ssc *model.SubSubCategory) error {
	if ssc.ID == 0 {
		return errors.New("id is required")
	}
	return s.repo.UpdateSubSub(ctx, ssc)
}
func (s *CategoryService) DeleteSubSub(ctx context.Context, id int64) error {
	return s.repo.DeleteSubSub(ctx, id)
}
func (s *CategoryService) GetSubSubByID(ctx context.Context, id int64) (*model.SubSubCategory, error) {
	return s.repo.GetSubSubByID(ctx, id)
}
