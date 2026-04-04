package model

import "time"

type HeroSection struct {
	ID int64 `json:"id"`

	// Main Banner (Left)
	// Changed from MainBannerImg to MainBanner to match repo usage
	MainBanner   string `json:"main_banner"`
	MainTitle    string `json:"main_title"`
	MainSubtitle string `json:"main_subtitle"`

	// Side Top Banner (Right Top)
	SideTopBanner string `json:"side_top_banner"` // Changed to match repo
	SideTopTitle  string `json:"side_top_title"`
	SideTopTag    string `json:"side_top_tag"`

	// Mini Banner 1
	MiniBanner1      string `json:"mini_banner_1"` // Changed to match repo
	MiniBanner1Title string `json:"mini_banner_1_title"`

	// Mini Banner 2
	MiniBanner2      string `json:"mini_banner_2"` // Changed to match repo
	MiniBanner2Title string `json:"mini_banner_2_title"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
