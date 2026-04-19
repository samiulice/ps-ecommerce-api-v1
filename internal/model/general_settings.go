package model

import "time"

type GeneralSettings struct {
	ID                int64     `json:"id" db:"id"`
	CompanyName       string    `json:"company_name" db:"company_name"`
	CompanyLogo       string    `json:"company_logo" db:"company_logo"`
	CompanyLogoMobile string    `json:"company_logo_mobile" db:"company_logo_mobile"`
	CompanyLogoFooter string    `json:"company_logo_footer" db:"company_logo_footer"`
	CompanyEmail      string    `json:"company_email" db:"company_email"`
	CompanyPhone      string    `json:"company_phone" db:"company_phone"`
	CompanyAddress    string    `json:"company_address" db:"company_address"`
	AboutUsShort      string    `json:"about_us_short" db:"about_us_short"`
	AboutUsLong       string    `json:"about_us_long" db:"about_us_long"`
	CurrencySymbol    string    `json:"currency_symbol" db:"currency_symbol"`
	CurrencyCode      string    `json:"currency_code" db:"currency_code"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}
