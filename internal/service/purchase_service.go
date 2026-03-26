package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

var (
	ErrPurchaseNotFound          = errors.New("purchase not found")
	ErrPurchaseDateRequired      = errors.New("purchase date is required")
	ErrPurchaseCodeRequired      = errors.New("purchase code is required")
	ErrPurchaseSupplierRequired  = errors.New("supplier is required")
	ErrPurchaseItemsRequired     = errors.New("at least one purchase item is required")
	ErrPurchaseItemInvalid       = errors.New("purchase item is invalid")
	ErrPurchaseGrandTotalInvalid = errors.New("grand total cannot be negative")
	ErrPurchasePaidAmountInvalid = errors.New("paid amount cannot be negative")
	ErrPurchaseShippingInvalid   = errors.New("shipping charge cannot be negative")
	ErrPurchaseRoundOffInvalid   = errors.New("round off cannot be negative")
	ErrPurchaseExchangeInvalid   = errors.New("exchange rate cannot be negative")
	ErrPurchaseCodeDuplicate     = errors.New("purchase code already exists")
)

type PurchaseService struct {
	repo *repository.PurchaseRepo
}

func NewPurchaseService(repo *repository.PurchaseRepo) *PurchaseService {
	return &PurchaseService{repo: repo}
}

func (s *PurchaseService) Create(ctx context.Context, purchase *model.Purchase, attachmentHeaders []*multipart.FileHeader) (*model.Purchase, error) {
	if err := s.validatePurchase(purchase); err != nil {
		return nil, err
	}

	attachments, err := buildPurchaseAttachments(purchase.PurchaseCode, attachmentHeaders)
	if err != nil {
		return nil, err
	}
	purchase.Attachments = attachments

	exists, err := s.repo.ExistsByCode(ctx, purchase.PurchaseCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPurchaseCodeDuplicate
	}

	if err := s.repo.Create(ctx, purchase); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, mapPurchaseUniqueError(err)
		}
		return nil, err
	}

	if err := savePurchaseAttachments(attachmentHeaders, attachments); err != nil {
		return nil, err
	}

	created, err := s.repo.FindByID(ctx, purchase.ID)
	if err != nil {
		return purchase, nil
	}

	return created, nil
}

func (s *PurchaseService) Update(ctx context.Context, purchase *model.Purchase, attachmentHeaders []*multipart.FileHeader) (*model.Purchase, error) {
	if purchase.ID <= 0 {
		return nil, ErrPurchaseNotFound
	}
	if err := s.validatePurchase(purchase); err != nil {
		return nil, err
	}

	attachments, err := buildPurchaseAttachments(purchase.PurchaseCode, attachmentHeaders)
	if err != nil {
		return nil, err
	}
	purchase.Attachments = attachments

	current, err := s.repo.FindByID(ctx, purchase.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrPurchaseNotFound
		}
		return nil, err
	}

	if purchase.PurchaseCode != current.PurchaseCode {
		exists, err := s.repo.ExistsByCode(ctx, purchase.PurchaseCode)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrPurchaseCodeDuplicate
		}
	}

	if err := s.repo.Update(ctx, purchase); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrPurchaseNotFound
		}
		if strings.Contains(err.Error(), "already exists") {
			return nil, mapPurchaseUniqueError(err)
		}
		return nil, err
	}

	if err := savePurchaseAttachments(attachmentHeaders, attachments); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(ctx, purchase.ID)
	if err != nil {
		return purchase, nil
	}

	return updated, nil
}

func (s *PurchaseService) GetByID(ctx context.Context, id int64) (*model.Purchase, error) {
	if id <= 0 {
		return nil, ErrPurchaseNotFound
	}

	purchase, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrPurchaseNotFound
		}
		return nil, err
	}

	return purchase, nil
}

func (s *PurchaseService) ListPurchases(ctx context.Context, filter model.PurchaseFilter) ([]*model.PurchaseResponse, int64, error) {
	purchases, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*model.PurchaseResponse, 0, len(purchases))
	for i := range purchases {
		responses = append(responses, purchases[i].ToResponse())
	}

	return responses, total, nil
}

func (s *PurchaseService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrPurchaseNotFound
	}

	err := s.repo.Delete(ctx, id)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return ErrPurchaseNotFound
	}
	return err
}

