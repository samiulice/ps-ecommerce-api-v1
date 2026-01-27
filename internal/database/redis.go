package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(dsn string) *redis.Client {
	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	// Verify Redis connectivity
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Redis ping failed:", err)
	}

	return client
}
