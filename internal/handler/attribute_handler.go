package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type AttributeHandler struct {
	svc *service.AttributeService
}

func NewAttributeHandler(svc *service.AttributeService) *AttributeHandler {
	return &AttributeHandler{svc: svc}
}

// Create Attribute
func (h *AttributeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u model.Attribute
	if err := utils.ReadJSON(w, r, &u); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if err := h.svc.Create(r.Context(), &u); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, u)
}

// Update Attribute
func (h *AttributeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var u model.Attribute
	if err := utils.ReadJSON(w, r, &u); err != nil {
		utils.BadRequest(w, err)
		return
	}
	u.ID = id

	if err := h.svc.Update(r.Context(), &u); err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, u)
}

// Delete Attribute
func (h *AttributeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	if err := h.svc.Delete(r.Context(), id); err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "attribute deleted successfully"})
}

// Get By ID
func (h *AttributeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	u, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, u)
}

// List All Attributes
func (h *AttributeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	attributes, err := h.svc.GetAll(r.Context())
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	// Wrapping in a standard response format if you prefer
	var response struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Attributes   []*model.Attribute `json:"attributes"`
	}
	response.Error = false
	response.Message = "Attributes retrieved"
	response.Attributes = attributes

	utils.WriteJSON(w, http.StatusOK, response)
}
