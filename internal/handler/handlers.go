package handler

import (
	"github.com/projuktisheba/pse-api-v1/internal/service"
)

type HandlerRepository struct {
	AuthHandler         *AuthHandler
	RoleHandler         *RoleHandler
	EmployeeHandler     *EmployeeHandler
	CategoryHandler     *CategoryHandler
	BrandHandler        *BrandHandler
	ProductHandler      *ProductHandler
	CustomerHandler     *CustomerHandler
	SupplierHandler     *SupplierHandler
	PurchaseHandler     *PurchaseHandler
	OrderHandler        *OrderHandler
	SiteSettingsHandler *SiteSettingsHandler
	BranchHandler       *BranchHandler
	UnitHandler         *UnitHandler
	AttributeHandler    *AttributeHandler
	POSHandler          *POSHandler
	ReportHandler       *ReportHandler
}

func NewHandlerRepository(svc *service.ServiceRepository) *HandlerRepository {
	return &HandlerRepository{
		AuthHandler:         NewAuthHandler(svc.AuthService),
		RoleHandler:         NewRoleHandler(svc.RoleService),
		EmployeeHandler:     NewEmployeeHandler(svc.EmployeeService),
		CategoryHandler:     NewCategoryHandler(svc.CategoryService),
		BrandHandler:        NewBrandHandler(svc.BrandService),
		ProductHandler:      NewProductHandler(svc.ProductService),
		CustomerHandler:     NewCustomerHandler(svc.CustomerService),
		SupplierHandler:     NewSupplierHandler(svc.SupplierService),
		PurchaseHandler:     NewPurchaseHandler(svc.PurchaseService),
		OrderHandler:        NewOrderHandler(svc.OrderService),
		SiteSettingsHandler: NewSiteSettingsHandler(svc.SiteSettingsService),
		BranchHandler:       NewBranchHandler(svc.BranchService),
		UnitHandler:         NewUnitHandler(svc.UnitService),
		AttributeHandler:    NewAttributeHandler(svc.AttributeService),
		POSHandler:          NewPOSHandler(svc.POSService),
		ReportHandler:       NewReportHandler(svc.ReportService),
	}
}
