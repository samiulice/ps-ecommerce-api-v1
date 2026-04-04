package service

import (
	"context"
	"errors"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

type UnitService struct {
	repo *repository.UnitRepo
}

func NewUnitService(repo *repository.UnitRepo) *UnitService {
	return &UnitService{repo: repo}
}

func (s *UnitService) Create(ctx context.Context, u *model.Unit) error {
	if u.Name == "" {
		return errors.New("unit name is required")
	}
	if u.Symbol == "" {
		return errors.New("unit symbol is required")
	}
	return s.repo.Create(ctx, u)
}

func (s *UnitService) Update(ctx context.Context, u *model.Unit) error {
	if u.ID == 0 {
		return errors.New("unit id is required")
	}
	if u.Name == "" {
		return errors.New("unit name is required")
	}
	return s.repo.Update(ctx, u)
}

func (s *UnitService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *UnitService) GetByID(ctx context.Context, id int) (*model.Unit, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UnitService) GetAll(ctx context.Context) ([]*model.Unit, error) {
	return s.repo.GetAll(ctx)
}
