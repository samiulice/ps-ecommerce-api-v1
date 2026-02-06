package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type OrderService struct {
	orderRepo    *repository.OrderRepo
	customerRepo *repository.CustomerRepository
}

func NewOrderService(orderRepo *repository.OrderRepo, customerRepo *repository.CustomerRepository) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
	}
}

// PlaceOrder handles the complete order placement flow
func (s *OrderService) PlaceOrder(ctx context.Context, req model.CreateOrderRequest) (*model.OrderResponse, error) {
	// Validate request
	if err := s.validateOrderRequest(req); err != nil {
		return nil, err
	}

	// Handle customer creation if requested
	var customerID sql.NullInt64
	if req.Customer.ID > 0 {
		customerID = sql.NullInt64{Int64: req.Customer.ID, Valid: true}
	} else if req.CreateAccount && req.Customer.Email != "" && req.Customer.Password != "" {
		// Create new customer account
		newCustomer := &model.Customer{
			Name:     model.ToNullString(req.Customer.Name),
			FName:    model.ToNullString(req.Customer.Name),
			Phone:    req.Customer.Mobile,
			Email:    model.ToNullString(req.Customer.Email),
			City:     model.ToNullString(req.Customer.City),
			IsActive: true,
		}

		//Generate password hash
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Customer.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		newCustomer.Password = string(hashedPassword)
		err = s.customerRepo.Create(ctx, newCustomer)
		if err != nil {
			// Don't fail the order if account creation fails, just log it
			fmt.Printf("Warning: Failed to create customer account: %v\n", err)
		} else {
			customerID = sql.NullInt64{Int64: newCustomer.ID, Valid: true}
		}
	}

	// Build order model
	order := &model.Order{
		CustomerID:     customerID,
		CustomerName:   req.Customer.Name,
		CustomerMobile: req.Customer.Mobile,
		CustomerEmail:  model.ToNullString(req.Customer.Email),
		CustomerArea:   model.ToNullString(req.Customer.Area),
		CustomerCity:   model.ToNullString(req.Customer.City),
		PaymentMethod:  req.PaymentMethod,
		PaymentStatus:  "pending",
		OrderStatus:    "pending",
		Subtotal:       req.Subtotal,
		ShippingCost:   req.ShippingCost,
		Discount:       req.Discount,
		Tax:            req.Tax,
		Total:          req.Total,
		OrderNote:      model.ToNullString(req.OrderNote),
	}

	// Build order items
	var items []model.OrderItem
	for _, item := range req.Items {
		orderItem := model.OrderItem{
			ProductID:   sql.NullInt64{Int64: item.ProductID, Valid: item.ProductID > 0},
			ProductName: item.Name,
			Quantity:    item.Qty,
			UnitPrice:   item.Price,
			TotalPrice:  item.Price * float64(item.Qty),
			Discount:    0,
			Tax:         0,
		}
		items = append(items, orderItem)
	}

	// Create order in database
	createdOrder, err := s.orderRepo.Create(ctx, order, items)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Build response
	response := createdOrder.ToOrderResponse()

	// Add items to response
	var itemResponses []model.OrderItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToOrderItemResponse())
	}
	response.Items = itemResponses

	return &response, nil
}

// validateOrderRequest validates the order request
func (s *OrderService) validateOrderRequest(req model.CreateOrderRequest) error {
	if req.Customer.Name == "" {
		return errors.New("customer name is required")
	}
	if req.Customer.Mobile == "" {
		return errors.New("customer mobile is required")
	}
	if len(req.Items) == 0 {
		return errors.New("order must have at least one item")
	}
	if req.PaymentMethod == "" {
		return errors.New("payment method is required")
	}
	if req.Total <= 0 {
		return errors.New("order total must be greater than zero")
	}

	// Validate create account requirements
	if req.CreateAccount {
		if req.Customer.Email == "" {
			return errors.New("email is required for account creation")
		}
		if req.Customer.Password == "" {
			return errors.New("password is required for account creation")
		}
	}

	return nil
}

// GetOrderByID retrieves an order by ID with items
func (s *OrderService) GetOrderByID(ctx context.Context, id int64) (*model.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	// Get order items
	items, err := s.orderRepo.GetOrderItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	response := order.ToOrderResponse()

	// Add items
	var itemResponses []model.OrderItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToOrderItemResponse())
	}
	response.Items = itemResponses

	return &response, nil
}

// GetOrderByOrderNumber retrieves an order by order number
func (s *OrderService) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*model.OrderResponse, error) {
	order, err := s.orderRepo.FindByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	items, err := s.orderRepo.GetOrderItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s  %s\n\n", order.OrderStatus, order.PaymentStatus)
	response := order.ToOrderResponse()
	fmt.Printf("%s  %s\n\n", response.OrderStatus, response.PaymentStatus)
	var itemResponses []model.OrderItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToOrderItemResponse())
	}
	response.Items = itemResponses

	return &response, nil
}

// ListOrders retrieves orders with filters
func (s *OrderService) ListOrders(ctx context.Context, filter model.OrderFilter) ([]model.OrderResponse, int64, error) {
	orders, total, err := s.orderRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.OrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToOrderResponse())
	}

	return responses, total, nil
}

// GetCustomerOrders retrieves orders for a specific customer
func (s *OrderService) GetCustomerOrders(ctx context.Context, customerID int64, page, limit int) ([]model.OrderResponse, int64, error) {
	orders, total, err := s.orderRepo.GetCustomerOrders(ctx, customerID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.OrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToOrderResponse())
	}

	return responses, total, nil
}

// UpdateOrderStatus updates order status
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id int64, status, reason string) error {
	// Validate status
	validStatuses := map[string]bool{
		"pending": true, "confirmed": true, "processing": true,
		"shipped": true, "delivered": true, "cancelled": true, "returned": true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid order status: %s", status)
	}

	return s.orderRepo.UpdateStatus(ctx, id, status, reason)
}

// UpdatePaymentStatus updates payment status
func (s *OrderService) UpdatePaymentStatus(ctx context.Context, id int64, status string) error {
	validStatuses := map[string]bool{
		"pending": true, "paid": true, "failed": true, "refunded": true, "partially_refunded": true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid payment status: %s", status)
	}

	return s.orderRepo.UpdatePaymentStatus(ctx, id, status)
}

// DeleteOrder deletes an order
func (s *OrderService) DeleteOrder(ctx context.Context, id int64) error {
	return s.orderRepo.Delete(ctx, id)
}

// GetOrderStats retrieves order statistics
func (s *OrderService) GetOrderStats(ctx context.Context) (*model.OrderStats, error) {
	return s.orderRepo.GetStats(ctx)
}
