package service

import (
	"github.com/projuktisheba/pse-api-v1/internal/config"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"github.com/redis/go-redis/v9"
)

// ServiceRepository contains all individual service
type ServiceRepository struct {
	AuthService         *AuthService
	RoleService         *RoleService
	EmployeeService     *EmployeeService
	CategoryService     *CategoryService
	BrandService        *BrandService
	ProductService      *ProductService
	CustomerService     *CustomerService
	SupplierService     *SupplierService
	PurchaseService     *PurchaseService
	OrderService        *OrderService
	SiteSettingsService *SiteSettingsService
	BranchService       *BranchService
	UnitService         *UnitService
	AttributeService    *AttributeService
}

// NewServiceRepository initializes all repositories with a shared connection pool
func NewServiceRepository(dbrepo *repository.DBRepository, rdb *redis.Client, config *config.Config) *ServiceRepository {
	return &ServiceRepository{
		AuthService:         NewAuthService(dbrepo.EmployeeRepository, dbrepo.RoleRepo, dbrepo.CustomerRepository, dbrepo.RedisTokenRepository, config.JWT.Access.SecretKey),
		RoleService:         NewRoleService(dbrepo.RoleRepo),
		EmployeeService:     NewEmployeeService(dbrepo.EmployeeRepository, dbrepo.RoleRepo),
		CategoryService:     NewCategoryService(dbrepo.CategoryRepo),
		ProductService:      NewProductService(dbrepo.ProductRepo),
		BrandService:        NewBrandService(dbrepo.BrandRepo),
		CustomerService:     NewCustomerService(dbrepo.CustomerRepository),
		SupplierService:     NewSupplierService(dbrepo.SupplierRepo),
		PurchaseService:     NewPurchaseService(dbrepo.PurchaseRepo),
		OrderService:        NewOrderService(dbrepo.OrderRepo, dbrepo.CustomerRepository),
		SiteSettingsService: NewSiteSettingsService(dbrepo.SiteSettingsRepo),
		BranchService:       NewBranchService(dbrepo.BranchRepo),
		UnitService:         NewUnitService(dbrepo.UnitRepo),
		AttributeService:    NewAttributeService(dbrepo.AttributeRepo),
	}
}
