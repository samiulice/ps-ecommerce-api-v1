package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (r *ProductRepo) Create(ctx context.Context, p *model.Product) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert into products table
	// Note: We insert with the provided p.SKU (which might be empty)
	// Ensure unit_id is set to retail_unit_id for backward compatibility
	if p.RetailUnitID != nil {
		p.UnitID = p.RetailUnitID
	}

	query := `
        INSERT INTO products (
            name, description, category_id, sub_category_id, sub_sub_category_id, brand_id, sku, status, unit_id, retail_unit_id, wholesale_unit_id, tags, thumbnail, gallery_images, retail_price, wholesale_price, purchase_price, min_retail_order_qty, min_wholesale_order_qty, current_stock_qty, stock_alert_qty, discount_type, discount_amount, tax_amount, tax_type, shipping_cost, shipping_type, has_variation, variation_attributes
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
            $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29
        ) RETURNING id, created_at, updated_at
    `

	err = tx.QueryRow(ctx, query,
		p.Name, p.Description, p.CategoryID, p.SubCategoryID, p.SubSubCategoryID,
		p.BrandID, p.SKU, p.Status, p.UnitID, p.RetailUnitID, p.WholesaleUnitID, p.Tags, p.Thumbnail, p.GalleryImages,
		p.RetailPrice, p.WholesalePrice, p.PurchasePrice, p.MinRetailOrderQty, p.MinWholesaleOrderQty, p.CurrentStockQty,
		p.StockAlertQty, p.DiscountType, p.DiscountAmount, p.TaxAmount, p.TaxType,
		p.ShippingCost, p.ShippingType, p.HasVariation, p.VariationAttributes,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return fmt.Errorf("product with SKU '%s' already exists", p.SKU)
		}
		return err
	}

	// 2. Generate SKU if missing using the Database ID
	// Format: PRD-26-CATID-DBID (e.g., PRD-26-5-102)
	// This is much shorter and guaranteed unique.
	if p.SKU == "" {
		yearShort := time.Now().Format("06") // "26"
		p.SKU = fmt.Sprintf("PRD-%s-%d-%d", yearShort, p.CategoryID, p.ID)

		_, err = tx.Exec(ctx, "UPDATE products SET sku = $1 WHERE id = $2", p.SKU, p.ID)
		if err != nil {
			return fmt.Errorf("failed to update generated sku: %w", err)
		}
	}

	// 3. Insert Variations
	if p.HasVariation && len(p.Variations) > 0 {
		varQuery := `
            INSERT INTO product_variations (
                product_id, variation_attributes, sku, price, stock_qty, thumbnail
            ) VALUES ($1, $2, $3, $4, $5, $6)
        `
		for i, v := range p.Variations {
			// If variation SKU is missing, derive it from parent SKU
			// Result: PRD-26-5-102-V1
			if v.SKU == "" {
				v.SKU = fmt.Sprintf("%s-V%d", p.SKU, i+1)
			}

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

// UpdateImageURLs update the existing product thumbnail and gallery images
func (r *ProductRepo) UpdateImageURLs(ctx context.Context, thumbnail string, galleryImages []string, productId int64) error {

	// 1. Update Product images urls
	query := `
        UPDATE products SET
            thumbnail = $1, gallery_images = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3
        RETURNING created_at, updated_at
    `
	_, err := r.db.Exec(ctx, query, thumbnail, galleryImages, productId)

	return err
}

// Update modifies an existing product and recreates variations
func (r *ProductRepo) Update(ctx context.Context, p *model.Product) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Ensure SKU is not lost or is regenerated if cleared
	if p.SKU == "" {
		yearShort := time.Now().Format("06") // "26"
		p.SKU = fmt.Sprintf("PRD-%s-%d-%d", yearShort, p.CategoryID, p.ID)
	}

	// Ensure unit_id is set to retail_unit_id for backward compatibility
	if p.RetailUnitID != nil {
		p.UnitID = p.RetailUnitID
	}

	// 2. Update Product Table
	query := `
        UPDATE products SET
            name = $1, description = $2, category_id = $3, sub_category_id = $4, sub_sub_category_id = $5, brand_id = $6, sku = $7, status = $8, unit_id = $9, retail_unit_id = $10, wholesale_unit_id = $11, tags = $12, thumbnail = $13, gallery_images = $14, retail_price = $15, wholesale_price = $16, purchase_price = $17, min_retail_order_qty = $18, min_wholesale_order_qty = $19, current_stock_qty = $20, stock_alert_qty = $21, discount_type = $22, discount_amount = $23, tax_amount = $24, tax_type = $25, shipping_cost = $26, shipping_type = $27, has_variation = $28, variation_attributes = $29, updated_at = CURRENT_TIMESTAMP
        WHERE id = $30
        RETURNING created_at, updated_at
    `
	err = tx.QueryRow(ctx, query,
		p.Name, p.Description, p.CategoryID, p.SubCategoryID, p.SubSubCategoryID,
		p.BrandID, p.SKU, p.Status, p.UnitID, p.RetailUnitID, p.WholesaleUnitID, p.Tags, p.Thumbnail, p.GalleryImages,
		p.RetailPrice, p.WholesalePrice, p.PurchasePrice, p.MinRetailOrderQty, p.MinWholesaleOrderQty, p.CurrentStockQty,
		p.StockAlertQty, p.DiscountType, p.DiscountAmount, p.TaxAmount, p.TaxType,
		p.ShippingCost, p.ShippingType, p.HasVariation, p.VariationAttributes,
		p.ID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return fmt.Errorf("sku '%s' is already taken by another product", p.SKU)
		}
		return err
	}

	// 3. Handle Variations (Delete and Re-insert)
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
		for i, v := range p.Variations {
			// Re-apply naming convention if variation SKU is missing
			if v.SKU == "" {
				v.SKU = fmt.Sprintf("%s-V%d", p.SKU, i+1)
			}

			_, err := tx.Exec(ctx, varQuery,
				p.ID, v.VariationAttributes, v.SKU, v.Price, v.StockQty, v.Thumbnail,
			)
			if err != nil {
				return fmt.Errorf("failed to save variation %s: %w", v.SKU, err)
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
		SELECT p.id, p.name, p.description, p.category_id, p.sub_category_id, p.sub_sub_category_id, p.brand_id, p.sku, p.status, p.unit_id, p.retail_unit_id, p.wholesale_unit_id, p.tags, p.thumbnail, p.gallery_images, p.retail_price, p.wholesale_price, p.purchase_price, p.min_retail_order_qty, p.min_wholesale_order_qty, p.current_stock_qty, p.stock_alert_qty, p.total_sold, p.discount_type, p.discount_amount, p.tax_amount, p.tax_type, p.shipping_cost, p.shipping_type, p.has_variation, p.variation_attributes, p.total_reviews, p.avg_rating, p.five_star_count, p.four_star_count, p.three_star_count, p.two_star_count, p.one_star_count, p.created_at, p.updated_at,
			r_u.symbol AS retail_unit_symbol, w_u.symbol AS wholesale_unit_symbol
		FROM products p
		LEFT JOIN units r_u ON p.retail_unit_id = r_u.id
		LEFT JOIN units w_u ON p.wholesale_unit_id = w_u.id
		WHERE p.id = $1
	`
	var p model.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.SubCategoryID, &p.SubSubCategoryID,
		&p.BrandID, &p.SKU, &p.Status, &p.UnitID, &p.RetailUnitID, &p.WholesaleUnitID, &p.Tags, &p.Thumbnail, &p.GalleryImages,
		&p.RetailPrice, &p.WholesalePrice, &p.PurchasePrice, &p.MinRetailOrderQty, &p.MinWholesaleOrderQty, &p.CurrentStockQty, &p.StockAlertQty,
		&p.TotalSold, &p.DiscountType, &p.DiscountAmount, &p.TaxAmount, &p.TaxType,
		&p.ShippingCost, &p.ShippingType, &p.HasVariation, &p.VariationAttributes,
		&p.TotalReviews, &p.AvgRating, &p.FiveStarCount, &p.FourStarCount, &p.ThreeStarCount,
		&p.TwoStarCount, &p.OneStarCount, &p.CreatedAt, &p.UpdatedAt,
		&p.RetailUnitSymbol, &p.WholesaleUnitSymbol,
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

// GetProductVariationsByProductID retrieves a product variations
func (r *ProductRepo) GetProductVariationsByProductID(ctx context.Context, id int64) ([]*model.ProductVariation, error) {
	var variations []*model.ProductVariation

	query := `
			SELECT id, product_id, variation_attributes, sku, price, stock_qty, thumbnail
			FROM product_variations WHERE product_id = $1
		`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v model.ProductVariation
		if err := rows.Scan(&v.ID, &v.ProductID, &v.VariationAttributes, &v.SKU, &v.Price, &v.StockQty, &v.Thumbnail); err != nil {
			return nil, err
		}
		variations = append(variations, &v)
	}

	return variations, nil
}

// GetProducts retrieves products with filters
func (r *ProductRepo) GetProducts(ctx context.Context, filter model.ProductFilter) ([]*model.Product, int64, error) {
	baseQuery := `
		SELECT p.id, p.name, p.description, p.category_id, p.sub_category_id, p.sub_sub_category_id,
			p.brand_id, p.sku, p.status, p.unit_id, p.retail_unit_id, p.wholesale_unit_id, p.tags, p.thumbnail, p.gallery_images,
			p.retail_price, p.wholesale_price, p.purchase_price, p.min_retail_order_qty, p.min_wholesale_order_qty, p.current_stock_qty, p.stock_alert_qty, p.total_sold, p.discount_type, p.discount_amount, p.tax_amount, p.tax_type,
			p.shipping_cost, p.shipping_type, p.has_variation, p.variation_attributes,
			p.total_reviews, p.avg_rating, p.five_star_count, p.four_star_count, p.three_star_count,
			p.two_star_count, p.one_star_count, p.created_at, p.updated_at,
			r_u.symbol AS retail_unit_symbol, w_u.symbol AS wholesale_unit_symbol
		FROM products p
		LEFT JOIN units r_u ON p.retail_unit_id = r_u.id
		LEFT JOIN units w_u ON p.wholesale_unit_id = w_u.id
	`
	// 1. Add alias 'p' to countQuery
	countQuery := `SELECT COUNT(*) FROM products p`

	if filter.Compact {
		// 2. Add alias 'p' to compact mode baseQuery
		baseQuery = `
			SELECT p.id, p.name, p.thumbnail, p.retail_price, p.wholesale_price
			FROM products p
		`
	}

	var conditions []string
	var args []any
	argPos := 1

	// Status filter (1 = active)
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argPos))
		statusInt := 1
		if filter.Status == "inactive" {
			statusInt = 0
		}
		args = append(args, statusInt)
		argPos++
	}

	// Category filter (MULTIPLE)
	if len(filter.CategoryIDs) > 0 {
		placeholders := []string{}

		for _, id := range filter.CategoryIDs {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argPos))
			args = append(args, id)
			argPos++
		}

		conditions = append(conditions,
			fmt.Sprintf("p.category_id IN (%s)", strings.Join(placeholders, ",")),
		)
	}

	// Search (Tags or Name)
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(p.name ILIKE $%d OR p.tags ILIKE $%d)", argPos, argPos))
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int64
	if !filter.SkipCount {
		err := r.db.QueryRow(ctx, countQuery+whereClause, args...).Scan(&totalCount)
		if err != nil {
			return nil, 0, err
		}
	}

	if (filter.Sort == "ASC" || filter.Sort == "DESC") && (filter.PriceType == "retail" || filter.PriceType == "wholesale") {
		baseQuery += whereClause + fmt.Sprintf(" ORDER BY p.%s_price %s", filter.PriceType, filter.Sort)
	} else if filter.Sort == "LT" {
		baseQuery += whereClause + " ORDER BY p.created_at DESC"
	} else {
		baseQuery += whereClause
	}

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
		var err error
		if filter.Compact {
			err = rows.Scan(&p.ID, &p.Name, &p.Thumbnail, &p.RetailPrice, &p.WholesalePrice)
		} else {
			err = rows.Scan(
				&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.SubCategoryID, &p.SubSubCategoryID,
				&p.BrandID, &p.SKU, &p.Status, &p.UnitID, &p.RetailUnitID, &p.WholesaleUnitID, &p.Tags, &p.Thumbnail, &p.GalleryImages,
				&p.RetailPrice, &p.WholesalePrice, &p.PurchasePrice, &p.MinRetailOrderQty, &p.MinWholesaleOrderQty, &p.CurrentStockQty, &p.StockAlertQty,
				&p.TotalSold, &p.DiscountType, &p.DiscountAmount, &p.TaxAmount, &p.TaxType,
				&p.ShippingCost, &p.ShippingType, &p.HasVariation, &p.VariationAttributes,
				&p.TotalReviews, &p.AvgRating, &p.FiveStarCount, &p.FourStarCount, &p.ThreeStarCount,
				&p.TwoStarCount, &p.OneStarCount, &p.CreatedAt, &p.UpdatedAt,
				&p.RetailUnitSymbol, &p.WholesaleUnitSymbol,
			)
		}
		if err != nil {
			return nil, 0, err
		}
		products = append(products, &p)
	}

	if filter.SkipCount {
		totalCount = int64(len(products))
	}

	return products, totalCount, nil
}
