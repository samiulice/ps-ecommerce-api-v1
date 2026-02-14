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
        &h.MainBanner, &h.MainTitle, &h.MainSubtitle,       // Main Banner
        &h.SideTopBanner, &h.SideTopTitle, &h.SideTopTag,   // Side Top Banner
        &h.MiniBanner1, &h.MiniBanner1Title,                // Mini Banner 1
        &h.MiniBanner2, &h.MiniBanner2Title,                // Mini Banner 2
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
		h.MainBanner, h.MainTitle, h.MainSubtitle,       // $1, $2, $3
		h.SideTopBanner, h.SideTopTitle, h.SideTopTag,   // $4, $5, $6
		h.MiniBanner1, h.MiniBanner1Title,               // $7, $8
		h.MiniBanner2, h.MiniBanner2Title,               // $9, $10
	).Scan(&h.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update hero section: %w", err)
	}
	return nil
}