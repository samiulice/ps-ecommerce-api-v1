package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func roleRoutes(h *handler.RoleHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.JWTAuth(secretKey))
	mux.Use(middleware.RequireEmployee)

	mux.With(middleware.RequirePermission("role.view")).Get("/list", h.List)
	mux.With(middleware.RequirePermission("role.view")).Get("/permissions", h.ListPermissions)
	mux.With(middleware.RequirePermission("role.view")).Get("/get/{id}", h.GetByID)
	mux.With(middleware.RequirePermission("role.create")).Post("/new", h.Create)
	mux.With(middleware.RequirePermission("role.edit")).Put("/update/{id}", h.Update)
	mux.With(middleware.RequirePermission("role.delete")).Delete("/delete/{id}", h.Delete)

	return mux
}
