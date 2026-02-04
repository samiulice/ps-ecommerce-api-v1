package model

import (
	"database/sql"
	"time"
)

// Order represents the main order entity
type Order struct {
	ID              int64          `json:"id"`
	OrderNumber     string         `json:"order_number"`
	CustomerID      sql.NullInt64  `json:"customer_id"`
	CustomerName    string         `json:"customer_name"`
	CustomerMobile  string         `json:"customer_mobile"`
	CustomerEmail   sql.NullString `json:"customer_email"`
	CustomerArea    sql.NullString `json:"customer_area"`
	CustomerCity    sql.NullString `json:"customer_city"`
	PaymentMethod   string         `json:"payment_method"`
	PaymentStatus   string         `json:"payment_status"`
	OrderStatus     string         `json:"order_status"`
	Subtotal        float64        `json:"subtotal"`
	ShippingCost    float64        `json:"shipping_cost"`
	Discount        float64        `json:"discount"`
	Tax             float64        `json:"tax"`
	Total           float64        `json:"total"`
	OrderNote       sql.NullString `json:"order_note"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeliveredAt     sql.NullTime   `json:"delivered_at"`
	CancelledAt     sql.NullTime   `json:"cancelled_at"`
	CancelledReason sql.NullString `json:"cancelled_reason"`
}

// OrderItem represents a line item in an order
type OrderItem struct {
	ID            int64          `json:"id"`
	OrderID       int64          `json:"order_id"`
	ProductID     sql.NullInt64  `json:"product_id"`
	ProductName   string         `json:"product_name"`
	Quantity      int            `json:"quantity"`
	UnitPrice     float64        `json:"unit_price"`
	TotalPrice    float64        `json:"total_price"`
	Discount      float64        `json:"discount"`
	Tax           float64        `json:"tax"`
	VariationInfo sql.NullString `json:"variation_info"`
	CreatedAt     time.Time      `json:"created_at"`
}

// --- Request DTOs ---

// OrderCustomer represents customer info in checkout request
type OrderCustomer struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	Area     string `json:"area"`
	City     string `json:"city"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// OrderItemRequest represents an item in checkout request
type OrderItemRequest struct {
	ProductID int64   `json:"product_id,omitempty"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Qty       int     `json:"qty"`
}

// CreateOrderRequest is the payload from checkout
type CreateOrderRequest struct {
	Customer      OrderCustomer      `json:"customer"`
	CreateAccount bool               `json:"create_account"`
	PaymentMethod string             `json:"payment_method"`
	Items         []OrderItemRequest `json:"items"`
	Subtotal      float64            `json:"subtotal"`
	ShippingCost  float64            `json:"shipping_cost,omitempty"`
	Discount      float64            `json:"discount,omitempty"`
	Tax           float64            `json:"tax,omitempty"`
	Total         float64            `json:"total"`
	OrderNote     string             `json:"order_note,omitempty"`
}

// UpdateOrderStatusRequest for status updates
type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// --- Response DTOs ---

// OrderItemResponse for API response
type OrderItemResponse struct {
	ID            int64   `json:"id"`
	ProductID     *int64  `json:"product_id,omitempty"`
	ProductName   string  `json:"product_name"`
	Quantity      int     `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	TotalPrice    float64 `json:"total_price"`
	Discount      float64 `json:"discount"`
	Tax           float64 `json:"tax"`
	VariationInfo *string `json:"variation_info,omitempty"`
}

// OrderResponse for API response
type OrderResponse struct {
	ID              int64               `json:"id"`
	OrderNumber     string              `json:"order_number"`
	CustomerID      *int64              `json:"customer_id,omitempty"`
	CustomerName    string              `json:"customer_name"`
	CustomerMobile  string              `json:"customer_mobile"`
	CustomerEmail   *string             `json:"customer_email,omitempty"`
	CustomerArea    *string             `json:"customer_area,omitempty"`
	CustomerCity    *string             `json:"customer_city,omitempty"`
	PaymentMethod   string              `json:"payment_method"`
	PaymentStatus   string              `json:"payment_status"`
	OrderStatus     string              `json:"order_status"`
	Subtotal        float64             `json:"subtotal"`
	ShippingCost    float64             `json:"shipping_cost"`
	Discount        float64             `json:"discount"`
	Tax             float64             `json:"tax"`
	Total           float64             `json:"total"`
	OrderNote       *string             `json:"order_note,omitempty"`
	Items           []OrderItemResponse `json:"items,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	DeliveredAt     *time.Time          `json:"delivered_at,omitempty"`
	CancelledAt     *time.Time          `json:"cancelled_at,omitempty"`
	CancelledReason *string             `json:"cancelled_reason,omitempty"`
}

// OrderFilter for listing orders with filters
type OrderFilter struct {
	Status        string // pending, confirmed, processing, shipped, delivered, cancelled
	PaymentStatus string // pending, paid, failed, refunded
	PaymentMethod string // COD, ONLINE, etc.
	CustomerID    int64
	FromDate      string
	ToDate        string
	Search        string // order number or customer name/mobile
	Page          int
	Limit         int
}

// OrderStats for dashboard
type OrderStats struct {
	TotalOrders      int64   `json:"total_orders"`
	PendingOrders    int64   `json:"pending_orders"`
	ProcessingOrders int64   `json:"processing_orders"`
	DeliveredOrders  int64   `json:"delivered_orders"`
	CancelledOrders  int64   `json:"cancelled_orders"`
	TotalRevenue     float64 `json:"total_revenue"`
	TodayOrders      int64   `json:"today_orders"`
	TodayRevenue     float64 `json:"today_revenue"`
}

// --- Helper Functions ---

// ToOrderResponse converts Order model to OrderResponse
func (o *Order) ToOrderResponse() OrderResponse {
	resp := OrderResponse{
		ID:             o.ID,
		OrderNumber:    o.OrderNumber,
		CustomerName:   o.CustomerName,
		CustomerMobile: o.CustomerMobile,
		PaymentMethod:  o.PaymentMethod,
		PaymentStatus:  o.PaymentStatus,
		OrderStatus:    o.OrderStatus,
		Subtotal:       o.Subtotal,
		ShippingCost:   o.ShippingCost,
		Discount:       o.Discount,
		Tax:            o.Tax,
		Total:          o.Total,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
	}

	if o.CustomerID.Valid {
		resp.CustomerID = &o.CustomerID.Int64
	}
	if o.CustomerEmail.Valid {
		resp.CustomerEmail = &o.CustomerEmail.String
	}
	if o.CustomerArea.Valid {
		resp.CustomerArea = &o.CustomerArea.String
	}
	if o.CustomerCity.Valid {
		resp.CustomerCity = &o.CustomerCity.String
	}
	if o.OrderNote.Valid {
		resp.OrderNote = &o.OrderNote.String
	}
	if o.DeliveredAt.Valid {
		resp.DeliveredAt = &o.DeliveredAt.Time
	}
	if o.CancelledAt.Valid {
		resp.CancelledAt = &o.CancelledAt.Time
	}
	if o.CancelledReason.Valid {
		resp.CancelledReason = &o.CancelledReason.String
	}

	return resp
}

// ToOrderItemResponse converts OrderItem model to OrderItemResponse
func (oi *OrderItem) ToOrderItemResponse() OrderItemResponse {
	resp := OrderItemResponse{
		ID:          oi.ID,
		ProductName: oi.ProductName,
		Quantity:    oi.Quantity,
		UnitPrice:   oi.UnitPrice,
		TotalPrice:  oi.TotalPrice,
		Discount:    oi.Discount,
		Tax:         oi.Tax,
	}

	if oi.ProductID.Valid {
		resp.ProductID = &oi.ProductID.Int64
	}
	if oi.VariationInfo.Valid {
		resp.VariationInfo = &oi.VariationInfo.String
	}

	return resp
}
