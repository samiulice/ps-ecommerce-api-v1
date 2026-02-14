package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

// CustomerHandler handles HTTP requests for customer operations.
type CustomerHandler struct {
	svc *service.CustomerService
}

// NewCustomerHandler creates a new CustomerHandler.
func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

// ---------------- Helpers for Nullable Types ----------------

func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func stringToNullFloat64(s string) sql.NullFloat64 {
	if s == "" {
		return sql.NullFloat64{Valid: false}
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}

func stringToNullInt32(s string) sql.NullInt32 {
	if s == "" {
		return sql.NullInt32{Valid: false}
	}
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(i), Valid: true}
}

// ---------------- Handlers ----------------

// Create handles POST /customers - creates a new customer using Multipart Form.
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Parse Multipart Form (10MB limit)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.BadRequest(w, err)
		return
	}

	// 2. Parse Boolean helpers
	parseBool := func(key string) bool {
		val := r.FormValue(key)
		b, _ := strconv.ParseBool(val)
		return b
	}

	// 3. Construct Model using Helpers
	// Note: We only map fields present in your struct and the form
	customer := &model.Customer{
		FName:           stringToNullString(r.FormValue("f_name")),
		LName:           stringToNullString(r.FormValue("l_name")),
		Name:            stringToNullString(r.FormValue("name")), // Full Name
		Phone:           r.FormValue("phone"),                    // Required string
		Email:           stringToNullString(r.FormValue("email")),
		Password:        r.FormValue("password"),
		StreetAddress:   stringToNullString(r.FormValue("street_address")),
		City:            stringToNullString(r.FormValue("city")),
		Country:         stringToNullString(r.FormValue("country")),
		Zip:             stringToNullString(r.FormValue("zip")),
		ReferralCode:    stringToNullString(r.FormValue("referral_code")),
		ReferredBy:      stringToNullInt32(r.FormValue("refer_by")), // Assumes ID is sent
		WalletBalance:   stringToNullFloat64(r.FormValue("wallet_balance")),
		LoyaltyPoint:    stringToNullFloat64(r.FormValue("loyalty_point")),
		IsActive:        parseBool("is_active"),
		IsPhoneVerified: parseBool("is_phone_verified"),
		IsEmailVerified: parseBool("is_email_verified"),
		IsTempBlocked:   parseBool("is_temp_blocked"),
	}

	// 4. Handle Image File
	file, header, _ := r.FormFile("image")

	// 5. Call Service
	// Update your Service.Create to accept (ctx, customer, file, header)
	resCustomer, err := h.svc.Create(r.Context(), customer, file, header)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// 6. Response
	response := struct {
		Error    bool                    `json:"error"`
		Message  string                  `json:"message"`
		Customer *model.CustomerResponse `json:"customer"`
	}{
		Error:    false,
		Message:  "Customer created successfully",
		Customer: resCustomer.ToResponse(),
	}
	utils.WriteJSON(w, http.StatusCreated, response)
}

