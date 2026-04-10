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
	return r.db.QueryRow(ctx, query,
		d.OrderID, d.DeliveryManID, d.DeliveryStatus,
		d.DeliveryFeeCollected, d.DeliveryManEarning,
	).Scan(&d.ID, &d.AssignedAt, &d.CreatedAt, &d.UpdatedAt)
}

// UpdateOrderDeliveryStatus updates the rider's progress for an order.
func (r *DeliveryRepository) UpdateOrderDeliveryStatus(ctx context.Context, orderID int64, payload *model.OrderDelivery) error {
	query := `
		UPDATE order_deliveries
		SET delivery_status = $1, delivered_at = CASE WHEN $1 = 'delivered' THEN NOW() ELSE delivered_at END
		WHERE order_id = $2
		RETURNING id, delivery_man_id, delivery_status, delivery_man_earning, delivered_at
	`
	return r.db.QueryRow(ctx, query, payload.DeliveryStatus, orderID).Scan(
		&payload.ID, &payload.DeliveryManID, &payload.DeliveryStatus,
		&payload.DeliveryManEarning, &payload.DeliveredAt,
	)
}

// CreditWallet earns the delivery man a commision
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
		       COALESCE(e.name, ''), COALESCE(e.mobile, '')
		FROM delivery_men d
		JOIN employees e ON e.id = d.employee_id
		WHERE d.employee_id = $1 LIMIT 1
	`
	m := &model.DeliveryMan{}
	err := r.db.QueryRow(ctx, query, customerID).Scan(
		&m.ID, &m.EmployeeID, &m.IsActive, &m.IsOnline, &m.CreatedAt,
		&m.EmployeeName, &m.EmployeeMobile,
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
                LEFT JOIN employees e ON dm.employee_id = c.id
                ORDER BY od.created_at DESC
                LIMIT $1 OFFSET $2
        `
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []OrderDeliveryHistory
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
SELECT id, order_id, delivery_man_id, delivery_status, delivery_fee_collected, delivery_man_earning, assigned_at, delivered_at
FROM order_deliveries
WHERE delivery_man_id = $1
ORDER BY assigned_at DESC
`
rows, err := r.db.Query(ctx, query, dmID)
if err != nil {
return nil, err
}
defer rows.Close()

var history []model.OrderDelivery
for rows.Next() {
var h model.OrderDelivery
if err := rows.Scan(&h.ID, &h.OrderID, &h.DeliveryManID, &h.DeliveryStatus, &h.DeliveryFeeCollected, &h.DeliveryManEarning, &h.AssignedAt, &h.DeliveredAt); err != nil {
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
return &model.DeliveryWallet{DeliveryManID: dmID, TotalEarned: 0, TotalWithdrawn: 0, CurrentBalance: 0}, nil
}
return &w, err
}
