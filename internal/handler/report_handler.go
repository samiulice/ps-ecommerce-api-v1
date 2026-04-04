package handler

import (
	"net/http"
	"strconv"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func parseReportFilter(r *http.Request) model.ReportFilter {
	var filter model.ReportFilter
	filter.SaleType = r.URL.Query().Get("sale_type")
	filter.OrderBy = r.URL.Query().Get("order_by")
	filter.Search = r.URL.Query().Get("search")

	pageStr := r.URL.Query().Get("page")
	if page, err := strconv.Atoi(pageStr); err == nil {
		filter.Page = page
	} else {
		filter.Page = 1
	}

	limitStr := r.URL.Query().Get("limit")
	if limit, err := strconv.Atoi(limitStr); err == nil {
		filter.Limit = limit
	} else {
		filter.Limit = 20
	}

	return filter
}

func (h *ReportHandler) GetPOSSalesReport(w http.ResponseWriter, r *http.Request) {
	filter := parseReportFilter(r)
	resp, err := h.reportService.GetPOSSalesReport(r.Context(), filter)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ReportHandler) GetOrdersReport(w http.ResponseWriter, r *http.Request) {
	filter := parseReportFilter(r)
	resp, err := h.reportService.GetOrdersReport(r.Context(), filter)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ReportHandler) GetCustomerDueReport(w http.ResponseWriter, r *http.Request) {
	filter := parseReportFilter(r)
	resp, err := h.reportService.GetCustomerDueReport(r.Context(), filter)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ReportHandler) GetSupplierDueReport(w http.ResponseWriter, r *http.Request) {
	filter := parseReportFilter(r)
	resp, err := h.reportService.GetSupplierDueReport(r.Context(), filter)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ReportHandler) GetLowStockReport(w http.ResponseWriter, r *http.Request) {
	filter := parseReportFilter(r)
	resp, err := h.reportService.GetLowStockReport(r.Context(), filter)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}
