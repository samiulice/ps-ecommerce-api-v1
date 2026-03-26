package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type PurchaseHandler struct {
	svc *service.PurchaseService
}

func NewPurchaseHandler(svc *service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{svc: svc}
}

func (h *PurchaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	purchase, headers, err := parseCreatePurchaseRequest(w, r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	resPurchase, err := h.svc.Create(r.Context(), purchase, headers)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error    bool                    `json:"error"`
		Message  string                  `json:"message"`
		Purchase *model.PurchaseResponse `json:"purchase"`
	}{
		Error:    false,
		Message:  "Purchase created successfully",
		Purchase: resPurchase.ToResponse(),
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}

func (h *PurchaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid purchase ID"))
		return
	}

	purchase, headers, err := parseUpdatePurchaseRequest(w, r, id)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	resPurchase, err := h.svc.Update(r.Context(), purchase, headers)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error    bool                    `json:"error"`
		Message  string                  `json:"message"`
		Purchase *model.PurchaseResponse `json:"purchase"`
	}{
		Error:    false,
		Message:  "Purchase updated successfully",
		Purchase: resPurchase.ToResponse(),
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *PurchaseHandler) ListPurchases(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filter := model.PurchaseFilter{
		Search:   query.Get("search"),
		FromDate: query.Get("from_date"),
		ToDate:   query.Get("to_date"),
	}

	if partyID := query.Get("party_id"); partyID != "" {
		id, err := strconv.ParseInt(partyID, 10, 64)
		if err != nil || id <= 0 {
			utils.BadRequest(w, errors.New("invalid supplier ID"))
			return
		}
		filter.PartyID = id
		filter.HasSupplierSet = true
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	filter.Page = page

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	filter.Limit = limit

	purchases, total, err := h.svc.ListPurchases(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	utils.OK(w, "Purchases retrieved successfully", map[string]interface{}{
		"purchases":   purchases,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

func (h *PurchaseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid purchase ID"))
		return
	}

	purchase, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error    bool                    `json:"error"`
		Message  string                  `json:"message"`
		Purchase *model.PurchaseResponse `json:"purchase"`
	}{
		Error:    false,
		Message:  "Purchase retrieved successfully",
		Purchase: purchase.ToResponse(),
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *PurchaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(w, errors.New("invalid purchase ID"))
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Purchase deleted successfully",
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *PurchaseHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrPurchaseNotFound):
		utils.NotFound(w, err)
	case errors.Is(err, service.ErrPurchaseDateRequired),
		errors.Is(err, service.ErrPurchaseCodeRequired),
		errors.Is(err, service.ErrPurchaseSupplierRequired),
		errors.Is(err, service.ErrPurchaseItemsRequired),
		errors.Is(err, service.ErrPurchaseItemInvalid),
		errors.Is(err, service.ErrPurchaseGrandTotalInvalid),
		errors.Is(err, service.ErrPurchasePaidAmountInvalid),
		errors.Is(err, service.ErrPurchaseShippingInvalid),
		errors.Is(err, service.ErrPurchaseRoundOffInvalid),
		errors.Is(err, service.ErrPurchaseExchangeInvalid):
		utils.BadRequest(w, err)
	case errors.Is(err, service.ErrPurchaseCodeDuplicate):
		utils.WriteJSON(w, http.StatusConflict, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
	default:
		fmt.Println("purchase error:", err)
		utils.ServerError(w, err)
	}
}

func mapCreateRequestToPurchase(req model.PurchaseCreateRequest) (*model.Purchase, error) {
	purchaseDate, err := parsePurchaseDate(req.PurchaseDate)
	if err != nil {
		return nil, err
	}

	return &model.Purchase{
		PurchaseDate:                purchaseDate,
		PrefixCode:                  model.ToNullString(strings.TrimSpace(req.PrefixCode)),
		CountID:                     model.ToNullString(strings.TrimSpace(req.CountID)),
		PurchaseCode:                strings.TrimSpace(req.PurchaseCode),
		ReferenceNo:                 model.ToNullString(strings.TrimSpace(req.ReferenceNo)),
		PurchaseOrderID:             toNullInt64Ptr(req.PurchaseOrderID),
		PartyID:                     req.PartyID,
		StateID:                     toNullInt64Ptr(req.StateID),
		CarrierID:                   toNullInt64Ptr(req.CarrierID),
		Note:                        model.ToNullString(strings.TrimSpace(req.Note)),
		ShippingCharge:              toSQLFloat64Purchase(req.ShippingCharge),
		IsShippingChargeDistributed: req.IsShippingChargeDistributed,
		RoundOff:                    toSQLFloat64Purchase(req.RoundOff),
		GrandTotal:                  toSQLFloat64Purchase(req.GrandTotal),
		ChangeReturn:                toNullInt32Ptr(req.ChangeReturn),
		PaidAmount:                  toSQLFloat64Purchase(req.PaidAmount),
		CurrencyID:                  toNullInt64Ptr(req.CurrencyID),
		ExchangeRate:                toSQLFloat64Purchase(req.ExchangeRate),
		CreatedBy:                   toNullInt64Ptr(req.CreatedBy),
		UpdatedBy:                   toNullInt64Ptr(req.UpdatedBy),
		Items:                       mapPurchaseItems(req.Items),
	}, nil
}

func mapUpdateRequestToPurchase(id int64, req model.PurchaseUpdateRequest) (*model.Purchase, error) {
	purchaseDate, err := parsePurchaseDate(req.PurchaseDate)
	if err != nil {
		return nil, err
	}

	return &model.Purchase{
		ID:                          id,
		PurchaseDate:                purchaseDate,
		PrefixCode:                  model.ToNullString(strings.TrimSpace(req.PrefixCode)),
		CountID:                     model.ToNullString(strings.TrimSpace(req.CountID)),
		PurchaseCode:                strings.TrimSpace(req.PurchaseCode),
		ReferenceNo:                 model.ToNullString(strings.TrimSpace(req.ReferenceNo)),
		PurchaseOrderID:             toNullInt64Ptr(req.PurchaseOrderID),
		PartyID:                     req.PartyID,
		StateID:                     toNullInt64Ptr(req.StateID),
		CarrierID:                   toNullInt64Ptr(req.CarrierID),
		Note:                        model.ToNullString(strings.TrimSpace(req.Note)),
		ShippingCharge:              toSQLFloat64Purchase(req.ShippingCharge),
		IsShippingChargeDistributed: req.IsShippingChargeDistributed,
		RoundOff:                    toSQLFloat64Purchase(req.RoundOff),
		GrandTotal:                  toSQLFloat64Purchase(req.GrandTotal),
		ChangeReturn:                toNullInt32Ptr(req.ChangeReturn),
		PaidAmount:                  toSQLFloat64Purchase(req.PaidAmount),
		CurrencyID:                  toNullInt64Ptr(req.CurrencyID),
		ExchangeRate:                toSQLFloat64Purchase(req.ExchangeRate),
		CreatedBy:                   toNullInt64Ptr(req.CreatedBy),
		UpdatedBy:                   toNullInt64Ptr(req.UpdatedBy),
		Items:                       mapPurchaseItems(req.Items),
		RemoveAttachmentIDs:         req.RemoveAttachmentIDs,
	}, nil
}

func parsePurchaseDate(v string) (time.Time, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return time.Time{}, errors.New("purchase date is required")
	}

	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return time.Time{}, errors.New("purchase date must be in YYYY-MM-DD format")
	}
	return t, nil
}

func toNullInt64Ptr(v *int64) sql.NullInt64 {
	if v == nil || *v <= 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *v, Valid: true}
}

