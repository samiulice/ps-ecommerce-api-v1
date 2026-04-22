package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/projuktisheba/pse-api-v1/internal/middleware"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type ExpenseHandler struct {
	svc *service.ExpenseService
}

func NewExpenseHandler(svc *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{svc: svc}
}

func (h *ExpenseHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var c model.ExpenseCategory
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.BadRequest(w, err)
		return
	}

	// We assume branch_id is passed in JSON if applicable.

	if err := h.svc.CreateCategory(r.Context(), &c); err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.OK(w, "Expense category created", c)
}

func (h *ExpenseHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	// The client can optionally send branch_id as query param to filter
	var bID *int64
	// parsing query parameter
	// branchIDStr := r.URL.Query().Get("branch_id")
	// if branchIDStr != "" {
	// 	id, _ := strconv.ParseInt(branchIDStr, 10, 64)
	// 	bID = &id
	// }

	cats, err := h.svc.GetCategories(r.Context(), bID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.OK(w, "Expense categories retrieved", cats)
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var e model.Expense
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if uid, ok := middleware.AuthIDFromContext(r.Context()); ok {
		e.CreatedBy = int64(uid)
	}

	if err := h.svc.CreateExpense(r.Context(), &e); err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.OK(w, "Expense recorded successfully", e)
}

func (h *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	// The client can optionally send branch_id as query param to filter
	var bID int64
	branchIDStr := r.URL.Query().Get("branch_id")
	if branchIDStr != "" {
		id, _ := strconv.ParseInt(branchIDStr, 10, 64)
		bID = id
	}

	expenses, err := h.svc.GetExpenses(r.Context(), bID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.OK(w, "Expenses retrieved", expenses)
}
