package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

// DeliveryHandler handles HTTP requests for delivery module.
type DeliveryHandler struct {
	svc *service.DeliveryService
}

// NewDeliveryHandler creates a new DeliveryHandler.
func NewDeliveryHandler(svc *service.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{svc: svc}
}

// AddDeliveryMethod handles POST /delivery-methods
func (h *DeliveryHandler) AddDeliveryMethod(w http.ResponseWriter, r *http.Request) {
	var m model.DeliveryMethod
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if err := h.svc.AddDeliveryMethod(r.Context(), &m); err != nil {
		utils.ServerError(w, err)
		return
	}

	response := struct {
		Error          bool                  `json:"error"`
		Message        string                `json:"message"`
		DeliveryMethod *model.DeliveryMethod `json:"delivery_method"`
	}{
		Error:          false,
		Message:        "Delivery method created successfully",
		DeliveryMethod: &m,
	}
	utils.WriteJSON(w, http.StatusCreated, response)
}

// RegisterDeliveryMan handles POST /delivery-men
func (h *DeliveryHandler) RegisterDeliveryMan(w http.ResponseWriter, r *http.Request) {
	var m model.DeliveryMan
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if err := h.svc.RegisterDeliveryMan(r.Context(), &m); err != nil {
		utils.BadRequest(w, err) // using BadRequest for validation issues
		return
	}

	response := struct {
		Error       bool               `json:"error"`
		Message     string             `json:"message"`
		DeliveryMan *model.DeliveryMan `json:"delivery_man"`
	}{
		Error:       false,
		Message:     "Delivery man registered successfully",
		DeliveryMan: &m,
	}
	utils.WriteJSON(w, http.StatusCreated, response)
}

// AssignDelivery handles POST /orders/{id}/assign-delivery
func (h *DeliveryHandler) AssignDelivery(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	var d model.OrderDelivery
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		utils.BadRequest(w, err)
		return
	}
	d.OrderID = orderID

	if err := h.svc.AssignDelivery(r.Context(), &d); err != nil {
		utils.ServerError(w, err)
		return
	}

	response := struct {
		Error   bool                 `json:"error"`
		Message string               `json:"message"`
		Details *model.OrderDelivery `json:"assignment_details"`
	}{
		Error:   false,
		Message: "Order assigned to delivery man successfully",
		Details: &d,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

// UpdateDeliveryStatus handles PUT /orders/{id}/delivery-status
// It allows riders to mark orders as delivered, processing wallet credit.
func (h *DeliveryHandler) UpdateDeliveryStatus(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if err := h.svc.UpdateDeliveryStatus(r.Context(), orderID, payload.Status); err != nil {
		utils.ServerError(w, err)
		return
	}

	response := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Delivery status updated successfully",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

// RequestWithdrawal handles POST /delivery-portal/withdraw
func (h *DeliveryHandler) RequestWithdrawal(w http.ResponseWriter, r *http.Request) {
	var req model.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if err := h.svc.RequestWithdrawal(r.Context(), &req); err != nil {
		utils.ServerError(w, err)
		return
	}

	response := struct {
		Error   bool                   `json:"error"`
		Message string                 `json:"message"`
		Request *model.WithdrawRequest `json:"request"`
	}{
		Error:   false,
		Message: "Withdrawal request submitted successfully",
		Request: &req,
	}
	utils.WriteJSON(w, http.StatusCreated, response)
}

// ListDeliveryMen handles GET /delivery/men
func (h *DeliveryHandler) ListDeliveryMen(w http.ResponseWriter, r *http.Request) {
        men, err := h.svc.ListDeliveryMen(r.Context())
        if err != nil {
                utils.ServerError(w, err)
                return
        }
        utils.OK(w, "Delivery men retrieved", men)
}

// GetDeliveryHistory handles GET /delivery/history
func (h *DeliveryHandler) GetDeliveryHistory(w http.ResponseWriter, r *http.Request) {
        pageStr := r.URL.Query().Get("page")
        limitStr := r.URL.Query().Get("limit")

        page, _ := strconv.Atoi(pageStr)
        limit, _ := strconv.Atoi(limitStr)

        if page <= 0 {
                page = 1
        }
        if limit <= 0 {
                limit = 20
        }

        history, err := h.svc.GetDeliveryHistory(r.Context(), page, limit)
        if err != nil {
                utils.ServerError(w, err)
                return
        }
        utils.OK(w, "Delivery history retrieved", history)
}
