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
func (r *POSRepo) GetSaleByReference(ctx context.Context, referenceNo string) (*model.POSSale, error) {
	sale := &model.POSSale{}
	query := `
                SELECT id, reference_no, customer_id, branch_id, sale_type, subtotal, discount,
                       tax, total, amount_paid, payment_method, payment_status, sale_date, sale_note, created_at, updated_at
                FROM pos_sales
                WHERE reference_no = $1
        `
	err := r.db.QueryRow(ctx, query, referenceNo).Scan(
		&sale.ID, &sale.ReferenceNo, &sale.CustomerID, &sale.BranchID, &sale.SaleType, &sale.Subtotal, &sale.Discount,
		&sale.Tax, &sale.Total, &sale.AmountPaid, &sale.PaymentMethod, &sale.PaymentStatus, &sale.SaleDate, &sale.SaleNote, &sale.CreatedAt, &sale.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("pos sale not found")
		}
		return nil, fmt.Errorf("failed to get pos sale: %w", err)
	}

	itemsQuery := `
                SELECT id, pos_sale_id, product_id, product_variation_id, product_name, quantity, unit_price, subtotal, tax_amount, total
                FROM pos_sale_items
                WHERE pos_sale_id = $1
        `
	rows, err := r.db.Query(ctx, itemsQuery, sale.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pos sale items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.POSSaleItem
		err := rows.Scan(
			&item.ID, &item.POSSaleID, &item.ProductID, &item.ProductVariationID, &item.ProductName, &item.Quantity, &item.UnitPrice, &item.Subtotal, &item.TaxAmount, &item.Total,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pos sale item: %w", err)
		}
		sale.Items = append(sale.Items, item)
	}

	return sale, nil
}