// Update handles PUT /customers/{id} - updates an existing customer using Multipart Form.
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid customer ID"))
		return
	}

	// 1. Parse Multipart Form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.BadRequest(w, err)
		return
	}

	// 2. Parse Boolean helpers
	parseBool := func(key string) bool {
		val := r.FormValue(key)
		b, _ := strconv.ParseBool(val)
		return b
	}

	// 3. Construct Model
	customer := &model.Customer{
		ID:              id,
		FName:           stringToNullString(r.FormValue("f_name")),
		LName:           stringToNullString(r.FormValue("l_name")),
		Name:            stringToNullString(r.FormValue("name")),
		Phone:           r.FormValue("phone"),
		Email:           stringToNullString(r.FormValue("email")),
		StreetAddress:   stringToNullString(r.FormValue("street_address")),
		City:            stringToNullString(r.FormValue("city")),
		Country:         stringToNullString(r.FormValue("country")),
		Zip:             stringToNullString(r.FormValue("zip")),
		ReferralCode:    stringToNullString(r.FormValue("referral_code")),
		ReferredBy:      stringToNullInt32(r.FormValue("refer_by")),
		WalletBalance:   stringToNullFloat64(r.FormValue("wallet_balance")),
		LoyaltyPoint:    stringToNullFloat64(r.FormValue("loyalty_point")),
		IsActive:        parseBool("is_active"),
		IsPhoneVerified: parseBool("is_phone_verified"),
		IsEmailVerified: parseBool("is_email_verified"),
		IsTempBlocked:   parseBool("is_temp_blocked"),
	}

	// Handle Password only if provided (Optional in Edit)
	if pwd := r.FormValue("password"); pwd != "" {
		customer.Password = pwd
	}

	// 4. Handle Image File
	file, header, _ := r.FormFile("image")

	// 5. Call Service
	resCustomer, err := h.svc.Update(r.Context(), customer, file, header)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// 6. Response
	response := struct {
		Error    bool                    `json:"error"`
		Message  string                  `json:"message"`
		Customer *model.CustomerResponse `json:"customer"`
	}{
		Error:    false,
		Message:  "Customer updated successfully",
		Customer: resCustomer.ToResponse(),
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

// ... [ListCustomers, Delete, GetByID, UpdateAccountStatus, handleServiceError remain unchanged] ...
// ListCustomers handles GET /customers
func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Initialize filter with string parameters
	filter := model.CustomerFilter{
		Search: query.Get("search"),
	}

	// Parse 'is_active'
	// if error occurs, don't include this in filter
	status, err := strconv.ParseBool(r.URL.Query().Get("is_active"))
	if err != nil {
		fmt.Println(err)
		filter.CheckAccountStatus = false
	} else {
		filter.CheckAccountStatus = true
		filter.IsActive = status
	}

	// Parse 'page' (default to 1)
	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	filter.Page = page

	// Parse 'limit' (default to 20)
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	filter.Limit = limit

	// Call service layer
	customers, total, err := h.svc.ListCustomers(r.Context(), filter)
	if err != nil {
		fmt.Println(err)
		h.handleServiceError(w, err)
		return
	}

	// Calculate total pages
	// Logic: (total + limit - 1) / limit is an integer math trick for ceiling division
	totalPages := (int(total) + limit - 1) / limit

	// Return response
	utils.OK(w, "Customers retrieved successfully", map[string]interface{}{
		"customers":   customers,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// GetByID handles GET /customers/{id} - retrieves a customer by ID.
func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid customer ID"))
		return
	}

	customer, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error    bool                    `json:"error"`
		Message  string                  `json:"message"`
		Customer *model.CustomerResponse `json:"customer"`
	}{
		Error:    false,
		Message:  "Customer retrieved successfully",
		Customer: customer.ToResponse(),
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

// Update handles PUT /customers/update/account/status/{id} - updates an existing customer.
func (h *CustomerHandler) UpdateAccountStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid customer ID"))
		return
	}

	status, err := strconv.ParseBool(r.URL.Query().Get("is_active"))
	if err != nil {
		utils.BadRequest(w, errors.New("invalid account status"))
		return
	}

	err = h.svc.UpdateAccountStatus(r.Context(), status, id)
	if err != nil {
		fmt.Println(err)
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Account Status updated successfully",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

// Delete handles DELETE /customers/{id} - removes a customer.
func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid customer ID"))
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
		Message: "Customer deleted successfully",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

// handleServiceError maps service errors to HTTP responses.
func (h *CustomerHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrCustomerNotFound):
		utils.NotFound(w, err)
	case errors.Is(err, service.ErrInvalidPhone),
		errors.Is(err, service.ErrInvalidEmail),
		errors.Is(err, service.ErrInvalidPassword):
		utils.BadRequest(w, err)
	case errors.Is(err, service.ErrEmailAlreadyExists),
		errors.Is(err, service.ErrPhoneAlreadyExists):
		utils.WriteJSON(w, http.StatusConflict, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
	case errors.Is(err, service.ErrCustomerInactive):
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error":   true,
			"message": "customer account is inactive",
		})
	case errors.Is(err, service.ErrCustomerBlocked):
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error":   true,
			"message": "customer account is temporarily blocked",
		})
	default:
		utils.ServerError(w, err)
	}
}
