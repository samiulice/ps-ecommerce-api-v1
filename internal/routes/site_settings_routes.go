package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

func siteSettingsRoutes(h *handler.SiteSettingsHandler) *chi.Mux {
	mux := chi.NewRouter()

	// Mount path: /api/v1/site-settings (or similar)

	mux.Route("/hero", func(r chi.Router) {
		r.Get("/", h.GetHeroSection)
		r.Post("/update", h.UpdateHeroSection)
	})

	return mux
}