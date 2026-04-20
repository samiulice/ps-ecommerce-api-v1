package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

// DeliveryRepository handles data access for delivery features.
type DeliveryRepository struct {
	db *pgxpool.Pool
}

// NewDeliveryRepo creates a new DeliveryRepository.
func NewDeliveryRepo(db *pgxpool.Pool) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

// CreateDeliveryMethod creates a new delivery method
func (r *DeliveryRepository) CreateDeliveryMethod(ctx context.Context, m *model.DeliveryMethod) error {
	query := `
		INSERT INTO delivery_methods (name, base_cost, estimated_days, is_active)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		m.Name, m.BaseCost, m.EstimatedDays, m.IsActive,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

// CreateDeliveryMan promotes a customer to delivery man
func (r *DeliveryRepository) CreateDeliveryMan(ctx context.Context, m *model.DeliveryMan) error {
	query := `
		INSERT INTO delivery_men (
			employee_id, identity_type, identity_number, identity_image,
			vehicle_type, vehicle_number, bank_name, account_no, account_holder_name,
			is_active, is_online
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		m.EmployeeID, m.IdentityType, m.IdentityNumber, m.IdentityImage,
		m.VehicleType, m.VehicleNumber, m.BankName, m.AccountNo, m.AccountHolderName,
		m.IsActive, m.IsOnline,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)

	if err != nil {
		return err
	}

	// Create an empty wallet right away for the delivery man
	walletQuery := `INSERT INTO delivery_wallets (delivery_man_id) VALUES ($1)`
	_, err = r.db.Exec(ctx, walletQuery, m.ID)
	return err
}

