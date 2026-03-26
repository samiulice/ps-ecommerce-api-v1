package model

import (
	"time"
)

type Employee struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	BranchID    int64     `json:"branch_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Mobile      string    `json:"mobile"`
	RoleID      int64     `json:"role_id"`
	Role        string    `json:"role"`
	RoleName    string    `json:"role_name,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
	IsActive    bool      `json:"isActive"`
	IsVerified  bool      `json:"isVerified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