func (s *PurchaseService) validatePurchase(purchase *model.Purchase) error {
	if purchase.PurchaseDate.IsZero() {
		return ErrPurchaseDateRequired
	}
	if strings.TrimSpace(purchase.PurchaseCode) == "" {
		return ErrPurchaseCodeRequired
	}
	if purchase.PartyID <= 0 {
		return ErrPurchaseSupplierRequired
	}
	if len(purchase.Items) == 0 {
		return ErrPurchaseItemsRequired
	}
	if err := normalizePurchaseItems(purchase); err != nil {
		return err
	}
	if purchase.GrandTotal.Valid && purchase.GrandTotal.Float64 < 0 {
		return ErrPurchaseGrandTotalInvalid
	}
	if purchase.PaidAmount.Valid && purchase.PaidAmount.Float64 < 0 {
		return ErrPurchasePaidAmountInvalid
	}
	if purchase.ShippingCharge.Valid && purchase.ShippingCharge.Float64 < 0 {
		return ErrPurchaseShippingInvalid
	}
	if purchase.RoundOff.Valid && purchase.RoundOff.Float64 < 0 {
		return ErrPurchaseRoundOffInvalid
	}
	if purchase.ExchangeRate.Valid && purchase.ExchangeRate.Float64 < 0 {
		return ErrPurchaseExchangeInvalid
	}

	return nil
}

func normalizePurchaseItems(purchase *model.Purchase) error {
	items := make([]model.PurchaseItem, 0, len(purchase.Items))
	var subtotal float64

	for _, item := range purchase.Items {
		itemType := strings.ToLower(strings.TrimSpace(item.ItemType))
		itemName := strings.TrimSpace(item.ItemName)

		switch itemType {
		case "product":
			if !item.ProductID.Valid || item.ProductID.Int64 <= 0 {
				return fmt.Errorf("%w: product item must have product_id", ErrPurchaseItemInvalid)
			}
			if itemName == "" {
				return fmt.Errorf("%w: product item name is required", ErrPurchaseItemInvalid)
			}
		case "material":
			item.ProductID.Valid = false
			if itemName == "" {
				return fmt.Errorf("%w: material name is required", ErrPurchaseItemInvalid)
			}
		default:
			return fmt.Errorf("%w: item_type must be product or material", ErrPurchaseItemInvalid)
		}

		if item.Quantity <= 0 {
			return fmt.Errorf("%w: quantity must be greater than zero", ErrPurchaseItemInvalid)
		}
		if item.UnitPrice < 0 {
			return fmt.Errorf("%w: unit price cannot be negative", ErrPurchaseItemInvalid)
		}

		item.ItemType = itemType
		item.ItemName = itemName
		item.TotalPrice = roundPurchaseAmount(item.Quantity * item.UnitPrice)
		item.Note = model.ToNullString(strings.TrimSpace(item.Note.String))

		subtotal += item.TotalPrice
		items = append(items, item)
	}

	purchase.Items = items
	shipping := 0.0
	if purchase.ShippingCharge.Valid {
		shipping = purchase.ShippingCharge.Float64
	}
	roundOff := 0.0
	if purchase.RoundOff.Valid {
		roundOff = purchase.RoundOff.Float64
	}
	purchase.GrandTotal = toNullPurchaseFloat(roundPurchaseAmount(subtotal + shipping + roundOff))

	return nil
}

func roundPurchaseAmount(v float64) float64 {
	return math.Round(v*100) / 100
}

func toNullPurchaseFloat(v float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: v, Valid: true}
}

func mapPurchaseUniqueError(err error) error {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "purchase_code"):
		return ErrPurchaseCodeDuplicate
	default:
		return err
	}
}

func buildPurchaseAttachments(purchaseCode string, headers []*multipart.FileHeader) ([]model.PurchaseAttachment, error) {
	if len(headers) == 0 {
		return nil, nil
	}

	attachments := make([]model.PurchaseAttachment, 0, len(headers))
	now := time.Now().UnixNano()

	for idx, header := range headers {
		if header == nil {
			continue
		}

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext == "" {
			return nil, fmt.Errorf("attachment extension is required for %s", header.Filename)
		}

		filename := fmt.Sprintf("%s-attachment-%d-%d", purchaseCode, now, idx+1)
		attachments = append(attachments, model.PurchaseAttachment{
			FileURL:  utils.GetPurchaseAttachmentURL(filename, ext),
			FileName: header.Filename,
			FileExt:  ext,
			MimeType: header.Header.Get("Content-Type"),
			FileSize: header.Size,
		})
	}

	return attachments, nil
}

func savePurchaseAttachments(headers []*multipart.FileHeader, attachments []model.PurchaseAttachment) error {
	if len(headers) == 0 || len(attachments) == 0 {
		return nil
	}

	for idx, header := range headers {
		if header == nil || idx >= len(attachments) {
			continue
		}

		file, err := header.Open()
		if err != nil {
			return err
		}

		baseName := strings.TrimSuffix(filepath.Base(attachments[idx].FileURL), filepath.Ext(attachments[idx].FileURL))
		_, saveErr := utils.SaveMultipartDocument(file, header, utils.GetPurchaseAttachmentFolderPath(""), baseName)
		file.Close()
		if saveErr != nil {
			return saveErr
		}
	}

	return nil
}
