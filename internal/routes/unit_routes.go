package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

func UnitRoutes(h *handler.UnitHandler) *chi.Mux {
	mux := chi.NewRouter()

	mux.Post("/new", h.Create)
	mux.Get("/list", h.GetAll)
	mux.Get("/get/{id}", h.GetByID)
	mux.Put("/update/{id}", h.Update)
	mux.Delete("/delete/{id}", h.Delete)

	return mux
}