package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

func orderRoutes(h *handler.OrderHandler) chi.Router {
	r := chi.NewRouter()

	// Public routes (no auth required for placing orders)
	r.Post("/new", h.PlaceOrder) // POST /api/v1/orders/new

	// Public order tracking
	r.Get("/track/{orderNumber}", h.GetOrderByNumber) // GET /api/v1/orders/track/{orderNumber}

	// Protected routes (require authentication)
	// These can be wrapped with auth middleware in the main routes file

	// Order listing and stats (admin)
	r.Get("/", h.ListOrders)         // GET /api/v1/orders
	r.Get("/stats", h.GetOrderStats) // GET /api/v1/orders/stats

	// Single order operations
	r.Get("/{id}", h.GetOrder)                           // GET /api/v1/orders/{id}
	r.Put("/{id}/status", h.UpdateOrderStatus)           // PUT /api/v1/orders/{id}/status
	r.Put("/{id}/payment-status", h.UpdatePaymentStatus) // PUT /api/v1/orders/{id}/payment-status
	r.Delete("/{id}", h.DeleteOrder)                     // DELETE /api/v1/orders/{id}

	// Customer orders
	r.Get("/customer/{customerId}", h.GetCustomerOrders) // GET /api/v1/orders/customer/{customerId}

	return r
}
