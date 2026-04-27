package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func dashboardRoutes(h *handler.DashboardHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()
	mux.With(employeeAuth(secretKey), middleware.RequireEmployee,).Get("/stats", h.GetDashboardStats)

	return mux
}
