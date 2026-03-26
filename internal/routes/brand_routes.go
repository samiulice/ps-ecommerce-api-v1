package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func brandRoutes(h *handler.BrandHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/list", h.GetBrands)

	// Level 1
	mux.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("brand.create")).Post("/new", h.Create)
	mux.Get("/get/{id}", h.GetByID)
	mux.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("brand.edit")).Put("/update/{id}", h.Update)
	mux.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("brand.delete")).Delete("/delete/{id}", h.Delete)
	return mux
}
