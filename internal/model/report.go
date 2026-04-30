package model

import "time"

// POS Sale Report
type POSSaleReportItem struct {
	ID           int64     `json:"id" db:"id"`
	ReferenceNo  string    `json:"reference_no" db:"reference_no"`
	CustomerName *string   `json:"customer_name" db:"customer_name"`
	SaleType     string    `json:"sale_type" db:"sale_type"`
	Total        float64   `json:"total" db:"total"`
	AmountPaid   float64   `json:"amount_paid" db:"amount_paid"`
	SaleDate     time.Time `json:"sale_date" db:"sale_date"`
}

type POSSaleReportResponse struct {
	TotalHits   int                 `json:"total_hits"`
	TotalAmount float64             `json:"total_amount"`
	TotalPaid   float64             `json:"total_paid"`
	Data        []POSSaleReportItem `json:"data"`
}

// Order Report
type OrderReportItem struct {
	ID            int64     `json:"id" db:"id"`
	OrderNumber   string    `json:"order_number" db:"order_number"`
	CustomerName  string    `json:"customer_name" db:"customer_name"`
	SaleType      string    `json:"sale_type" db:"sale_type"`
	OrderStatus   string    `json:"order_status" db:"order_status"`
	PaymentStatus string    `json:"payment_status" db:"payment_status"`
	Total         float64   `json:"total" db:"total"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type OrderReportResponse struct {
	TotalHits   int               `json:"total_hits"`
	TotalAmount float64           `json:"total_amount"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
	Items        []OrderReportItem `json:"items"`
}

// Customer Due
type CustomerDueReportItem struct {
	CustomerID int64   `json:"customer_id" db:"customer_id"`
	Name       string  `json:"name" db:"name"`
	Phone      string  `json:"phone" db:"phone"`
	TotalDue   float64 `json:"total_due" db:"total_due"`
}

type CustomerDueReportResponse struct {
	TotalHits int                     `json:"total_hits"`
	TotalDue  float64                 `json:"total_due"`
	Data      []CustomerDueReportItem `json:"data"`
}

// Supplier Due
type SupplierDueReportItem struct {
	SupplierID   int64   `json:"supplier_id" db:"supplier_id"`
	SupplierCode string  `json:"supplier_code" db:"supplier_code"`
	Name         string  `json:"name" db:"name"`
	Phone        string  `json:"phone" db:"phone"`
	TotalDue     float64 `json:"total_due" db:"total_due"`
}

type SupplierDueReportResponse struct {
	TotalHits int                     `json:"total_hits"`
	TotalDue  float64                 `json:"total_due"`
	Data      []SupplierDueReportItem `json:"data"`
}

// Low Stock
type LowStockReportItem struct {
	ProductID         int64   `json:"product_id" db:"id"`
	Name              string  `json:"name" db:"name"`
	SKU               string  `json:"sku" db:"sku"`
	CurrentStockQty   float64 `json:"stock_quantity" db:"current_stock_qty"`
	AlertQty float64 `json:"alert_stock" db:"stock_alert_qty"`
	ImageURL          *string `json:"image_url" db:"thumbnail"`
}

type LowStockReportResponse struct {
	TotalItems int                  `json:"total_items"`
	TotalPages int                  `json:"total_pages"`
	Items      []LowStockReportItem `json:"items"`
}

// Financial Report

type FinancialReportItem struct {
	TransactionDate   time.Time `json:"transaction_date" db:"transaction_date"`
	TotalSalesIncome  float64   `json:"total_sales_income" db:"total_sales_income"`
	TotalSalesRefunds float64   `json:"total_sales_refunds" db:"total_sales_refunds"`
	TotalPurchases    float64   `json:"total_purchases" db:"total_purchases"`
	TotalExpenses     float64   `json:"total_expenses" db:"total_expenses"`
	NetProfit         float64   `json:"net_profit" db:"net_profit"`
}

type FinancialReportResponse struct {
	TotalHits int                   `json:"total_hits"`
	Data      []FinancialReportItem `json:"data"`
}

type ReportFilter struct {
	SaleType string // 'retail', 'wholesale', ''

	OrderBy string // 'price_asc', 'price_desc', 'date_desc'
	Search  string
	Page    int
	Limit   int
}
