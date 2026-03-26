package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type SupplierHandler struct {
	svc *service.SupplierService
}

func NewSupplierHandler(svc *service.SupplierService) *SupplierHandler {
	return &SupplierHandler{svc: svc}
}

func (h *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.SupplierCreateRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	supplier := mapCreateRequestToSupplier(req)

	resSupplier, err := h.svc.Create(r.Context(), supplier)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error    bool                    `json:"error"`
		Message  string                  `json:"message"`
		Supplier *model.SupplierResponse `json:"supplier"`
	}{
		Error:    false,
		Message:  "Supplier created successfully",
		Supplier: resSupplier.ToResponse(),
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}

func (h *SupplierHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid supplier ID"))
		return
	}

	var req model.SupplierUpdateRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	supplier := mapUpdateRequestToSupplier(id, req)

	resSupplier, err := h.svc.Update(r.Context(), supplier)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error    bool                    `json:"error"`
		Message  string                  `json:"message"`
		Supplier *model.SupplierResponse `json:"supplier"`
	}{
		Error:    false,
		Message:  "Supplier updated successfully",
		Supplier: resSupplier.ToResponse(),
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *SupplierHandler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filter := model.SupplierFilter{
		Search: query.Get("search"),
	}

	status, err := strconv.ParseBool(query.Get("is_active"))
	if err == nil {
		filter.CheckAccountStatus = true
		filter.IsActive = status
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	filter.Page = page

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	filter.Limit = limit

	suppliers, total, err := h.svc.ListSuppliers(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	utils.OK(w, "Suppliers retrieved successfully", map[string]interface{}{
		"suppliers":   suppliers,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

func (h *SupplierHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid supplier ID"))
		return
	}

	supplier, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error    bool                    `json:"error"`
		Message  string                  `json:"message"`
		Supplier *model.SupplierResponse `json:"supplier"`
	}{
		Error:    false,
		Message:  "Supplier retrieved successfully",
		Supplier: supplier.ToResponse(),
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *SupplierHandler) UpdateAccountStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid supplier ID"))
		return
	}

	status, err := strconv.ParseBool(r.URL.Query().Get("is_active"))
	if err != nil {
		utils.BadRequest(w, errors.New("invalid account status"))
		return
	}

	if err := h.svc.UpdateAccountStatus(r.Context(), status, id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Account status updated successfully",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *SupplierHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid supplier ID"))
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Supplier deleted successfully",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *SupplierHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrSupplierNotFound):
		utils.NotFound(w, err)
	case errors.Is(err, service.ErrSupplierNameRequired),
		errors.Is(err, service.ErrSupplierCodeRequired),
		errors.Is(err, service.ErrSupplierPhoneRequired),
		errors.Is(err, service.ErrSupplierPhoneInvalid),
		errors.Is(err, service.ErrSupplierCreditInvalid),
		errors.Is(err, service.ErrSupplierBalanceInvalid),
		errors.Is(err, service.ErrSupplierLeadTimeInvalid),
		errors.Is(err, service.ErrSupplierRatingInvalid):
		utils.BadRequest(w, err)
	case errors.Is(err, service.ErrSupplierCodeDuplicate),
		errors.Is(err, service.ErrSupplierPhoneDuplicate),
		errors.Is(err, service.ErrSupplierEmailDuplicate):
		utils.WriteJSON(w, http.StatusConflict, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
	default:
		fmt.Println("supplier error:", err)
		utils.ServerError(w, err)
	}
}

func mapCreateRequestToSupplier(req model.SupplierCreateRequest) *model.Supplier {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return &model.Supplier{
		SupplierCode:       strings.TrimSpace(req.SupplierCode),
		Name:               strings.TrimSpace(req.Name),
		CompanyName:        model.ToNullString(req.CompanyName),
		ContactPerson:      model.ToNullString(req.ContactPerson),
		Phone:              strings.TrimSpace(req.Phone),
		Email:              model.ToNullString(req.Email),
		Website:            model.ToNullString(req.Website),
		TaxID:              model.ToNullString(req.TaxID),
		TradeLicenseNo:     model.ToNullString(req.TradeLicenseNo),
		PaymentTerms:       model.ToNullString(req.PaymentTerms),
		CreditLimit:        toSQLFloat64(req.CreditLimit),
		OutstandingBalance: toSQLFloat64(req.OutstandingBalance),
		LeadTimeDays:       toSQLInt32(req.LeadTimeDays),
		Rating:             toSQLFloat64(req.Rating),
		StreetAddress:      model.ToNullString(req.StreetAddress),
		Country:            model.ToNullString(req.Country),
		City:               model.ToNullString(req.City),
		Zip:                model.ToNullString(req.Zip),
		Notes:              model.ToNullString(req.Notes),
		IsActive:           isActive,
	}
}

func mapUpdateRequestToSupplier(id int64, req model.SupplierUpdateRequest) *model.Supplier {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return &model.Supplier{
		ID:                 id,
		SupplierCode:       strings.TrimSpace(req.SupplierCode),
		Name:               strings.TrimSpace(req.Name),
		CompanyName:        model.ToNullString(req.CompanyName),
		ContactPerson:      model.ToNullString(req.ContactPerson),
		Phone:              strings.TrimSpace(req.Phone),
		Email:              model.ToNullString(req.Email),
		Website:            model.ToNullString(req.Website),
		TaxID:              model.ToNullString(req.TaxID),
		TradeLicenseNo:     model.ToNullString(req.TradeLicenseNo),
		PaymentTerms:       model.ToNullString(req.PaymentTerms),
		CreditLimit:        toSQLFloat64(req.CreditLimit),
		OutstandingBalance: toSQLFloat64(req.OutstandingBalance),
		LeadTimeDays:       toSQLInt32(req.LeadTimeDays),
		Rating:             toSQLFloat64(req.Rating),
		StreetAddress:      model.ToNullString(req.StreetAddress),
		Country:            model.ToNullString(req.Country),
		City:               model.ToNullString(req.City),
		Zip:                model.ToNullString(req.Zip),
		Notes:              model.ToNullString(req.Notes),
		IsActive:           isActive,
	}
}

func toSQLFloat64(v float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: v, Valid: true}
}

func toSQLInt32(v int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(v), Valid: true}
}
