package database

import (
    "context"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

func ConnectRedis(dsn, password string) *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:	  dsn,      // e.g., "203.161.48.179:6379"
        Password: password, // your "u3FJR..." password
        DB:       0,        // default DB
        // Add a timeout for remote connections
        DialTimeout: 5 * time.Second, 
    })

    // Use a context with a timeout for the initial Ping
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Verify Redis connectivity
    if err := client.Ping(ctx).Err(); err != nil {
        // Log the error specifically to help debugging
        log.Fatalf("Could not connect to Redis at %s: %v", dsn, err)
    }

    log.Printf("Successfully connected to Redis at %s", dsn)
    return client
}