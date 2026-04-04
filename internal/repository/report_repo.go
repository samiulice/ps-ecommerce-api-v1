package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type ReportRepo struct {
	db *pgxpool.Pool
}

func NewReportRepo(db *pgxpool.Pool) *ReportRepo {
	return &ReportRepo{db: db}
}

func (r *ReportRepo) GetPOSSalesReport(ctx context.Context, filter model.ReportFilter) (*model.POSSaleReportResponse, error) {
	var resp model.POSSaleReportResponse

	offset := (filter.Page - 1) * filter.Limit

	whereClauses := []string{"1=1"}
	var args []interface{}
	argID := 1

	if filter.SaleType != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.sale_type = $%d", argID))
		args = append(args, filter.SaleType)
		argID++
	}
	if filter.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(p.reference_no ILIKE $%d)", argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	whereSql := strings.Join(whereClauses, " AND ")

	orderSql := "p.sale_date DESC"
	if filter.OrderBy == "price_asc" {
		orderSql = "p.total ASC"
	} else if filter.OrderBy == "price_desc" {
		orderSql = "p.total DESC"
	} else if filter.OrderBy == "date_desc" {
		orderSql = "p.sale_date DESC"
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(p.id), COALESCE(SUM(p.total), 0), COALESCE(SUM(p.amount_paid), 0) FROM pos_sales p WHERE %s`, whereSql)
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&resp.TotalHits, &resp.TotalAmount, &resp.TotalPaid)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
SELECT p.id, p.reference_no, p.sale_type, p.total, p.amount_paid, p.sale_date, c.name
FROM pos_sales p
LEFT JOIN customers c ON p.customer_id = c.id
WHERE %s
ORDER BY %s
LIMIT $%d OFFSET $%d
`, whereSql, orderSql, argID, argID+1)

	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.POSSaleReportItem
		err := rows.Scan(&item.ID, &item.ReferenceNo, &item.SaleType, &item.Total, &item.AmountPaid, &item.SaleDate, &item.CustomerName)
		if err != nil {
			return nil, err
		}
		resp.Data = append(resp.Data, item)
	}

	if resp.Data == nil {
		resp.Data = []model.POSSaleReportItem{}
	}

	return &resp, nil
}

func (r *ReportRepo) GetOrdersReport(ctx context.Context, filter model.ReportFilter) (*model.OrderReportResponse, error) {
	var resp model.OrderReportResponse

	offset := (filter.Page - 1) * filter.Limit

	whereClauses := []string{"1=1"}
	var args []interface{}
	argID := 1

	if filter.SaleType != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("sale_type = $%d", argID))
		args = append(args, filter.SaleType)
		argID++
	}
	if filter.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(order_number ILIKE $%d OR customer_name ILIKE $%d)", argID, argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	whereSql := strings.Join(whereClauses, " AND ")

	orderSql := "created_at DESC"
	if filter.OrderBy == "price_asc" {
		orderSql = "total ASC"
	} else if filter.OrderBy == "price_desc" {
		orderSql = "total DESC"
	} else if filter.OrderBy == "date_desc" {
		orderSql = "created_at DESC"
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(id), COALESCE(SUM(total), 0) FROM orders WHERE %s`, whereSql)
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&resp.TotalHits, &resp.TotalAmount)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
SELECT id, order_number, customer_name, sale_type, order_status, payment_status, total, created_at
FROM orders
WHERE %s
ORDER BY %s
LIMIT $%d OFFSET $%d
`, whereSql, orderSql, argID, argID+1)

	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderReportItem
		err := rows.Scan(&item.ID, &item.OrderNumber, &item.CustomerName, &item.SaleType, &item.OrderStatus, &item.PaymentStatus, &item.Total, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		resp.Data = append(resp.Data, item)
	}

	if resp.Data == nil {
		resp.Data = []model.OrderReportItem{}
	}

	return &resp, nil
}

