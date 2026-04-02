package routes

import (
"github.com/go-chi/chi/v5"
"github.com/projuktisheba/pse-api-v1/internal/handler"
"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func reportRoutes(h *handler.ReportHandler, jwtSecret string) *chi.Mux {
r := chi.NewRouter()
r.Use(employeeAuth(jwtSecret))
r.Use(middleware.RequireEmployee)
r.Use(middleware.RequirePermission("report.view"))

r.Get("/pos-sales", h.GetPOSSalesReport)
r.Get("/orders", h.GetOrdersReport)
r.Get("/customer-dues", h.GetCustomerDueReport)
r.Get("/supplier-dues", h.GetSupplierDueReport)
r.Get("/low-stock", h.GetLowStockReport)

return r
}
