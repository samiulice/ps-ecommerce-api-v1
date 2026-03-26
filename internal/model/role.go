package model

import "time"

type Permission struct {
	ID          int64     `json:"id"`
	Key         string    `json:"key"`
	DisplayName string    `json:"display_name"`
	Module      string    `json:"module"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Role struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Slug        string       `json:"slug"`
	Description string       `json:"description,omitempty"`
	IsActive    bool         `json:"is_active"`
	Permissions []Permission `json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type RoleCreateRequest struct {
	Name           string   `json:"name"`
	Slug           string   `json:"slug"`
	Description    string   `json:"description"`
	IsActive       *bool    `json:"is_active"`
	PermissionKeys []string `json:"permission_keys"`
}

type RoleUpdateRequest struct {
	Name           string   `json:"name"`
	Slug           string   `json:"slug"`
	Description    string   `json:"description"`
	IsActive       *bool    `json:"is_active"`
	PermissionKeys []string `json:"permission_keys"`
}

type EmployeeAdminUpdateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	BranchID int64  `json:"branch_id"`
	RoleID   int64  `json:"role_id"`
	IsActive *bool  `json:"is_active"`
}