func toNullInt32Ptr(v *int) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*v), Valid: true}
}

func toSQLFloat64Purchase(v float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: v, Valid: true}
}

func parseCreatePurchaseRequest(w http.ResponseWriter, r *http.Request) (*model.Purchase, []*multipart.FileHeader, error) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			return nil, nil, err
		}
		return mapCreateFormToPurchase(r)
	}

	var req model.PurchaseCreateRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		return nil, nil, err
	}
	p, err := mapCreateRequestToPurchase(req)
	return p, nil, err
}

func parseUpdatePurchaseRequest(w http.ResponseWriter, r *http.Request, id int64) (*model.Purchase, []*multipart.FileHeader, error) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			return nil, nil, err
		}
		return mapUpdateFormToPurchase(r, id)
	}

	var req model.PurchaseUpdateRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		return nil, nil, err
	}
	p, err := mapUpdateRequestToPurchase(id, req)
	return p, nil, err
}

func mapCreateFormToPurchase(r *http.Request) (*model.Purchase, []*multipart.FileHeader, error) {
	date, err := parsePurchaseDate(r.FormValue("purchase_date"))
	if err != nil {
		return nil, nil, err
	}

	partyID, err := parseInt64Required(r.FormValue("party_id"), "supplier is required")
	if err != nil {
		return nil, nil, err
	}

	attachments := getAttachmentHeaders(r)

	return &model.Purchase{
		PurchaseDate:                date,
		PrefixCode:                  model.ToNullString(strings.TrimSpace(r.FormValue("prefix_code"))),
		CountID:                     model.ToNullString(strings.TrimSpace(r.FormValue("count_id"))),
		PurchaseCode:                strings.TrimSpace(r.FormValue("purchase_code")),
		ReferenceNo:                 model.ToNullString(strings.TrimSpace(r.FormValue("reference_no"))),
		PurchaseOrderID:             toNullInt64String(r.FormValue("purchase_order_id")),
		PartyID:                     partyID,
		StateID:                     toNullInt64String(r.FormValue("state_id")),
		CarrierID:                   toNullInt64String(r.FormValue("carrier_id")),
		Note:                        model.ToNullString(strings.TrimSpace(r.FormValue("note"))),
		ShippingCharge:              toSQLFloat64Purchase(parseFloat64Default(r.FormValue("shipping_charge"), 0)),
		IsShippingChargeDistributed: parseBoolDefault(r.FormValue("is_shipping_charge_distributed"), false),
		RoundOff:                    toSQLFloat64Purchase(parseFloat64Default(r.FormValue("round_off"), 0)),
		GrandTotal:                  toSQLFloat64Purchase(parseFloat64Default(r.FormValue("grand_total"), 0)),
		ChangeReturn:                toNullInt32String(r.FormValue("change_return")),
		PaidAmount:                  toSQLFloat64Purchase(parseFloat64Default(r.FormValue("paid_amount"), 0)),
		CurrencyID:                  toNullInt64String(r.FormValue("currency_id")),
		ExchangeRate:                toSQLFloat64Purchase(parseFloat64Default(r.FormValue("exchange_rate"), 0)),
		CreatedBy:                   toNullInt64String(r.FormValue("created_by")),
		UpdatedBy:                   toNullInt64String(r.FormValue("updated_by")),
		Items:                       parsePurchaseItemsForm(r.FormValue("items")),
	}, attachments, nil
}

