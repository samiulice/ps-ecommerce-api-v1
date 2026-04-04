package service

import (
	"context"
	"fmt"
	"time"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

type POSService struct {
	posRepo *repository.POSRepo
}

func NewPOSService(posRepo *repository.POSRepo) *POSService {
	return &POSService{posRepo: posRepo}
}
func (s *POSService) GetPOSSaleByReference(ctx context.Context, referenceNo string) (*model.POSSale, error) {
	return s.posRepo.GetSaleByReference(ctx, referenceNo)
}
func (s *POSService) CreatePOSSale(ctx context.Context, req model.CreatePOSSaleRequest) (*model.POSSale, error) {
	sale := &model.POSSale{
		ReferenceNo:   fmt.Sprintf("POS-%d", time.Now().Unix()),
		CustomerID:    req.CustomerID,
		BranchID:      req.BranchID,
		SaleType:      req.SaleType,
		Discount:      req.Discount,
		AmountPaid:    req.AmountPaid,
		PaymentMethod: req.PaymentMethod,
		SaleNote:      req.SaleNote,
	}

	var totalSubtotal float64
	var totalTax float64

	for _, itemReq := range req.Items {
		itemSubtotal := float64(itemReq.Quantity) * itemReq.UnitPrice
		itemTax := 0.0 // Implement tax logic if needed
		itemTotal := itemSubtotal + itemTax

		totalSubtotal += itemSubtotal
		totalTax += itemTax

		var pVariationID *int64
		if itemReq.ProductVariationID != nil {
			vID := int64(*itemReq.ProductVariationID)
			pVariationID = &vID
		}

		saleItem := model.POSSaleItem{
			ProductID:          int64(itemReq.ProductID),
			ProductVariationID: pVariationID,
			ProductName:        "Product", // Ideally fetch name from DB
			Quantity:           itemReq.Quantity,
			UnitPrice:          itemReq.UnitPrice,
			Subtotal:           itemSubtotal,
			TaxAmount:          itemTax,
			Total:              itemTotal,
		}
		sale.Items = append(sale.Items, saleItem)
	}

	sale.Subtotal = totalSubtotal
	sale.Tax = totalTax
	sale.Total = totalSubtotal + totalTax - req.Discount

	return s.posRepo.CreateSale(ctx, sale)
}
