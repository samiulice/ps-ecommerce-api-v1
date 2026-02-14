package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

func branchRoutes(h *handler.BranchHandler) *chi.Mux {
	mux := chi.NewRouter()

	// List all branches
	mux.Get("/list", h.GetBranches)

	// Single branch operations
	mux.Post("/new", h.Create)
	mux.Get("/get/{id}", h.GetByID)
	mux.Put("/update/{id}", h.Update)
	mux.Delete("/delete/{id}", h.Delete)

	return mux
}