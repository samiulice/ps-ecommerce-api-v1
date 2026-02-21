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

// Create inserts a new product
func (r *ProductRepo) Create(ctx context.Context, p *model.Product) error {
	query := `
		INSERT INTO products (
			added_by, user_id, name, slug, product_type, category_ids, category_id,
			sub_category_id, sub_sub_category_id, brand_id, unit, min_qty, refundable,
			digital_product_type, digital_file_ready, digital_file_ready_storage_type,
			images, color_image, thumbnail, thumbnail_storage_type, preview_file,
			preview_file_storage_type, featured, flash_deal, video_provider, video_url,
			colors, variant_product, attributes, choice_options, variation,
			digital_product_file_types, digital_product_extensions, published,
			unit_price, purchase_price, tax, tax_type, tax_model, discount, discount_type,
			current_stock, minimum_order_qty, details, free_shipping, attachment,
			status, featured_status, meta_title, meta_description, meta_image,
			request_status, denied_note, shipping_cost, multiply_qty, temp_shipping_cost,
			is_shipping_cost_updated, code
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44,
			$45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58
		) RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		p.AddedBy, p.UserID, p.Name, p.Slug, p.ProductType, p.CategoryIDs, p.CategoryID,
		p.SubCategoryID, p.SubSubCategoryID, p.BrandID, p.Unit, p.MinQty, p.Refundable,
		p.DigitalProductType, p.DigitalFileReady, p.DigitalFileReadyStorageType,
		p.Images, p.ColorImage, p.Thumbnail, p.ThumbnailStorageType, p.PreviewFile,
		p.PreviewFileStorageType, p.Featured, p.FlashDeal, p.VideoProvider, p.VideoURL,
		p.Colors, p.VariantProduct, p.Attributes, p.ChoiceOptions, p.Variation,
		p.DigitalProductFileTypes, p.DigitalProductExtensions, p.Published,
		p.UnitPrice, p.PurchasePrice, p.Tax, p.TaxType, p.TaxModel, p.Discount, p.DiscountType,
		p.CurrentStock, p.MinimumOrderQty, p.Details, p.FreeShipping, p.Attachment,
		p.Status, p.FeaturedStatus, p.MetaTitle, p.MetaDescription, p.MetaImage,
		p.RequestStatus, p.DeniedNote, p.ShippingCost, p.MultiplyQty, p.TempShippingCost,
		p.IsShippingCostUpdated, p.Code,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if isUniqueViolation(err) {
		return fmt.Errorf("product with slug '%s' already exists", *p.Slug)
	}
	return err
}

// Update modifies an existing product
func (r *ProductRepo) Update(ctx context.Context, p *model.Product) error {
	query := `
		UPDATE products SET
			added_by = $1, user_id = $2, name = $3, slug = $4, product_type = $5,
			category_ids = $6, category_id = $7, sub_category_id = $8, sub_sub_category_id = $9,
			brand_id = $10, unit = $11, min_qty = $12, refundable = $13, digital_product_type = $14,
			digital_file_ready = $15, digital_file_ready_storage_type = $16, images = $17,
			color_image = $18, thumbnail = $19, thumbnail_storage_type = $20, preview_file = $21,
			preview_file_storage_type = $22, featured = $23, flash_deal = $24, video_provider = $25,
			video_url = $26, colors = $27, variant_product = $28, attributes = $29,
			choice_options = $30, variation = $31, digital_product_file_types = $32,
			digital_product_extensions = $33, published = $34, unit_price = $35, purchase_price = $36,
			tax = $37, tax_type = $38, tax_model = $39, discount = $40, discount_type = $41,
			current_stock = $42, minimum_order_qty = $43, details = $44, free_shipping = $45,
			attachment = $46, status = $47, featured_status = $48, meta_title = $49,
			meta_description = $50, meta_image = $51, request_status = $52, denied_note = $53,
			shipping_cost = $54, multiply_qty = $55, temp_shipping_cost = $56,
			is_shipping_cost_updated = $57, code = $58, updated_at = CURRENT_TIMESTAMP
		WHERE id = $59
		RETURNING created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		p.AddedBy, p.UserID, p.Name, p.Slug, p.ProductType,
		p.CategoryIDs, p.CategoryID, p.SubCategoryID, p.SubSubCategoryID,
		p.BrandID, p.Unit, p.MinQty, p.Refundable, p.DigitalProductType,
		p.DigitalFileReady, p.DigitalFileReadyStorageType, p.Images,
		p.ColorImage, p.Thumbnail, p.ThumbnailStorageType, p.PreviewFile,
		p.PreviewFileStorageType, p.Featured, p.FlashDeal, p.VideoProvider,
		p.VideoURL, p.Colors, p.VariantProduct, p.Attributes,
		p.ChoiceOptions, p.Variation, p.DigitalProductFileTypes,
		p.DigitalProductExtensions, p.Published, p.UnitPrice, p.PurchasePrice,
		p.Tax, p.TaxType, p.TaxModel, p.Discount, p.DiscountType,
		p.CurrentStock, p.MinimumOrderQty, p.Details, p.FreeShipping,
		p.Attachment, p.Status, p.FeaturedStatus, p.MetaTitle,
		p.MetaDescription, p.MetaImage, p.RequestStatus, p.DeniedNote,
		p.ShippingCost, p.MultiplyQty, p.TempShippingCost,
		p.IsShippingCostUpdated, p.Code, p.ID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)

	if isUniqueViolation(err) {
		return fmt.Errorf("product with slug '%s' already exists", *p.Slug)
	}
	return err
}

