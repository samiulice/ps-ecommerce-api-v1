package routes

import (
	"log"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/projuktisheba/pse-api-v1/internal/config"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
	"github.com/redis/go-redis/v9"
)

func Routes(cfg *config.Config, rdb *redis.Client, handlers *handler.HandlerRepository, infoLogger, errorLogger *log.Logger) http.Handler {
	mux := chi.NewRouter()

	// 1. Setup Global Middleware
	setupMiddlewares(mux, cfg, rdb)

	// 2. Static file serving (Images)
	setupStaticFiles(mux)

	// 3. Health Check
	setupHealthCheck(mux, cfg)

	// 4. API Routes (v1)
	mux.Route("/api/v1", func(r chi.Router) {
		// Public Routes
		r.Mount("/auth", authRoutes(handlers.AuthHandler, cfg.JWT.Access.SecretKey))
		r.Mount("/categories", categoryRoutes(handlers.CategoryHandler))
		r.Mount("/products", productRoutes(handlers.ProductHandler))
		r.Mount("/customers", customerRoutes(handlers.CustomerHandler))
		r.Mount("/orders", orderRoutes(handlers.OrderHandler))
		r.Mount("/site-settings", siteSettingsRoutes(handlers.SiteSettingsHandler))
		r.Mount("/branches", branchRoutes(handlers.BranchHandler))
		r.Mount("/units", UnitRoutes(handlers.UnitHandler))
		r.Mount("/attributes", AttributeRoutes(handlers.AttributeHandler))
	})

	return mux
}

// setupMiddleware configures CORS, logging, recovery, and rate limiting
func setupMiddlewares(mux *chi.Mux, cfg *config.Config, rdb *redis.Client) {
	// CORS
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: Change to specific domains in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Branch-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Standard Chi Middleware
	mux.Use(chimiddleware.RequestID)
	mux.Use(chimiddleware.RealIP)
	mux.Use(chimiddleware.Logger)
	mux.Use(chimiddleware.Recoverer)
	mux.Use(chimiddleware.Timeout(cfg.Server.ReadTimeout))
	log.Println("Global middleware enabled")
	// Rate Limiter
	middleware.AttachRateLimiter(mux, middleware.RateLimiterConfig{
		Requests: cfg.Server.MaxRequests,
		Duration: cfg.Server.RequestWindow,
		Redis:    rdb,
	})
	rps := float64(cfg.Server.MaxRequests) / cfg.Server.RequestWindow.Seconds()

	log.Printf(
		"rate limiter active: %d requests per %s (%.2f req/sec)",
		cfg.Server.MaxRequests,
		cfg.Server.RequestWindow,
		rps,
	)

}

// setupStaticFiles serves static files from ./assets/public
func setupStaticFiles(mux *chi.Mux) {
	imageDir := filepath.Join(".", "assets", "public")
	// Note: We use /api/v1/images/ here to match the StripPrefix logic
	fs := http.StripPrefix("/api/v1/public/", http.FileServer(http.Dir(imageDir)))
	mux.Handle("/api/v1/public/*", fs)
}

// setupHealthCheck adds a simple ping endpoint
func setupHealthCheck(mux *chi.Mux, cfg *config.Config) {
	mux.Get("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		ip := "unknown"
		if conn, err := net.Dial("udp", "1.1.1.1:80"); err == nil {
			defer conn.Close()
			ip = conn.LocalAddr().(*net.UDPAddr).IP.String()
		}

		resp := map[string]any{
			"status":    "active",
			"env":       cfg.App.Env,
			"server_ip": ip,
			"timestamp": time.Now().UTC(),
		}
		utils.WriteJSON(w, http.StatusOK, resp)
	})
}
