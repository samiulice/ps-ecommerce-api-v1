package handler

import (
	"net/http"
	"strings"

	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) getPeriod(r *http.Request) string {
	period := strings.ToLower(r.URL.Query().Get("period"))
	if period != "today" && period != "weekly" && period != "monthly" && period != "yearly" {
		return "today"
	}
	return period
}

func (h *DashboardHandler) GetStatsCards(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetStatsCards(r.Context(), h.getPeriod(r))
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}

func (h *DashboardHandler) GetSaleComparison(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetSaleComparison(r.Context(), h.getPeriod(r))
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}

func (h *DashboardHandler) GetOperationalExpenseGraph(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetOperationalExpenseGraph(r.Context(), h.getPeriod(r))
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}

func (h *DashboardHandler) GetNetCashFlowGraph(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetNetCashFlowGraph(r.Context(), h.getPeriod(r))
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}

func (h *DashboardHandler) GetPopularProducts(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetPopularProducts(r.Context())
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}

func (h *DashboardHandler) GetLowStockProducts(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetLowStockProducts(r.Context())
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}
