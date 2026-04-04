package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type BranchHandler struct {
	svc *service.BranchService
}

func NewBranchHandler(svc *service.BranchService) *BranchHandler {
	return &BranchHandler{svc: svc}
}

func (h *BranchHandler) handleErr(w http.ResponseWriter, err error) {
	fmt.Println("Error: ", err)
	if strings.Contains(err.Error(), "already exists") {
		utils.WriteJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	} else {
		utils.ServerError(w, err)
	}
}

func (h *BranchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var branch model.Branch

	// Read JSON Body
	if err := utils.ReadJSON(w, r, &branch); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if err := h.svc.Create(r.Context(), &branch); err != nil {
		h.handleErr(w, err)
		return
	}

	var response struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Branch  *model.Branch `json:"branch"`
	}
	response.Error = false
	response.Message = "Branch added successfully"
	response.Branch = &branch

	utils.WriteJSON(w, http.StatusCreated, response)
}

func (h *BranchHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var branch model.Branch
	if err := utils.ReadJSON(w, r, &branch); err != nil {
		utils.BadRequest(w, err)
		return
	}
	branch.ID = id

	if err := h.svc.Update(r.Context(), &branch); err != nil {
		h.handleErr(w, err)
		return
	}

	var response struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Branch  *model.Branch `json:"branch"`
	}
	response.Error = false
	response.Message = "Branch updated successfully"
	response.Branch = &branch

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *BranchHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := h.svc.Delete(r.Context(), id); err != nil {
		utils.ServerError(w, err)
		return
	}

	var response struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	response.Error = false
	response.Message = "Branch deleted successfully"

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *BranchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	branch, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, branch)
}

func (h *BranchHandler) GetBranches(w http.ResponseWriter, r *http.Request) {
	branches, err := h.svc.GetBranches(r.Context())
	if err != nil {
		utils.NotFound(w, err)
		return
	}

	var response struct {
		Error    bool            `json:"error"`
		Message  string          `json:"message"`
		Branches []*model.Branch `json:"branches"`
	}
	response.Error = false
	response.Message = "Branches retrieved successfully"
	response.Branches = branches

	utils.WriteJSON(w, http.StatusOK, response)
}
