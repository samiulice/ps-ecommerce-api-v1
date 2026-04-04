package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// DBRepository contains all individual repository
type DBRepository struct {
	EmployeeRepository   *EmployeeRepository
	RoleRepo             *RoleRepository
	CustomerRepository   *CustomerRepository
	SupplierRepo         *SupplierRepo
	PurchaseRepo         *PurchaseRepo
	RedisTokenRepository *RedisTokenRepo
	CategoryRepo         *CategoryRepo
	BrandRepo            *BrandRepo
	ProductRepo          *ProductRepo
	OrderRepo            *OrderRepo
	SiteSettingsRepo     *SiteSettingsRepo
	BranchRepo           *BranchRepo
	UnitRepo             *UnitRepo
	AttributeRepo        *AttributeRepo
	POSRepo              *POSRepo
	ReportRepo           *ReportRepo
	DeliveryRepo         *DeliveryRepository
}

// NewDBRepository initializes all repositories with a shared connection pool
func NewDBRepository(db *pgxpool.Pool, rdb *redis.Client) *DBRepository {
	return &DBRepository{
		EmployeeRepository:   NewEmployeeRepo(db),
		RoleRepo:             NewRoleRepo(db),
		CustomerRepository:   NewCustomerRepo(db),
		SupplierRepo:         NewSupplierRepo(db),
		PurchaseRepo:         NewPurchaseRepo(db),
		RedisTokenRepository: NewRedisTokenRepo(rdb),
		CategoryRepo:         NewCategoryRepo(db),
		BrandRepo:            NewBrandRepo(db),
		ProductRepo:          NewProductRepo(db),
		OrderRepo:            NewOrderRepo(db),
		SiteSettingsRepo:     NewSiteSettingsRepo(db),
		BranchRepo:           NewBranchRepo(db),
		UnitRepo:             NewUnitRepo(db),
		AttributeRepo:        NewAttributeRepo(db),
		POSRepo:              NewPOSRepo(db),
		ReportRepo:           NewReportRepo(db),
		DeliveryRepo:         NewDeliveryRepo(db),
	}
}
