// Package handler contains HTTP handlers.
package handler

import (
	"fmt"
	"net/http"

	"github.com/projuktisheba/pse-api-v1/internal/middleware"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

// AuthHandler handles authentication HTTP endpoints.
type AuthHandler struct {
	svc *service.AuthService
}

// NewAuthHandler constructs a new AuthHandler.
func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

// ==================== EMPLOYEE (ADMIN) AUTH ====================

// EmployeeRegister handles POST /auth/admin/register.
func (h *AuthHandler) EmployeeRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Mobile   string `json:"mobile"`
		Role     string `json:"role"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.BadRequest(w, nil)
		return
	}

	if req.Role == "" {
		req.Role = "staff" // default role
	}

	if err := h.svc.EmployeeRegister(r.Context(), req.Email, req.Password, req.Name, req.Mobile, req.Role); err != nil {
		utils.BadRequest(w, err)
		return
	}

	utils.Created(w, "Employee registration successful", nil)
}

// EmployeeLogin handles POST /auth/admin/login.
func (h *AuthHandler) EmployeeLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	employee, access, refresh, err := h.svc.EmployeeLogin(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	response := struct {
		Error        bool            `json:"error"`
		Message      string          `json:"message"`
		Employee     *model.Employee `json:"employee"`
		AccessToken  string          `json:"access_token"`
		RefreshToken string          `json:"refresh_token"`
	}{
		Error:        false,
		Message:      "Login successful",
		Employee:     employee,
		AccessToken:  access,
		RefreshToken: refresh,
	}
	_ = utils.WriteJSON(w, http.StatusOK, response)
}

// EmployeeRefresh handles POST /auth/admin/refresh.
func (h *AuthHandler) EmployeeRefresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"refresh_token"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	access, err := h.svc.EmployeeRefresh(r.Context(), req.Token)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, map[string]any{
		"error":        false,
		"access_token": access,
	})
}

// EmployeeMe handles GET /auth/admin/me (protected endpoint).
func (h *AuthHandler) EmployeeMe(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.CustomerIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error":   true,
			"message": "unauthorized",
		})
		return
	}

	customerType, _ := middleware.CustomerTypeFromContext(r.Context())
	role, _ := middleware.RoleFromContext(r.Context())

	_ = utils.WriteJSON(w, http.StatusOK, map[string]any{
		"error":         false,
		"customer_id":   uid,
		"customer_type": customerType,
		"role":          role,
	})
}

// ==================== CUSTOMER (CUSTOMER) AUTH ====================

// CustomerRegister handles POST /auth/customer/register.
func (h *AuthHandler) CustomerRegister(w http.ResponseWriter, r *http.Request) {
	var req model.CustomerCreateRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if req.Phone == "" || req.Password == "" {
		utils.BadRequest(w, nil)
		return
	}

	fmt.Printf("%+v\n", req)
	customer, err := h.svc.CustomerRegister(r.Context(), &req)
	if err != nil {
		utils.WriteJSON(w, http.StatusConflict, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	response := struct {
		Error    bool                    `json:"error"`
		Message  string                  `json:"message"`
		Customer *model.CustomerResponse `json:"customer"`
	}{
		Error:    false,
		Message:  "Customer registration successful",
		Customer: customer.ToResponse(),
	}
	utils.WriteJSON(w, http.StatusCreated, response)
}

// CustomerLogin handles POST /auth/customer/login.
func (h *AuthHandler) CustomerLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmailOrPhone string `json:"email_or_phone"`
		Password     string `json:"password"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	customer, access, refresh, err := h.svc.CustomerLogin(r.Context(), req.EmailOrPhone, req.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	response := struct {
		Error        bool                    `json:"error"`
		Message      string                  `json:"message"`
		Customer     *model.CustomerResponse `json:"customer"`
		AccessToken  string                  `json:"access_token"`
		RefreshToken string                  `json:"refresh_token"`
	}{
		Error:        false,
		Message:      "Login successful",
		Customer:     customer.ToResponse(),
		AccessToken:  access,
		RefreshToken: refresh,
	}
	_ = utils.WriteJSON(w, http.StatusOK, response)
}

// CustomerRefresh handles POST /auth/customer/refresh.
func (h *AuthHandler) CustomerRefresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"refresh_token"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	access, err := h.svc.CustomerRefresh(r.Context(), req.Token)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, map[string]any{
		"error":        false,
		"access_token": access,
	})
}

// CustomerMe handles GET /auth/customer/me (protected endpoint).
func (h *AuthHandler) CustomerMe(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.CustomerIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error":   true,
			"message": "unauthorized",
		})
		return
	}

	customerType, _ := middleware.CustomerTypeFromContext(r.Context())

	_ = utils.WriteJSON(w, http.StatusOK, map[string]any{
		"error":         false,
		"customer_id":   uid,
		"customer_type": customerType,
	})
}

// ==================== LEGACY HANDLERS (for backward compatibility) ====================

// Register handles POST /auth/register (legacy - maps to customer register).
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.CustomerRegister(w, r)
}

// Login handles POST /auth/login (legacy - maps to customer login).
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.CustomerLogin(w, r)
}

// Refresh handles POST /auth/refresh (legacy - maps to customer refresh).
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	h.CustomerRefresh(w, r)
}

// Me handles GET /auth/me (legacy - maps to customer me).
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	h.CustomerMe(w, r)
}
