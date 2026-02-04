package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// DBRepository contains all individual repository
type DBRepository struct {
	EmployeeRepository   *EmployeeRepository
	CustomerRepository   *CustomerRepository
	RedisTokenRepository *RedisTokenRepo
	CategoryRepo         *CategoryRepo
	ProductRepo          *ProductRepo
	OrderRepo            *OrderRepo
}

// NewDBRepository initializes all repositories with a shared connection pool
func NewDBRepository(db *pgxpool.Pool, rdb *redis.Client) *DBRepository {
	return &DBRepository{
		EmployeeRepository:   NewEmployeeRepo(db),
		CustomerRepository:   NewCustomerRepo(db),
		RedisTokenRepository: NewRedisTokenRepo(rdb),
		CategoryRepo:         NewCategoryRepo(db),
		ProductRepo:          NewProductRepo(db),
		OrderRepo:            NewOrderRepo(db),
	}
}
