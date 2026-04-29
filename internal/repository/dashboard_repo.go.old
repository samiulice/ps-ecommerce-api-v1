package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetDashboardStats(ctx context.Context, period string) (*model.DashboardStats, error) {
	stats := &model.DashboardStats{}

	// Define time filters based on period
	// We use transaction_date for financial data and created_at for entities
	interval := "'1 month'"
	switch period {
	case "weekly":
		interval = "'7 days'"
	case "yearly":
		interval = "'1 year'"
	}

	// 1. High-Level Cards (Summary Stats)
	summaryQuery := fmt.Sprintf(`
		SELECT 
			COALESCE(SUM(total_sales_income), 0) as total_revenue,
			COALESCE(SUM(total_expenses + total_purchases), 0) as total_outflow,
			COALESCE(SUM(net_profit), 0) as net_profit
		FROM view_daily_financial_report 
		WHERE transaction_date >= NOW() - INTERVAL %s`, interval)

	err := r.db.QueryRow(ctx, summaryQuery).Scan(
		&stats.TotalRevenue, 
		&stats.TotalExpenses, 
		&stats.NetProfit,
	)
	if err != nil {
		return nil, fmt.Errorf("summary stats error: %w", err)
	}

	// 2. Operational Stats (Counts)
	countsQuery := fmt.Sprintf(`
		SELECT 
			(SELECT COUNT(*) FROM orders WHERE created_at >= NOW() - INTERVAL %s AND order_status != 'cancelled') as online_orders,
			(SELECT COUNT(*) FROM pos_sales WHERE created_at >= NOW() - INTERVAL %s) as pos_sales,
			(SELECT COUNT(*) FROM customers WHERE created_at >= NOW() - INTERVAL %s) as new_customers,
			(SELECT COUNT(*) FROM products WHERE current_stock_qty <= stock_alert_qty) as low_stock_products
		`, interval, interval, interval)

	err = r.db.QueryRow(ctx, countsQuery).Scan(
		&stats.TotalOrders,
		&stats.TotalPOSSales,
		&stats.NewCustomers,
		&stats.LowStockAlerts,
	)
	if err != nil {
		return nil, fmt.Errorf("counts error: %w", err)
	}

	// 3. Financial Chart Data (Revenue vs Expense vs Profit)
	// Using the view to get a clean time-series
	dateTrunc := "day"
	if period == "yearly" {
		dateTrunc = "month"
	}

	chartQuery := fmt.Sprintf(`
		SELECT 
			TO_CHAR(DATE_TRUNC('%s', transaction_date), 'YYYY-MM-DD') as label,
			SUM(total_sales_income) as revenue,
			SUM(total_expenses + total_purchases) as expenses,
			SUM(net_profit) as profit
		FROM view_daily_financial_report
		WHERE transaction_date >= NOW() - INTERVAL %s
		GROUP BY 1 ORDER BY 1 ASC`, dateTrunc, interval)

	rows, err := r.db.Query(ctx, chartQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.FinancialChartData
		if err := rows.Scan(&item.Label, &item.Revenue, &item.Expenses, &item.Profit); err != nil {
			return nil, err
		}
		stats.FinancialChart = append(stats.FinancialChart, item)
	}

	// 4. Top Selling Products (Insight for the frontend "Top Products" card)
	topProductsQuery := `
		SELECT product_name, SUM(quantity) as total_qty
		FROM order_items
		GROUP BY product_name
		ORDER BY total_qty DESC
		LIMIT 5`

	pRows, err := r.db.Query(ctx, topProductsQuery)
	if err == nil {
		defer pRows.Close()
		for pRows.Next() {
			var p model.ChartData
			if err := pRows.Scan(&p.Label, &p.Value); err == nil {
				stats.TopProducts = append(stats.TopProducts, p)
			}
		}
	}

	return stats, nil
}