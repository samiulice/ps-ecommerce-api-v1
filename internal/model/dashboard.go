package model

// StatsCards for the high-level dashboard metrics
type StatsCards struct {
	TotalOrders      int     `json:"total_orders"`
	TotalPOSSales    int     `json:"total_pos_sales"`
	TotalCustomers   int     `json:"total_customers"`
	InventoryPurchaseCosts   float64 `json:"total_purchases"`
	OperationalExpenses    float64 `json:"total_expenses"`
	NetCashFlow        float64 `json:"net_profit"`
	TotalProducts    int     `json:"total_products"`
	LowStockProducts int     `json:"low_stock_products"`
}

// ChartPoint generic structure for time-series data
type ChartPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// SaleComparisonData for the Bar Chart (POS vs Online and Sales vs Purchases)
type SaleComparisonData struct {
	Label       string  `json:"label"`
	POSSales    float64 `json:"pos_sales"`
	OnlineSales float64 `json:"online_sales"`
	Purchases   float64 `json:"purchases"`
}

// ProductSummary for popular and low stock lists
type ProductSummary struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	CurrentStock float64 `json:"current_stock"`
	AlertStock   float64 `json:"alert_stock"`
	TotalSold    float64 `json:"total_sold,omitempty"`
}