// AssignOrderToDelivery assigns an order to a delivery man
func (r *DeliveryRepository) AssignOrderToDelivery(ctx context.Context, d *model.OrderDelivery) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO order_deliveries (
			order_id, delivery_man_id, delivery_status,
			delivery_fee_collected, delivery_man_earning, assigned_at
		) VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (order_id) DO UPDATE SET
			delivery_man_id = EXCLUDED.delivery_man_id,
			delivery_status = EXCLUDED.delivery_status,
			assigned_at = NOW()
		RETURNING id, assigned_at, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		d.OrderID, d.DeliveryManID, d.DeliveryStatus,
		d.DeliveryFeeCollected, d.DeliveryManEarning,
	).Scan(&d.ID, &d.AssignedAt, &d.CreatedAt, &d.UpdatedAt)

	if err != nil {
		return err
	}

	// Update the order status to shipped since a delivery man is assigned
	orderQuery := `UPDATE orders SET order_status = 'shipped', updated_at = NOW() WHERE id = $1 AND order_status != 'delivered'`
	_, err = tx.Exec(ctx, orderQuery, d.OrderID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// UpdateOrderDeliveryStatus updates the rider's progress for an order.
func (r *DeliveryRepository) UpdateOrderDeliveryStatus(ctx context.Context, orderID int64, payload *model.OrderDelivery) error {
	var query string
	if payload.DeliveryStatus == "delivered" {
		query = `
                        UPDATE order_deliveries
                        SET delivery_status = $1, delivered_at = NOW()
                        WHERE order_id = $2
                        RETURNING id, delivery_man_id, delivery_status, delivery_man_earning, delivered_at
                `
	} else {
		query = `
                        UPDATE order_deliveries
                        SET delivery_status = $1
                        WHERE order_id = $2
                        RETURNING id, delivery_man_id, delivery_status, delivery_man_earning, delivered_at
                `
	}
	err := r.db.QueryRow(ctx, query, payload.DeliveryStatus, orderID).Scan(
		&payload.ID, &payload.DeliveryManID, &payload.DeliveryStatus,
		&payload.DeliveryManEarning, &payload.DeliveredAt,
	)
	if err != nil {
		return err
	}

	// Implicitly sync the core orders table status and payment
	if payload.DeliveryStatus == "delivered" {
		orderQuery := `
			UPDATE orders 
			SET order_status = 'delivered', 
			    delivered_at = NOW(),
			    payment_status = CASE WHEN payment_method = 'COD' THEN 'paid'::payment_status ELSE payment_status END
			WHERE id = $1
		`
		_, _ = r.db.Exec(ctx, orderQuery, orderID)
	} else if payload.DeliveryStatus == "failed" || payload.DeliveryStatus == "cancelled" {
		orderQuery := `
			UPDATE orders 
			SET order_status = 'cancelled', 
			    cancelled_at = NOW(),
			    cancelled_reason = 'Delivery reported as ' || $2
			WHERE id = $1 AND order_status != 'cancelled'
		`
		_, _ = r.db.Exec(ctx, orderQuery, orderID, payload.DeliveryStatus)
	} else if payload.DeliveryStatus == "out_for_delivery" {
		orderQuery := `UPDATE orders SET order_status = 'shipped', updated_at = NOW() WHERE id = $1 AND order_status != 'delivered'`
		_, _ = r.db.Exec(ctx, orderQuery, orderID)
	}

	return nil
}
func (r *DeliveryRepository) CreditWallet(ctx context.Context, deliveryManID int64, amount float64) error {
	query := `
		UPDATE delivery_wallets 
		SET total_earned = total_earned + $1, current_balance = current_balance + $1 
		WHERE delivery_man_id = $2
	`
	_, err := r.db.Exec(ctx, query, amount, deliveryManID)
	return err
}

// CreateWithdrawRequest requests a cashout to bank
func (r *DeliveryRepository) CreateWithdrawRequest(ctx context.Context, wr *model.WithdrawRequest) error {
	query := `
		INSERT INTO withdraw_requests (delivery_man_id, amount)
		VALUES ($1, $2) RETURNING id, status, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, wr.DeliveryManID, wr.Amount).Scan(
		&wr.ID, &wr.Status, &wr.CreatedAt, &wr.UpdatedAt,
	)
}

// GetDeliveryManByEmployeeID finds a rider by their primary customer ID
func (r *DeliveryRepository) GetDeliveryManByEmployeeID(ctx context.Context, customerID int64) (*model.DeliveryMan, error) {
	query := `
		SELECT d.id, d.employee_id, d.is_active, d.is_online, d.created_at,
		       COALESCE(e.name, ''), COALESCE(e.mobile, ''),
		       d.identity_type, d.identity_number, d.identity_image, d.vehicle_type, d.vehicle_number,
		       d.bank_name, d.account_no, d.account_holder_name
		FROM delivery_men d
		JOIN employees e ON e.id = d.employee_id
		WHERE d.employee_id = $1 LIMIT 1
	`
	m := &model.DeliveryMan{}
	err := r.db.QueryRow(ctx, query, customerID).Scan(
		&m.ID, &m.EmployeeID, &m.IsActive, &m.IsOnline, &m.CreatedAt,
		&m.EmployeeName, &m.EmployeeMobile,
		&m.IdentityType, &m.IdentityNumber, &m.IdentityImage, &m.VehicleType, &m.VehicleNumber,
		&m.BankName, &m.AccountNo, &m.AccountHolderName,
	)
	if err == pgx.ErrNoRows {
		return nil, nil // Return clear nil on not found
	}
	return m, err
}

func (r *DeliveryRepository) ListDeliveryMen(ctx context.Context) ([]model.DeliveryMan, error) {
	// Auto-sync employees who have delivery roles into delivery_men
	syncQuery := `
		INSERT INTO delivery_men (employee_id, is_active, is_online)
		SELECT e.id, true, true 
		FROM employees e
		LEFT JOIN roles r ON r.id = e.role_id
		WHERE (r.slug = 'delivery_man' OR e.role ILIKE '%delivery%')
		ON CONFLICT (employee_id) DO NOTHING;
	`
	r.db.Exec(ctx, syncQuery)

	walletSyncQuery := `
		INSERT INTO delivery_wallets (delivery_man_id, total_earned, total_withdrawn, current_balance)
		SELECT id, 0, 0, 0 FROM delivery_men
		ON CONFLICT (delivery_man_id) DO NOTHING;
	`
	r.db.Exec(ctx, walletSyncQuery)

	query := `
                SELECT
                        dm.id, dm.employee_id, COALESCE(e.name, ''), COALESCE(e.mobile, ''), dm.is_active, dm.is_online, dm.vehicle_type, dm.vehicle_number
                FROM delivery_men dm
                JOIN employees e ON dm.employee_id = e.id
                ORDER BY e.name ASC
        `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var men []model.DeliveryMan
	for rows.Next() {
		var d model.DeliveryMan
		if err := rows.Scan(&d.ID, &d.EmployeeID, &d.EmployeeName, &d.EmployeeMobile, &d.IsActive, &d.IsOnline, &d.VehicleType, &d.VehicleNumber); err != nil {
		}
		men = append(men, d)
	}
	return men, nil
}

type OrderDeliveryHistory struct {
	model.OrderDelivery
	DeliveryManName  *string `json:"delivery_man_name,omitempty"`
	DeliveryManPhone *string `json:"delivery_man_phone,omitempty"`
}

func (r *DeliveryRepository) GetDeliveryHistory(ctx context.Context, limit, offset int) ([]OrderDeliveryHistory, error) {
	query := `
                SELECT 
                        od.id, od.order_id, od.delivery_man_id, COALESCE(e.name, ''), COALESCE(e.mobile, ''), od.delivery_status, 
                        od.delivery_fee_collected, od.delivery_man_earning, od.assigned_at, od.delivered_at, od.created_at, od.updated_at
                FROM order_deliveries od
                LEFT JOIN delivery_men dm ON od.delivery_man_id = dm.id
                LEFT JOIN employees e ON dm.employee_id = e.id
                ORDER BY od.created_at DESC
                LIMIT $1 OFFSET $2
        `
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]OrderDeliveryHistory, 0)
	for rows.Next() {
		var d OrderDeliveryHistory
		if err := rows.Scan(&d.ID, &d.OrderID, &d.DeliveryManID, &d.DeliveryManName, &d.DeliveryManPhone, &d.DeliveryStatus,
			&d.DeliveryFeeCollected, &d.DeliveryManEarning, &d.AssignedAt, &d.DeliveredAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		history = append(history, d)
	}
	return history, nil
}

func (r *DeliveryRepository) ListDeliveryMethods(ctx context.Context) ([]model.DeliveryMethod, error) {
	query := `SELECT id, name, base_cost, is_active FROM delivery_methods ORDER BY id ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []model.DeliveryMethod
	for rows.Next() {
		var m model.DeliveryMethod
		if err := rows.Scan(&m.ID, &m.Name, &m.BaseCost, &m.IsActive); err != nil {
			return nil, err
		}
		methods = append(methods, m)
	}
	return methods, nil
}

func (r *DeliveryRepository) GetOrdersByDeliveryMan(ctx context.Context, dmID int64) ([]model.OrderDelivery, error) {
	query := `
SELECT 
    od.id, od.order_id, od.delivery_man_id, od.delivery_status, 
    od.delivery_fee_collected, od.delivery_man_earning, od.assigned_at, od.delivered_at,
    o.order_number, o.customer_name, o.customer_mobile, 
    COALESCE(o.customer_area, ''), COALESCE(o.customer_city, ''), o.total
FROM order_deliveries od
LEFT JOIN orders o ON od.order_id = o.id
WHERE od.delivery_man_id = $1
ORDER BY od.assigned_at DESC
`
	rows, err := r.db.Query(ctx, query, dmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]model.OrderDelivery, 0)
	for rows.Next() {
		var h model.OrderDelivery
		if err := rows.Scan(
			&h.ID, &h.OrderID, &h.DeliveryManID, &h.DeliveryStatus,
			&h.DeliveryFeeCollected, &h.DeliveryManEarning, &h.AssignedAt, &h.DeliveredAt,
			&h.OrderNumber, &h.CustomerName, &h.CustomerMobile,
			&h.CustomerArea, &h.CustomerCity, &h.OrderTotal,
		); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func (r *DeliveryRepository) GetWalletByDeliveryMan(ctx context.Context, dmID int64) (*model.DeliveryWallet, error) {
	query := `SELECT id, delivery_man_id, total_earned, total_withdrawn, current_balance FROM delivery_wallets WHERE delivery_man_id = $1 LIMIT 1`
	var w model.DeliveryWallet
	err := r.db.QueryRow(ctx, query, dmID).Scan(&w.ID, &w.DeliveryManID, &w.TotalEarned, &w.TotalWithdrawn, &w.CurrentBalance)
	if err != nil {
		w = model.DeliveryWallet{DeliveryManID: dmID, TotalEarned: 0, TotalWithdrawn: 0, CurrentBalance: 0, CompletedDeliveries: 0}
	}

	// Fetch completed deliveries count
	countQuery := `SELECT COUNT(*) FROM order_deliveries WHERE delivery_man_id = $1 AND delivery_status = 'delivered'`
	var count int64
	_ = r.db.QueryRow(ctx, countQuery, dmID).Scan(&count)
	w.CompletedDeliveries = count
	var assignedCount int64
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM order_deliveries WHERE delivery_man_id = $1`, dmID).Scan(&assignedCount)
	w.TotalAssigned = assignedCount

	// Get earnings for the last 15 days (grouped by absolute day, using arbitrary sum)
	// Example query to fill the curve backwards from today
	earningsQuery := `
		SELECT COALESCE(SUM(delivery_man_earning), 0)
		FROM (
			SELECT CURRENT_DATE - generate_series(0, 14) AS d
		) dates
		LEFT JOIN order_deliveries od 
			ON od.delivery_man_id = $1 
			AND od.delivery_status = 'delivered'
			AND DATE(od.delivered_at) = dates.d
		GROUP BY dates.d
		ORDER BY dates.d ASC
	`
	rows, err := r.db.Query(ctx, earningsQuery, dmID)
	w.RecentEarnings = make([]float64, 0, 15)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e float64
			if err := rows.Scan(&e); err == nil {
				w.RecentEarnings = append(w.RecentEarnings, e)
			}
		}
	}

	// Fetch delivery counts and cancellation counts per day
	countsQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN od.delivery_status = 'delivered' THEN 1 ELSE 0 END), 0) as delivered_cnt,
			COALESCE(SUM(CASE WHEN od.delivery_status = 'cancelled' OR od.delivery_status = 'failed' THEN 1 ELSE 0 END), 0) as cancelled_cnt
		FROM (
			SELECT CURRENT_DATE - generate_series(0, 14) AS d
		) dates
		LEFT JOIN order_deliveries od 
			ON od.delivery_man_id = $1 
			AND (DATE(od.delivered_at) = dates.d OR DATE(od.updated_at) = dates.d)
		GROUP BY dates.d
		ORDER BY dates.d ASC
	`
	rows2, err := r.db.Query(ctx, countsQuery, dmID)
	w.RecentDeliveries = make([]int64, 0, 15)
	w.RecentCancellations = make([]int64, 0, 15)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var dCnt, cCnt int64
			if err := rows2.Scan(&dCnt, &cCnt); err == nil {
				w.RecentDeliveries = append(w.RecentDeliveries, dCnt)
				w.RecentCancellations = append(w.RecentCancellations, cCnt)
			}
		}
	}

	return &w, nil
}
