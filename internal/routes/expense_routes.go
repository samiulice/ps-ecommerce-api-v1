package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
)

func expenseRoutes(h *handler.ExpenseHandler, secretKey string) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(employeeAuth(secretKey)) // Make sure user is logged in
	// Ideally use middleware.RequirePermission("manage_expenses") here if you have RBAC set up

	// Expense categories
	mux.Post("/categories", h.CreateCategory)
	mux.Get("/categories", h.GetCategories)

	// Actual expenses
	mux.Post("/", h.CreateExpense)
	mux.Get("/", h.GetExpenses)

	return mux
}
