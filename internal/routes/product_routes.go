package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

// root path: /products
func productRoutes(h *handler.ProductHandler) *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/list", h.GetProducts)

	mux.Post("/new", h.Create)
	mux.Get("/get/{id}", h.GetByID)
	mux.Get("/get-variations/{id}", h.GetProductVariationsByProductID)
	mux.Put("/update/{id}", h.Update)
	mux.Delete("/delete/{id}", h.Delete)
	mux.Delete("/delete-gallery-image/{id}", h.DeleteGalleryImage)

	return mux
}
