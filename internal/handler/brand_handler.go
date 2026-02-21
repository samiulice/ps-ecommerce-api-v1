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

type BrandHandler struct {
	svc *service.BrandService
}

func NewBrandHandler(svc *service.BrandService) *BrandHandler {
	return &BrandHandler{svc: svc}
}

func (h *BrandHandler) handleErr(w http.ResponseWriter, err error) {
	fmt.Println("Error: ", err)
	if strings.Contains(err.Error(), "already exists") {
		utils.WriteJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	} else {
		utils.ServerError(w, err)
	}
}

// Create Brand (Multipart)
func (h *BrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(5 << 20); err != nil { // 5MB max
		utils.BadRequest(w, err)
		return
	}

	name := r.FormValue("name")
	priority, _ := strconv.Atoi(r.FormValue("priority"))

	brand := &model.Brand{
		Name:     name,
		Priority: int16(priority),
		IsActive: true,
	}

	// Image
	file, header, _ := r.FormFile("thumbnail")
	if file != nil {
		fmt.Println("Image exist")
	}
	err := h.svc.Create(r.Context(), brand, file, header)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	var response struct {
		Error   bool         `json:"error"`
		Message string       `json:"message"`
		Brand   *model.Brand `json:"brand"`
	}
	response.Error = false
	response.Message = "Brand added successfully"
	response.Brand = brand
	utils.WriteJSON(w, http.StatusOK, response)
}

// Update Brand (Multipart) - Handles Photo Replacement
func (h *BrandHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	// 1. Parse Form
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.BadRequest(w, err)
		return
	}

	var brand model.Brand
	// 2. Update Text Fields
	brand.ID = id
	brand.Name = r.FormValue("name")
	brand.Priority = int16(utils.ParseInt(r.FormValue("priority")))
	if val := r.FormValue("is_active"); val != "" {
		brand.IsActive = val == "1"
	}

	// 4. Handle New Image
	file, header, _ := r.FormFile("thumbnail")

	// 5. Save
	if err := h.svc.Update(r.Context(), &brand, file, header); err != nil {
		h.handleErr(w, err)
		return
	}
	var response struct {
		Error   bool        `json:"error"`
		Message string      `json:"message"`
		Brand   model.Brand `json:"brand"`
	}
	response.Error = false
	response.Message = "Brand updated successfully"
	response.Brand = brand
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *BrandHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	response.Message = "Brand deleted successfully"
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *BrandHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	ssc, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, ssc)
}
func (h *BrandHandler) GetBrands(w http.ResponseWriter, r *http.Request) {
	status := strings.TrimSpace(r.URL.Query().Get("status"))

	brand, err := h.svc.GetBrands(r.Context(), status)
	if err != nil {
		utils.NotFound(w, err)
		return
	}
	var response struct {
		Error   bool           `json:"error"`
		Message string         `json:"message"`
		Brand   []*model.Brand `json:"brands"`
	}
	response.Error = false
	response.Message = "Brand retrieved"
	response.Brand = brand
	utils.WriteJSON(w, http.StatusOK, response)
}
