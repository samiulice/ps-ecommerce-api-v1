package model

import (
	"database/sql"
	"time"
)

type Purchase struct {
	ID                          int64                `json:"id" db:"id"`
	PurchaseDate                time.Time            `json:"purchase_date" db:"purchase_date"`
	PrefixCode                  sql.NullString       `json:"prefix_code,omitempty" db:"prefix_code"`
	CountID                     sql.NullString       `json:"count_id,omitempty" db:"count_id"`
	PurchaseCode                string               `json:"purchase_code" db:"purchase_code"`
	ReferenceNo                 sql.NullString       `json:"reference_no,omitempty" db:"reference_no"`
	PurchaseOrderID             sql.NullInt64        `json:"purchase_order_id,omitempty" db:"purchase_order_id"`
	PartyID                     int64                `json:"party_id" db:"party_id"`
	SupplierName                sql.NullString       `json:"supplier_name,omitempty" db:"supplier_name"`
	StateID                     sql.NullInt64        `json:"state_id,omitempty" db:"state_id"`
	CarrierID                   sql.NullInt64        `json:"carrier_id,omitempty" db:"carrier_id"`
	Note                        sql.NullString       `json:"note,omitempty" db:"note"`
	ShippingCharge              sql.NullFloat64      `json:"shipping_charge,omitempty" db:"shipping_charge"`
	IsShippingChargeDistributed bool                 `json:"is_shipping_charge_distributed" db:"is_shipping_charge_distributed"`
	RoundOff                    sql.NullFloat64      `json:"round_off,omitempty" db:"round_off"`
	GrandTotal                  sql.NullFloat64      `json:"grand_total,omitempty" db:"grand_total"`
	ChangeReturn                sql.NullInt32        `json:"change_return,omitempty" db:"change_return"`
	PaidAmount                  sql.NullFloat64      `json:"paid_amount,omitempty" db:"paid_amount"`
	CurrencyID                  sql.NullInt64        `json:"currency_id,omitempty" db:"currency_id"`
	ExchangeRate                sql.NullFloat64      `json:"exchange_rate,omitempty" db:"exchange_rate"`
	CreatedBy                   sql.NullInt64        `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy                   sql.NullInt64        `json:"updated_by,omitempty" db:"updated_by"`
	Items                       []PurchaseItem       `json:"items,omitempty"`
	Attachments                 []PurchaseAttachment `json:"attachments,omitempty"`
	RemoveAttachmentIDs         []int64              `json:"remove_attachment_ids,omitempty"`
	CreatedAt                   sql.NullTime         `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt                   sql.NullTime         `json:"updated_at,omitempty" db:"updated_at"`
}

