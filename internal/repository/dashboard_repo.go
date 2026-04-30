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
			(SELECT COUNT(*) FROM orders WHERE created_at >= CURRENT_DATE - INTERVAL %s AND order_status != 'cancelled') as online_orders,
			(SELECT COUNT(*) FROM pos_sales WHERE created_at >= CURRENT_DATE - INTERVAL %s) as pos_sales,
			(SELECT COUNT(*) FROM customers WHERE created_at >= CURRENT_DATE - INTERVAL %s) as customers,
			(SELECT COUNT(*) FROM products) as total_products,
			(SELECT COUNT(*) FROM products WHERE current_stock_qty <= stock_alert_qty) as low_stock,
			COALESCE((SELECT SUM(inventory_purchase_costs) FROM view_daily_financial_report WHERE transaction_date >= CURRENT_DATE - INTERVAL %s), 0) as inventory_purchase_costs,
			COALESCE((SELECT SUM(operational_expenses) FROM view_daily_financial_report WHERE transaction_date >= CURRENT_DATE - INTERVAL %s), 0) as operational_expenses,
			COALESCE((SELECT SUM(net_cash_flow) FROM view_daily_financial_report WHERE transaction_date >= CURRENT_DATE - INTERVAL %s), 0) as net_cash_flow
	`, interval, interval, interval, interval, interval, interval)

	err := r.db.QueryRow(ctx, query).Scan(
		&stats.TotalOrders, &stats.TotalPOSSales, &stats.TotalCustomers,
		&stats.TotalProducts, &stats.LowStockProducts,
		&stats.InventoryPurchaseCosts, &stats.OperationalExpenses, &stats.NetCashFlow,
	)
	return stats, err
}

func (r *DashboardRepository) GetSaleComparisonChart(ctx context.Context, period string) ([]model.SaleComparisonData, error) {
	interval := r.getInterval(period)
	dateTrunc := r.getDateTrunc(period)

	query := fmt.Sprintf(`
		WITH date_series AS (
			SELECT generate_series(
				DATE_TRUNC('%[1]s', CURRENT_DATE - INTERVAL %[2]s),
				DATE_TRUNC('%[1]s', CURRENT_DATE),
				'1 %[1]s'::interval
			)::date AS series_date
		),
		agg_transactions AS (
			SELECT 
				DATE_TRUNC('%[1]s', transaction_date)::date as trans_date,
				SUM(CASE WHEN type = 'sale_income' AND description LIKE 'POS%%' THEN credit_amount ELSE 0 END) as pos_val,
				SUM(CASE WHEN type = 'sale_income' AND description LIKE 'Income from Order%%' THEN credit_amount ELSE 0 END) as online_val
			FROM account_transactions
			WHERE transaction_date >= CURRENT_DATE - INTERVAL %[2]s
			GROUP BY 1
		),
		agg_report AS (
			SELECT 
				DATE_TRUNC('%[1]s', transaction_date)::date as rep_date,
				SUM(inventory_purchase_costs) as purchases
			FROM view_daily_financial_report
			WHERE transaction_date >= CURRENT_DATE - INTERVAL %[2]s
			GROUP BY 1
		)
		SELECT 
			TO_CHAR(ds.series_date, 'YYYY-MM-DD') as label,
			COALESCE(at.pos_val, 0) as pos_val,
			COALESCE(at.online_val, 0) as online_val,
			COALESCE(ar.purchases, 0) as purchases
		FROM date_series ds
		LEFT JOIN agg_transactions at ON ds.series_date = at.trans_date
		LEFT JOIN agg_report ar ON ds.series_date = ar.rep_date
		ORDER BY ds.series_date ASC`, dateTrunc, interval)

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
		WITH date_series AS (
			SELECT generate_series(
				DATE_TRUNC('%[1]s', CURRENT_DATE - INTERVAL %[2]s),
				DATE_TRUNC('%[1]s', CURRENT_DATE),
				'1 %[1]s'::interval
			)::date AS series_date
		),
		agg_report AS (
			SELECT 
				DATE_TRUNC('%[1]s', transaction_date)::date as rep_date,
				SUM(%[3]s) as value
			FROM view_daily_financial_report
			WHERE transaction_date >= CURRENT_DATE - INTERVAL %[2]s
			GROUP BY 1
		)
		SELECT 
			TO_CHAR(ds.series_date, 'YYYY-MM-DD') as label,
			COALESCE(ar.value, 0) as value
		FROM date_series ds
		LEFT JOIN agg_report ar ON ds.series_date = ar.rep_date
		ORDER BY ds.series_date ASC`, dateTrunc, interval, column)

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
                SELECT product_id, product_name, SUM(quantity)::float8 as total_qty
                FROM order_items
                GROUP BY 1, 2 ORDER BY 3 DESC LIMIT 10`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]model.ProductSummary, 0)
	for rows.Next() {
		var p model.ProductSummary
		// Use a pointer/sql.Null to ignore issues if a product got deleted or named null
		var pID *int64
		var pName *string
		var pSold *float64
		if err := rows.Scan(&pID, &pName, &pSold); err != nil {
			fmt.Printf("GetPopularProducts scan err: %v\n", err)
			continue
		}
		if pID != nil {
			p.ID = *pID
		}
		if pName != nil {
			p.Name = *pName
		}
		if pSold != nil {
			p.TotalSold = *pSold
		}

		products = append(products, p)
	}
	return products, nil
}

func (r *DashboardRepository) GetLowStockProducts(ctx context.Context) ([]model.ProductSummary, error) {
	query := `SELECT id, name, sku, current_stock_qty::float8, stock_alert_qty::float8 FROM products WHERE current_stock_qty <= stock_alert_qty`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]model.ProductSummary, 0)
	for rows.Next() {
		var p model.ProductSummary
		var sku *string
		if err := rows.Scan(&p.ID, &p.Name, &sku, &p.CurrentStock, &p.AlertStock); err != nil {
			fmt.Printf("GetLowStockProducts scan err: %v\n", err)
			continue
		}
		if sku != nil {
			p.SKU = *sku
		}
		products = append(products, p)
	}
	return products, nil
}

// Helpers
func (r *DashboardRepository) getInterval(p string) string {
    switch p {
    case "weekly":
        return "'7 days'"
    case "monthly":
        return "'1 month'"
    case "yearly":
        return "'1 year'"
    default:
        // Use 0 days because we will anchor to CURRENT_DATE for 'today'
        return "'0 days'"
    }
}

func (r *DashboardRepository) getDateTrunc(p string) string {
	if p == "yearly" {
		return "month"
	}
	return "day"
}
