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

// Create handles product creation (Multipart with Image)
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		utils.BadRequest(w, err)
		return
	}

	product := h.parseProductForm(r)

	// Handle thumbnail
	file, header, _ := r.FormFile("thumbnail")

	err := h.svc.Create(r.Context(), product, file, header)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	var response struct {
		Error   bool           `json:"error"`
		Message string         `json:"message"`
		Product *model.Product `json:"product"`
	}
	response.Error = false
	response.Message = "Product added successfully"
	response.Product = product
	utils.WriteJSON(w, http.StatusOK, response)
}

// Update handles product modification (Multipart with Image)
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.BadRequest(w, err)
		return
	}

	product := h.parseProductForm(r)
	product.ID = id

	// Handle new image
	file, header, _ := r.FormFile("thumbnail")

	if err := h.svc.Update(r.Context(), product, file, header); err != nil {
		h.handleErr(w, err)
		return
	}

	var response struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Product model.Product `json:"product"`
	}
	response.Error = false
	response.Message = "Product updated successfully"
	response.Product = *product
	utils.WriteJSON(w, http.StatusOK, response)
}

// Delete handles product removal
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	response.Message = "Product deleted successfully"
	utils.WriteJSON(w, http.StatusOK, response)
}

// GetByID retrieves a single product
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	product, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}

