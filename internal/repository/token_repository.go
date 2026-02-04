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
// The value can be a user identifier string like "employee:123" or "user:456".
func (r *RedisTokenRepo) Save(ctx context.Context, token string, value string, ttl time.Duration) error {
	return r.rdb.Set(ctx, token, value, ttl).Err()
}

// Get retrieves the value associated with a refresh token.
func (r *RedisTokenRepo) Get(ctx context.Context, token string) (string, error) {
	return r.rdb.Get(ctx, token).Result()
}

// Delete revokes a refresh token.
func (r *RedisTokenRepo) Delete(ctx context.Context, token string) error {
	return r.rdb.Del(ctx, token).Err()
}
