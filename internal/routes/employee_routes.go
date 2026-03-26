package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func employeeRoutes(h *handler.EmployeeHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.JWTAuth(secretKey))
	mux.Use(middleware.RequireEmployee)

	mux.With(middleware.RequirePermission("user.view")).Get("/list", h.List)
	mux.With(middleware.RequirePermission("user.view")).Get("/get/{id}", h.GetByID)
	mux.With(middleware.RequirePermission("user.edit")).Put("/update/{id}", h.Update)
	mux.With(middleware.RequirePermission("user.delete")).Delete("/delete/{id}", h.Delete)

	return mux
}