type PurchaseItem struct {
	ID         int64          `json:"id" db:"id"`
	PurchaseID int64          `json:"purchase_id" db:"purchase_id"`
	ItemType   string         `json:"item_type" db:"item_type"`
	ProductID  sql.NullInt64  `json:"product_id,omitempty" db:"product_id"`
	ItemName   string         `json:"item_name" db:"item_name"`
	Quantity   float64        `json:"quantity" db:"quantity"`
	UnitPrice  float64        `json:"unit_price" db:"unit_price"`
	TotalPrice float64        `json:"total_price" db:"total_price"`
	Note       sql.NullString `json:"note,omitempty" db:"note"`
	CreatedAt  sql.NullTime   `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  sql.NullTime   `json:"updated_at,omitempty" db:"updated_at"`
}

type PurchaseAttachment struct {
	ID         int64        `json:"id" db:"id"`
	PurchaseID int64        `json:"purchase_id" db:"purchase_id"`
	FileURL    string       `json:"file_url" db:"file_url"`
	FileName   string       `json:"file_name" db:"file_name"`
	FileExt    string       `json:"file_ext" db:"file_ext"`
	MimeType   string       `json:"mime_type" db:"mime_type"`
	FileSize   int64        `json:"file_size" db:"file_size"`
	CreatedAt  sql.NullTime `json:"created_at,omitempty" db:"created_at"`
}

type PurchaseFilter struct {
	Search         string `json:"search" query:"search"`
	PartyID        int64  `json:"party_id" query:"party_id"`
	FromDate       string `json:"from_date" query:"from_date"`
	ToDate         string `json:"to_date" query:"to_date"`
	Page           int    `json:"page" query:"page"`
	Limit          int    `json:"limit" query:"limit"`
	HasSupplierSet bool   `json:"-"`
}

type PurchaseCreateRequest struct {
	PurchaseDate                string              `json:"purchase_date"`
	PrefixCode                  string              `json:"prefix_code"`
	CountID                     string              `json:"count_id"`
	PurchaseCode                string              `json:"purchase_code"`
	ReferenceNo                 string              `json:"reference_no"`
	PurchaseOrderID             *int64              `json:"purchase_order_id"`
	PartyID                     int64               `json:"party_id"`
	StateID                     *int64              `json:"state_id"`
	CarrierID                   *int64              `json:"carrier_id"`
	Note                        string              `json:"note"`
	ShippingCharge              float64             `json:"shipping_charge"`
	IsShippingChargeDistributed bool                `json:"is_shipping_charge_distributed"`
	RoundOff                    float64             `json:"round_off"`
	GrandTotal                  float64             `json:"grand_total"`
	ChangeReturn                *int                `json:"change_return"`
	PaidAmount                  float64             `json:"paid_amount"`
	CurrencyID                  *int64              `json:"currency_id"`
	ExchangeRate                float64             `json:"exchange_rate"`
	CreatedBy                   *int64              `json:"created_by"`
	UpdatedBy                   *int64              `json:"updated_by"`
	Items                       []PurchaseItemInput `json:"items"`
}

type PurchaseUpdateRequest struct {
	PurchaseDate                string              `json:"purchase_date"`
	PrefixCode                  string              `json:"prefix_code"`
	CountID                     string              `json:"count_id"`
	PurchaseCode                string              `json:"purchase_code"`
	ReferenceNo                 string              `json:"reference_no"`
	PurchaseOrderID             *int64              `json:"purchase_order_id"`
	PartyID                     int64               `json:"party_id"`
	StateID                     *int64              `json:"state_id"`
	CarrierID                   *int64              `json:"carrier_id"`
	Note                        string              `json:"note"`
	ShippingCharge              float64             `json:"shipping_charge"`
	IsShippingChargeDistributed bool                `json:"is_shipping_charge_distributed"`
	RoundOff                    float64             `json:"round_off"`
	GrandTotal                  float64             `json:"grand_total"`
	ChangeReturn                *int                `json:"change_return"`
	PaidAmount                  float64             `json:"paid_amount"`
	CurrencyID                  *int64              `json:"currency_id"`
	ExchangeRate                float64             `json:"exchange_rate"`
	CreatedBy                   *int64              `json:"created_by"`
	UpdatedBy                   *int64              `json:"updated_by"`
	Items                       []PurchaseItemInput `json:"items"`
	RemoveAttachmentIDs         []int64             `json:"remove_attachment_ids"`
}

type PurchaseItemInput struct {
	ItemType   string  `json:"item_type"`
	ProductID  *int64  `json:"product_id,omitempty"`
	ItemName   string  `json:"item_name"`
	Quantity   float64 `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
	Note       string  `json:"note"`
}

type PurchaseResponse struct {
	ID                          int64                        `json:"id"`
	PurchaseDate                string                       `json:"purchase_date"`
	PrefixCode                  string                       `json:"prefix_code,omitempty"`
	CountID                     string                       `json:"count_id,omitempty"`
	PurchaseCode                string                       `json:"purchase_code"`
	ReferenceNo                 string                       `json:"reference_no,omitempty"`
	PurchaseOrderID             *int64                       `json:"purchase_order_id,omitempty"`
	PartyID                     int64                        `json:"party_id"`
	SupplierName                string                       `json:"supplier_name,omitempty"`
	StateID                     *int64                       `json:"state_id,omitempty"`
	CarrierID                   *int64                       `json:"carrier_id,omitempty"`
	Note                        string                       `json:"note,omitempty"`
	ShippingCharge              float64                      `json:"shipping_charge"`
	IsShippingChargeDistributed bool                         `json:"is_shipping_charge_distributed"`
	RoundOff                    float64                      `json:"round_off"`
	GrandTotal                  float64                      `json:"grand_total"`
	ChangeReturn                *int32                       `json:"change_return,omitempty"`
	PaidAmount                  float64                      `json:"paid_amount"`
	CurrencyID                  *int64                       `json:"currency_id,omitempty"`
	ExchangeRate                float64                      `json:"exchange_rate"`
	CreatedBy                   *int64                       `json:"created_by,omitempty"`
	UpdatedBy                   *int64                       `json:"updated_by,omitempty"`
	Items                       []PurchaseItemResponse       `json:"items,omitempty"`
	Attachments                 []PurchaseAttachmentResponse `json:"attachments,omitempty"`
	CreatedAt                   string                       `json:"created_at,omitempty"`
	UpdatedAt                   string                       `json:"updated_at,omitempty"`
}

