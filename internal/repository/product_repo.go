package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

// Create inserts a product and its variations.
func (r *ProductRepo) Create(ctx context.Context, p *model.Product) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert Main Product
	query := `
		INSERT INTO products (
			name, description, category_id, sub_category_id, sub_sub_category_id, brand_id, 
			sku, unit, search_tags, thumbnail, additional_thumbnails, unit_price, 
			min_order_qty, current_stock_qty, discount_type, discount_amount, 
			tax_amount, tax_calculation, shipping_cost, shipping_cost_type, has_variation
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		p.Name, p.Description, p.CategoryID, p.SubCategoryID, p.SubSubCategoryID, p.BrandID,
		p.SKU, p.Unit, p.SearchTags, p.Thumbnail, p.AdditionalThumbnails, p.UnitPrice,
		p.MinOrderQty, p.CurrentStockQty, p.DiscountType, p.DiscountAmount,
		p.TaxAmount, p.TaxCalculation, p.ShippingCost, p.ShippingCostType, p.HasVariation,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	// 2. Insert Variations if they exist
	if p.HasVariation && len(p.Variations) > 0 {
		vQuery := `
			INSERT INTO product_variations (product_id, name, price, sku, stock, thumbnail) 
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		for _, v := range p.Variations {
			if _, err := tx.Exec(ctx, vQuery, p.ID, v.Name, v.Price, v.SKU, v.Stock, v.Thumbnail); err != nil {
				return fmt.Errorf("failed to insert variation: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}

// GetByID fetches product and joins variations.
func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `
		SELECT id, name, description, category_id, sku, unit_price, has_variation, thumbnail, current_stock_qty
		FROM products WHERE id = $1
	`
	var p model.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.SKU, 
		&p.UnitPrice, &p.HasVariation, &p.Thumbnail, &p.CurrentStockQty,
	)
	if err != nil {
		return nil, err
	}

	// Fetch variations if flag is true
	if p.HasVariation {
		vRows, err := r.db.Query(ctx, "SELECT id, name, price, sku, stock, thumbnail FROM product_variations WHERE product_id=$1", id)
		if err != nil {
			return nil, err
		}
		defer vRows.Close()

		for vRows.Next() {
			var v model.ProductVariation
			if err := vRows.Scan(&v.ID, &v.Name, &v.Price, &v.SKU, &v.Stock, &v.Thumbnail); err == nil {
				p.Variations = append(p.Variations, v)
			}
		}
	}

	return &p, nil
}

// Update modifies product details.
func (r *ProductRepo) Update(ctx context.Context, p *model.Product) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Update Main Product Table
	// We use COALESCE or dynamic query building to avoid wiping fields if they are empty strings (optional strategy)
	// But strictly speaking, for PUT, we overwrite.
	query := `
		UPDATE products 
		SET name=$1, description=$2, category_id=$3, unit_price=$4, 
			current_stock_qty=$5, has_variation=$6, thumbnail=$7, updated_at=NOW() 
		WHERE id=$8`
	
	_, err = tx.Exec(ctx, query, 
		p.Name, p.Description, p.CategoryID, p.UnitPrice, 
		p.CurrentStockQty, p.HasVariation, p.Thumbnail, p.ID)
	if err != nil {
		return err
	}

	// 2. Handle Variations (Full Replace Strategy)
	// This is the safest way to ensure the DB matches the Frontend exactly.
	// A. Delete old variations
	if _, err := tx.Exec(ctx, "DELETE FROM product_variations WHERE product_id = $1", p.ID); err != nil {
		return err
	}

	// B. Insert new variations (if any)
	if p.HasVariation && len(p.Variations) > 0 {
		vQuery := `
			INSERT INTO product_variations (product_id, name, price, sku, stock, thumbnail) 
			VALUES ($1, $2, $3, $4, $5, $6)`
		
		for _, v := range p.Variations {
			// Ensure we use the Product ID
			_, err := tx.Exec(ctx, vQuery, p.ID, v.Name, v.Price, v.SKU, v.Stock, v.Thumbnail)
			if err != nil {
				return fmt.Errorf("failed to update variation: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}

// Delete removes product (cascades to variations).
func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	return err
}