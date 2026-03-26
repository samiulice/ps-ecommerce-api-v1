package model

import "time"

type SocialLink struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Icon         string    `json:"icon"`
	Link         string    `json:"link"`
	AltText      string    `json:"alt_text"`
	IsActive     bool      `json:"is_active"`
	ShowInTopbar bool      `json:"show_in_topbar"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SocialLinkRequest struct {
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	Link         string `json:"link"`
	AltText      string `json:"alt_text"`
	IsActive     *bool  `json:"is_active"`
	ShowInTopbar *bool  `json:"show_in_topbar"`
	DisplayOrder int    `json:"display_order"`
}
