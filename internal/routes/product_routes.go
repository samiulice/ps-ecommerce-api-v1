package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

// root path: /products
func productRoutes(productHandler *handler.ProductHandler) *chi.Mux {
	mux := chi.NewRouter()

	mux.Post("/", productHandler.Create)
	mux.Get("/{id}", productHandler.GetByID)
	mux.Put("/{id}", productHandler.Update) 
	mux.Delete("/{id}", productHandler.Delete)

	return mux
}
