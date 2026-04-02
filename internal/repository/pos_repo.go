package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type POSRepo struct {
	db   *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

func NewPOSRepo(db *pgxpool.Pool) *POSRepo {
	return &POSRepo{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *POSRepo) CreateSale(ctx context.Context, sale *model.POSSale) (*model.POSSale, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// insert pos sale
	query := `
		INSERT INTO pos_sales (
			reference_no, customer_id, branch_id, sale_type, subtotal, discount,
			tax, total, amount_paid, payment_method, sale_date, sale_note
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, $11)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		sale.ReferenceNo, sale.CustomerID, sale.BranchID, sale.SaleType, sale.Subtotal, sale.Discount,
		sale.Tax, sale.Total, sale.AmountPaid, sale.PaymentMethod, sale.SaleNote,
	).Scan(&sale.ID, &sale.CreatedAt, &sale.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert pos sale: %w", err)
	}

	// insert pos sale items
	if len(sale.Items) > 0 {
		var values []string
		var args []interface{}

		for i, item := range sale.Items {
			offset := i * 9
			placeholders := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7, offset+8, offset+9)
			values = append(values, placeholders)
			args = append(args,
				sale.ID, item.ProductID, item.ProductVariationID, item.ProductName, item.Quantity,
				item.UnitPrice, item.Subtotal, item.TaxAmount, item.Total,
			)
		}

		itemsQuery := fmt.Sprintf(`
			INSERT INTO pos_sale_items (
				pos_sale_id, product_id, product_variation_id, product_name, quantity, 
				unit_price, subtotal, tax_amount, total
			) VALUES %s`, strings.Join(values, ","))

		_, err = tx.Exec(ctx, itemsQuery, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to insert pos sale items: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return sale, nil
}
