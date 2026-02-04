// Package handler implements HTTP handlers for customer endpoints.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

// CustomerHandler handles HTTP requests for customer operations.
type CustomerHandler struct {
	svc *service.CustomerService
}

// NewCustomerHandler creates a new CustomerHandler.
func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

// Create handles POST /customers - creates a new customer.
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CustomerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	customer, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error   bool                `json:"error"`
		Message string              `json:"message"`
		Customer    *model.CustomerResponse `json:"customer"`
	}{
		Error:   false,
		Message: "Customer created successfully",
		Customer:    customer.ToResponse(),
	}
	utils.WriteJSON(w, http.StatusCreated, response)
}

// GetByID handles GET /customers/{id} - retrieves a customer by ID.
func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid customer ID"))
		return
	}

	customer, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error   bool                `json:"error"`
		Message string              `json:"message"`
		Customer    *model.CustomerResponse `json:"customer"`
	}{
		Error:   false,
		Message: "Customer retrieved successfully",
		Customer:    customer.ToResponse(),
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

// Update handles PUT /customers/{id} - updates an existing customer.
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid customer ID"))
		return
	}

	var req model.CustomerUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	customer, err := h.svc.Update(r.Context(), id, &req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error   bool                `json:"error"`
		Message string              `json:"message"`
		Customer    *model.CustomerResponse `json:"customer"`
	}{
		Error:   false,
		Message: "Customer updated successfully",
		Customer:    customer.ToResponse(),
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

// Delete handles DELETE /customers/{id} - removes a customer.
func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid customer ID"))
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Customer deleted successfully",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

// handleServiceError maps service errors to HTTP responses.
func (h *CustomerHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrCustomerNotFound):
		utils.NotFound(w, err)
	case errors.Is(err, service.ErrInvalidPhone),
		errors.Is(err, service.ErrInvalidEmail),
		errors.Is(err, service.ErrInvalidPassword):
		utils.BadRequest(w, err)
	case errors.Is(err, service.ErrEmailAlreadyExists),
		errors.Is(err, service.ErrPhoneAlreadyExists):
		utils.WriteJSON(w, http.StatusConflict, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
	case errors.Is(err, service.ErrCustomerInactive):
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error":   true,
			"message": "customer account is inactive",
		})
	case errors.Is(err, service.ErrCustomerBlocked):
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error":   true,
			"message": "customer account is temporarily blocked",
		})
	default:
		utils.ServerError(w, err)
	}
}
