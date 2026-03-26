package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func orderRoutes(h *handler.OrderHandler, secretKey string) chi.Router {
	r := chi.NewRouter()

	// Public routes (no auth required for placing orders)
	r.Post("/new", h.PlaceOrder) // POST /api/v1/orders/new

	// Public order tracking
	r.Get("/track/{orderNumber}", h.GetOrderByNumber) // GET /api/v1/orders/track/{orderNumber}

	// Protected routes (require authentication)
	// These can be wrapped with auth middleware in the main routes file

	r.Group(func(admin chi.Router) {
		admin.Use(employeeAuth(secretKey))
		admin.Use(middleware.RequireEmployee)

		// Order listing and stats (admin)
		admin.With(middleware.RequirePermission("order.view")).Get("/", h.ListOrders)
		admin.With(middleware.RequirePermission("order.view")).Get("/stats", h.GetOrderStats)

		// Single order operations
		admin.With(middleware.RequirePermission("order.view")).Get("/{id}", h.GetOrder)
		admin.With(middleware.RequirePermission("order.edit")).Put("/{id}/status", h.UpdateOrderStatus)
		admin.With(middleware.RequirePermission("order.edit")).Put("/{id}/payment-status", h.UpdatePaymentStatus)
		admin.With(middleware.RequirePermission("order.delete")).Delete("/{id}", h.DeleteOrder)

		// Customer orders
		admin.With(middleware.RequirePermission("order.view")).Get("/customer/{customerId}", h.GetCustomerOrders)
	})

	return r
}
