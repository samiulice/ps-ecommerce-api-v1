package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"github.com/redis/go-redis/v9"
)

type ProductService struct {
	repo  *repository.ProductRepo
	redis *redis.Client // Add Redis client here
}

func NewProductService(repo *repository.ProductRepo, rdb *redis.Client) *ProductService {
	return &ProductService{repo: repo, redis: rdb}
}

func (s *ProductService) Create(ctx context.Context, p *model.Product) error {
	if p.Name == "" || p.SKU == "" {
		return errors.New("name and SKU are required")
	}
	return s.repo.Create(ctx, p)
}

// GetByID implements Redis Caching
func (s *ProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)

	// 1. Check Redis
	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var p model.Product
		if json.Unmarshal([]byte(val), &p) == nil {
			return &p, nil // Return cached data
		}
	}

	// 2. Fetch from DB
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. Set Cache (1 Hour TTL)
	if data, err := json.Marshal(p); err == nil {
		s.redis.Set(ctx, cacheKey, data, time.Hour)
	}

	return p, nil
}

func (s *ProductService) Update(ctx context.Context, p *model.Product) error {
	if p.ID == 0 {
		return errors.New("product ID required")
	}
	err := s.repo.Update(ctx, p)
	if err == nil {
		// Invalidate cache
		s.redis.Del(ctx, fmt.Sprintf("product:%d", p.ID))
	}
	return err
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		// Invalidate cache
		s.redis.Del(ctx, fmt.Sprintf("product:%d", id))
	}
	return err
}