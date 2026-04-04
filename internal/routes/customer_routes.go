// Package routes defines HTTP routes for customer endpoints.
package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

// customerRoutes registers all customer-related routes.
func customerRoutes(h *handler.CustomerHandler) *chi.Mux {
	mux := chi.NewRouter()

	// GET /customers/suggestions - lightweight customer search for POS autocomplete
	mux.Get("/suggestions", h.SuggestCustomers)

	// GET /customers/list - get all customers
	mux.Get("/list", h.ListCustomers)

	// POST /customers/new - Create a new customer
	mux.Post("/new", h.Create)

	// GET /customers/profile/{id} - Get customer by ID
	mux.Get("/profile/{id}", h.GetByID)

	// PUT /customers/update/{id} - Update customer by ID
	mux.Put("/update/{id}", h.Update)

	// PUT /customers/update/account/status/{id} - Update customer by ID, query parameter {is_active=true}
	mux.Put("/update/account/status/{id}", h.UpdateAccountStatus)

	// DELETE /customers/delete/{id} - Delete customer by ID
	mux.Delete("/delete/{id}", h.Delete)

	return mux
}
