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
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&resp.TotalItems, &resp.TotalAmount)
	if err != nil {
		return nil, err
	}

	resp.TotalPages = (resp.TotalItems + filter.Limit - 1) / filter.Limit
	if resp.TotalPages <= 0 {
		resp.TotalPages = 1
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
		resp.Items = append(resp.Items, item)
	}

	if resp.Items == nil {
		resp.Items = []model.OrderReportItem{}
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
	whereClauses := []string{"(current_stock_qty <= COALESCE(stock_alert_qty, 0))"}
	var args []interface{}
	argID := 1

	if filter.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(name ILIKE $%d OR sku ILIKE $%d)", argID, argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	whereSql := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(id) FROM products WHERE %s`, whereSql)
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&resp.TotalItems)
	if err != nil {
		return nil, err
	}
	
	resp.TotalPages = (resp.TotalItems + filter.Limit - 1) / filter.Limit
	if resp.TotalPages <= 0 {
		resp.TotalPages = 1
	}

	query := fmt.Sprintf(`
SELECT id, name, sku, current_stock_qty, stock_alert_qty, thumbnail
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
		var sku *string
		err := rows.Scan(&item.ProductID, &item.Name, &sku, &item.CurrentStockQty, &item.AlertQty, &item.ImageURL)
		if err != nil {
			return nil, err
		}
		if sku != nil {
			item.SKU = *sku
		}
		
		resp.Items = append(resp.Items, item)
	}

	if resp.Items == nil {
		resp.Items = []model.LowStockReportItem{}
	}

	return &resp, nil
}

func (r *ReportRepo) GetFinancialReport(ctx context.Context, filter model.ReportFilter) (*model.FinancialReportResponse, error) {
	var resp model.FinancialReportResponse

	offset := (filter.Page - 1) * filter.Limit

	// Build optional date filter clause for each CTE
	orderDateWhere := "WHERE order_status = 'delivered'"
	posDateWhere := "WHERE 1=1"
	purchaseDateWhere := "WHERE 1=1"
	expenseDateWhere := "WHERE 1=1"

	var args []interface{}
	argID := 1

	if filter.StartDate != "" {
		orderDateWhere += fmt.Sprintf(" AND created_at >= $%d::date", argID)
		posDateWhere += fmt.Sprintf(" AND sale_date >= $%d::date", argID)
		purchaseDateWhere += fmt.Sprintf(" AND purchase_date >= $%d::date", argID)
		expenseDateWhere += fmt.Sprintf(" AND expense_date >= $%d::date", argID)
		args = append(args, filter.StartDate)
		argID++
	}
	if filter.EndDate != "" {
		orderDateWhere += fmt.Sprintf(" AND created_at < ($%d::date + interval '1 day')", argID)
		posDateWhere += fmt.Sprintf(" AND sale_date < ($%d::date + interval '1 day')", argID)
		purchaseDateWhere += fmt.Sprintf(" AND purchase_date <= $%d::date", argID)
		expenseDateWhere += fmt.Sprintf(" AND expense_date < ($%d::date + interval '1 day')", argID)
		args = append(args, filter.EndDate)
		argID++
	}

	baseCTE := fmt.Sprintf(`
WITH order_daily AS (
    SELECT DATE(created_at) AS d, COALESCE(SUM(total), 0) AS amount
    FROM orders %s
    GROUP BY DATE(created_at)
),
pos_daily AS (
    SELECT DATE(sale_date) AS d, COALESCE(SUM(total), 0) AS amount
    FROM pos_sales %s
    GROUP BY DATE(sale_date)
),
purchase_daily AS (
    SELECT purchase_date AS d, COALESCE(SUM(grand_total), 0) AS amount
    FROM purchases %s
    GROUP BY purchase_date
),
expense_daily AS (
    SELECT DATE(expense_date) AS d, COALESCE(SUM(amount), 0) AS amount
    FROM expenses %s
    GROUP BY DATE(expense_date)
),
all_dates AS (
    SELECT d FROM order_daily
    UNION SELECT d FROM pos_daily
    UNION SELECT d FROM purchase_daily
    UNION SELECT d FROM expense_daily
),
daily_report AS (
    SELECT
        ad.d AS transaction_date,
        COALESCE(od.amount, 0) + COALESCE(pd.amount, 0) AS total_sales_income,
        0::numeric AS total_sales_refunds,
        COALESCE(pud.amount, 0) AS total_purchases,
        COALESCE(ed.amount, 0) AS total_expenses,
        (COALESCE(od.amount, 0) + COALESCE(pd.amount, 0) - COALESCE(pud.amount, 0) - COALESCE(ed.amount, 0)) AS net_profit
    FROM all_dates ad
    LEFT JOIN order_daily od ON ad.d = od.d
    LEFT JOIN pos_daily pd ON ad.d = pd.d
    LEFT JOIN purchase_daily pud ON ad.d = pud.d
    LEFT JOIN expense_daily ed ON ad.d = ed.d
)`, orderDateWhere, posDateWhere, purchaseDateWhere, expenseDateWhere)

	// Count total days
	countQuery := baseCTE + ` SELECT COUNT(*) FROM daily_report`
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&resp.TotalHits)
	if err != nil {
		return nil, fmt.Errorf("financial count query: %w", err)
	}

	// Fetch summary totals
	summaryQuery := baseCTE + `
SELECT COALESCE(SUM(total_sales_income), 0), COALESCE(SUM(total_sales_refunds), 0),
       COALESCE(SUM(total_purchases), 0), COALESCE(SUM(total_expenses), 0), COALESCE(SUM(net_profit), 0)
FROM daily_report`
	err = r.db.QueryRow(ctx, summaryQuery, args...).Scan(
		&resp.TotalSalesIncome, &resp.TotalSalesRefunds,
		&resp.TotalPurchases, &resp.TotalExpenses, &resp.TotalNetProfit,
	)
	if err != nil {
		return nil, fmt.Errorf("financial summary query: %w", err)
	}

	// Fetch paginated rows
	dataQuery := fmt.Sprintf(baseCTE+`
SELECT transaction_date, total_sales_income, total_sales_refunds, total_purchases, total_expenses, net_profit
FROM daily_report
ORDER BY transaction_date DESC
LIMIT $%d OFFSET $%d`, argID, argID+1)

	dataArgs := append(args, filter.Limit, offset)
	rows, err := r.db.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, fmt.Errorf("financial data query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.FinancialReportItem
		err := rows.Scan(
			&item.TransactionDate,
			&item.TotalSalesIncome,
			&item.TotalSalesRefunds,
			&item.TotalPurchases,
			&item.TotalExpenses,
			&item.NetProfit,
		)
		if err != nil {
			return nil, fmt.Errorf("financial row scan: %w", err)
		}
		resp.Data = append(resp.Data, item)
	}

	if resp.Data == nil {
		resp.Data = []model.FinancialReportItem{}
	}

	return &resp, nil
}

func (r *ReportRepo) GetIncomeStatement(ctx context.Context, filter model.ReportFilter) (*model.IncomeStatementResponse, error) {
	var resp model.IncomeStatementResponse
	resp.StartDate = filter.StartDate
	resp.EndDate = filter.EndDate

	// ── 1. Online Order Revenue (delivered orders only) ──
	orderQuery := `
		SELECT
			COALESCE(SUM(subtotal), 0),
			COALESCE(SUM(shipping_cost), 0),
			COALESCE(SUM(discount), 0),
			COALESCE(SUM(tax), 0),
			COUNT(*)
		FROM orders
		WHERE order_status = 'delivered'`
	orderArgs := []interface{}{}
	oArgID := 1
	if filter.StartDate != "" {
		orderQuery += fmt.Sprintf(" AND created_at >= $%d::date", oArgID)
		orderArgs = append(orderArgs, filter.StartDate)
		oArgID++
	}
	if filter.EndDate != "" {
		orderQuery += fmt.Sprintf(" AND created_at < ($%d::date + interval '1 day')", oArgID)
		orderArgs = append(orderArgs, filter.EndDate)
		oArgID++
	}
	var orderSubtotal, orderShipping, orderDiscount, orderTax float64
	err := r.db.QueryRow(ctx, orderQuery, orderArgs...).Scan(
		&orderSubtotal, &orderShipping, &orderDiscount, &orderTax, &resp.OnlineOrderCount,
	)
	if err != nil {
		return nil, fmt.Errorf("order revenue query: %w", err)
	}
	resp.OnlineOrderRevenue = orderSubtotal
	resp.ShippingIncome = orderShipping

	// ── 2. POS Sales Revenue ──
	posQuery := `
		SELECT
			COALESCE(SUM(subtotal), 0),
			COALESCE(SUM(discount), 0),
			COALESCE(SUM(tax), 0),
			COUNT(*)
		FROM pos_sales
		WHERE 1=1`
	posArgs := []interface{}{}
	pArgID := 1
	if filter.StartDate != "" {
		posQuery += fmt.Sprintf(" AND sale_date >= $%d::date", pArgID)
		posArgs = append(posArgs, filter.StartDate)
		pArgID++
	}
	if filter.EndDate != "" {
		posQuery += fmt.Sprintf(" AND sale_date < ($%d::date + interval '1 day')", pArgID)
		posArgs = append(posArgs, filter.EndDate)
		pArgID++
	}
	var posSubtotal, posDiscount, posTax float64
	err = r.db.QueryRow(ctx, posQuery, posArgs...).Scan(
		&posSubtotal, &posDiscount, &posTax, &resp.POSSalesCount,
	)
	if err != nil {
		return nil, fmt.Errorf("pos revenue query: %w", err)
	}
	resp.POSSalesRevenue = posSubtotal

	// ── Compute Revenue Totals ──
	resp.TotalDiscounts = orderDiscount + posDiscount
	resp.TaxCollected = orderTax + posTax
	resp.GrossRevenue = resp.OnlineOrderRevenue + resp.POSSalesRevenue + resp.ShippingIncome
	resp.NetRevenue = resp.GrossRevenue - resp.TotalDiscounts

	// ── 3. Purchases (COGS) ──
	purchaseQuery := `
		SELECT COALESCE(SUM(grand_total), 0), COALESCE(SUM(shipping_charge), 0), COUNT(*)
		FROM purchases
		WHERE 1=1`
	purchaseArgs := []interface{}{}
	puArgID := 1
	if filter.StartDate != "" {
		purchaseQuery += fmt.Sprintf(" AND purchase_date >= $%d::date", puArgID)
		purchaseArgs = append(purchaseArgs, filter.StartDate)
		puArgID++
	}
	if filter.EndDate != "" {
		purchaseQuery += fmt.Sprintf(" AND purchase_date <= $%d::date", puArgID)
		purchaseArgs = append(purchaseArgs, filter.EndDate)
		puArgID++
	}
	err = r.db.QueryRow(ctx, purchaseQuery, purchaseArgs...).Scan(
		&resp.TotalPurchases, &resp.PurchaseShipping, &resp.PurchaseCount,
	)
	if err != nil {
		return nil, fmt.Errorf("purchases query: %w", err)
	}

	// ── 4. Purchase Returns ──
	returnQuery := `
		SELECT COALESCE(SUM(grand_total), 0)
		FROM purchase_return
		WHERE 1=1`
	returnArgs := []interface{}{}
	rArgID := 1
	if filter.StartDate != "" {
		returnQuery += fmt.Sprintf(" AND return_date >= $%d::date", rArgID)
		returnArgs = append(returnArgs, filter.StartDate)
		rArgID++
	}
	if filter.EndDate != "" {
		returnQuery += fmt.Sprintf(" AND return_date <= $%d::date", rArgID)
		returnArgs = append(returnArgs, filter.EndDate)
		rArgID++
	}
	err = r.db.QueryRow(ctx, returnQuery, returnArgs...).Scan(&resp.PurchaseReturns)
	if err != nil {
		return nil, fmt.Errorf("purchase returns query: %w", err)
	}

	// ── Compute COGS ──
	resp.NetCOGS = resp.TotalPurchases - resp.PurchaseReturns + resp.PurchaseShipping

	// ── Gross Profit ──
	resp.GrossProfit = resp.NetRevenue - resp.NetCOGS
	if resp.NetRevenue > 0 {
		resp.GrossProfitMargin = (resp.GrossProfit / resp.NetRevenue) * 100
	}

	// ── 5. Operating Expenses by Category ──
	expenseQuery := `
		SELECT ec.name, COALESCE(SUM(e.amount), 0)
		FROM expenses e
		JOIN expense_categories ec ON e.category_id = ec.id
		WHERE 1=1`
	expenseArgs := []interface{}{}
	eArgID := 1
	if filter.StartDate != "" {
		expenseQuery += fmt.Sprintf(" AND e.expense_date >= $%d::date", eArgID)
		expenseArgs = append(expenseArgs, filter.StartDate)
		eArgID++
	}
	if filter.EndDate != "" {
		expenseQuery += fmt.Sprintf(" AND e.expense_date < ($%d::date + interval '1 day')", eArgID)
		expenseArgs = append(expenseArgs, filter.EndDate)
		eArgID++
	}
	expenseQuery += " GROUP BY ec.name ORDER BY SUM(e.amount) DESC"

	rows, err := r.db.Query(ctx, expenseQuery, expenseArgs...)
	if err != nil {
		return nil, fmt.Errorf("expenses query: %w", err)
	}
	defer rows.Close()

	resp.ExpenseBreakdown = []model.ExpenseCategoryItem{}
	resp.TotalExpenses = 0
	for rows.Next() {
		var item model.ExpenseCategoryItem
		if err := rows.Scan(&item.CategoryName, &item.Amount); err != nil {
			return nil, fmt.Errorf("expense scan: %w", err)
		}
		resp.TotalExpenses += item.Amount
		resp.ExpenseBreakdown = append(resp.ExpenseBreakdown, item)
	}

	// ── Net Income ──
	resp.NetIncome = resp.GrossProfit - resp.TotalExpenses
	if resp.NetRevenue > 0 {
		resp.NetIncomeMargin = (resp.NetIncome / resp.NetRevenue) * 100
	}

	return &resp, nil
}

