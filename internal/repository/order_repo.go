package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type OrderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{db: db}
}

// GenerateOrderNumber generates a unique order number
func (r *OrderRepo) GenerateOrderNumber(ctx context.Context) (string, error) {
	var orderNumber string
	err := r.db.QueryRow(ctx, "SELECT generate_order_number()").Scan(&orderNumber)
	if err != nil {
		// Fallback to manual generation if function fails
		prefix := "ORD"
		datePart := time.Now().Format("20060102")
		var count int
		err = r.db.QueryRow(ctx,
			"SELECT COUNT(*) + 1 FROM orders WHERE DATE(created_at) = CURRENT_DATE").Scan(&count)
		if err != nil {
			count = 1
		}
		orderNumber = fmt.Sprintf("%s-%s-%04d", prefix, datePart, count)
	}
	return orderNumber, nil
}

// Create creates a new order with items (within a transaction)
func (r *OrderRepo) Create(ctx context.Context, order *model.Order, items []model.OrderItem) (*model.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Generate order number
	orderNumber, err := r.GenerateOrderNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate order number: %w", err)
	}
	order.OrderNumber = orderNumber

	// Insert order
	query := `
		INSERT INTO orders (
			order_number, customer_id, customer_name, customer_mobile, customer_email,
			customer_area, customer_city, payment_method, payment_status, order_status,
			subtotal, shipping_cost, discount, tax, total, order_note
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		order.OrderNumber,
		order.CustomerID,
		order.CustomerName,
		order.CustomerMobile,
		order.CustomerEmail,
		order.CustomerArea,
		order.CustomerCity,
		order.PaymentMethod,
		order.PaymentStatus,
		order.OrderStatus,
		order.Subtotal,
		order.ShippingCost,
		order.Discount,
		order.Tax,
		order.Total,
		order.OrderNote,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	// Insert order items
	for i := range items {
		items[i].OrderID = order.ID
		itemQuery := `
			INSERT INTO order_items (
				order_id, product_id, product_name, quantity, unit_price,
				total_price, discount, tax, variation_info
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, created_at
		`
		err = tx.QueryRow(ctx, itemQuery,
			items[i].OrderID,
			items[i].ProductID,
			items[i].ProductName,
			items[i].Quantity,
			items[i].UnitPrice,
			items[i].TotalPrice,
			items[i].Discount,
			items[i].Tax,
			items[i].VariationInfo,
		).Scan(&items[i].ID, &items[i].CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

// FindByID retrieves an order by ID
func (r *OrderRepo) FindByID(ctx context.Context, id int64) (*model.Order, error) {
	query := `
		SELECT id, order_number, customer_id, customer_name, customer_mobile, customer_email,
			customer_area, customer_city, payment_method, payment_status, order_status,
			subtotal, shipping_cost, discount, tax, total, order_note,
			created_at, updated_at, delivered_at, cancelled_at, cancelled_reason
		FROM orders
		WHERE id = $1
	`

	order := &model.Order{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID, &order.OrderNumber, &order.CustomerID, &order.CustomerName,
		&order.CustomerMobile, &order.CustomerEmail, &order.CustomerArea, &order.CustomerCity,
		&order.PaymentMethod, &order.PaymentStatus, &order.OrderStatus,
		&order.Subtotal, &order.ShippingCost, &order.Discount, &order.Tax, &order.Total,
		&order.OrderNote, &order.CreatedAt, &order.UpdatedAt, &order.DeliveredAt,
		&order.CancelledAt, &order.CancelledReason,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	return order, nil
}

// FindByOrderNumber retrieves an order by order number
func (r *OrderRepo) FindByOrderNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	query := `
		SELECT id, order_number, customer_id, customer_name, customer_mobile, customer_email,
			customer_area, customer_city, payment_method, payment_status, order_status,
			subtotal, shipping_cost, discount, tax, total, order_note,
			created_at, updated_at, delivered_at, cancelled_at, cancelled_reason
		FROM orders
		WHERE order_number = $1
	`

	order := &model.Order{}
	err := r.db.QueryRow(ctx, query, orderNumber).Scan(
		&order.ID, &order.OrderNumber, &order.CustomerID, &order.CustomerName,
		&order.CustomerMobile, &order.CustomerEmail, &order.CustomerArea, &order.CustomerCity,
		&order.PaymentMethod, &order.PaymentStatus, &order.OrderStatus,
		&order.Subtotal, &order.ShippingCost, &order.Discount, &order.Tax, &order.Total,
		&order.OrderNote, &order.CreatedAt, &order.UpdatedAt, &order.DeliveredAt,
		&order.CancelledAt, &order.CancelledReason,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	return order, nil
}

// GetOrderItems retrieves all items for an order
func (r *OrderRepo) GetOrderItems(ctx context.Context, orderID int64) ([]model.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, product_name, quantity, unit_price,
			total_price, discount, tax, variation_info, created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	var items []model.OrderItem
	for rows.Next() {
		var item model.OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.ProductName,
			&item.Quantity, &item.UnitPrice, &item.TotalPrice,
			&item.Discount, &item.Tax, &item.VariationInfo, &item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// List retrieves orders with filters and pagination
func (r *OrderRepo) List(ctx context.Context, filter model.OrderFilter) ([]model.Order, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("order_status = $%d", argIndex))
		args = append(args, filter.Status)
		argIndex++
	}
	if filter.PaymentStatus != "" {
		conditions = append(conditions, fmt.Sprintf("payment_status = $%d", argIndex))
		args = append(args, filter.PaymentStatus)
		argIndex++
	}
	if filter.PaymentMethod != "" {
		conditions = append(conditions, fmt.Sprintf("payment_method = $%d", argIndex))
		args = append(args, filter.PaymentMethod)
		argIndex++
	}
	if filter.CustomerID > 0 {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", argIndex))
		args = append(args, filter.CustomerID)
		argIndex++
	}
	if filter.FromDate != "" {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, filter.FromDate)
		argIndex++
	}
	if filter.ToDate != "" {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, filter.ToDate+" 23:59:59")
		argIndex++
	}
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		conditions = append(conditions, fmt.Sprintf(
			"(order_number ILIKE $%d OR customer_name ILIKE $%d OR customer_mobile ILIKE $%d)",
			argIndex, argIndex, argIndex,
		))
		args = append(args, searchPattern)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders %s", whereClause)
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Pagination
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (filter.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Query orders
	query := fmt.Sprintf(`
		SELECT id, order_number, customer_id, customer_name, customer_mobile, customer_email,
			customer_area, customer_city, payment_method, payment_status, order_status,
			subtotal, shipping_cost, discount, tax, total, order_note,
			created_at, updated_at, delivered_at, cancelled_at, cancelled_reason
		FROM orders
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID, &order.OrderNumber, &order.CustomerID, &order.CustomerName,
			&order.CustomerMobile, &order.CustomerEmail, &order.CustomerArea, &order.CustomerCity,
			&order.PaymentMethod, &order.PaymentStatus, &order.OrderStatus,
			&order.Subtotal, &order.ShippingCost, &order.Discount, &order.Tax, &order.Total,
			&order.OrderNote, &order.CreatedAt, &order.UpdatedAt, &order.DeliveredAt,
			&order.CancelledAt, &order.CancelledReason,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// GetCustomerOrders retrieves orders for a specific customer
func (r *OrderRepo) GetCustomerOrders(ctx context.Context, customerID int64, page, limit int) ([]model.Order, int64, error) {
	filter := model.OrderFilter{
		CustomerID: customerID,
		Page:       page,
		Limit:      limit,
	}
	return r.List(ctx, filter)
}

// UpdateStatus updates the order status
func (r *OrderRepo) UpdateStatus(ctx context.Context, id int64, status string, reason string) error {
	var query string
	var args []interface{}

	if status == "cancelled" {
		query = `
			UPDATE orders 
			SET order_status = $1, cancelled_at = NOW(), cancelled_reason = $2, updated_at = NOW()
			WHERE id = $3
		`
		args = []interface{}{status, sql.NullString{String: reason, Valid: reason != ""}, id}
	} else if status == "delivered" {
		query = `
			UPDATE orders 
			SET order_status = $1, delivered_at = NOW(), updated_at = NOW()
			WHERE id = $2
		`
		args = []interface{}{status, id}
	} else {
		query = `
			UPDATE orders 
			SET order_status = $1, updated_at = NOW()
			WHERE id = $2
		`
		args = []interface{}{status, id}
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// UpdatePaymentStatus updates the payment status
func (r *OrderRepo) UpdatePaymentStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE orders SET payment_status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// Delete deletes an order (soft delete recommended, but this is hard delete)
func (r *OrderRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// GetStats retrieves order statistics for dashboard
func (r *OrderRepo) GetStats(ctx context.Context) (*model.OrderStats, error) {
	stats := &model.OrderStats{}

	// Total orders and revenue
	err := r.db.QueryRow(ctx, `
		SELECT 
			COUNT(*),
			COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN total ELSE 0 END), 0)
		FROM orders
	`).Scan(&stats.TotalOrders, &stats.TotalRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get total stats: %w", err)
	}

	// Orders by status
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE order_status = 'pending'`).Scan(&stats.PendingOrders)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE order_status IN ('confirmed', 'processing', 'shipped')`).Scan(&stats.ProcessingOrders)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE order_status = 'delivered'`).Scan(&stats.DeliveredOrders)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE order_status = 'cancelled'`).Scan(&stats.CancelledOrders)
	if err != nil {
		return nil, err
	}

	// Today's orders and revenue
	err = r.db.QueryRow(ctx, `
		SELECT 
			COUNT(*),
			COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN total ELSE 0 END), 0)
		FROM orders
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&stats.TodayOrders, &stats.TodayRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get today stats: %w", err)
	}

	return stats, nil
}
