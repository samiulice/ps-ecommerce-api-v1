// Package repository contains Redis-backed token persistence.
package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisTokenRepo implements TokenRepository using Redis.
type RedisTokenRepo struct {
	rdb *redis.Client
}

// NewRedisTokenRepo creates a new RedisTokenRepo.
func NewRedisTokenRepo(rdb *redis.Client) *RedisTokenRepo {
	return &RedisTokenRepo{rdb: rdb}
}

// Save stores a refresh token with expiration.
func (r *RedisTokenRepo) Save(ctx context.Context, token string, userID int, ttl time.Duration) error {
	return r.rdb.Set(ctx, token, userID, ttl).Err()
}

// Get retrieves the user ID associated with a refresh token.
func (r *RedisTokenRepo) Get(ctx context.Context, token string) (string, error) {
	return r.rdb.Get(ctx, token).Result()
}

// Delete revokes a refresh token.
func (r *RedisTokenRepo) Delete(ctx context.Context, token string) error {
	return r.rdb.Del(ctx, token).Err()
}
