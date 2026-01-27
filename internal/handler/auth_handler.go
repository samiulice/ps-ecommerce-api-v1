// Package handler contains HTTP handlers.
package handler

import (
	"encoding/json"
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

// Register handles POST /auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.User
	_ = utils.ReadJSON(w, r, &req)
	fmt.Printf("%+v", req)
	if err := h.svc.Register(r.Context(), req.Email, req.Password); err != nil {
		utils.BadRequest(w, err)
		return
	}

	utils.Created(w, "User registration successful", nil)
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.User
	_ = utils.ReadJSON(w, r, &req)

	user, access, refresh, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//sanitize output
	user.Password = ""
	var response struct {
		Error   bool        `json:"error"`
		Message string      `json:"message"`
		User    *model.User `json:"user"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	response.Error = false
	response.Message = "Login Successful"
	response.User = user
	response.AccessToken = access
	response.RefreshToken = refresh
	fmt.Printf("%+v\n", response)
	_ = utils.WriteJSON(w, http.StatusOK, response)
}

// Refresh handles POST /auth/refresh.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"refresh_token"`
	}
	_ = utils.ReadJSON(w, r, &req)

	access, err := h.svc.Refresh(r.Context(), req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token": access,
	})
}

// Me handles GET /api/me (protected endpoint).
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"user_id": uid,
	})
}
