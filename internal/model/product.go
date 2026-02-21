package model

import "time"

// Product represents the product entity based on the new schema
type Product struct {
	ID                  int64                  `json:"id"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`      // Stored as Markdown
	DescriptionHTML     string                 `json:"description_html"` // Computed: Markdown converted to sanitized HTML (not stored in DB)
	CategoryID          int64                  `json:"category_id"`
	SubCategoryID       *int64                 `json:"sub_category_id"`
	SubSubCategoryID    *int64                 `json:"sub_sub_category_id"`
	BrandID             *int64                 `json:"brand_id"`
	SKU                 string                 `json:"sku"`
	Status              int                    `json:"status"` // 1: Active, 0: Inactive
	UnitID              *int                   `json:"unit_id"`
	Tags                string                 `json:"tags"`
	Thumbnail           string                 `json:"thumbnail"`
	GalleryImages       []string               `json:"gallery_images"`
	UnitPrice           float64                `json:"unit_price"`
	PurchasePrice       float64                `json:"purchase_price"`
	MinOrderQty         float64                `json:"min_order_qty"`
	CurrentStockQty     float64                `json:"current_stock_qty"`
	StockAlertQty       float64                `json:"stock_alert_qty"`
	TotalSold           float64                `json:"total_sold"`
	DiscountType        string                 `json:"discount_type"` // "percentage" or "flat"
	DiscountAmount      float64                `json:"discount_amount"`
	TaxAmount           float64                `json:"tax_amount"`
	TaxType             string                 `json:"tax_type"` // "inclusive" or "exclusive"
	ShippingCost        float64                `json:"shipping_cost"`
	ShippingType        string                 `json:"shipping_type"`
	HasVariation        bool                   `json:"has_variation"`
	VariationAttributes map[string]interface{} `json:"variation_attributes"` // JSONB column
	TotalReviews        int64                  `json:"total_reviews"`
	AvgRating           float64                `json:"avg_rating"`
	FiveStarCount       int64                  `json:"five_star_count"`
	FourStarCount       int64                  `json:"four_star_count"`
	ThreeStarCount      int64                  `json:"three_star_count"`
	TwoStarCount        int64                  `json:"two_star_count"`
	OneStarCount        int64                  `json:"one_star_count"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`

	// Relations
	Variations []ProductVariation `json:"variations,omitempty"`
}

// ProductVariation represents the product_variations table
type ProductVariation struct {
	ID                  int64                  `json:"id"`
	ProductID           int64                  `json:"product_id"`
	VariationAttributes map[string]interface{} `json:"variation_attributes"` // JSONB
	SKU                 string                 `json:"sku"`
	Price               float64                `json:"price"`
	StockQty            float64                `json:"stock_qty"`
	Thumbnail           string                 `json:"thumbnail"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// ProductFilter holds query parameters for listing products
type ProductFilter struct {
	Status     string // "active", "inactive"
	CategoryID string // Filter by category
	Search     string // Search by name or tags
	Page       int    // Pagination page
	Limit      int    // Pagination limit
}
