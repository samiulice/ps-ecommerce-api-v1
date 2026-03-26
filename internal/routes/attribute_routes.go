package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func AttributeRoutes(h *handler.AttributeHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(employeeAuth(secretKey))
	mux.Use(middleware.RequireEmployee)

	mux.With(middleware.RequirePermission("attribute.create")).Post("/new", h.Create)
	mux.With(middleware.RequirePermission("attribute.view")).Get("/list", h.GetAll)
	mux.With(middleware.RequirePermission("attribute.view")).Get("/get/{id}", h.GetByID)
	mux.With(middleware.RequirePermission("attribute.edit")).Put("/update/{id}", h.Update)
	mux.With(middleware.RequirePermission("attribute.delete")).Delete("/delete/{id}", h.Delete)

	return mux
}
