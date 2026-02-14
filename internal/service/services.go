package service

import (
	"github.com/projuktisheba/pse-api-v1/internal/config"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"github.com/redis/go-redis/v9"
)

// ServiceRepository contains all individual service
type ServiceRepository struct {
	AuthService     *AuthService
	CategoryService *CategoryService
	ProductService  *ProductService
	CustomerService *CustomerService
	OrderService    *OrderService
	SiteSettingsService *SiteSettingsService
	BranchService *BranchService
}

// NewServiceRepository initializes all repositories with a shared connection pool
func NewServiceRepository(dbrepo *repository.DBRepository, rdb *redis.Client, config *config.Config) *ServiceRepository {
	return &ServiceRepository{
		AuthService:     NewAuthService(dbrepo.EmployeeRepository, dbrepo.CustomerRepository, dbrepo.RedisTokenRepository, config.JWT.Access.SecretKey),
		CategoryService: NewCategoryService(dbrepo.CategoryRepo),
		ProductService:  NewProductService(dbrepo.ProductRepo),
		CustomerService: NewCustomerService(dbrepo.CustomerRepository),
		OrderService:    NewOrderService(dbrepo.OrderRepo, dbrepo.CustomerRepository),
		SiteSettingsService: NewSiteSettingsService(dbrepo.SiteSettingsRepo),
		BranchService: NewBranchService(dbrepo.BranchRepo),
	}
}
