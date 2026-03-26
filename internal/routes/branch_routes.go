package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func branchRoutes(h *handler.BranchHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(employeeAuth(secretKey))
	mux.Use(middleware.RequireEmployee)

	// List all branches
	mux.With(middleware.RequirePermission("branch.view")).Get("/list", h.GetBranches)

	// Single branch operations
	mux.With(middleware.RequirePermission("branch.create")).Post("/new", h.Create)
	mux.With(middleware.RequirePermission("branch.view")).Get("/get/{id}", h.GetByID)
	mux.With(middleware.RequirePermission("branch.edit")).Put("/update/{id}", h.Update)
	mux.With(middleware.RequirePermission("branch.delete")).Delete("/delete/{id}", h.Delete)

	return mux
}
