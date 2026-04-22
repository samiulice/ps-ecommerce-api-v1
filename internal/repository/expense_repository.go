package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type ExpenseRepository struct {
	db *pgxpool.Pool
}

func NewExpenseRepository(db *pgxpool.Pool) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// category
func (r *ExpenseRepository) CreateCategory(ctx context.Context, c *model.ExpenseCategory) error {
	query := `INSERT INTO expense_categories (name, branch_id) VALUES ($1, $2) RETURNING id, total_amount, created_at`
	return r.db.QueryRow(ctx, query, c.Name, c.BranchID).Scan(&c.ID, &c.TotalAmount, &c.CreatedAt)
}

func (r *ExpenseRepository) GetCategories(ctx context.Context, branchID *int64) ([]model.ExpenseCategory, error) {
	query := `SELECT id, name, branch_id, total_amount, created_at FROM expense_categories WHERE branch_id = $1 OR branch_id IS NULL`
	rows, err := r.db.Query(ctx, query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.ExpenseCategory
	for rows.Next() {
		var c model.ExpenseCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.BranchID, &c.TotalAmount, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

// expense
func (r *ExpenseRepository) CreateExpense(ctx context.Context, e *model.Expense) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO expenses (category_id, branch_id, amount, expense_date, description, created_by) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	err = tx.QueryRow(ctx, query, e.CategoryID, e.BranchID, e.Amount, e.ExpenseDate, e.Description, e.CreatedBy).Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return err
	}

	updateCatQuery := `UPDATE expense_categories SET total_amount = total_amount + $1 WHERE id = $2`
	_, err = tx.Exec(ctx, updateCatQuery, e.Amount, e.CategoryID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *ExpenseRepository) GetExpenses(ctx context.Context, branchID int64) ([]model.Expense, error) {
	query := `
		SELECT e.id, e.category_id, c.name, e.branch_id, e.amount, e.expense_date, e.description, e.created_by, e.created_at
		FROM expenses e JOIN expense_categories c ON e.category_id = c.id
		WHERE e.branch_id = $1
		ORDER BY e.expense_date DESC
	`
	rows, err := r.db.Query(ctx, query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var e model.Expense
		if err := rows.Scan(&e.ID, &e.CategoryID, &e.CategoryName, &e.BranchID, &e.Amount, &e.ExpenseDate, &e.Description, &e.CreatedBy, &e.CreatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}
