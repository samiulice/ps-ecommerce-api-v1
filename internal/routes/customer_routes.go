// Package routes defines HTTP routes for customer endpoints.
package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

// customerRoutes registers all customer-related routes.
func customerRoutes(h *handler.CustomerHandler) *chi.Mux {
	mux := chi.NewRouter()

	// POST /customers - Create a new customer
	mux.Post("/", h.Create)

	// GET /customers/{id} - Get customer by ID
	mux.Get("/{id}", h.GetByID)

	// PUT /customers/{id} - Update customer by ID
	mux.Put("/{id}", h.Update)

	// DELETE /customers/{id} - Delete customer by ID
	mux.Delete("/{id}", h.Delete)

	return mux
}
