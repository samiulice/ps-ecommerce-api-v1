package service

import (
	"context"
	"errors"
	"strings"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

var (
	ErrEmployeeNotFound = errors.New("employee not found")
	ErrEmployeeEmail    = errors.New("employee email is required")
	ErrEmployeeRole     = errors.New("employee role is required")
)

type EmployeeService struct {
	repo  *repository.EmployeeRepository
	roles *repository.RoleRepository
}

func NewEmployeeService(repo *repository.EmployeeRepository, roles *repository.RoleRepository) *EmployeeService {
	return &EmployeeService{repo: repo, roles: roles}
}

func (s *EmployeeService) List(ctx context.Context) ([]model.Employee, error) {
	return s.repo.List(ctx)
}

func (s *EmployeeService) GetByID(ctx context.Context, id int) (*model.Employee, error) {
	employee, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}
	employee.Password = ""
	return employee, nil
}

func (s *EmployeeService) Update(ctx context.Context, id int, req model.EmployeeAdminUpdateRequest) (*model.Employee, error) {
	if strings.TrimSpace(req.Email) == "" {
		return nil, ErrEmployeeEmail
	}
	if req.RoleID <= 0 {
		return nil, ErrEmployeeRole
	}

	current, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	role, err := s.roles.FindByID(ctx, req.RoleID)
	if err != nil {
		return nil, ErrEmployeeRole
	}

	exists, err := s.repo.ExistsByEmailExcludingID(ctx, strings.TrimSpace(req.Email), id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	current.Name = strings.TrimSpace(req.Name)
	current.Email = strings.TrimSpace(req.Email)
	current.Mobile = strings.TrimSpace(req.Mobile)
	if req.BranchID > 0 {
		current.BranchID = req.BranchID
	}
	current.RoleID = role.ID
	current.Role = role.Slug
	current.RoleName = role.Name
	current.Permissions = permissionKeys(role.Permissions)
	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateAdmin(ctx, current); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *EmployeeService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return ErrEmployeeNotFound
		}
		return err
	}
	return nil
}
