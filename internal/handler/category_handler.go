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

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) handleErr(w http.ResponseWriter, err error) {
	if strings.Contains(err.Error(), "already exists") {
		utils.WriteJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	} else {
		utils.ServerError(w, err)
	}
}

// ---------------- LEVEL 1 (Multipart with Image) ----------------

// Create Category (Multipart)
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(5 << 20); err != nil { // 5MB max
		utils.BadRequest(w, err)
		return
	}

	name := r.FormValue("name")
	priority, _ := strconv.Atoi(r.FormValue("priority"))

	cat := &model.Category{
		Name:     name,
		Priority: int16(priority),
		IsActive: true,
	}

	// Save Image
	file, header, _ := r.FormFile("logo")
	if file != nil {
		defer file.Close()
		path, err := utils.SaveMultipartImage(file, header, "categories", name)
		if err != nil {
			utils.ServerError(w, err)
			return
		}
		cat.LogoURL = path
	}

	if err := h.svc.Create(r.Context(), cat); err != nil {
		h.handleErr(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, cat)
}

// Update Category (Multipart) - Handles Photo Replacement
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	// 1. Fetch Existing Data to preserve old Logo if not updating
	existingCat, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err)
		return
	}

	// 2. Parse Form
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.BadRequest(w, err)
		return
	}

	// 3. Update Text Fields
	existingCat.Name = r.FormValue("name")
	existingCat.Priority = int16(utils.ParseInt(r.FormValue("priority")))
	if val := r.FormValue("is_active"); val != "" {
		existingCat.IsActive = (val == "true")
	}

	// 4. Handle New Image
	file, header, _ := r.FormFile("logo")
	if file != nil {
		defer file.Close()
		// Optional: Delete old image here using os.Remove(existingCat.LogoURL)

		path, err := utils.SaveMultipartImage(file, header, "categories", existingCat.Name)
		if err != nil {
			utils.ServerError(w, err)
			return
		}
		existingCat.LogoURL = path
	}
	// If file == nil, existingCat.LogoURL remains what it was from DB

	// 5. Save
	if err := h.svc.Update(r.Context(), existingCat); err != nil {
		h.handleErr(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, existingCat)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.Delete(r.Context(), id); err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	ssc, err := h.svc.GetSubSubByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, ssc)
}

// ---------------- LEVEL 2 (JSON - No Image) ----------------

func (h *CategoryHandler) CreateSub(w http.ResponseWriter, r *http.Request) {
	var sub model.SubCategory
	// Read JSON Body
	if err := utils.ReadJSON(w, r, &sub); err != nil {
		utils.BadRequest(w, err)
		return
	}
	sub.IsActive = true // Default

	if err := h.svc.CreateSub(r.Context(), &sub); err != nil {
		h.handleErr(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, sub)
}

func (h *CategoryHandler) UpdateSub(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var sub model.SubCategory
	if err := utils.ReadJSON(w, r, &sub); err != nil {
		utils.BadRequest(w, err)
		return
	}
	sub.ID = id

	if err := h.svc.UpdateSub(r.Context(), &sub); err != nil {
		h.handleErr(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, sub)
}

func (h *CategoryHandler) DeleteSub(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.DeleteSub(r.Context(), id); err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
func (h *CategoryHandler) GetSubByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	sub, err := h.svc.GetSubByID(r.Context(), id)
	if err != nil {
		// Differentiate between 404 and 500
		utils.NotFound(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, sub)
}

// ---------------- LEVEL 3 (JSON - No Image) ----------------

func (h *CategoryHandler) CreateSubSub(w http.ResponseWriter, r *http.Request) {
	var ssc model.SubSubCategory
	if err := utils.ReadJSON(w, r, &ssc); err != nil {
		utils.BadRequest(w, err)
		return
	}
	ssc.IsActive = true

	if err := h.svc.CreateSubSub(r.Context(), &ssc); err != nil {
		h.handleErr(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, ssc)
}

func (h *CategoryHandler) UpdateSubSub(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var ssc model.SubSubCategory
	if err := utils.ReadJSON(w, r, &ssc); err != nil {
		utils.BadRequest(w, err)
		return
	}
	ssc.ID = id

	if err := h.svc.UpdateSubSub(r.Context(), &ssc); err != nil {
		h.handleErr(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, ssc)
}

func (h *CategoryHandler) DeleteSubSub(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.DeleteSubSub(r.Context(), id); err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// Tree View
func (h *CategoryHandler) GetTree(w http.ResponseWriter, r *http.Request) {
	onlyActive := r.URL.Query().Get("active") == "true"
	tree, err := h.svc.GetTree(r.Context(), onlyActive)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, tree)
}

func (h *CategoryHandler) GetSubSubByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	ssc, err := h.svc.GetSubSubByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, ssc)
}
