package service

import (
	"context"
	"errors"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

var (
	ErrSocialLinkNotFound     = errors.New("social link not found")
	ErrSocialLinkNameRequired = errors.New("social link name is required")
	ErrSocialLinkIconRequired = errors.New("social link icon is required")
	ErrSocialLinkURLRequired  = errors.New("social link URL is required")
	ErrSocialLinkURLInvalid   = errors.New("social link URL must be a valid http or https URL")
	ErrSocialLinkNameTooLong  = errors.New("social link name is too long")
	ErrSocialLinkIconTooLong  = errors.New("social link icon is too long")
	ErrSocialLinkURLTooLong   = errors.New("social link URL is too long")
	ErrSocialLinkAltTooLong   = errors.New("social link alt text is too long")
	ErrSocialLinkIconInvalid  = errors.New("social link icon contains invalid characters")
)

var socialIconClassRe = regexp.MustCompile(`^[a-zA-Z0-9 _:-]+$`)

type SiteSettingsService struct {
	repo *repository.SiteSettingsRepo
}

func NewSiteSettingsService(repo *repository.SiteSettingsRepo) *SiteSettingsService {
	return &SiteSettingsService{repo: repo}
}

func (s *SiteSettingsService) GetHeroSection(ctx context.Context) (*model.HeroSection, error) {
	return s.repo.GetHeroSection(ctx)
}

func (s *SiteSettingsService) ListSocialLinks(ctx context.Context) ([]model.SocialLink, error) {
	return s.repo.ListSocialLinks(ctx)
}

func (s *SiteSettingsService) GetSocialLinkByID(ctx context.Context, id int64) (*model.SocialLink, error) {
	link, err := s.repo.GetSocialLinkByID(ctx, id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, ErrSocialLinkNotFound
		}
		return nil, err
	}
	return link, nil
}

func (s *SiteSettingsService) CreateSocialLink(ctx context.Context, link *model.SocialLink) (*model.SocialLink, error) {
	if err := validateSocialLink(link); err != nil {
		return nil, err
	}
	if err := s.repo.CreateSocialLink(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *SiteSettingsService) UpdateSocialLink(ctx context.Context, link *model.SocialLink) (*model.SocialLink, error) {
	if link.ID <= 0 {
		return nil, ErrSocialLinkNotFound
	}
	if err := validateSocialLink(link); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateSocialLink(ctx, link); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, ErrSocialLinkNotFound
		}
		return nil, err
	}
	return s.repo.GetSocialLinkByID(ctx, link.ID)
}

func (s *SiteSettingsService) DeleteSocialLink(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrSocialLinkNotFound
	}
	if err := s.repo.DeleteSocialLink(ctx, id); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return ErrSocialLinkNotFound
		}
		return err
	}
	return nil
}

func (s *SiteSettingsService) ListTopbarSocialLinks(ctx context.Context) ([]model.SocialLink, error) {
	return s.repo.ListTopbarSocialLinks(ctx, 5)
}

