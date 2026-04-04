package model

import "time"

type Branch struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Address   string    `json:"address"`
	Mobile    string    `json:"mobile"`
	Telephone string    `json:"telephone"`
	Email     string    `json:"email"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
