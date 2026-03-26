package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func purchaseRoutes(h *handler.PurchaseHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(employeeAuth(secretKey))
	mux.Use(middleware.RequireEmployee)

	mux.With(middleware.RequirePermission("purchase.view")).Get("/list", h.ListPurchases)
	mux.With(middleware.RequirePermission("purchase.view")).Get("/profile/{id}", h.GetByID)
	mux.With(middleware.RequirePermission("purchase.create")).Post("/new", h.Create)
	mux.With(middleware.RequirePermission("purchase.edit")).Put("/update/{id}", h.Update)
	mux.With(middleware.RequirePermission("purchase.delete")).Delete("/delete/{id}", h.Delete)

	return mux
}
