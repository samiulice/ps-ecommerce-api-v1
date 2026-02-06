// Package routes defines HTTP routes for customer endpoints.
package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

// customerRoutes registers all customer-related routes.
func customerRoutes(h *handler.CustomerHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()

	// Routes accessible by both customers and admin (authenticated users)
	mux.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(secretKey))

		// POST /customers - Create a new customer
		r.Post("/", h.Create)

		// GET /customers/{id} - Get customer by ID
		r.Get("/{id}", h.GetByID)
	})

	// Routes accessible by employee only
	mux.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(secretKey))
		r.Use(middleware.RequireEmployee)

		// PUT /customers/{id} - Update customer by ID
		r.Put("/{id}", h.Update)

		// DELETE /customers/{id} - Delete customer by ID
		r.Delete("/{id}", h.Delete)
	})

	return mux
}
