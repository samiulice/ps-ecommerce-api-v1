package handler

import (
	"github.com/projuktisheba/pse-api-v1/internal/service"
)

type HandlerRepository struct {
	AuthHandler         *AuthHandler
	CategoryHandler     *CategoryHandler
	ProductHandler      *ProductHandler
	CustomerHandler     *CustomerHandler
	OrderHandler        *OrderHandler
	SiteSettingsHandler *SiteSettingsHandler
	BranchHandler       *BranchHandler
	UnitHandler         *UnitHandler
	AttributeHandler    *AttributeHandler
}

func NewHandlerRepository(svc *service.ServiceRepository) *HandlerRepository {
	return &HandlerRepository{
		AuthHandler:         NewAuthHandler(svc.AuthService),
		CategoryHandler:     NewCategoryHandler(svc.CategoryService),
		ProductHandler:      NewProductHandler(svc.ProductService),
		CustomerHandler:     NewCustomerHandler(svc.CustomerService),
		OrderHandler:        NewOrderHandler(svc.OrderService),
		SiteSettingsHandler: NewSiteSettingsHandler(svc.SiteSettingsService),
		BranchHandler:       NewBranchHandler(svc.BranchService),
		UnitHandler:         NewUnitHandler(svc.UnitService),
		AttributeHandler: NewAttributeHandler(svc.AttributeService),
	}
}
