package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleNameRequired  = errors.New("role name is required")
	ErrRoleSlugRequired  = errors.New("role slug is required")
	ErrRoleDeleteBlocked = errors.New("role cannot be deleted")
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) List(ctx context.Context) ([]model.Role, error) {
	return s.repo.List(ctx)
}

func (s *RoleService) ListPermissions(ctx context.Context) ([]model.Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *RoleService) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	role, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	return role, err
}

func (s *RoleService) Create(ctx context.Context, req model.RoleCreateRequest) (*model.Role, error) {
	role, keys, err := normalizeRoleRequest(req.Name, req.Slug, req.Description, req.IsActive, req.PermissionKeys)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, role, keys); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, role.ID)
}

func (s *RoleService) Update(ctx context.Context, id int64, req model.RoleUpdateRequest) (*model.Role, error) {
	role, keys, err := normalizeRoleRequest(req.Name, req.Slug, req.Description, req.IsActive, req.PermissionKeys)
	if err != nil {
		return nil, err
	}
	role.ID = id
	if err := s.repo.Update(ctx, role, keys); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *RoleService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return ErrRoleNotFound
		}
		if strings.Contains(strings.ToLower(err.Error()), "assigned to") {
			return ErrRoleDeleteBlocked
		}
		return err
	}
	return nil
}

func normalizeRoleRequest(name, slug, description string, isActive *bool, permissionKeys []string) (*model.Role, []string, error) {
	name = strings.TrimSpace(name)
	slug = strings.ToLower(strings.TrimSpace(slug))
	if name == "" {
		return nil, nil, ErrRoleNameRequired
	}
	if slug == "" {
		return nil, nil, ErrRoleSlugRequired
	}
	active := true
	if isActive != nil {
		active = *isActive
	}
	return &model.Role{
		Name:        name,
		Slug:        slug,
		Description: strings.TrimSpace(description),
		IsActive:    active,
	}, permissionKeys, nil
}
