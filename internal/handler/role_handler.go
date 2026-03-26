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

type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (h *RoleHandler) List(w http.ResponseWriter, r *http.Request) {
	roles, err := h.svc.List(r.Context())
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.OK(w, "Roles retrieved", map[string]any{"roles": roles})
}

func (h *RoleHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.svc.ListPermissions(r.Context())
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.OK(w, "Permissions retrieved", map[string]any{"permissions": permissions})
}

func (h *RoleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid role id"))
		return
	}
	role, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.handleErr(w, err)
		return
	}
	utils.OK(w, "Role retrieved", map[string]any{"role": role})
}

func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.RoleCreateRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}
	role, err := h.svc.Create(r.Context(), req)
	if err != nil {
		h.handleErr(w, err)
		return
	}
	utils.Created(w, "Role created successfully", map[string]any{"role": role})
}

func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid role id"))
		return
	}
	var req model.RoleUpdateRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}
	role, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		h.handleErr(w, err)
		return
	}
	utils.OK(w, "Role updated successfully", map[string]any{"role": role})
}

func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid role id"))
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.handleErr(w, err)
		return
	}
	utils.OK(w, "Role deleted successfully", nil)
}

func (h *RoleHandler) handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrRoleNotFound):
		utils.NotFound(w, err)
	case errors.Is(err, service.ErrRoleNameRequired), errors.Is(err, service.ErrRoleSlugRequired), errors.Is(err, service.ErrRoleDeleteBlocked):
		utils.BadRequest(w, err)
	default:
		utils.ServerError(w, err)
	}
}
