package model

import "time"

type Product struct {
	ID                          int64      `json:"id"`
	AddedBy                     *string    `json:"added_by"`
	UserID                      *int64     `json:"user_id"`
	Name                        *string    `json:"name"`
	Slug                        *string    `json:"slug"`
	ProductType                 string     `json:"product_type"`
	CategoryIDs                 *string    `json:"category_ids"`
	CategoryID                  *string    `json:"category_id"`
	SubCategoryID               *string    `json:"sub_category_id"`
	SubSubCategoryID            *string    `json:"sub_sub_category_id"`
	BrandID                     *int64     `json:"brand_id"`
	Unit                        *string    `json:"unit"`
	MinQty                      int        `json:"min_qty"`
	Refundable                  bool       `json:"refundable"`
	DigitalProductType          *string    `json:"digital_product_type"`
	DigitalFileReady            *string    `json:"digital_file_ready"`
	DigitalFileReadyStorageType *string    `json:"digital_file_ready_storage_type"`
	Images                      *string    `json:"images"`
	ColorImage                  string     `json:"color_image"`
	Thumbnail                   *string    `json:"thumbnail"`
	ThumbnailStorageType        *string    `json:"thumbnail_storage_type"`
	PreviewFile                 *string    `json:"preview_file"`
	PreviewFileStorageType      *string    `json:"preview_file_storage_type"`
	Featured                    *string    `json:"featured"`
	FlashDeal                   *string    `json:"flash_deal"`
	VideoProvider               *string    `json:"video_provider"`
	VideoURL                    *string    `json:"video_url"`
	Colors                      *string    `json:"colors"`
	VariantProduct              bool       `json:"variant_product"`
	Attributes                  *string    `json:"attributes"`
	ChoiceOptions               *string    `json:"choice_options"`
	Variation                   *string    `json:"variation"`
	DigitalProductFileTypes     *string    `json:"digital_product_file_types"`
	DigitalProductExtensions    *string    `json:"digital_product_extensions"`
	Published                   bool       `json:"published"`
	UnitPrice                   float64    `json:"unit_price"`
	PurchasePrice               float64    `json:"purchase_price"`
	Tax                         string     `json:"tax"`
	TaxType                     *string    `json:"tax_type"`
	TaxModel                    string     `json:"tax_model"`
	Discount                    string     `json:"discount"`
	DiscountType                *string    `json:"discount_type"`
	CurrentStock                *int       `json:"current_stock"`
	MinimumOrderQty             int        `json:"minimum_order_qty"`
	Details                     *string    `json:"details"`
	FreeShipping                bool       `json:"free_shipping"`
	Attachment                  *string    `json:"attachment"`
	CreatedAt                   *time.Time `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at"`
	Status                      bool       `json:"status"`
	FeaturedStatus              bool       `json:"featured_status"`
	MetaTitle                   *string    `json:"meta_title"`
	MetaDescription             *string    `json:"meta_description"`
	MetaImage                   *string    `json:"meta_image"`
	RequestStatus               bool       `json:"request_status"`
	DeniedNote                  *string    `json:"denied_note"`
	ShippingCost                *float64   `json:"shipping_cost"`
	MultiplyQty                 *bool      `json:"multiply_qty"`
	TempShippingCost            *float64   `json:"temp_shipping_cost"`
	IsShippingCostUpdated       *bool      `json:"is_shipping_cost_updated"`
	Code                        *string    `json:"code"`
}

// ProductFilter holds optional query parameters for listing products
type ProductFilter struct {
	Status     string // "active", "inactive", or empty for all
	Published  string // "true", "false", or empty for all
	CategoryID string // filter by category_id
	Page       int    // pagination
	Limit      int    // pagination
}
