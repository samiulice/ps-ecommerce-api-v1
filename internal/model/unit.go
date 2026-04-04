package model

import "time"

type Unit struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`   // e.g., "Kilogram"
	Symbol    string    `json:"symbol"` // e.g., "kg"
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
