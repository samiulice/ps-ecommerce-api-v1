package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type EmployeeHandler struct {
	svc *service.EmployeeService
}

func NewEmployeeHandler(svc *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{svc: svc}
}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	employees, err := h.svc.List(r.Context())
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.OK(w, "Employees retrieved", map[string]any{"employees": employees})
}

func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid employee id"))
		return
	}
	employee, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.handleErr(w, err)
		return
	}
	utils.OK(w, "Employee retrieved", map[string]any{"employee": employee})
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid employee id"))
		return
	}
	var req model.EmployeeAdminUpdateRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}
	employee, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		h.handleErr(w, err)
		return
	}
	utils.OK(w, "Employee updated successfully", map[string]any{"employee": employee})
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid employee id"))
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.handleErr(w, err)
		return
	}
	utils.OK(w, "Employee deleted successfully", nil)
}

func (h *EmployeeHandler) handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrEmployeeNotFound):
		utils.NotFound(w, err)
	case errors.Is(err, service.ErrEmployeeEmail), errors.Is(err, service.ErrEmployeeRole):
		utils.BadRequest(w, err)
	default:
		utils.ServerError(w, err)
	}
}
