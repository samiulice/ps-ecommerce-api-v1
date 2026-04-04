package service

import (
	"context"
	"errors"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

// DeliveryService handles business logic for delivery management.
type DeliveryService struct {
	repo *repository.DeliveryRepository
}

// NewDeliveryService creates a DeliveryService.
func NewDeliveryService(repo *repository.DeliveryRepository) *DeliveryService {
	return &DeliveryService{repo: repo}
}

// AddDeliveryMethod adds a new delivery method to the platform (e.g., standard, express)
func (s *DeliveryService) AddDeliveryMethod(ctx context.Context, m *model.DeliveryMethod) error {
	if m.Name == "" {
		return errors.New("delivery method name is required")
	}
	if m.BaseCost < 0 {
		return errors.New("base cost cannot be negative")
	}
	return s.repo.CreateDeliveryMethod(ctx, m)
}

// RegisterDeliveryMan converts an existing customer into a platform rider
func (s *DeliveryService) RegisterDeliveryMan(ctx context.Context, dm *model.DeliveryMan) error {
	if dm.CustomerID <= 0 {
		return errors.New("valid customer ID is required to register as delivery man")
	}

	// Check if already registered
	existing, err := s.repo.GetDeliveryManByCustomerID(ctx, dm.CustomerID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("customer is already registered as a delivery man")
	}

	return s.repo.CreateDeliveryMan(ctx, dm)
}

// AssignDelivery routes an order to a delivery man.
func (s *DeliveryService) AssignDelivery(ctx context.Context, assignment *model.OrderDelivery) error {
	if assignment.OrderID <= 0 || !assignment.DeliveryManID.Valid {
		return errors.New("order ID and delivery man ID are required")
	}

	assignment.DeliveryStatus = "assigned"
	return s.repo.AssignOrderToDelivery(ctx, assignment)
}

// UpdateDeliveryStatus allows a rider to update the status of their delivery task.
// Features logic for crediting the delivery man's wallet when order status changes to 'Delivered'.
func (s *DeliveryService) UpdateDeliveryStatus(ctx context.Context, orderID int64, newStatus string) error {
	if orderID <= 0 {
		return errors.New("valid order ID required")
	}

	validStatuses := map[string]bool{
		"accepted":         true,
		"out_for_delivery": true,
		"delivered":        true,
		"failed":           true,
		"cancelled":        true,
	}
	if !validStatuses[newStatus] {
		return errors.New("invalid delivery status")
	}

	payload := &model.OrderDelivery{DeliveryStatus: newStatus}
	err := s.repo.UpdateOrderDeliveryStatus(ctx, orderID, payload)
	if err != nil {
		return err
	}

	// Logic: Credit the delivery man's wallet when an order status changes to 'Delivered'
	if newStatus == "delivered" && payload.DeliveryManID.Valid && payload.DeliveryManEarning > 0 {
		err := s.repo.CreditWallet(ctx, payload.DeliveryManID.Int64, payload.DeliveryManEarning)
		if err != nil {
			// In a real production system, this should likely be in a DB Transaction
			return errors.New("delivery successful, but failed to credit wallet: " + err.Error())
		}
	}

	return nil
}

// RequestWithdrawal initiates a withdrawal request against available balance.
func (s *DeliveryService) RequestWithdrawal(ctx context.Context, wr *model.WithdrawRequest) error {
	if wr.DeliveryManID <= 0 {
		return errors.New("delivery man ID required")
	}
	if wr.Amount <= 0 {
		return errors.New("withdrawal amount must be greater than zero")
	}

	return s.repo.CreateWithdrawRequest(ctx, wr)
}



func (s *DeliveryService) ListDeliveryMen(ctx context.Context) ([]model.DeliveryMan, error) {
        return s.repo.ListDeliveryMen(ctx)
}

func (s *DeliveryService) ListDeliveryMethods(ctx context.Context) ([]model.DeliveryMethod, error) {
        return s.repo.ListDeliveryMethods(ctx)
}

func (s *DeliveryService) GetDeliveryHistory(ctx context.Context, page, limit int) ([]repository.OrderDeliveryHistory, error) {
        if page < 1 {
                page = 1
        }
        if limit < 1 || limit > 100 {
                limit = 20
        }
        offset := (page - 1) * limit
        return s.repo.GetDeliveryHistory(ctx, limit, offset)
}
