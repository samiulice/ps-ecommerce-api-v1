package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

// Assuming hRepo, cfg, and errorLog are available to this function.
func authRoutes(authHandler *handler.AuthHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()
	// ---- Public Auth Routes ----
	mux.Post("/register", authHandler.Register)
	mux.Post("/login", authHandler.Login)
	mux.Post("/refresh", authHandler.Refresh)

	// ---- Protected Routes ----
	mux.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(secretKey))
		r.Get("/me", authHandler.Me)
	})

	return mux
}
