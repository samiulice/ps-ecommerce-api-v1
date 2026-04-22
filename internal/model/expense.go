package model

import "time"

type ExpenseCategory struct {
	ID          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	BranchID    *int64     `json:"branch_id,omitempty" db:"branch_id"`
	TotalAmount float64    `json:"total_amount" db:"total_amount"`
	CreatedAt   *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type Expense struct {
	ID          int64      `json:"id" db:"id"`
	CategoryID  int64      `json:"category_id" db:"category_id"`
	BranchID    int64      `json:"branch_id" db:"branch_id"`
	Amount      float64    `json:"amount" db:"amount"`
	ExpenseDate string     `json:"expense_date" db:"expense_date"`
	Description string     `json:"description,omitempty" db:"description"`
	CreatedBy   int64      `json:"created_by,omitempty" db:"created_by"`
	CreatedAt   *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`

	CategoryName string `json:"category_name,omitempty" db:"category_name"`
}
