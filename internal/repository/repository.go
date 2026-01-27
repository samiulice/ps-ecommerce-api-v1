package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// DBRepository contains all individual repository
type DBRepository struct {
	UserRepository       *UserRepository
	RedisTokenRepository *RedisTokenRepo
	CategoryRepo         *CategoryRepo
	ProductRepo          *ProductRepo
}

// NewDBRepository initializes all repositories with a shared connection pool
func NewDBRepository(db *pgxpool.Pool, rdb *redis.Client) *DBRepository {
	return &DBRepository{
		UserRepository:       NewUserRepo(db),
		RedisTokenRepository: NewRedisTokenRepo(rdb),
		CategoryRepo:         NewCategoryRepo(db),
		ProductRepo:          NewProductRepo(db),
	}
}
