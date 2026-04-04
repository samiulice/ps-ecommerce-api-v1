package model

import (
	"database/sql"
	"time"
)

type Supplier struct {
	ID                 int64           `json:"id" db:"id"`
	SupplierCode       string          `json:"supplier_code" db:"supplier_code"`
	Name               string          `json:"name" db:"name"`
	CompanyName        sql.NullString  `json:"company_name,omitempty" db:"company_name"`
	ContactPerson      sql.NullString  `json:"contact_person,omitempty" db:"contact_person"`
	Phone              string          `json:"phone" db:"phone"`
	Email              sql.NullString  `json:"email,omitempty" db:"email"`
	Website            sql.NullString  `json:"website,omitempty" db:"website"`
	TaxID              sql.NullString  `json:"tax_id,omitempty" db:"tax_id"`
	TradeLicenseNo     sql.NullString  `json:"trade_license_no,omitempty" db:"trade_license_no"`
	PaymentTerms       sql.NullString  `json:"payment_terms,omitempty" db:"payment_terms"`
	CreditLimit        sql.NullFloat64 `json:"credit_limit,omitempty" db:"credit_limit"`
	OutstandingBalance sql.NullFloat64 `json:"outstanding_balance,omitempty" db:"outstanding_balance"`
	LeadTimeDays       sql.NullInt32   `json:"lead_time_days,omitempty" db:"lead_time_days"`
	Rating             sql.NullFloat64 `json:"rating,omitempty" db:"rating"`
	StreetAddress      sql.NullString  `json:"street_address,omitempty" db:"street_address"`
	Country            sql.NullString  `json:"country,omitempty" db:"country"`
	City               sql.NullString  `json:"city,omitempty" db:"city"`
	Zip                sql.NullString  `json:"zip,omitempty" db:"zip"`
	Notes              sql.NullString  `json:"notes,omitempty" db:"notes"`
	IsActive           bool            `json:"is_active" db:"is_active"`
	CreatedAt          sql.NullTime    `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt          sql.NullTime    `json:"updated_at,omitempty" db:"updated_at"`
}

type SupplierFilter struct {
	Search             string `json:"search" query:"search"`
	CheckAccountStatus bool   `json:"check_account_status" query:"check_account_status"`
	IsActive           bool   `json:"is_active" query:"is_active"`
	Page               int    `json:"page" query:"page"`
	Limit              int    `json:"limit" query:"limit"`
}

type SupplierCreateRequest struct {
	SupplierCode       string  `json:"supplier_code"`
	Name               string  `json:"name"`
	CompanyName        string  `json:"company_name"`
	ContactPerson      string  `json:"contact_person"`
	Phone              string  `json:"phone"`
	Email              string  `json:"email"`
	Website            string  `json:"website"`
	TaxID              string  `json:"tax_id"`
	TradeLicenseNo     string  `json:"trade_license_no"`
	PaymentTerms       string  `json:"payment_terms"`
	CreditLimit        float64 `json:"credit_limit"`
	OutstandingBalance float64 `json:"outstanding_balance"`
	LeadTimeDays       int     `json:"lead_time_days"`
	Rating             float64 `json:"rating"`
	StreetAddress      string  `json:"street_address"`
	Country            string  `json:"country"`
	City               string  `json:"city"`
	Zip                string  `json:"zip"`
	Notes              string  `json:"notes"`
	IsActive           *bool   `json:"is_active"`
}

type SupplierUpdateRequest struct {
	SupplierCode       string  `json:"supplier_code"`
	Name               string  `json:"name"`
	CompanyName        string  `json:"company_name"`
	ContactPerson      string  `json:"contact_person"`
	Phone              string  `json:"phone"`
	Email              string  `json:"email"`
	Website            string  `json:"website"`
	TaxID              string  `json:"tax_id"`
	TradeLicenseNo     string  `json:"trade_license_no"`
	PaymentTerms       string  `json:"payment_terms"`
	CreditLimit        float64 `json:"credit_limit"`
	OutstandingBalance float64 `json:"outstanding_balance"`
	LeadTimeDays       int     `json:"lead_time_days"`
	Rating             float64 `json:"rating"`
	StreetAddress      string  `json:"street_address"`
	Country            string  `json:"country"`
	City               string  `json:"city"`
	Zip                string  `json:"zip"`
	Notes              string  `json:"notes"`
	IsActive           *bool   `json:"is_active"`
}

type SupplierResponse struct {
	ID                 int64   `json:"id"`
	SupplierCode       string  `json:"supplier_code"`
	Name               string  `json:"name"`
	CompanyName        string  `json:"company_name,omitempty"`
	ContactPerson      string  `json:"contact_person,omitempty"`
	Phone              string  `json:"phone"`
	Email              string  `json:"email,omitempty"`
	Website            string  `json:"website,omitempty"`
	TaxID              string  `json:"tax_id,omitempty"`
	TradeLicenseNo     string  `json:"trade_license_no,omitempty"`
	PaymentTerms       string  `json:"payment_terms,omitempty"`
	CreditLimit        float64 `json:"credit_limit"`
	OutstandingBalance float64 `json:"outstanding_balance"`
	LeadTimeDays       int32   `json:"lead_time_days"`
	Rating             float64 `json:"rating"`
	StreetAddress      string  `json:"street_address,omitempty"`
	Country            string  `json:"country,omitempty"`
	City               string  `json:"city,omitempty"`
	Zip                string  `json:"zip,omitempty"`
	Notes              string  `json:"notes,omitempty"`
	IsActive           bool    `json:"is_active"`
	CreatedAt          string  `json:"created_at,omitempty"`
}

func (s *Supplier) ToResponse() *SupplierResponse {
	resp := &SupplierResponse{
		ID:           s.ID,
		SupplierCode: s.SupplierCode,
		Name:         s.Name,
		Phone:        s.Phone,
		IsActive:     s.IsActive,
	}

	if s.CompanyName.Valid {
		resp.CompanyName = s.CompanyName.String
	}
	if s.ContactPerson.Valid {
		resp.ContactPerson = s.ContactPerson.String
	}
	if s.Email.Valid {
		resp.Email = s.Email.String
	}
	if s.Website.Valid {
		resp.Website = s.Website.String
	}
	if s.TaxID.Valid {
		resp.TaxID = s.TaxID.String
	}
	if s.TradeLicenseNo.Valid {
		resp.TradeLicenseNo = s.TradeLicenseNo.String
	}
	if s.PaymentTerms.Valid {
		resp.PaymentTerms = s.PaymentTerms.String
	}
	if s.CreditLimit.Valid {
		resp.CreditLimit = s.CreditLimit.Float64
	}
	if s.OutstandingBalance.Valid {
		resp.OutstandingBalance = s.OutstandingBalance.Float64
	}
	if s.LeadTimeDays.Valid {
		resp.LeadTimeDays = s.LeadTimeDays.Int32
	}
	if s.Rating.Valid {
		resp.Rating = s.Rating.Float64
	}
	if s.StreetAddress.Valid {
		resp.StreetAddress = s.StreetAddress.String
	}
	if s.Country.Valid {
		resp.Country = s.Country.String
	}
	if s.City.Valid {
		resp.City = s.City.String
	}
	if s.Zip.Valid {
		resp.Zip = s.Zip.String
	}
	if s.Notes.Valid {
		resp.Notes = s.Notes.String
	}
	if s.CreatedAt.Valid {
		resp.CreatedAt = s.CreatedAt.Time.Format(time.RFC3339)
	}

	return resp
}
