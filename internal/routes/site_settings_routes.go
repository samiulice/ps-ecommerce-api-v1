package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func siteSettingsRoutes(h *handler.SiteSettingsHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()

	mux.Route("/general", func(r chi.Router) {
		r.Get("/", h.GetGeneralSettings)
		r.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("settings.edit")).Post("/update", h.UpdateGeneralSettings)
	})

	mux.Route("/hero", func(r chi.Router) {
		r.Get("/", h.GetHeroSection)
		r.With(employeeAuth(secretKey), middleware.RequireEmployee, middleware.RequirePermission("settings.edit")).Post("/update", h.UpdateHeroSection)
	})

	mux.Route("/social-links", func(r chi.Router) {
		r.Get("/topbar", h.GetTopbarSocialLinks)

		r.Group(func(admin chi.Router) {
			admin.Use(employeeAuth(secretKey))
			admin.Use(middleware.RequireEmployee)

			admin.With(middleware.RequirePermission("settings.view")).Get("/", h.ListSocialLinks)
			admin.With(middleware.RequirePermission("settings.view")).Get("/{id}", h.GetSocialLinkByID)
			admin.With(middleware.RequirePermission("settings.edit")).Post("/new", h.CreateSocialLink)
			admin.With(middleware.RequirePermission("settings.edit")).Put("/update/{id}", h.UpdateSocialLink)
			admin.With(middleware.RequirePermission("settings.edit")).Delete("/delete/{id}", h.DeleteSocialLink)
		})
	})

	return mux
}