// Delete removes a product by ID
func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if err == nil && tag.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return err
}

// GetByID retrieves a single product by ID
func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `
		SELECT id, added_by, user_id, name, slug, product_type, category_ids, category_id,
			sub_category_id, sub_sub_category_id, brand_id, unit, min_qty, refundable,
			digital_product_type, digital_file_ready, digital_file_ready_storage_type,
			images, color_image, thumbnail, thumbnail_storage_type, preview_file,
			preview_file_storage_type, featured, flash_deal, video_provider, video_url,
			colors, variant_product, attributes, choice_options, variation,
			digital_product_file_types, digital_product_extensions, published,
			unit_price, purchase_price, tax, tax_type, tax_model, discount, discount_type,
			current_stock, minimum_order_qty, details, free_shipping, attachment,
			created_at, updated_at, status, featured_status, meta_title, meta_description,
			meta_image, request_status, denied_note, shipping_cost, multiply_qty,
			temp_shipping_cost, is_shipping_cost_updated, code
		FROM products WHERE id = $1
	`
	var p model.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.AddedBy, &p.UserID, &p.Name, &p.Slug, &p.ProductType, &p.CategoryIDs, &p.CategoryID,
		&p.SubCategoryID, &p.SubSubCategoryID, &p.BrandID, &p.Unit, &p.MinQty, &p.Refundable,
		&p.DigitalProductType, &p.DigitalFileReady, &p.DigitalFileReadyStorageType,
		&p.Images, &p.ColorImage, &p.Thumbnail, &p.ThumbnailStorageType, &p.PreviewFile,
		&p.PreviewFileStorageType, &p.Featured, &p.FlashDeal, &p.VideoProvider, &p.VideoURL,
		&p.Colors, &p.VariantProduct, &p.Attributes, &p.ChoiceOptions, &p.Variation,
		&p.DigitalProductFileTypes, &p.DigitalProductExtensions, &p.Published,
		&p.UnitPrice, &p.PurchasePrice, &p.Tax, &p.TaxType, &p.TaxModel, &p.Discount, &p.DiscountType,
		&p.CurrentStock, &p.MinimumOrderQty, &p.Details, &p.FreeShipping, &p.Attachment,
		&p.CreatedAt, &p.UpdatedAt, &p.Status, &p.FeaturedStatus, &p.MetaTitle, &p.MetaDescription,
		&p.MetaImage, &p.RequestStatus, &p.DeniedNote, &p.ShippingCost, &p.MultiplyQty,
		&p.TempShippingCost, &p.IsShippingCostUpdated, &p.Code,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.New("product not found")
	}
	return &p, err
}

