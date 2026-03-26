package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func categoryRoutes(h *handler.CategoryHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/tree", h.GetTree) // The Full Efficient Tree

	mux.Get("/list", h.GetCategories)

	// Level 1
	mux.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("category.create")).Post("/new", h.Create)
	mux.Get("/get/{id}", h.GetByID)
	mux.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("category.edit")).Put("/update/{id}", h.Update)
	mux.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("category.delete")).Delete("/delete/{id}", h.Delete)

	mux.Route("/sub-categories", func(r chi.Router) {
		// Level 2
		r.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("category.create")).Post("/new", h.CreateSub)
		r.Get("/get/{id}", h.GetSubByID)
		r.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("category.edit")).Put("/update/{id}", h.UpdateSub)
		r.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("category.delete")).Delete("/delete/{id}", h.DeleteSub)
		r.Get("/list", h.GetSubCategories)
	})

	mux.Route("/sub-sub-categories", func(r chi.Router) {
		// Level 3
		r.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("category.create")).Post("/new", h.CreateSubSub)
		r.Get("/get/{id}", h.GetSubSubByID)
		r.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("category.edit")).Put("/update/{id}", h.UpdateSubSub)
		r.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("category.delete")).Delete("/delete/{id}", h.DeleteSubSub)
		r.Get("/list", h.GetSubSubCategories)
	})
	return mux
}
