package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

// authRoutes configures all authentication routes for both admin and customer.
func authRoutes(authHandler *handler.AuthHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()

	// ==================== ADMIN (Employee) AUTH ====================
	mux.Route("/admin", func(r chi.Router) {
		// Public admin routes
		r.Post("/register", authHandler.EmployeeRegister)
		r.Post("/login", authHandler.EmployeeLogin)
		r.Post("/refresh", authHandler.EmployeeRefresh)

		// Protected admin routes
		r.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTAuth(secretKey))
			protected.Get("/me", authHandler.EmployeeMe)
		})
	})

	// ==================== CUSTOMER (Customer) AUTH ====================
	mux.Route("/customer", func(r chi.Router) {
		// Public customer routes
		r.Post("/register", authHandler.CustomerRegister)
		r.Post("/login", authHandler.CustomerLogin)
		r.Post("/refresh", authHandler.CustomerRefresh)

		// Protected customer routes
		r.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTAuth(secretKey))
			protected.Get("/me", authHandler.CustomerMe)
		})
	})

	return mux
}
