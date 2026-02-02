package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

func categoryRoutes(h *handler.CategoryHandler) *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/tree", h.GetTree) // The Full Efficient Tree

	// Level 1
	mux.Post("/new", h.Create)
	mux.Get("/get/{id}", h.GetByID)
	mux.Put("/update/{id}", h.Update)
	mux.Delete("/delete/{id}", h.Delete)

	mux.Route("/sub-categories", func(r chi.Router) {
		// Level 2
		r.Post("/", h.CreateSub)
		r.Get("/", h.GetSubByID)
		r.Put("/{id}", h.UpdateSub)
		r.Delete("/{id}", h.DeleteSub)
	})

	mux.Route("/sub-sub-categories", func(r chi.Router) {
		// Level 3
		r.Post("/", h.CreateSubSub)
		r.Get("/{id}", h.GetSubSubByID)
		r.Put("/{id}", h.UpdateSubSub)
		r.Delete("/{id}", h.DeleteSubSub)
	})
	return mux
}
