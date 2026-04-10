package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

// deliveryRoutes initializes routes related to delivery management.
func deliveryRoutes(h *handler.DeliveryHandler, secretKey string) *chi.Mux {
	r := chi.NewRouter()

	// Apply authentication middleware for routes that require an authenticated user
	r.Use(middleware.JWTAuth(secretKey))

	// Routes that require admin or special permission
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("delivery.manage")) // Assume permission requirement

		// Admin delivery methods and assignment management
		r.Post("/methods", h.AddDeliveryMethod)
		r.Post("/men", h.RegisterDeliveryMan)
		r.Get("/men", h.ListDeliveryMen)
		r.Get("/history", h.GetDeliveryHistory)
	})

	// Platform delivery assignments endpoint via order sub-routes?
	// It's cleaner to mount these on orders routes typically, but we map here as requested
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("delivery.assign"))
		r.Post("/orders/{id}/assign", h.AssignDelivery)
	})

	// Portal Routes (For the delivery driver)
	r.Group(func(r chi.Router) {
		// Update status as being delivered
		r.Put("/portal/orders/{id}/status", h.UpdateDeliveryStatus)
 	// Get assigned orders
		r.Get("/portal/orders", h.GetPortalOrders)
		// Get portal wallet
		r.Get("/portal/wallet", h.GetPortalWallet)

		// Request withdrawal
		r.Post("/portal/withdraw", h.RequestWithdrawal)
	})

	return r
}