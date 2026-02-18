package service

import (
	"context"
	"errors"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

type AttributeService struct {
	repo *repository.AttributeRepo
}

func NewAttributeService(repo *repository.AttributeRepo) *AttributeService {
	return &AttributeService{repo: repo}
}

func (s *AttributeService) Create(ctx context.Context, u *model.Attribute) error {
	if u.Name == "" {
		return errors.New("attribute name is required")
	}
	return s.repo.Create(ctx, u)
}

func (s *AttributeService) Update(ctx context.Context, u *model.Attribute) error {
	if u.ID == 0 {
		return errors.New("attribute id is required")
	}
	if u.Name == "" {
		return errors.New("attribute name is required")
	}
	return s.repo.Update(ctx, u)
}

func (s *AttributeService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *AttributeService) GetByID(ctx context.Context, id int) (*model.Attribute, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AttributeService) GetAll(ctx context.Context) ([]*model.Attribute, error) {
	return s.repo.GetAll(ctx)
}