package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func supplierRoutes(h *handler.SupplierHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(employeeAuth(secretKey))
	mux.Use(middleware.RequireEmployee)

	mux.With(middleware.RequirePermission("supplier.view")).Get("/list", h.ListSuppliers)
	mux.With(middleware.RequirePermission("supplier.view")).Get("/profile/{id}", h.GetByID)
	mux.With(middleware.RequirePermission("supplier.create")).Post("/new", h.Create)
	mux.With(middleware.RequirePermission("supplier.edit")).Put("/update/{id}", h.Update)
	mux.With(middleware.RequirePermission("supplier.edit")).Put("/update/account/status/{id}", h.UpdateAccountStatus)
	mux.With(middleware.RequirePermission("supplier.delete")).Delete("/delete/{id}", h.Delete)

	return mux
}
