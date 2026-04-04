package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func posRoutes(h *handler.POSHandler, secretKey string) chi.Router {
	r := chi.NewRouter()

	r.Group(func(admin chi.Router) {
		admin.Use(employeeAuth(secretKey))
		admin.Use(middleware.RequireEmployee)

		// Create a pos sale
		admin.With(middleware.RequirePermission("pos.create")).Post("/sale", h.CreateSale)
                admin.With(middleware.RequirePermission("pos.view")).Get("/sale/{reference}", h.GetSaleByReference)
        })

        return r
}
