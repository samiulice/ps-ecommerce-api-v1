package model

type ChartData struct {
    Label string  `json:"label"`
    Value float64 `json:"value"`
}
type DashboardStats struct {
	TotalRevenue    float64            `json:"total_revenue"`
	TotalExpenses   float64            `json:"total_expenses"`
	NetProfit       float64            `json:"net_profit"`
	TotalOrders     int                `json:"total_orders"`     // Online
	TotalPOSSales   int                `json:"total_pos_sales"` // In-store
	NewCustomers    int                `json:"new_customers"`
	LowStockAlerts  int                `json:"low_stock_alerts"`
	FinancialChart  []FinancialChartData `json:"financial_chart"`
	TopProducts     []ChartData        `json:"top_products"`
}

type FinancialChartData struct {
	Label    string  `json:"label"`
	Revenue  float64 `json:"revenue"`
	Expenses float64 `json:"expenses"`
	Profit   float64 `json:"profit"`
}