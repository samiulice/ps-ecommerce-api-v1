package model

import (
	"database/sql"
)

// DeliveryMethod represents a shipping/delivery option.
type DeliveryMethod struct {
	ID            int64          `json:"id" db:"id"`
	Name          string         `json:"name" db:"name"`
	BaseCost      float64        `json:"base_cost" db:"base_cost"`
	EstimatedDays sql.NullString `json:"estimated_days,omitempty" db:"estimated_days"`
	IsActive      bool           `json:"is_active" db:"is_active"`
	CreatedAt     sql.NullTime   `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt     sql.NullTime   `json:"updated_at,omitempty" db:"updated_at"`
}

// DeliveryMan represents an employee who delivers orders.
type DeliveryMan struct {
	ID                int64          `json:"id" db:"id"`
	CustomerID        int64          `json:"customer_id" db:"customer_id"`
	IdentityType      sql.NullString `json:"identity_type,omitempty" db:"identity_type"`
	IdentityNumber    sql.NullString `json:"identity_number,omitempty" db:"identity_number"`
	IdentityImage     sql.NullString `json:"identity_image,omitempty" db:"identity_image"`
	VehicleType       sql.NullString `json:"vehicle_type,omitempty" db:"vehicle_type"`
	VehicleNumber     sql.NullString `json:"vehicle_number,omitempty" db:"vehicle_number"`
	BankName          sql.NullString `json:"bank_name,omitempty" db:"bank_name"`
	AccountNo         sql.NullString `json:"account_no,omitempty" db:"account_no"`
	AccountHolderName sql.NullString `json:"account_holder_name,omitempty" db:"account_holder_name"`
	IsActive          bool           `json:"is_active" db:"is_active"`
	IsOnline          bool           `json:"is_online" db:"is_online"`
	CreatedAt         sql.NullTime   `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt         sql.NullTime   `json:"updated_at,omitempty" db:"updated_at"`

	// Nested/Joined Customer profile details
	CustomerName  sql.NullString `json:"customer_name,omitempty" db:"name"`
	CustomerPhone sql.NullString `json:"customer_phone,omitempty" db:"phone"`
}

// OrderDelivery maps an order to a delivery man and tracks status.
type OrderDelivery struct {
	ID                   int64         `json:"id" db:"id"`
	OrderID              int64         `json:"order_id" db:"order_id"`
	DeliveryManID        sql.NullInt64 `json:"delivery_man_id,omitempty" db:"delivery_man_id"`
	DeliveryStatus       string        `json:"delivery_status" db:"delivery_status"`
	DeliveryFeeCollected float64       `json:"delivery_fee_collected" db:"delivery_fee_collected"`
	DeliveryManEarning   float64       `json:"delivery_man_earning" db:"delivery_man_earning"`
	AssignedAt           sql.NullTime  `json:"assigned_at,omitempty" db:"assigned_at"`
	DeliveredAt          sql.NullTime  `json:"delivered_at,omitempty" db:"delivered_at"`
	CreatedAt            sql.NullTime  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt            sql.NullTime  `json:"updated_at,omitempty" db:"updated_at"`
}

// DeliveryWallet stores the earnings and withdrawals.
type DeliveryWallet struct {
	ID             int64        `json:"id" db:"id"`
	DeliveryManID  int64        `json:"delivery_man_id" db:"delivery_man_id"`
	TotalEarned    float64      `json:"total_earned" db:"total_earned"`
	TotalWithdrawn float64      `json:"total_withdrawn" db:"total_withdrawn"`
	CurrentBalance float64      `json:"current_balance" db:"current_balance"`
	CreatedAt      sql.NullTime `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      sql.NullTime `json:"updated_at,omitempty" db:"updated_at"`
}

// WithdrawRequest tracks a delivery man's request to transfer wallet funds to bank.
type WithdrawRequest struct {
	ID            int64          `json:"id" db:"id"`
	DeliveryManID int64          `json:"delivery_man_id" db:"delivery_man_id"`
	Amount        float64        `json:"amount" db:"amount"`
	Status        string         `json:"status" db:"status"` // pending, approved, rejected
	AdminNote     sql.NullString `json:"admin_note,omitempty" db:"admin_note"`
	CreatedAt     sql.NullTime   `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt     sql.NullTime   `json:"updated_at,omitempty" db:"updated_at"`
}