func mapUpdateFormToPurchase(r *http.Request, id int64) (*model.Purchase, []*multipart.FileHeader, error) {
	p, headers, err := mapCreateFormToPurchase(r)
	if err != nil {
		return nil, nil, err
	}
	p.ID = id
	p.RemoveAttachmentIDs = parseInt64CSV(r.FormValue("remove_attachment_ids"))
	return p, headers, nil
}

func getAttachmentHeaders(r *http.Request) []*multipart.FileHeader {
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil
	}
	return r.MultipartForm.File["attachments"]
}

func parseInt64Required(v string, errMsg string) (int64, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, errors.New(errMsg)
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n <= 0 {
		return 0, errors.New(errMsg)
	}
	return n, nil
}

func toNullInt64String(v string) sql.NullInt64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return sql.NullInt64{Valid: false}
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n <= 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: n, Valid: true}
}

func toNullInt32String(v string) sql.NullInt32 {
	v = strings.TrimSpace(v)
	if v == "" {
		return sql.NullInt32{Valid: false}
	}
	n, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(n), Valid: true}
}

func parseFloat64Default(v string, fallback float64) float64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return n
}

func parseBoolDefault(v string, fallback bool) bool {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func parseInt64CSV(csv string) []int64 {
	csv = strings.TrimSpace(csv)
	if csv == "" {
		return nil
	}

	parts := strings.Split(csv, ",")
	ids := make([]int64, 0, len(parts))
	for _, part := range parts {
		v := strings.TrimSpace(part)
		if v == "" {
			continue
		}
		n, err := strconv.ParseInt(v, 10, 64)
		if err == nil && n > 0 {
			ids = append(ids, n)
		}
	}
	return ids
}

func mapPurchaseItems(inputs []model.PurchaseItemInput) []model.PurchaseItem {
	items := make([]model.PurchaseItem, 0, len(inputs))
	for _, input := range inputs {
		item := model.PurchaseItem{
			ItemType:   strings.TrimSpace(input.ItemType),
			ItemName:   strings.TrimSpace(input.ItemName),
			Quantity:   input.Quantity,
			UnitPrice:  input.UnitPrice,
			TotalPrice: input.TotalPrice,
			Note:       model.ToNullString(strings.TrimSpace(input.Note)),
		}
		if input.ProductID != nil && *input.ProductID > 0 {
			item.ProductID = sql.NullInt64{Int64: *input.ProductID, Valid: true}
		}
		items = append(items, item)
	}
	return items
}

func parsePurchaseItemsForm(raw string) []model.PurchaseItem {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	var inputs []model.PurchaseItemInput
	if err := json.Unmarshal([]byte(raw), &inputs); err != nil {
		return nil
	}

	return mapPurchaseItems(inputs)
}
