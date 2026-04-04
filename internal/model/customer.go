package model

import (
	"database/sql"
	"time"
)

// Customer represents a customer customer in the system.
type Customer struct {
	ID                    int64           `json:"id" db:"id"`
	Name                  sql.NullString  `json:"name,omitempty" db:"name"`
	FName                 sql.NullString  `json:"f_name,omitempty" db:"f_name"`
	LName                 sql.NullString  `json:"l_name,omitempty" db:"l_name"`
	Phone                 string          `json:"phone" db:"phone"`
	Image                 string          `json:"image" db:"image"`
	Email                 sql.NullString  `json:"email,omitempty" db:"email"`
	IsRetailer            bool            `json:"is_retailer" db:"is_retailer"`
	EmailVerifiedAt       sql.NullTime    `json:"email_verified_at,omitempty" db:"email_verified_at"`
	Password              string          `json:"-" db:"password"`
	RememberToken         sql.NullString  `json:"-" db:"remember_token"`
	CreatedAt             sql.NullTime    `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt             sql.NullTime    `json:"updated_at,omitempty" db:"updated_at"`
	StreetAddress         sql.NullString  `json:"street_address,omitempty" db:"street_address"`
	Country               sql.NullString  `json:"country,omitempty" db:"country"`
	City                  sql.NullString  `json:"city,omitempty" db:"city"`
	Zip                   sql.NullString  `json:"zip,omitempty" db:"zip"`
	HouseNo               sql.NullString  `json:"house_no,omitempty" db:"house_no"`
	ApartmentNo           sql.NullString  `json:"apartment_no,omitempty" db:"apartment_no"`
	CMFirebaseToken       sql.NullString  `json:"cm_firebase_token,omitempty" db:"cm_firebase_token"`
	IsActive              bool            `json:"is_active" db:"is_active"`
	PaymentCardLastFour   sql.NullString  `json:"payment_card_last_four,omitempty" db:"payment_card_last_four"`
	PaymentCardBrand      sql.NullString  `json:"payment_card_brand,omitempty" db:"payment_card_brand"`
	PaymentCardFawryToken sql.NullString  `json:"-" db:"payment_card_fawry_token"`
	LoginMedium           sql.NullString  `json:"login_medium,omitempty" db:"login_medium"`
	SocialID              sql.NullString  `json:"social_id,omitempty" db:"social_id"`
	IsPhoneVerified       bool            `json:"is_phone_verified" db:"is_phone_verified"`
	TemporaryToken        sql.NullString  `json:"-" db:"temporary_token"`
	IsEmailVerified       bool            `json:"is_email_verified" db:"is_email_verified"`
	WalletBalance         sql.NullFloat64 `json:"wallet_balance,omitempty" db:"wallet_balance"`
	LoyaltyPoint          sql.NullFloat64 `json:"loyalty_point,omitempty" db:"loyalty_point"`
	LoginHitCount         int16           `json:"login_hit_count" db:"login_hit_count"`
	IsTempBlocked         bool            `json:"is_temp_blocked" db:"is_temp_blocked"`
	TempBlockTime         sql.NullTime    `json:"temp_block_time,omitempty" db:"temp_block_time"`
	ReferralCode          sql.NullString  `json:"referral_code,omitempty" db:"referral_code"`
	ReferredBy            sql.NullInt32   `json:"referred_by,omitempty" db:"referred_by"`
	AppLanguage           string          `json:"app_language" db:"app_language"`
}

// CustomerFilter defines the criteria for querying customers.
type CustomerFilter struct {
	// Search matches against name, phone, or email
	Search string `json:"search" query:"search"`

	// Status matches against the account status
	CheckAccountStatus bool `json:"check_account_status" query:"check_account_status"`
	IsActive           bool `json:"is_active" query:"is_active"`

	// Page number for pagination (starts at 1)
	Page int `json:"page" query:"page"`

	// Limit is the number of items per page
	Limit int `json:"limit" query:"limit"`
}

// CustomerCreateRequest represents the payload for creating a new customer.
type CustomerCreateRequest struct {
	Name       string `json:"name"`
	FName      string `json:"f_name"`
	LName      string `json:"l_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsRetailer bool   `json:"is_retailer"`
}

