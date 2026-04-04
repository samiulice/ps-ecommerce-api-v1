package model

import "time"

type POSSale struct {
	ID            int64         `json:"id" db:"id"`
	ReferenceNo   string        `json:"reference_no" db:"reference_no"`
	CustomerID    *int64        `json:"customer_id,omitempty" db:"customer_id"`
	BranchID      *int64        `json:"branch_id,omitempty" db:"branch_id"`
	SaleType      string        `json:"sale_type" db:"sale_type"`
	Subtotal      float64       `json:"subtotal" db:"subtotal"`
	Discount      float64       `json:"discount" db:"discount"`
	Tax           float64       `json:"tax" db:"tax"`
	Total         float64       `json:"total" db:"total"`
	AmountPaid    float64       `json:"amount_paid" db:"amount_paid"`
	PaymentMethod string        `json:"payment_method" db:"payment_method"`
	PaymentStatus string        `json:"payment_status" db:"payment_status"`
	SaleDate      time.Time     `json:"sale_date" db:"sale_date"`
	SaleNote      string        `json:"sale_note" db:"sale_note"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
	Items         []POSSaleItem `json:"items,omitempty"`
}

type POSSaleItem struct {
	ID                 int64   `json:"id" db:"id"`
	POSSaleID          int64   `json:"pos_sale_id" db:"pos_sale_id"`
	ProductID          int64   `json:"product_id" db:"product_id"`
	ProductVariationID *int64  `json:"product_variation_id,omitempty" db:"product_variation_id"`
	ProductName        string  `json:"product_name" db:"product_name"`
	Quantity           int     `json:"quantity" db:"quantity"`
	UnitPrice          float64 `json:"unit_price" db:"unit_price"`
	Subtotal           float64 `json:"subtotal" db:"subtotal"`
	TaxAmount          float64 `json:"tax_amount" db:"tax_amount"`
	Total              float64 `json:"total" db:"total"`
}

type CreatePOSSaleRequest struct {
	CustomerID    *int64                 `json:"customer_id"`
	BranchID      *int64                 `json:"branch_id"`
	SaleType      string                 `json:"sale_type"` // retail or wholesale
	Discount      float64                `json:"discount"`
	AmountPaid    float64                `json:"amount_paid"`
	PaymentMethod string                 `json:"payment_method"` // cash, card, mobile_banking
	SaleNote      string                 `json:"sale_note"`
	Items         []CreatePOSSaleItemReq `json:"items"`
}

type CreatePOSSaleItemReq struct {
	ProductID          int     `json:"product_id"`
	ProductVariationID *int    `json:"product_variation_id,omitempty"`
	Quantity           int     `json:"quantity"`
	UnitPrice          float64 `json:"unit_price"`
}
