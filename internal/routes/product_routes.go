package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

// root path: /products
func productRoutes(h *handler.ProductHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/list", h.GetProducts)

	mux.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("product.create")).Post("/new", h.Create)
	mux.Get("/get/{id}", h.GetByID)
	mux.Get("/get-variations/{id}", h.GetProductVariationsByProductID)
	mux.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("product.edit")).Put("/update/{id}", h.Update)
	mux.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("product.delete")).Delete("/delete/{id}", h.Delete)
	mux.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("product.edit")).Delete("/delete-gallery-image/{id}", h.DeleteGalleryImage)

	return mux
}
