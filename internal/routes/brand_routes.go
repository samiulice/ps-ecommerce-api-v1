package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

func brandRoutes(h *handler.BrandHandler) *chi.Mux {
	mux := chi.NewRouter()

	

	mux.Get("/list", h.GetBrands)

	// Level 1
	mux.Post("/new", h.Create)
	mux.Get("/get/{id}", h.GetByID)
	mux.Put("/update/{id}", h.Update)
	mux.Delete("/delete/{id}", h.Delete)
	return mux
}
