package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func (h *ProductHandler) handleErr(w http.ResponseWriter, err error) {
	fmt.Println("Error: ", err)
	if strings.Contains(err.Error(), "already exists") {
		utils.WriteJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	} else {
		utils.ServerError(w, err)
	}
}

// Create handles product creation
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.BadRequest(w, err)
		return
	}

	product, err := h.parseProductForm(r)
	if err != nil {
		fmt.Println(err)
		utils.BadRequest(w, err)
		return
	}

	// Handle thumbnail
	thumbFile, thumbHeader, _ := r.FormFile("thumbnail")
	// Handle gallery
	form := r.MultipartForm
	galleryFiles := form.File["gallery_images"]

	if err := h.svc.Create(r.Context(), product, thumbFile, thumbHeader, galleryFiles); err != nil {
		fmt.Println(err)
		h.handleErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"error":   false,
		"message": "Product added successfully",
		"product": product,
	})
}

// Update handles product modification
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.BadRequest(w, err)
		return
	}

	product, err := h.parseProductForm(r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	product.ID = id

	thumbFile, thumbHeader, _ := r.FormFile("thumbnail")
	form := r.MultipartForm
	galleryFiles := form.File["gallery_images"]

	if err := h.svc.Update(r.Context(), product, thumbFile, thumbHeader, galleryFiles); err != nil {
		h.handleErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"error":   false,
		"message": "Product updated successfully",
		"product": product,
	})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.Delete(r.Context(), id); err != nil {
		utils.ServerError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"error": false, "message": "Product deleted successfully"})
}

// DeleteGalleryImage removes a specific gallery image from a product
func (h *ProductHandler) DeleteGalleryImage(w http.ResponseWriter, r *http.Request) {
	productID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	imagePath := r.URL.Query().Get("image_path")

	if productID == 0 || imagePath == "" {
		utils.BadRequest(w, fmt.Errorf("product_id and image_path are required"))
		return
	}

	if err := h.svc.DeleteGalleryImage(r.Context(), productID, imagePath); err != nil {
		h.handleErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"error":   false,
		"message": "Gallery image deleted successfully",
	})
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	product, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err)
		return
	}
	// Convert Markdown description to sanitized HTML for web rendering
	product.DescriptionHTML = utils.MarkdownToHTML(product.Description)
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"error":   false,
		"message": "Product details & variations successfully retrieved retrieved",
		"product": product,
	})
}

// GetProductVariationsByProductID read the id from the path and call service to get the product variations
func (h *ProductHandler) GetProductVariationsByProductID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	productVariations, err := h.svc.GetProductVariationsByProductID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"error":              false,
		"message":            "Product variations retrieved",
		"product_variations": productVariations,
	})
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	search := strings.TrimSpace(r.URL.Query().Get("search_text"))
	sort := strings.TrimSpace(r.URL.Query().Get("sort"))
	priceType := strings.TrimSpace(r.URL.Query().Get("price_type"))
	searchMode := strings.TrimSpace(r.URL.Query().Get("search_mode"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	compactMode := searchMode == "suggestion"

	// Multiple category support
	categoryParams := r.URL.Query()["category_id"]

	var categoryIDs []int64
	for _, idStr := range categoryParams {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			categoryIDs = append(categoryIDs, id)
		}
	}

	// sanitize values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	filter := model.ProductFilter{
		Status:      status,
		CategoryIDs: categoryIDs,
		Search:      search,
		Sort:        sort,
		PriceType:   priceType,
		Page:        page,
		Limit:       limit,
		Compact:     compactMode,
		SkipCount:   compactMode,
	}

	products, total, err := h.svc.GetProducts(r.Context(), filter)
	if err != nil {
		utils.NotFound(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"error":    false,
		"message":  "Products retrieved",
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// parseProductForm extracts product fields from multipart form based on new SQL Schema
func (h *ProductHandler) parseProductForm(r *http.Request) (*model.Product, error) {
	p := &model.Product{
		VariationAttributes: make(map[string]interface{}), // Default to empty map to avoid NULL
	}

	p.Name = r.FormValue("name")
	p.Description = r.FormValue("description")
	p.SKU = r.FormValue("sku")
	p.Tags = r.FormValue("tags")
	p.DiscountType = r.FormValue("discount_type")
	p.TaxType = r.FormValue("tax_type")
	p.ShippingType = r.FormValue("shipping_type")

	// Integers
	if v := r.FormValue("category_id"); v != "" {
		p.CategoryID, _ = strconv.ParseInt(v, 10, 64)
	}
	if v := r.FormValue("sub_category_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			p.SubCategoryID = &id
		}
	}
	if v := r.FormValue("sub_sub_category_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			p.SubSubCategoryID = &id
		}
	}
	if v := r.FormValue("brand_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			p.BrandID = &id
		}
	}
	if v := r.FormValue("unit_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			p.UnitID = &id
		}
	}
	if v := r.FormValue("status"); v != "" {
		if s, err := strconv.Atoi(v); err == nil {
			p.Status = s
		} else {
			p.Status = 1 // Default active
		}
	}

	// Floats
	p.RetailPrice, _ = strconv.ParseFloat(r.FormValue("retail_price"), 64)
	p.WholesalePrice, _ = strconv.ParseFloat(r.FormValue("wholesale_price"), 64)
	p.PurchasePrice, _ = strconv.ParseFloat(r.FormValue("purchase_price"), 64)
	p.MinRetailOrderQty, _ = strconv.ParseFloat(r.FormValue("min_retail_order_qty"), 64)
	p.MinWholesaleOrderQty, _ = strconv.ParseFloat(r.FormValue("min_wholesale_order_qty"), 64)
	p.CurrentStockQty, _ = strconv.ParseFloat(r.FormValue("current_stock_qty"), 64)
	p.StockAlertQty, _ = strconv.ParseFloat(r.FormValue("stock_alert_qty"), 64)
	p.DiscountAmount, _ = strconv.ParseFloat(r.FormValue("discount_amount"), 64)
	p.TaxAmount, _ = strconv.ParseFloat(r.FormValue("tax_amount"), 64)
	p.ShippingCost, _ = strconv.ParseFloat(r.FormValue("shipping_cost"), 64)

	// Boolean
	p.HasVariation = r.FormValue("has_variation") == "true" || r.FormValue("has_variation") == "1"

	// JSONB Fields (Variation Attributes)
	// Expecting JSON string from frontend e.g. '{"color": "red", "size": "XL"}'
	if v := r.FormValue("variation_attributes"); v != "" {
		if err := json.Unmarshal([]byte(v), &p.VariationAttributes); err != nil {
			return nil, fmt.Errorf("invalid variation_attributes json")
		}
	}

	// Handle Variations List
	// Expecting a JSON string for the list of variations if has_variation is true
	if p.HasVariation {
		if v := r.FormValue("variations"); v != "" {
			if err := json.Unmarshal([]byte(v), &p.Variations); err != nil {
				return nil, fmt.Errorf("invalid variations json data")
			}
		}
	}

	return p, nil
}
