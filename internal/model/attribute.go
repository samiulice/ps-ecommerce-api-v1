package model

import "time"

type Attribute struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"` // e.g., "Color", "Size"
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
