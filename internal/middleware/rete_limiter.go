package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

// RateLimiterConfig holds settings for rate limiting.
type RateLimiterConfig struct {
	Requests int           // Max requests per duration
	Duration time.Duration // Window duration
	Redis    *redis.Client // Redis client
}

// RateLimiterMiddleware returns a Chi middleware for per-IP or per-user limiting.
// If a userID is in context, limits per user. Otherwise, limits per IP.
func RateLimiterMiddleware(cfg RateLimiterConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var key string
			// If userID is set in context (from AuthMiddleware)
			if uid := ctx.Value("userID"); uid != nil {
				key = "rate:user:" + strconv.Itoa(uid.(int))
			} else {
				ip := clientIP(r)
				key = "rate:ip:" + ip
			}

			// Increment request counter
			count, err := cfg.Redis.Incr(ctx, key).Result()
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Set expiration for first request
			if count == 1 {
				cfg.Redis.Expire(ctx, key, cfg.Duration)
			}

			if count > int64(cfg.Requests) {
				// Optional: set Retry-After header
				ttl, _ := cfg.Redis.TTL(ctx, key).Result()
				w.Header().Set("Retry-After", strconv.Itoa(int(ttl.Seconds())))
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Helper: attach to router with optional per-route limits
func AttachRateLimiter(r chi.Router, cfg RateLimiterConfig) {
	r.Use(RateLimiterMiddleware(cfg))
}

// clientIP extracts the real client IP address from an HTTP request.
// It checks common reverse-proxy headers first (e.g., Nginx, Cloudflare),
// and falls back to r.RemoteAddr if no headers are present.
func clientIP(r *http.Request) string {
	// Check X-Forwarded-For header (may contain multiple IPs: client, proxy1, proxy2)
	// Format: "clientIP, proxy1, proxy2"
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// Take the first IP which is the original client
		return strings.Split(ip, ",")[0]
	}

	// Check X-Real-IP header (used by some proxies like Nginx)
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fallback to RemoteAddr (format: "IP:PORT")
	// Example: "192.168.1.10:52341"
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If parsing fails, return the whole RemoteAddr as a last resort
		return r.RemoteAddr
	}

	return host
}