// GetProducts retrieves products with optional filters
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	published := strings.TrimSpace(r.URL.Query().Get("published"))
	categoryID := strings.TrimSpace(r.URL.Query().Get("category_id"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	filter := model.ProductFilter{
		Status:     status,
		Published:  published,
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
	}

	products, total, err := h.svc.GetProducts(r.Context(), filter)
	if err != nil {
		utils.NotFound(w, err)
		return
	}

	var response struct {
		Error    bool             `json:"error"`
		Message  string           `json:"message"`
		Products []*model.Product `json:"products"`
		Total    int64            `json:"total"`
		Page     int              `json:"page"`
		Limit    int              `json:"limit"`
	}
	response.Error = false
	response.Message = "Products retrieved"
	response.Products = products
	response.Total = total
	response.Page = page
	response.Limit = limit
	utils.WriteJSON(w, http.StatusOK, response)
}

// parseProductForm extracts product fields from multipart form
func (h *ProductHandler) parseProductForm(r *http.Request) *model.Product {
	product := &model.Product{}

	// String pointers
	if v := r.FormValue("added_by"); v != "" {
		product.AddedBy = &v
	}
	if v := r.FormValue("name"); v != "" {
		product.Name = &v
	}
	if v := r.FormValue("category_ids"); v != "" {
		product.CategoryIDs = &v
	}
	if v := r.FormValue("category_id"); v != "" {
		product.CategoryID = &v
	}
	if v := r.FormValue("sub_category_id"); v != "" {
		product.SubCategoryID = &v
	}
	if v := r.FormValue("sub_sub_category_id"); v != "" {
		product.SubSubCategoryID = &v
	}
	if v := r.FormValue("unit"); v != "" {
		product.Unit = &v
	}
	if v := r.FormValue("digital_product_type"); v != "" {
		product.DigitalProductType = &v
	}
	if v := r.FormValue("digital_file_ready"); v != "" {
		product.DigitalFileReady = &v
	}
	if v := r.FormValue("digital_file_ready_storage_type"); v != "" {
		product.DigitalFileReadyStorageType = &v
	}
	if v := r.FormValue("images"); v != "" {
		product.Images = &v
	}
	if v := r.FormValue("thumbnail_storage_type"); v != "" {
		product.ThumbnailStorageType = &v
	}
	if v := r.FormValue("preview_file"); v != "" {
		product.PreviewFile = &v
	}
	if v := r.FormValue("preview_file_storage_type"); v != "" {
		product.PreviewFileStorageType = &v
	}
	if v := r.FormValue("featured"); v != "" {
		product.Featured = &v
	}
	if v := r.FormValue("flash_deal"); v != "" {
		product.FlashDeal = &v
	}
	if v := r.FormValue("video_provider"); v != "" {
		product.VideoProvider = &v
	}
	if v := r.FormValue("video_url"); v != "" {
		product.VideoURL = &v
	}
	if v := r.FormValue("colors"); v != "" {
		product.Colors = &v
	}
	if v := r.FormValue("attributes"); v != "" {
		product.Attributes = &v
	}
	if v := r.FormValue("choice_options"); v != "" {
		product.ChoiceOptions = &v
	}
	if v := r.FormValue("variation"); v != "" {
		product.Variation = &v
	}
	if v := r.FormValue("digital_product_file_types"); v != "" {
		product.DigitalProductFileTypes = &v
	}
	if v := r.FormValue("digital_product_extensions"); v != "" {
		product.DigitalProductExtensions = &v
	}
	if v := r.FormValue("tax_type"); v != "" {
		product.TaxType = &v
	}
	if v := r.FormValue("discount_type"); v != "" {
		product.DiscountType = &v
	}
	if v := r.FormValue("details"); v != "" {
		product.Details = &v
	}
	if v := r.FormValue("attachment"); v != "" {
		product.Attachment = &v
	}
	if v := r.FormValue("meta_title"); v != "" {
		product.MetaTitle = &v
	}
	if v := r.FormValue("meta_description"); v != "" {
		product.MetaDescription = &v
	}
	if v := r.FormValue("meta_image"); v != "" {
		product.MetaImage = &v
	}
	if v := r.FormValue("denied_note"); v != "" {
		product.DeniedNote = &v
	}
	if v := r.FormValue("code"); v != "" {
		product.Code = &v
	}

	// Int64 pointers
	if v := r.FormValue("user_id"); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			product.UserID = &i
		}
	}
	if v := r.FormValue("brand_id"); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			product.BrandID = &i
		}
	}

	// Int fields
	if v := r.FormValue("min_qty"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			product.MinQty = i
		}
	} else {
		product.MinQty = 1
	}
	if v := r.FormValue("minimum_order_qty"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			product.MinimumOrderQty = i
		}
	} else {
		product.MinimumOrderQty = 1
	}
	if v := r.FormValue("current_stock"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			product.CurrentStock = &i
		}
	}

	// Float64 fields
	if v := r.FormValue("unit_price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			product.UnitPrice = f
		}
	}
	if v := r.FormValue("purchase_price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			product.PurchasePrice = f
		}
	}
	if v := r.FormValue("shipping_cost"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			product.ShippingCost = &f
		}
	}
	if v := r.FormValue("temp_shipping_cost"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			product.TempShippingCost = &f
		}
	}

	// String fields (non-pointer)
	product.ProductType = r.FormValue("product_type")
	if product.ProductType == "" {
		product.ProductType = "physical"
	}
	product.ColorImage = r.FormValue("color_image")
	product.Tax = r.FormValue("tax")
	if product.Tax == "" {
		product.Tax = "0.00"
	}
	product.TaxModel = r.FormValue("tax_model")
	if product.TaxModel == "" {
		product.TaxModel = "exclude"
	}
	product.Discount = r.FormValue("discount")
	if product.Discount == "" {
		product.Discount = "0.00"
	}

	// Boolean fields
	product.Refundable = r.FormValue("refundable") != "0" && r.FormValue("refundable") != "false"
	product.VariantProduct = r.FormValue("variant_product") == "1" || r.FormValue("variant_product") == "true"
	product.Published = r.FormValue("published") == "1" || r.FormValue("published") == "true"
	product.FreeShipping = r.FormValue("free_shipping") == "1" || r.FormValue("free_shipping") == "true"
	product.Status = r.FormValue("status") != "0" && r.FormValue("status") != "false"
	product.FeaturedStatus = r.FormValue("featured_status") != "0" && r.FormValue("featured_status") != "false"
	product.RequestStatus = r.FormValue("request_status") == "1" || r.FormValue("request_status") == "true"

	// Boolean pointers
	if v := r.FormValue("multiply_qty"); v != "" {
		b := v == "1" || v == "true"
		product.MultiplyQty = &b
	}
	if v := r.FormValue("is_shipping_cost_updated"); v != "" {
		b := v == "1" || v == "true"
		product.IsShippingCostUpdated = &b
	}

	return product
}
