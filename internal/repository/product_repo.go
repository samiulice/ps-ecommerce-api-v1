package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

// Create inserts a new product and its variations transactionally
func (r *ProductRepo) Create(ctx context.Context, p *model.Product) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert into products table
	query := `
		INSERT INTO products (
			name, description, category_id, sub_category_id, sub_sub_category_id,
			brand_id, sku, status, unit_id, tags, thumbnail, gallery_images,
			unit_price, purchase_price, min_order_qty, current_stock_qty,
			stock_alert_qty, discount_type, discount_amount, tax_amount, tax_type,
			shipping_cost, shipping_type, has_variation, variation_attributes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25
		) RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		p.Name, p.Description, p.CategoryID, p.SubCategoryID, p.SubSubCategoryID,
		p.BrandID, p.SKU, p.Status, p.UnitID, p.Tags, p.Thumbnail, p.GalleryImages,
		p.UnitPrice, p.PurchasePrice, p.MinOrderQty, p.CurrentStockQty,
		p.StockAlertQty, p.DiscountType, p.DiscountAmount, p.TaxAmount, p.TaxType,
		p.ShippingCost, p.ShippingType, p.HasVariation, p.VariationAttributes,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return fmt.Errorf("product with SKU '%s' already exists", p.SKU)
		}
		return err
	}

	// 2. Insert Variations if they exist and has_variation is true
	if p.HasVariation && len(p.Variations) > 0 {
		varQuery := `
			INSERT INTO product_variations (
				product_id, variation_attributes, sku, price, stock_qty, thumbnail
			) VALUES ($1, $2, $3, $4, $5, $6)
		`
		for _, v := range p.Variations {
			_, err := tx.Exec(ctx, varQuery,
				p.ID, v.VariationAttributes, v.SKU, v.Price, v.StockQty, v.Thumbnail,
			)
			if err != nil {
				return fmt.Errorf("failed to save variation sku %s: %w", v.SKU, err)
			}
		}
	}

	return tx.Commit(ctx)
}

// Update modifies an existing product and recreates variations
func (r *ProductRepo) Update(ctx context.Context, p *model.Product) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Update Product Table
	query := `
		UPDATE products SET
			name = $1, description = $2, category_id = $3, sub_category_id = $4,
			sub_sub_category_id = $5, brand_id = $6, sku = $7, status = $8,
			unit_id = $9, tags = $10, thumbnail = $11, gallery_images = $12,
			unit_price = $13, purchase_price = $14, min_order_qty = $15,
			current_stock_qty = $16, stock_alert_qty = $17, discount_type = $18,
			discount_amount = $19, tax_amount = $20, tax_type = $21,
			shipping_cost = $22, shipping_type = $23, has_variation = $24,
			variation_attributes = $25, updated_at = CURRENT_TIMESTAMP
		WHERE id = $26
		RETURNING created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		p.Name, p.Description, p.CategoryID, p.SubCategoryID, p.SubSubCategoryID,
		p.BrandID, p.SKU, p.Status, p.UnitID, p.Tags, p.Thumbnail, p.GalleryImages,
		p.UnitPrice, p.PurchasePrice, p.MinOrderQty, p.CurrentStockQty,
		p.StockAlertQty, p.DiscountType, p.DiscountAmount, p.TaxAmount, p.TaxType,
		p.ShippingCost, p.ShippingType, p.HasVariation, p.VariationAttributes,
		p.ID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return err
	}

	// 2. Handle Variations (Delete all and Re-insert strategy for simplicity)
	// In a high-traffic production system, you might want to upsert instead.
	_, err = tx.Exec(ctx, "DELETE FROM product_variations WHERE product_id = $1", p.ID)
	if err != nil {
		return err
	}

	if p.HasVariation && len(p.Variations) > 0 {
		varQuery := `
			INSERT INTO product_variations (
				product_id, variation_attributes, sku, price, stock_qty, thumbnail
			) VALUES ($1, $2, $3, $4, $5, $6)
		`
		for _, v := range p.Variations {
			_, err := tx.Exec(ctx, varQuery,
				p.ID, v.VariationAttributes, v.SKU, v.Price, v.StockQty, v.Thumbnail,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

// Delete removes a product by ID
func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	// Constraints are set to CASCADE in DB, so this deletes variations too
	tag, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if err == nil && tag.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return err
}

// GetByID retrieves a single product and its variations
func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `
		SELECT id, name, description, category_id, sub_category_id, sub_sub_category_id,
			brand_id, sku, status, unit_id, tags, thumbnail, gallery_images,
			unit_price, purchase_price, min_order_qty, current_stock_qty, stock_alert_qty,
			total_sold, discount_type, discount_amount, tax_amount, tax_type,
			shipping_cost, shipping_type, has_variation, variation_attributes,
			total_reviews, avg_rating, five_star_count, four_star_count, three_star_count,
			two_star_count, one_star_count, created_at, updated_at
		FROM products WHERE id = $1
	`
	var p model.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.SubCategoryID, &p.SubSubCategoryID,
		&p.BrandID, &p.SKU, &p.Status, &p.UnitID, &p.Tags, &p.Thumbnail, &p.GalleryImages,
		&p.UnitPrice, &p.PurchasePrice, &p.MinOrderQty, &p.CurrentStockQty, &p.StockAlertQty,
		&p.TotalSold, &p.DiscountType, &p.DiscountAmount, &p.TaxAmount, &p.TaxType,
		&p.ShippingCost, &p.ShippingType, &p.HasVariation, &p.VariationAttributes,
		&p.TotalReviews, &p.AvgRating, &p.FiveStarCount, &p.FourStarCount, &p.ThreeStarCount,
		&p.TwoStarCount, &p.OneStarCount, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.New("product not found")
	} else if err != nil {
		return nil, err
	}

	// Fetch variations if they exist
	if p.HasVariation {
		varQuery := `
			SELECT id, product_id, variation_attributes, sku, price, stock_qty, thumbnail
			FROM product_variations WHERE product_id = $1
		`
		rows, err := r.db.Query(ctx, varQuery, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var v model.ProductVariation
			if err := rows.Scan(&v.ID, &v.ProductID, &v.VariationAttributes, &v.SKU, &v.Price, &v.StockQty, &v.Thumbnail); err != nil {
				return nil, err
			}
			p.Variations = append(p.Variations, v)
		}
	}

	return &p, nil
}