// UpdateHeroSection handles multiple image uploads and removals
// filesMap maps the key (e.g., "main_banner") to the uploaded file data
func (s *SiteSettingsService) UpdateHeroSection(ctx context.Context, h *model.HeroSection, filesMap map[string]*multipart.FileHeader, filesData map[string]multipart.File) error {

	// 1. Fetch Existing Data to preserve old images if not updating, or to delete later
	existingHero, err := s.repo.GetHeroSection(ctx)
	if err != nil {
		return err
	}

	// 2. Prepare new URLs based on uploaded files
	// Logic: If file exists in map, update h.Field with new URL.
	// If not, check if it was marked for deletion (empty string).
	// If neither, keep the existing URL from DB.

	// --- Main Banner ---
	if header, ok := filesMap["main_banner"]; ok {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		h.MainBanner = utils.GetHeroBannerURL("main", ext)
	} else if h.MainBanner == "" {
		// If empty, it means user sent "delete" flag (handled in handler), so we leave it empty
	} else {
		// Otherwise, keep existing
		h.MainBanner = existingHero.MainBanner
	}

	// --- Side Top Banner ---
	if header, ok := filesMap["side_top_banner"]; ok {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		h.SideTopBanner = utils.GetHeroBannerURL("side-top", ext)
	} else if h.SideTopBanner == "" {
		// User requested delete
	} else {
		h.SideTopBanner = existingHero.SideTopBanner
	}

	// --- Mini Banner 1 ---
	if header, ok := filesMap["mini_banner_1"]; ok {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		h.MiniBanner1 = utils.GetHeroBannerURL("mini-1", ext)
	} else if h.MiniBanner1 == "" {
		// User requested delete
	} else {
		h.MiniBanner1 = existingHero.MiniBanner1
	}

	// --- Mini Banner 2 ---
	if header, ok := filesMap["mini_banner_2"]; ok {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		h.MiniBanner2 = utils.GetHeroBannerURL("mini-2", ext)
	} else if h.MiniBanner2 == "" {
		// User requested delete
	} else {
		h.MiniBanner2 = existingHero.MiniBanner2
	}

	// 3. Update Database
	err = s.repo.UpdateHeroSection(ctx, h)
	if err != nil {
		return err
	}

	// 4. Save Files & Delete Old (Post-DB Success)
	// Helper to reduce repetition in the saving phase
	saveAndClean := func(key string, newURL string, oldURL string, namePrefix string) {
		if file, ok := filesData[key]; ok && file != nil {
			defer file.Close()
			header := filesMap[key]

			// Save new image
			_, err := utils.SaveMultipartImage(file, header, utils.GetHeroFolderPath(""), namePrefix)
			if err == nil {
				// Delete old image if it's different and not empty
				if oldURL != "" && oldURL != newURL {
					utils.DeleteFile(utils.GetHeroFolderPath(filepath.Base(oldURL)))
				}
			}
		} else if newURL == "" && oldURL != "" {
			// Case: User requested delete (no new file, but URL is empty)
			utils.DeleteFile(utils.GetHeroFolderPath(filepath.Base(oldURL)))
		}
	}

	// Execute file operations
	saveAndClean("main_banner", h.MainBanner, existingHero.MainBanner, "main")
	saveAndClean("side_top_banner", h.SideTopBanner, existingHero.SideTopBanner, "side-top")
	saveAndClean("mini_banner_1", h.MiniBanner1, existingHero.MiniBanner1, "mini-1")
	saveAndClean("mini_banner_2", h.MiniBanner2, existingHero.MiniBanner2, "mini-2")

	return nil
}

// Helper for reflection-less property access (mock logic)
func getattr(h *model.HeroSection, field string) string {
	switch field {
	case "main_banner":
		return h.MainBanner
	case "side_top_banner":
		return h.SideTopBanner
	default:
		return ""
	}
}

func validateSocialLink(link *model.SocialLink) error {
	link.Name = strings.TrimSpace(link.Name)
	link.Icon = strings.TrimSpace(link.Icon)
	link.Link = strings.TrimSpace(link.Link)
	link.AltText = strings.TrimSpace(link.AltText)

	if link.Name == "" {
		return ErrSocialLinkNameRequired
	}
	if len(link.Name) > 120 {
		return ErrSocialLinkNameTooLong
	}
	if link.Icon == "" {
		return ErrSocialLinkIconRequired
	}
	if len(link.Icon) > 160 {
		return ErrSocialLinkIconTooLong
	}
	if !socialIconClassRe.MatchString(link.Icon) {
		return ErrSocialLinkIconInvalid
	}
	if link.Link == "" {
		return ErrSocialLinkURLRequired
	}
	if len(link.Link) > 500 {
		return ErrSocialLinkURLTooLong
	}
	parsed, err := url.Parse(link.Link)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ErrSocialLinkURLInvalid
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrSocialLinkURLInvalid
	}
	if link.AltText == "" {
		link.AltText = link.Name
	}
	if len(link.AltText) > 255 {
		return ErrSocialLinkAltTooLong
	}
	if link.DisplayOrder < 0 {
		link.DisplayOrder = 0
	}

	return nil
}