func (r *ReportRepo) GetCustomerDueReport(ctx context.Context, filter model.ReportFilter) (*model.CustomerDueReportResponse, error) {
	var resp model.CustomerDueReportResponse

	offset := (filter.Page - 1) * filter.Limit

	whereClauses := []string{"1=1"}
	var args []interface{}
	argID := 1

	if filter.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(name ILIKE $%d OR phone ILIKE $%d)", argID, argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	whereSql := strings.Join(whereClauses, " AND ")

	// Define a CTE to compute due per customer from pos_sales and orders
	dueCTE := `
WITH customer_dues AS (
SELECT c.id as customer_id, c.name, c.phone, 
(
(SELECT COALESCE(SUM(total - amount_paid), 0) FROM pos_sales WHERE customer_id = c.id) +
(SELECT COALESCE(SUM(total), 0) FROM orders WHERE customer_id = c.id AND payment_status = 'pending')
) as total_due
FROM customers c
WHERE %s
)
`

	countQuery := fmt.Sprintf(dueCTE+` SELECT COUNT(*), COALESCE(SUM(total_due), 0) FROM customer_dues WHERE total_due > 0`, whereSql)
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&resp.TotalHits, &resp.TotalDue)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(dueCTE+`
SELECT customer_id, name, phone, total_due 
FROM customer_dues 
WHERE total_due > 0 
ORDER BY total_due DESC 
LIMIT $%d OFFSET $%d
`, whereSql, argID, argID+1)

	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.CustomerDueReportItem
		err := rows.Scan(&item.CustomerID, &item.Name, &item.Phone, &item.TotalDue)
		if err != nil {
			return nil, err
		}
		resp.Data = append(resp.Data, item)
	}

	if resp.Data == nil {
		resp.Data = []model.CustomerDueReportItem{}
	}

	return &resp, nil
}

func (r *ReportRepo) GetSupplierDueReport(ctx context.Context, filter model.ReportFilter) (*model.SupplierDueReportResponse, error) {
	var resp model.SupplierDueReportResponse

	offset := (filter.Page - 1) * filter.Limit

	whereClauses := []string{"1=1"}
	var args []interface{}
	argID := 1

	if filter.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(s.name ILIKE $%d OR s.supplier_code ILIKE $%d OR s.phone ILIKE $%d)", argID, argID, argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	whereSql := strings.Join(whereClauses, " AND ")

	dueCTE := `
WITH supplier_dues AS (
SELECT s.id as supplier_id, s.supplier_code, s.name, s.phone,
(SELECT COALESCE(SUM(grand_total - paid_amount), 0) FROM purchases WHERE party_id = s.id) as total_due
FROM suppliers s
WHERE %s
)
`

	countQuery := fmt.Sprintf(dueCTE+` SELECT COUNT(*), COALESCE(SUM(total_due), 0) FROM supplier_dues WHERE total_due > 0`, whereSql)
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&resp.TotalHits, &resp.TotalDue)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(dueCTE+`
SELECT supplier_id, supplier_code, name, phone, total_due 
FROM supplier_dues 
WHERE total_due > 0 
ORDER BY total_due DESC 
LIMIT $%d OFFSET $%d
`, whereSql, argID, argID+1)

	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.SupplierDueReportItem
		err := rows.Scan(&item.SupplierID, &item.SupplierCode, &item.Name, &item.Phone, &item.TotalDue)
		if err != nil {
			return nil, err
		}
		resp.Data = append(resp.Data, item)
	}

	if resp.Data == nil {
		resp.Data = []model.SupplierDueReportItem{}
	}

	return &resp, nil
}

func (r *ReportRepo) GetLowStockReport(ctx context.Context, filter model.ReportFilter) (*model.LowStockReportResponse, error) {
	var resp model.LowStockReportResponse

	offset := (filter.Page - 1) * filter.Limit

	// Ensure we only retrieve items where stock is low (e.g. <= min_retail_order_qty or < 10)
	whereClauses := []string{"(current_stock_qty <= COALESCE(min_retail_order_qty, 5))"}
	var args []interface{}
	argID := 1

	if filter.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(name ILIKE $%d OR sku ILIKE $%d)", argID, argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	whereSql := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(id) FROM products WHERE %s`, whereSql)
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&resp.TotalHits)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
SELECT id, name, sku, current_stock_qty, min_retail_order_qty, retail_price
FROM products
WHERE %s
ORDER BY current_stock_qty ASC
LIMIT $%d OFFSET $%d
`, whereSql, argID, argID+1)

	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.LowStockReportItem
		var minRetail *float64
		err := rows.Scan(&item.ProductID, &item.Name, &item.SKU, &item.CurrentStockQty, &minRetail, &item.RetailPrice)
		if err != nil {
			return nil, err
		}
		if minRetail != nil {
			item.MinRetailOrderQty = *minRetail
		}
		resp.Data = append(resp.Data, item)
	}

	if resp.Data == nil {
		resp.Data = []model.LowStockReportItem{}
	}

	return &resp, nil
}
