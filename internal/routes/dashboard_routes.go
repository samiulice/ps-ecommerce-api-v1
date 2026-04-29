package routes
import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func dashboardRoutes(h *handler.DashboardHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(employeeAuth(secretKey), middleware.RequireEmployee)

	mux.Get("/stats/cards", h.GetStatsCards)
	mux.Get("/stats/charts/sales", h.GetSaleComparison)
	mux.Get("/stats/charts/expenses", h.GetOperationalExpenseGraph)
	mux.Get("/stats/charts/profit", h.GetNetCashFlowGraph)
	mux.Get("/stats/products/popular", h.GetPopularProducts)
	mux.Get("/stats/products/low-stock", h.GetLowStockProducts)

	return mux
}