package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type SiteSettingsHandler struct {
	svc *service.SiteSettingsService
}

func NewSiteSettingsHandler(svc *service.SiteSettingsService) *SiteSettingsHandler {
	return &SiteSettingsHandler{svc: svc}
}

func (h *SiteSettingsHandler) handleErr(w http.ResponseWriter, err error) {
	fmt.Println("Error: ", err)
	utils.ServerError(w, err)
}

func (h *SiteSettingsHandler) handleSocialLinkErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrSocialLinkNotFound):
		utils.NotFound(w, err)
	case errors.Is(err, service.ErrSocialLinkNameRequired),
		errors.Is(err, service.ErrSocialLinkIconRequired),
		errors.Is(err, service.ErrSocialLinkURLRequired),
		errors.Is(err, service.ErrSocialLinkURLInvalid),
		errors.Is(err, service.ErrSocialLinkNameTooLong),
		errors.Is(err, service.ErrSocialLinkIconTooLong),
		errors.Is(err, service.ErrSocialLinkURLTooLong),
		errors.Is(err, service.ErrSocialLinkAltTooLong),
		errors.Is(err, service.ErrSocialLinkIconInvalid):
		utils.BadRequest(w, err)
	default:
		h.handleErr(w, err)
	}
}

// GetHeroSection returns the current settings
func (h *SiteSettingsHandler) GetHeroSection(w http.ResponseWriter, r *http.Request) {
	hero, err := h.svc.GetHeroSection(r.Context())
	if err != nil {
		utils.NotFound(w, err)
		return
	}

	response := struct {
		Error       bool               `json:"error"`
		HeroSection *model.HeroSection `json:"hero_section"`
	}{
		Error:       false,
		HeroSection: hero,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *SiteSettingsHandler) ListSocialLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.svc.ListSocialLinks(r.Context())
	if err != nil {
		h.handleErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"error":        false,
		"social_links": links,
	})
}

func (h *SiteSettingsHandler) GetSocialLinkByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid social link id"))
		return
	}

	link, err := h.svc.GetSocialLinkByID(r.Context(), id)
	if err != nil {
		h.handleSocialLinkErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"error":       false,
		"social_link": link,
	})
}

func (h *SiteSettingsHandler) CreateSocialLink(w http.ResponseWriter, r *http.Request) {
	link, err := decodeSocialLinkRequest(w, r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	created, err := h.svc.CreateSocialLink(r.Context(), link)
	if err != nil {
		h.handleSocialLinkErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"error":       false,
		"message":     "Social link created successfully",
		"social_link": created,
	})
}

func (h *SiteSettingsHandler) UpdateSocialLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid social link id"))
		return
	}

	link, err := decodeSocialLinkRequest(w, r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	link.ID = id

	updated, err := h.svc.UpdateSocialLink(r.Context(), link)
	if err != nil {
		h.handleSocialLinkErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"error":       false,
		"message":     "Social link updated successfully",
		"social_link": updated,
	})
}

func (h *SiteSettingsHandler) DeleteSocialLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid social link id"))
		return
	}

	if err := h.svc.DeleteSocialLink(r.Context(), id); err != nil {
		h.handleSocialLinkErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"error":   false,
		"message": "Social link deleted successfully",
	})
}

func (h *SiteSettingsHandler) GetTopbarSocialLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.svc.ListTopbarSocialLinks(r.Context())
	if err != nil {
		h.handleErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"error":        false,
		"social_links": links,
	})
}

// UpdateHeroSection handles multipart upload for images and text
func (h *SiteSettingsHandler) UpdateHeroSection(w http.ResponseWriter, r *http.Request) {
	// 1. Parse Multipart
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.BadRequest(w, err)
		return
	}

	// 2. Construct Model from Form Values (Updated to capture ALL text)
	hero := &model.HeroSection{
		// Main
		MainTitle:    r.FormValue("main_title"),
		MainSubtitle: r.FormValue("main_subtitle"),

		// Side Top
		SideTopTitle: r.FormValue("side_top_title"),
		SideTopTag:   r.FormValue("side_top_tag"),

		// Mini Banners
		MiniBanner1Title: r.FormValue("mini_banner_1_title"),
		MiniBanner2Title: r.FormValue("mini_banner_2_title"),

		// Images (Flag checking)
		MainBanner:    checkDelete(r, "main_banner"),
		SideTopBanner: checkDelete(r, "side_top_banner"),
		MiniBanner1:   checkDelete(r, "mini_banner_1"),
		MiniBanner2:   checkDelete(r, "mini_banner_2"),
	}

	// ... (Rest of the file extraction logic remains the same) ...

	filesMap := make(map[string]*multipart.FileHeader)
	filesData := make(map[string]multipart.File)
	keys := []string{"main_banner", "side_top_banner", "mini_banner_1", "mini_banner_2"}

	for _, key := range keys {
		file, header, _ := r.FormFile(key)
		if file != nil {
			filesMap[key] = header
			filesData[key] = file
		}
	}

	err := h.svc.UpdateHeroSection(r.Context(), hero, filesMap, filesData)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	response := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Hero section updated successfully",
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

// Helper to check if a delete flag was sent (e.g., "main_banner_delete" = "true")
func checkDelete(r *http.Request, key string) string {
	if r.FormValue(key+"_delete") == "true" {
		return "" // Return empty string to signify deletion
	}
	return "KEEP" // Sentinel value to tell Service to keep existing
}

func decodeSocialLinkRequest(w http.ResponseWriter, r *http.Request) (*model.SocialLink, error) {
	var req model.SocialLinkRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	showInTopbar := false
	if req.ShowInTopbar != nil {
		showInTopbar = *req.ShowInTopbar
	}

	return &model.SocialLink{
		Name:         req.Name,
		Icon:         req.Icon,
		Link:         req.Link,
		AltText:      req.AltText,
		IsActive:     isActive,
		ShowInTopbar: showInTopbar,
		DisplayOrder: req.DisplayOrder,
	}, nil
}
