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

func (r *DashboardRepository) GetStatsCards(ctx context.Context, period string) (*model.StatsCards, error) {
	stats := &model.StatsCards{}
	interval := r.getInterval(period)

	query := fmt.Sprintf(`
		SELECT 
			(SELECT COUNT(*) FROM orders WHERE created_at >= NOW() - INTERVAL %s AND order_status != 'cancelled') as online_orders,
			(SELECT COUNT(*) FROM pos_sales WHERE created_at >= NOW() - INTERVAL %s) as pos_sales,
			(SELECT COUNT(*) FROM customers WHERE created_at >= NOW() - INTERVAL %s) as customers,
			(SELECT COUNT(*) FROM products) as total_products,
			(SELECT COUNT(*) FROM products WHERE current_stock_qty <= stock_alert_qty) as low_stock,
			COALESCE((SELECT SUM(total_purchases) FROM view_daily_financial_report WHERE transaction_date >= NOW() - INTERVAL %s), 0) as total_purchases,
			COALESCE((SELECT SUM(total_expenses) FROM view_daily_financial_report WHERE transaction_date >= NOW() - INTERVAL %s), 0) as total_expenses,
			COALESCE((SELECT SUM(net_profit) FROM view_daily_financial_report WHERE transaction_date >= NOW() - INTERVAL %s), 0) as net_profit
	`, interval, interval, interval, interval, interval, interval)

	err := r.db.QueryRow(ctx, query).Scan(
		&stats.TotalOrders, &stats.TotalPOSSales, &stats.TotalCustomers,
		&stats.TotalProducts, &stats.LowStockProducts,
		&stats.TotalPurchases, &stats.TotalExpenses, &stats.NetProfit,
	)
	return stats, err
}

func (r *DashboardRepository) GetSaleComparisonChart(ctx context.Context, period string) ([]model.SaleComparisonData, error) {
	interval := r.getInterval(period)
	dateTrunc := r.getDateTrunc(period)

	// Note: We differentiate POS vs Online by joining the ledger with transaction descriptions
	query := fmt.Sprintf(`
		SELECT 
			TO_CHAR(DATE_TRUNC('%s', transaction_date), 'YYYY-MM-DD') as label,
			SUM(CASE WHEN type = 'sale_income' AND description LIKE 'POS%%' THEN credit_amount ELSE 0 END) as pos_val,
			SUM(CASE WHEN type = 'sale_income' AND description LIKE 'Income from Order%%' THEN credit_amount ELSE 0 END) as online_val,
			SUM(total_purchases) as purchases
		FROM view_daily_financial_report
		LEFT JOIN account_transactions USING (transaction_date)
		WHERE transaction_date >= NOW() - INTERVAL %s
		GROUP BY 1 ORDER BY 1 ASC`, dateTrunc, interval)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.SaleComparisonData
	for rows.Next() {
		var d model.SaleComparisonData
		if err := rows.Scan(&d.Label, &d.POSSales, &d.OnlineSales, &d.Purchases); err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	return data, nil
}

func (r *DashboardRepository) GetFinancialGraph(ctx context.Context, period, column string) ([]model.ChartPoint, error) {
	interval := r.getInterval(period)
	dateTrunc := r.getDateTrunc(period)

	query := fmt.Sprintf(`
		SELECT 
			TO_CHAR(DATE_TRUNC('%s', transaction_date), 'YYYY-MM-DD') as label,
			SUM(%s) as value
		FROM view_daily_financial_report
		WHERE transaction_date >= NOW() - INTERVAL %s
		GROUP BY 1 ORDER BY 1 ASC`, dateTrunc, column, interval)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.ChartPoint
	for rows.Next() {
		var p model.ChartPoint
		if err := rows.Scan(&p.Label, &p.Value); err != nil {
			return nil, err
		}
		data = append(data, p)
	}
	return data, nil
}

func (r *DashboardRepository) GetPopularProducts(ctx context.Context) ([]model.ProductSummary, error) {
	query := `
		SELECT product_id, product_name, SUM(quantity) as total_qty
		FROM order_items
		GROUP BY 1, 2 ORDER BY 3 DESC LIMIT 10`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var products []model.ProductSummary
	for rows.Next() {
		var p model.ProductSummary
		rows.Scan(&p.ID, &p.Name, &p.TotalSold)
		products = append(products, p)
	}
	return products, nil
}

func (r *DashboardRepository) GetLowStockProducts(ctx context.Context) ([]model.ProductSummary, error) {
	query := `SELECT id, name, current_stock_qty, stock_alert_qty FROM products WHERE current_stock_qty <= stock_alert_qty`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var products []model.ProductSummary
	for rows.Next() {
		var p model.ProductSummary
		rows.Scan(&p.ID, &p.Name, &p.CurrentStock, &p.AlertStock)
		products = append(products, p)
	}
	return products, nil
}

// Helpers
func (r *DashboardRepository) getInterval(p string) string {
	switch p {
	case "weekly": return "'7 days'"
	case "yearly": return "'1 year'"
	default: return "'1 month'"
	}
}

func (r *DashboardRepository) getDateTrunc(p string) string {
	if p == "yearly" { return "month" }
	return "day"
}