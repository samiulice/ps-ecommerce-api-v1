package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type POSHandler struct {
	posService *service.POSService
}

func NewPOSHandler(posService *service.POSService) *POSHandler {
	return &POSHandler{posService: posService}
}

func (h *POSHandler) CreateSale(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePOSSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	sale, err := h.posService.CreatePOSSale(r.Context(), req)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.Created(w, "POS sale created successfully", sale)
}
func (h *POSHandler) GetSaleByReference(w http.ResponseWriter, r *http.Request) {
	referenceNo := chi.URLParam(r, "reference")
	if referenceNo == "" {
		utils.BadRequest(w, fmt.Errorf("reference number is required"))
		return
	}

	sale, err := h.posService.GetPOSSaleByReference(r.Context(), referenceNo)
	if err != nil {
		utils.NotFound(w, err)
		return
	}

	utils.OK(w, "POS sale retrieved successfully", sale)
}