// CustomerUpdateRequest represents the payload for updating a customer.
type CustomerUpdateRequest struct {
	Name          string `json:"name,omitempty"`
	FName         string `json:"f_name,omitempty"`
	LName         string `json:"l_name,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Image         string `json:"image,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
	Country       string `json:"country,omitempty"`
	City          string `json:"city,omitempty"`
	Zip           string `json:"zip,omitempty"`
	HouseNo       string `json:"house_no,omitempty"`
	ApartmentNo   string `json:"apartment_no,omitempty"`
	AppLanguage   string `json:"app_language,omitempty"`
}

// CustomerResponse is a sanitized customer response (no sensitive fields).
type CustomerResponse struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name,omitempty"`
	FName           string  `json:"f_name,omitempty"`
	LName           string  `json:"l_name,omitempty"`
	Phone           string  `json:"phone"`
	Image           string  `json:"image"`
	Email           string  `json:"email,omitempty"`
	IsRetailer      bool    `json:"is_retailer"`
	IsActive        bool    `json:"is_active"`
	IsPhoneVerified bool    `json:"is_phone_verified"`
	IsEmailVerified bool    `json:"is_email_verified"`
	WalletBalance   float64 `json:"wallet_balance"`
	LoyaltyPoint    float64 `json:"loyalty_point"`
	ReferralCode    string  `json:"referral_code,omitempty"`
	AppLanguage     string  `json:"app_language"`
	CreatedAt       string  `json:"created_at,omitempty"`
	// Address info
	StreetAddress string `json:"street_address,omitempty"`
	Country       string `json:"country,omitempty"`
	City          string `json:"city,omitempty"`
	Zip           string `json:"zip,omitempty"`
	HouseNo       string `json:"house_no,omitempty"`
	ApartmentNo   string `json:"apartment_no,omitempty"`
}

// ToResponse converts Customer to a safe CustomerResponse.
func (u *Customer) ToResponse() *CustomerResponse {
	resp := &CustomerResponse{
		ID:              u.ID,
		Phone:           u.Phone,
		IsRetailer:      u.IsRetailer,
		Image:           u.Image,
		IsActive:        u.IsActive,
		IsPhoneVerified: u.IsPhoneVerified,
		IsEmailVerified: u.IsEmailVerified,
		AppLanguage:     u.AppLanguage,
	}
	if u.Name.Valid {
		resp.Name = u.Name.String
	}
	if u.FName.Valid {
		resp.FName = u.FName.String
	}
	if u.LName.Valid {
		resp.LName = u.LName.String
	}
	if u.Email.Valid {
		resp.Email = u.Email.String
	}
	if u.WalletBalance.Valid {
		resp.WalletBalance = u.WalletBalance.Float64
	}
	if u.LoyaltyPoint.Valid {
		resp.LoyaltyPoint = u.LoyaltyPoint.Float64
	}
	if u.ReferralCode.Valid {
		resp.ReferralCode = u.ReferralCode.String
	}
	if u.CreatedAt.Valid {
		resp.CreatedAt = u.CreatedAt.Time.Format(time.RFC3339)
	}
	// Address info
	if u.StreetAddress.Valid {
		resp.StreetAddress = u.StreetAddress.String
	}
	if u.Country.Valid {
		resp.Country = u.Country.String
	}
	if u.City.Valid {
		resp.City = u.City.String
	}
	if u.Zip.Valid {
		resp.Zip = u.Zip.String
	}
	if u.HouseNo.Valid {
		resp.HouseNo = u.HouseNo.String
	}
	if u.ApartmentNo.Valid {
		resp.ApartmentNo = u.ApartmentNo.String
	}
	return resp
}

// ToNullString converts a string to sql.NullString.
// Empty strings result in NULL in the database.
func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// ToNullInt32 converts an int to sql.NullInt32.
// Zero values result in NULL in the database.
func ToNullInt32(i int) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(i), Valid: true}
}

// ToNullFloat64 converts a float64 to sql.NullFloat64.
// Zero values result in NULL in the database.
func ToNullFloat64(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}
