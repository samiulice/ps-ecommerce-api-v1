package model

import "time"

type Product struct {
	ID                   int64              `json:"id"`
	Name                 string             `json:"name"`
	Description          string             `json:"description"`
	CategoryID           int64              `json:"categoryId"`
	SubCategoryID        int64              `json:"subCategoryId"`
	SubSubCategoryID     int64              `json:"subSubCategoryId"`
	BrandID              int64              `json:"brandId"`
	SKU                  string             `json:"sku"`
	Unit                 string             `json:"unit"`
	SearchTags           string             `json:"searchTags"`
	Thumbnail            string             `json:"thumbnail"`
	AdditionalThumbnails string             `json:"additionalThumbnails"`
	UnitPrice            float64            `json:"unitPrice"`
	MinOrderQty          float64            `json:"minOrderQty"`
	CurrentStockQty      float64            `json:"currentStockQty"`
	StockAlertQty        float64            `json:"stockAlertQty"`
	DiscountType         string             `json:"discountType"`
	DiscountAmount       float64            `json:"discountAmount"`
	TaxAmount            float64            `json:"taxAmount"`
	TaxCalculation       string             `json:"taxCalculation"`
	ShippingCost         float64            `json:"shippingCost"`
	ShippingCostType     string             `json:"shippingCostType"`
	HasVariation         bool               `json:"hasVariation"`
	Variations           []ProductVariation `json:"variations,omitempty"` // populated if HasVariation is true
	CreatedAt            time.Time          `json:"createdAt"`
	UpdatedAt            time.Time          `json:"updatedAt"`
}

type ProductVariation struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"productid"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	SKU       string  `json:"sku"`
	Stock     int     `json:"stock"`
	Thumbnail string  `json:"thumbnail"`
}
