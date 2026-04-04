package model

import "time"

// Category (Level 1)
type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Thumbnail string    `json:"thumbnail"`
	Priority  int16     `json:"priority"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relation: Nested for Tree View
	SubCategories []SubCategory `json:"subitems,omitempty"`
}

// SubCategory (Level 2)
type SubCategory struct {
	ID         int64     `json:"id"`
	CategoryID int64     `json:"category_id"`
	Name       string    `json:"name"`
	Priority   int16     `json:"priority"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relation: Nested for Tree View
	SubSubCategories []SubSubCategory `json:"subitems,omitempty"`
}

// SubSubCategory (Level 3)
type SubSubCategory struct {
	ID            int64     `json:"id"`
	SubCategoryID int64     `json:"sub_category_id"`
	Name          string    `json:"name"`
	Priority      int16     `json:"priority"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