// GetProducts retrieves products with filters
func (r *ProductRepo) GetProducts(ctx context.Context, filter model.ProductFilter) ([]*model.Product, int64, error) {
	baseQuery := `
		SELECT id, name, description, category_id, sub_category_id, sub_sub_category_id,
			brand_id, sku, status, unit_id, tags, thumbnail, gallery_images,
			unit_price, purchase_price, min_order_qty, current_stock_qty, stock_alert_qty,
			total_sold, discount_type, discount_amount, tax_amount, tax_type,
			shipping_cost, shipping_type, has_variation, variation_attributes,
			total_reviews, avg_rating, five_star_count, four_star_count, three_star_count,
			two_star_count, one_star_count, created_at, updated_at
		FROM products
	`
	countQuery := `SELECT COUNT(*) FROM products`

	var conditions []string
	var args []any
	argPos := 1

	// Status filter (1 = active)
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		statusInt := 1
		if filter.Status == "inactive" {
			statusInt = 0
		}
		args = append(args, statusInt)
		argPos++
	}

	// Category filter
	if filter.CategoryID != "" {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argPos))
		args = append(args, filter.CategoryID)
		argPos++
	}

	// Search (Tags or Name)
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR tags ILIKE $%d)", argPos, argPos))
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int64
	err := r.db.QueryRow(ctx, countQuery+whereClause, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	baseQuery += whereClause + " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filter.Limit)
		argPos++
		if filter.Page > 0 {
			offset := (filter.Page - 1) * filter.Limit
			baseQuery += fmt.Sprintf(" OFFSET $%d", argPos)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var p model.Product
		// Scanning a subset of fields for list view optimization
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.SubCategoryID, &p.SubSubCategoryID,
			&p.BrandID, &p.SKU, &p.Status, &p.UnitID, &p.Tags, &p.Thumbnail, &p.GalleryImages,
			&p.UnitPrice, &p.PurchasePrice, &p.MinOrderQty, &p.CurrentStockQty, &p.StockAlertQty,
			&p.TotalSold, &p.DiscountType, &p.DiscountAmount, &p.TaxAmount, &p.TaxType,
			&p.ShippingCost, &p.ShippingType, &p.HasVariation, &p.VariationAttributes,
			&p.TotalReviews, &p.AvgRating, &p.FiveStarCount, &p.FourStarCount, &p.ThreeStarCount,
			&p.TwoStarCount, &p.OneStarCount, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, &p)
	}

	return products, totalCount, nil
}