// GetProducts retrieves products with optional filters and pagination
func (r *ProductRepo) GetProducts(ctx context.Context, filter model.ProductFilter) ([]*model.Product, int64, error) {
	baseQuery := `
		SELECT id, added_by, user_id, name, slug, product_type, category_ids, category_id,
			sub_category_id, sub_sub_category_id, brand_id, unit, min_qty, refundable,
			digital_product_type, digital_file_ready, digital_file_ready_storage_type,
			images, color_image, thumbnail, thumbnail_storage_type, preview_file,
			preview_file_storage_type, featured, flash_deal, video_provider, video_url,
			colors, variant_product, attributes, choice_options, variation,
			digital_product_file_types, digital_product_extensions, published,
			unit_price, purchase_price, tax, tax_type, tax_model, discount, discount_type,
			current_stock, minimum_order_qty, details, free_shipping, attachment,
			created_at, updated_at, status, featured_status, meta_title, meta_description,
			meta_image, request_status, denied_note, shipping_cost, multiply_qty,
			temp_shipping_cost, is_shipping_cost_updated, code
		FROM products
	`
	countQuery := `SELECT COUNT(*) FROM products`

	var conditions []string
	var args []any
	argPos := 1

	// status filter
	if filter.Status == "active" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, true)
		argPos++
	} else if filter.Status == "inactive" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, false)
		argPos++
	}

	// published filter
	if filter.Published == "true" {
		conditions = append(conditions, fmt.Sprintf("published = $%d", argPos))
		args = append(args, true)
		argPos++
	} else if filter.Published == "false" {
		conditions = append(conditions, fmt.Sprintf("published = $%d", argPos))
		args = append(args, false)
		argPos++
	}

	// category_id filter
	if filter.CategoryID != "" {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argPos))
		args = append(args, filter.CategoryID)
		argPos++
	}

	// build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Get total count
	var totalCount int64
	err := r.db.QueryRow(ctx, countQuery+whereClause, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Add ordering
	baseQuery += whereClause + " ORDER BY created_at DESC"

	// Pagination
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
		err := rows.Scan(
			&p.ID, &p.AddedBy, &p.UserID, &p.Name, &p.Slug, &p.ProductType, &p.CategoryIDs, &p.CategoryID,
			&p.SubCategoryID, &p.SubSubCategoryID, &p.BrandID, &p.Unit, &p.MinQty, &p.Refundable,
			&p.DigitalProductType, &p.DigitalFileReady, &p.DigitalFileReadyStorageType,
			&p.Images, &p.ColorImage, &p.Thumbnail, &p.ThumbnailStorageType, &p.PreviewFile,
			&p.PreviewFileStorageType, &p.Featured, &p.FlashDeal, &p.VideoProvider, &p.VideoURL,
			&p.Colors, &p.VariantProduct, &p.Attributes, &p.ChoiceOptions, &p.Variation,
			&p.DigitalProductFileTypes, &p.DigitalProductExtensions, &p.Published,
			&p.UnitPrice, &p.PurchasePrice, &p.Tax, &p.TaxType, &p.TaxModel, &p.Discount, &p.DiscountType,
			&p.CurrentStock, &p.MinimumOrderQty, &p.Details, &p.FreeShipping, &p.Attachment,
			&p.CreatedAt, &p.UpdatedAt, &p.Status, &p.FeaturedStatus, &p.MetaTitle, &p.MetaDescription,
			&p.MetaImage, &p.RequestStatus, &p.DeniedNote, &p.ShippingCost, &p.MultiplyQty,
			&p.TempShippingCost, &p.IsShippingCostUpdated, &p.Code,
		)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, &p)
	}

	return products, totalCount, nil
}
