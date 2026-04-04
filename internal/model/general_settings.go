package model

import "time"

type GeneralSettings struct {
	ID             int64     `json:"id" db:"id"`
	CompanyName    string    `json:"company_name" db:"company_name"`
	CompanyLogo    string    `json:"company_logo" db:"company_logo"`
	CompanyAddress string    `json:"company_address" db:"company_address"`
	CurrencySymbol string    `json:"currency_symbol" db:"currency_symbol"`
	CurrencyCode   string    `json:"currency_code" db:"currency_code"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
