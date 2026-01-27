package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// Create handles product creation with multiple images (Main + Variations)
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Parse Multipart Form (10MB limit)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.BadRequest(w, fmt.Errorf("invalid form data: %v", err))
		return
	}

	// 2. Extract Basic Fields
	catID, _ := strconv.ParseInt(r.FormValue("categoryId"), 10, 64)
	subCatID, _ := strconv.ParseInt(r.FormValue("subCategoryId"), 10, 64)
	subSubCatID, _ := strconv.ParseInt(r.FormValue("subSubCategoryId"), 10, 64)
	brandID, _ := strconv.ParseInt(r.FormValue("brandId"), 10, 64)
	unitPrice, _ := strconv.ParseFloat(r.FormValue("unitPrice"), 64)
	minOrderQty, _ := strconv.ParseFloat(r.FormValue("minOrderQty"), 64)
	currentStockQty, _ := strconv.ParseFloat(r.FormValue("currentStockQty"), 64)
	stockAlertQty, _ := strconv.ParseFloat(r.FormValue("stockAlertQty"), 64)
	discountAmount, _ := strconv.ParseFloat(r.FormValue("discountAmount"), 64)
	taxAmount, _ := strconv.ParseFloat(r.FormValue("taxAmount"), 64)
	shippingCost, _ := strconv.ParseFloat(r.FormValue("shippingCost"), 64)
	hasVariation := r.FormValue("has_variation") == "true"

	product := &model.Product{
		Name:             r.FormValue("name"),
		Description:      r.FormValue("description"),
		CategoryID:       catID,
		SubCategoryID:    subCatID,
		SubSubCategoryID: subSubCatID,
		BrandID:          brandID,
		SKU:              r.FormValue("sku"),
		Unit:             r.FormValue("unit"),
		SearchTags:       r.FormValue("searchTags"),
		UnitPrice:        unitPrice,
		MinOrderQty:      minOrderQty,
		CurrentStockQty:  currentStockQty,
		StockAlertQty:    stockAlertQty,
		DiscountType:     r.FormValue("discountType"),
		DiscountAmount:   discountAmount,
		TaxAmount:        taxAmount,
		TaxCalculation:   r.FormValue("taxCalculation"),
		ShippingCost:     shippingCost,
		ShippingCostType: r.FormValue("shippingCostType"),
		HasVariation:     r.FormValue("hasVariation") == "true",
	}

	// 3. Handle Main Product Thumbnail
	file, header, err := r.FormFile("thumbnail")
	if err == nil && file != nil {
		defer file.Close()
		saveDir := utils.GetProductFolderPath("")
		// SaveMultipartImage is your utility function
		path, err := utils.SaveMultipartImage(file, header, saveDir, product.Name+"_main")
		if err != nil {
			utils.ServerError(w, fmt.Errorf("failed to save main thumbnail: %w", err))
			return
		}
		product.Thumbnail = path
	}

	// 4. Handle Variations & Their Thumbnails
	if hasVariation {
		var vars []model.ProductVariation
		// Decode JSON string: '[{"name":"Red", "price":10}, ...]'
		varsJSON := r.FormValue("variations")
		if err := json.Unmarshal([]byte(varsJSON), &vars); err != nil {
			utils.BadRequest(w, fmt.Errorf("invalid variations json: %v", err))
			return
		}

		// Iterate variations to check for matching image files
		for i, v := range vars {
			// Convention: variation_thumb_0, variation_thumb_1, etc.
			formKey := fmt.Sprintf("%s_%s_thumb", product.Name, v.SKU)

			vFile, vHeader, vErr := r.FormFile(formKey)
			if vErr == nil && vFile != nil {
				defer vFile.Close()

				// Naming convention: productName_variationSKU.jpg
				uniqueName := fmt.Sprintf("%s_%s_thumb", product.Name, v.SKU)
				saveDir := utils.GetProductFolderPath("")

				vPath, err := utils.SaveMultipartImage(vFile, vHeader, saveDir, uniqueName)
				if err == nil {
					vars[i].Thumbnail = vPath
				}
			}
		}
		product.Variations = vars
	}

	// 5. Call Service
	if err := h.svc.Create(r.Context(), product); err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		// Differentiate between "not found" and "db error" ideally
		utils.ServerError(w, err)
		return
	}
	if p == nil {
		utils.NotFound(w, fmt.Errorf("product not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, p)
}

// Update handles product modifications including images and variations.
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	// 1. Get Product ID from URL
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		utils.BadRequest(w, fmt.Errorf("invalid product id"))
		return
	}

	// 2. Parse Multipart Form (Max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.BadRequest(w, fmt.Errorf("invalid form data: %v", err))
		return
	}

	// 3. Extract Basic Fields
	// Note: We parse fields regardless of whether they changed.
	// Frontend must send all fields for a PUT (replace) or just changed ones for PATCH logic.
	// Here we assume a PUT-style update where the struct is fully populated.

	catID, _ := strconv.ParseInt(r.FormValue("categoryId"), 10, 64)
	subCatID, _ := strconv.ParseInt(r.FormValue("subCategoryId"), 10, 64)
	subSubCatID, _ := strconv.ParseInt(r.FormValue("subSubCategoryId"), 10, 64)
	brandID, _ := strconv.ParseInt(r.FormValue("brandId"), 10, 64)
	unitPrice, _ := strconv.ParseFloat(r.FormValue("unitPrice"), 64)
	minOrderQty, _ := strconv.ParseFloat(r.FormValue("minOrderQty"), 64)
	currentStockQty, _ := strconv.ParseFloat(r.FormValue("currentStockQty"), 64)
	stockAlertQty, _ := strconv.ParseFloat(r.FormValue("stockAlertQty"), 64)
	discountAmount, _ := strconv.ParseFloat(r.FormValue("discountAmount"), 64)
	taxAmount, _ := strconv.ParseFloat(r.FormValue("taxAmount"), 64)
	shippingCost, _ := strconv.ParseFloat(r.FormValue("shippingCost"), 64)
	hasVariation := r.FormValue("has_variation") == "true"

	product := &model.Product{
		Name:             r.FormValue("name"),
		Description:      r.FormValue("description"),
		CategoryID:       catID,
		SubCategoryID:    subCatID,
		SubSubCategoryID: subSubCatID,
		BrandID:          brandID,
		SKU:              r.FormValue("sku"),
		Unit:             r.FormValue("unit"),
		SearchTags:       r.FormValue("searchTags"),
		UnitPrice:        unitPrice,
		MinOrderQty:      minOrderQty,
		CurrentStockQty:  currentStockQty,
		StockAlertQty:    stockAlertQty,
		DiscountType:     r.FormValue("discountType"),
		DiscountAmount:   discountAmount,
		TaxAmount:        taxAmount,
		TaxCalculation:   r.FormValue("taxCalculation"),
		ShippingCost:     shippingCost,
		ShippingCostType: r.FormValue("shippingCostType"),
		HasVariation:     r.FormValue("hasVariation") == "true",
	}

	// 4. Handle Main Product Thumbnail (Optional)
	// Only update if a new file is provided.
	file, header, err := r.FormFile("thumbnail")
	if err == nil && file != nil {
		defer file.Close()

		// Use a unique name or overwrite existing SKU based name
		// Ideally append timestamp to avoid caching issues on frontend
		saveDir := utils.GetProductFolderPath("")
		path, err := utils.SaveMultipartImage(file, header, saveDir, product.Name+"_main")
		if err != nil {
			utils.ServerError(w, fmt.Errorf("failed to save new thumbnail: %w", err))
			return
		}
		product.Thumbnail = path
	} else {
		// If no new file, keep the old URL sent by frontend (hidden field)
		// or let the Repo handle "empty string means no change"
		product.Thumbnail = r.FormValue("existing_thumbnail")
	}

	// 5. Handle Variations (Optional)
	if hasVariation {
		var vars []model.ProductVariation
		varsJSON := r.FormValue("variations") // JSON string: '[{"id": 10, "name":"Red", ...}, ...]'

		if err := json.Unmarshal([]byte(varsJSON), &vars); err != nil {
			utils.BadRequest(w, fmt.Errorf("invalid variations json: %v", err))
			return
		}

		// Iterate through variations to check for NEW images
		for i, v := range vars {
			// Check if a file exists for this index: variation_thumb_0, variation_thumb_1...
			formKey := fmt.Sprintf("%s_%s_thumb", product.Name, v.SKU)

			vFile, vHeader, vErr := r.FormFile(formKey)
			if vErr == nil && vFile != nil {
				defer vFile.Close()

				uniqueName := fmt.Sprintf("%s_%s_thumb", product.Name, v.SKU)
				saveDir := utils.GetProductFolderPath("")

				vPath, err := utils.SaveMultipartImage(vFile, vHeader, saveDir, uniqueName)
				if err == nil {
					vars[i].Thumbnail = vPath
				}
			}
			// If no file uploaded, vars[i].Thumbnail remains whatever was in the JSON
			// (Client should send the existing URL if unchanged)
		}
		product.Variations = vars
	}

	// 6. Call Service Layer
	if err := h.svc.Update(r.Context(), product); err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.Delete(r.Context(), id); err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
