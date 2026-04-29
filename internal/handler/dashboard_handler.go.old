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

func (h *DashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	period := strings.ToLower(r.URL.Query().Get("period"))
	if period != "weekly" && period != "monthly" &&  period != "yearly"{
		period = "weekly"
	}
	Dashboard, err := h.svc.GetDashboardStats(r.Context(), period)
	if err != nil {
		utils.NotFound(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, Dashboard)
}