type PurchaseItemResponse struct {
	ID         int64   `json:"id"`
	ItemType   string  `json:"item_type"`
	ProductID  *int64  `json:"product_id,omitempty"`
	ItemName   string  `json:"item_name"`
	Quantity   float64 `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
	Note       string  `json:"note,omitempty"`
	CreatedAt  string  `json:"created_at,omitempty"`
	UpdatedAt  string  `json:"updated_at,omitempty"`
}

type PurchaseAttachmentResponse struct {
	ID        int64  `json:"id"`
	FileURL   string `json:"file_url"`
	FileName  string `json:"file_name"`
	FileExt   string `json:"file_ext"`
	MimeType  string `json:"mime_type"`
	FileSize  int64  `json:"file_size"`
	CreatedAt string `json:"created_at,omitempty"`
}

func (p *Purchase) ToResponse() *PurchaseResponse {
	resp := &PurchaseResponse{
		ID:                          p.ID,
		PurchaseDate:                p.PurchaseDate.Format("2006-01-02"),
		PurchaseCode:                p.PurchaseCode,
		PartyID:                     p.PartyID,
		IsShippingChargeDistributed: p.IsShippingChargeDistributed,
	}

	if p.PrefixCode.Valid {
		resp.PrefixCode = p.PrefixCode.String
	}
	if p.CountID.Valid {
		resp.CountID = p.CountID.String
	}
	if p.ReferenceNo.Valid {
		resp.ReferenceNo = p.ReferenceNo.String
	}
	if p.PurchaseOrderID.Valid {
		resp.PurchaseOrderID = &p.PurchaseOrderID.Int64
	}
	if p.SupplierName.Valid {
		resp.SupplierName = p.SupplierName.String
	}
	if p.StateID.Valid {
		resp.StateID = &p.StateID.Int64
	}
	if p.CarrierID.Valid {
		resp.CarrierID = &p.CarrierID.Int64
	}
	if p.Note.Valid {
		resp.Note = p.Note.String
	}
	if p.ShippingCharge.Valid {
		resp.ShippingCharge = p.ShippingCharge.Float64
	}
	if p.RoundOff.Valid {
		resp.RoundOff = p.RoundOff.Float64
	}
	if p.GrandTotal.Valid {
		resp.GrandTotal = p.GrandTotal.Float64
	}
	if p.ChangeReturn.Valid {
		resp.ChangeReturn = &p.ChangeReturn.Int32
	}
	if p.PaidAmount.Valid {
		resp.PaidAmount = p.PaidAmount.Float64
	}
	if p.CurrencyID.Valid {
		resp.CurrencyID = &p.CurrencyID.Int64
	}
	if p.ExchangeRate.Valid {
		resp.ExchangeRate = p.ExchangeRate.Float64
	}
	if p.CreatedBy.Valid {
		resp.CreatedBy = &p.CreatedBy.Int64
	}
	if p.UpdatedBy.Valid {
		resp.UpdatedBy = &p.UpdatedBy.Int64
	}
	if p.CreatedAt.Valid {
		resp.CreatedAt = p.CreatedAt.Time.Format(time.RFC3339)
	}
	if p.UpdatedAt.Valid {
		resp.UpdatedAt = p.UpdatedAt.Time.Format(time.RFC3339)
	}

	if len(p.Items) > 0 {
		resp.Items = make([]PurchaseItemResponse, 0, len(p.Items))
		for _, item := range p.Items {
			row := PurchaseItemResponse{
				ID:         item.ID,
				ItemType:   item.ItemType,
				ItemName:   item.ItemName,
				Quantity:   item.Quantity,
				UnitPrice:  item.UnitPrice,
				TotalPrice: item.TotalPrice,
			}
			if item.ProductID.Valid {
				row.ProductID = &item.ProductID.Int64
			}
			if item.Note.Valid {
				row.Note = item.Note.String
			}
			if item.CreatedAt.Valid {
				row.CreatedAt = item.CreatedAt.Time.Format(time.RFC3339)
			}
			if item.UpdatedAt.Valid {
				row.UpdatedAt = item.UpdatedAt.Time.Format(time.RFC3339)
			}
			resp.Items = append(resp.Items, row)
		}
	}

	if len(p.Attachments) > 0 {
		resp.Attachments = make([]PurchaseAttachmentResponse, 0, len(p.Attachments))
		for _, a := range p.Attachments {
			row := PurchaseAttachmentResponse{
				ID:       a.ID,
				FileURL:  a.FileURL,
				FileName: a.FileName,
				FileExt:  a.FileExt,
				MimeType: a.MimeType,
				FileSize: a.FileSize,
			}
			if a.CreatedAt.Valid {
				row.CreatedAt = a.CreatedAt.Time.Format(time.RFC3339)
			}
			resp.Attachments = append(resp.Attachments, row)
		}
	}

	return resp
}
