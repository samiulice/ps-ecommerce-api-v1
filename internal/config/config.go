// Package config loads and validates application configuration.
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv" // ✅ added
)

// Load loads config from environment variables.
func LoadConfig() *Config {
	// load .env file (non-fatal if missing)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "myapp"),
			Env:     getEnv("APP_ENV", "development"),
			Debug:   getEnvBool("APP_DEBUG", true),
			Version: getEnv("APP_VERSION", "0.1.0"),
		},
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			MaxRequests:  getEnvInt("MAX_REQUESTS", 100),
			RequestWindow:  getEnvDuration("REQUEST_WINDOW", 60*time.Second),
		},
		JWT: JWTConfig{
			Access: JWTAccessConfig{
				SecretKey: getEnv("JWT_ACCESS_SECRET", ""),
				Issuer:    getEnv("JWT_ISSUER", "myapp"),
				Audience:  getEnv("JWT_AUDIENCE", "myapp-users"),
				Algorithm: getEnv("JWT_ALGO", "HS256"),
				Expiry:    getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute),
			},
			Refresh: JWTRefreshConfig{
				SecretKey: getEnv("JWT_REFRESH_SECRET", ""),
				Expiry:    getEnvDuration("JWT_REFRESH_TTL", 7*24*time.Hour),
			},
		},
		DB: DBConfig{
			PostgresURL:    getEnv("DATABASE_URL", ""),
			PostgresDevURL: getEnv("DATABASE_DEV_URL", ""),
			RedisURL:       getEnv("REDIS_URL", ""),
		},
	}

	validate(cfg)
	return cfg
}

// --- helpers ---

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

// validate ensures critical config values exist.
func validate(cfg *Config) {
	if cfg.JWT.Access.SecretKey == "" {
		log.Fatal("JWT_ACCESS_SECRET is required")
	}
	if cfg.JWT.Refresh.SecretKey == "" {
		log.Fatal("JWT_REFRESH_SECRET is required")
	}
	if cfg.DB.PostgresURL == "" && cfg.DB.PostgresDevURL == "" {
		log.Fatal("DATABASE_URL or DATABASE_DEV_URL is required")
	}
}
