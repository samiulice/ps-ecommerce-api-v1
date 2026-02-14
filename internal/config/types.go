package config

import "time"

// JWTAccessConfig holds access-token settings.
type JWTAccessConfig struct {
	SecretKey string        `env:"JWT_ACCESS_SECRET"`
	Issuer    string        `env:"JWT_ISSUER"`
	Audience  string        `env:"JWT_AUDIENCE"`
	Algorithm string        `env:"JWT_ALGO"`
	Expiry    time.Duration `env:"JWT_ACCESS_TTL"`
}

// JWTRefreshConfig holds refresh-token settings.
type JWTRefreshConfig struct {
	SecretKey string        `env:"JWT_REFRESH_SECRET"`
	Expiry    time.Duration `env:"JWT_REFRESH_TTL"`
}

// JWTConfig groups all JWT-related settings.
type JWTConfig struct {
	Access  JWTAccessConfig
	Refresh JWTRefreshConfig
}

// DBConfig holds database connection settings.
type DBConfig struct {
	PostgresURL    string `env:"DATABASE_URL"`
	PostgresDevURL string `env:"DATABASE_DEV_URL"`
	RedisURL       string `env:"REDIS_URL"`
}

// ServerConfig holds HTTP server tuning options.
type ServerConfig struct {
	Port          int
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
	MaxRequests   int
	RequestWindow time.Duration
}

// AppConfig holds app metadata and feature flags.
type AppConfig struct {
	Name        string
	Env         string
	Debug       bool
	Version     string
	FeatureAuth bool
}

// Config is the root configuration object.
type Config struct {
	App    AppConfig
	Server ServerConfig
	JWT    JWTConfig
	DB     DBConfig
}
