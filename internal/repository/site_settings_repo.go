package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type SiteSettingsRepo struct {
	db *pgxpool.Pool
}

func NewSiteSettingsRepo(db *pgxpool.Pool) *SiteSettingsRepo {
	return &SiteSettingsRepo{db: db}
}

// GetHeroSection retrieves the single hero section row (assuming ID always 1)
func (r *SiteSettingsRepo) GetHeroSection(ctx context.Context) (*model.HeroSection, error) {
	// Assuming a singleton row where id = 1
	query := `
        SELECT id, 
               main_banner, main_title, main_subtitle,
               side_top_banner, side_top_title, side_top_tag,
               mini_banner_1, mini_banner_1_title,
               mini_banner_2, mini_banner_2_title,
               updated_at
        FROM hero_sections 
        WHERE id = 1
    `
	var h model.HeroSection
	err := r.db.QueryRow(ctx, query).Scan(
		&h.ID,
		&h.MainBanner, &h.MainTitle, &h.MainSubtitle, // Main Banner
		&h.SideTopBanner, &h.SideTopTitle, &h.SideTopTag, // Side Top Banner
		&h.MiniBanner1, &h.MiniBanner1Title, // Mini Banner 1
		&h.MiniBanner2, &h.MiniBanner2Title, // Mini Banner 2
		&h.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		// Optional: Auto-seed if missing, or return error
		return nil, errors.New("hero section settings not found")
	}
	return &h, err
}

// UpdateHeroSection updates images AND all text fields
func (r *SiteSettingsRepo) UpdateHeroSection(ctx context.Context, h *model.HeroSection) error {
	query := `
		UPDATE hero_sections 
		SET 
			main_banner=$1, main_title=$2, main_subtitle=$3,
			side_top_banner=$4, side_top_title=$5, side_top_tag=$6,
			mini_banner_1=$7, mini_banner_1_title=$8,
			mini_banner_2=$9, mini_banner_2_title=$10,
			updated_at=CURRENT_TIMESTAMP
		WHERE id = 1 
		RETURNING updated_at
	`
	// Make sure the order of arguments matches the $1, $2, $3... placeholders exactly
	err := r.db.QueryRow(ctx, query,
		h.MainBanner, h.MainTitle, h.MainSubtitle, // $1, $2, $3
		h.SideTopBanner, h.SideTopTitle, h.SideTopTag, // $4, $5, $6
		h.MiniBanner1, h.MiniBanner1Title, // $7, $8
		h.MiniBanner2, h.MiniBanner2Title, // $9, $10
	).Scan(&h.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update hero section: %w", err)
	}
	return nil
}

func (r *SiteSettingsRepo) ListSocialLinks(ctx context.Context) ([]model.SocialLink, error) {
	query := `
		SELECT id, name, icon, link, alt_text, is_active, show_in_topbar, display_order, created_at, updated_at
		FROM social_links
		ORDER BY display_order ASC, id ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list social links: %w", err)
	}
	defer rows.Close()

	links := make([]model.SocialLink, 0)
	for rows.Next() {
		var link model.SocialLink
		if err := rows.Scan(
			&link.ID,
			&link.Name,
			&link.Icon,
			&link.Link,
			&link.AltText,
			&link.IsActive,
			&link.ShowInTopbar,
			&link.DisplayOrder,
			&link.CreatedAt,
			&link.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan social link: %w", err)
		}
		links = append(links, link)
	}

	return links, rows.Err()
}

func (r *SiteSettingsRepo) GetSocialLinkByID(ctx context.Context, id int64) (*model.SocialLink, error) {
	query := `
		SELECT id, name, icon, link, alt_text, is_active, show_in_topbar, display_order, created_at, updated_at
		FROM social_links
		WHERE id = $1
	`

	var link model.SocialLink
	err := r.db.QueryRow(ctx, query, id).Scan(
		&link.ID,
		&link.Name,
		&link.Icon,
		&link.Link,
		&link.AltText,
		&link.IsActive,
		&link.ShowInTopbar,
		&link.DisplayOrder,
		&link.CreatedAt,
		&link.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.New("social link not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get social link: %w", err)
	}

	return &link, nil
}

func (r *SiteSettingsRepo) CreateSocialLink(ctx context.Context, link *model.SocialLink) error {
	query := `
		INSERT INTO social_links (name, icon, link, alt_text, is_active, show_in_topbar, display_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		link.Name,
		link.Icon,
		link.Link,
		link.AltText,
		link.IsActive,
		link.ShowInTopbar,
		link.DisplayOrder,
	).Scan(&link.ID, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create social link: %w", err)
	}

	return nil
}

func (r *SiteSettingsRepo) UpdateSocialLink(ctx context.Context, link *model.SocialLink) error {
	query := `
		UPDATE social_links
		SET name = $1,
			icon = $2,
			link = $3,
			alt_text = $4,
			is_active = $5,
			show_in_topbar = $6,
			display_order = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		link.Name,
		link.Icon,
		link.Link,
		link.AltText,
		link.IsActive,
		link.ShowInTopbar,
		link.DisplayOrder,
		link.ID,
	).Scan(&link.UpdatedAt)
	if err == pgx.ErrNoRows {
		return errors.New("social link not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update social link: %w", err)
	}

	return nil
}

func (r *SiteSettingsRepo) DeleteSocialLink(ctx context.Context, id int64) error {
	result, err := r.db.Exec(ctx, `DELETE FROM social_links WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete social link: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("social link not found")
	}

	return nil
}

func (r *SiteSettingsRepo) ListTopbarSocialLinks(ctx context.Context, limit int) ([]model.SocialLink, error) {
	if limit <= 0 {
		limit = 5
	}

	query := `
		SELECT id, name, icon, link, alt_text, is_active, show_in_topbar, display_order, created_at, updated_at
		FROM social_links
		WHERE is_active = TRUE AND show_in_topbar = TRUE
		ORDER BY display_order ASC, id ASC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list topbar social links: %w", err)
	}
	defer rows.Close()

	links := make([]model.SocialLink, 0)
	for rows.Next() {
		var link model.SocialLink
		if err := rows.Scan(
			&link.ID,
			&link.Name,
			&link.Icon,
			&link.Link,
			&link.AltText,
			&link.IsActive,
			&link.ShowInTopbar,
			&link.DisplayOrder,
			&link.CreatedAt,
			&link.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan topbar social link: %w", err)
		}
		links = append(links, link)
	}

	return links, rows.Err()
}
