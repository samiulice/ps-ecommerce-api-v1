package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

var (
	ErrSupplierNotFound        = errors.New("supplier not found")
	ErrSupplierNameRequired    = errors.New("supplier name is required")
	ErrSupplierCodeRequired    = errors.New("supplier code is required")
	ErrSupplierPhoneRequired   = errors.New("supplier phone is required")
	ErrSupplierPhoneInvalid    = errors.New("supplier phone format is invalid")
	ErrSupplierCreditInvalid   = errors.New("supplier credit limit cannot be negative")
	ErrSupplierBalanceInvalid  = errors.New("supplier outstanding balance cannot be negative")
	ErrSupplierLeadTimeInvalid = errors.New("supplier lead time cannot be negative")
	ErrSupplierRatingInvalid   = errors.New("supplier rating must be between 0 and 5")
	ErrSupplierEmailDuplicate  = errors.New("supplier email already exists")
	ErrSupplierPhoneDuplicate  = errors.New("supplier phone already exists")
	ErrSupplierCodeDuplicate   = errors.New("supplier code already exists")
)

type SupplierService struct {
	repo *repository.SupplierRepo
}

var supplierPhoneRe = regexp.MustCompile(`^[0-9+()\-\s]{7,25}$`)

func NewSupplierService(repo *repository.SupplierRepo) *SupplierService {
	return &SupplierService{repo: repo}
}

func (s *SupplierService) Create(ctx context.Context, supplier *model.Supplier) (*model.Supplier, error) {
	if err := s.validateSupplier(supplier); err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByCode(ctx, supplier.SupplierCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSupplierCodeDuplicate
	}

	exists, err = s.repo.ExistsByPhone(ctx, supplier.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSupplierPhoneDuplicate
	}

	if supplier.Email.Valid && supplier.Email.String != "" {
		exists, err = s.repo.ExistsByEmail(ctx, supplier.Email.String)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrSupplierEmailDuplicate
		}
	}

	if err := s.repo.Create(ctx, supplier); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, mapUniqueError(err)
		}
		return nil, err
	}
	return supplier, nil
}

func (s *SupplierService) Update(ctx context.Context, supplier *model.Supplier) (*model.Supplier, error) {
	if supplier.ID <= 0 {
		return nil, ErrSupplierNotFound
	}
	if err := s.validateSupplier(supplier); err != nil {
		return nil, err
	}

	current, err := s.repo.FindByID(ctx, supplier.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrSupplierNotFound
		}
		return nil, err
	}

	if supplier.SupplierCode != current.SupplierCode {
		exists, err := s.repo.ExistsByCode(ctx, supplier.SupplierCode)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrSupplierCodeDuplicate
		}
	}

	if supplier.Phone != current.Phone {
		exists, err := s.repo.ExistsByPhone(ctx, supplier.Phone)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrSupplierPhoneDuplicate
		}
	}

	if supplier.Email.Valid && supplier.Email.String != "" {
		currentEmail := ""
		if current.Email.Valid {
			currentEmail = current.Email.String
		}
		if supplier.Email.String != currentEmail {
			exists, err := s.repo.ExistsByEmail(ctx, supplier.Email.String)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrSupplierEmailDuplicate
			}
		}
	}

	if err := s.repo.Update(ctx, supplier); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrSupplierNotFound
		}
		if strings.Contains(err.Error(), "already exists") {
			return nil, mapUniqueError(err)
		}
		return nil, err
	}

	return supplier, nil
}

func (s *SupplierService) GetByID(ctx context.Context, id int64) (*model.Supplier, error) {
	if id <= 0 {
		return nil, ErrSupplierNotFound
	}

	supplier, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrSupplierNotFound
		}
		return nil, err
	}
	return supplier, nil
}

func (s *SupplierService) ListSuppliers(ctx context.Context, filter model.SupplierFilter) ([]*model.SupplierResponse, int64, error) {
	suppliers, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*model.SupplierResponse, 0, len(suppliers))
	for i := range suppliers {
		responses = append(responses, suppliers[i].ToResponse())
	}

	return responses, total, nil
}

func (s *SupplierService) UpdateAccountStatus(ctx context.Context, accountStatus bool, supplierID int64) error {
	if supplierID <= 0 {
		return ErrSupplierNotFound
	}
	err := s.repo.UpdateAccountStatus(ctx, accountStatus, supplierID)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return ErrSupplierNotFound
	}
	return err
}

func (s *SupplierService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrSupplierNotFound
	}
	err := s.repo.Delete(ctx, id)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return ErrSupplierNotFound
	}
	return err
}

func (s *SupplierService) validateSupplier(supplier *model.Supplier) error {
	supplier.SupplierCode = strings.TrimSpace(supplier.SupplierCode)
	supplier.Name = strings.TrimSpace(supplier.Name)
	supplier.Phone = normalizePhone(supplier.Phone)

	if supplier.Email.Valid {
		supplier.Email.String = strings.ToLower(strings.TrimSpace(supplier.Email.String))
		supplier.Email.Valid = supplier.Email.String != ""
	}
	if supplier.Website.Valid {
		supplier.Website.String = strings.TrimSpace(supplier.Website.String)
		supplier.Website.Valid = supplier.Website.String != ""
	}

	if supplier.SupplierCode == "" {
		return ErrSupplierCodeRequired
	}
	if supplier.Name == "" {
		return ErrSupplierNameRequired
	}
	if supplier.Phone == "" {
		return ErrSupplierPhoneRequired
	}
	if !supplierPhoneRe.MatchString(supplier.Phone) {
		return ErrSupplierPhoneInvalid
	}
	if supplier.CreditLimit.Valid && supplier.CreditLimit.Float64 < 0 {
		return ErrSupplierCreditInvalid
	}
	if supplier.OutstandingBalance.Valid && supplier.OutstandingBalance.Float64 < 0 {
		return ErrSupplierBalanceInvalid
	}
	if supplier.LeadTimeDays.Valid && supplier.LeadTimeDays.Int32 < 0 {
		return ErrSupplierLeadTimeInvalid
	}
	if supplier.Rating.Valid && (supplier.Rating.Float64 < 0 || supplier.Rating.Float64 > 5) {
		return ErrSupplierRatingInvalid
	}
	return nil
}

func normalizePhone(phone string) string {
	trimmed := strings.TrimSpace(phone)
	replacer := strings.NewReplacer("-", "", "(", "", ")", "", " ", "")
	return replacer.Replace(trimmed)
}

func mapUniqueError(err error) error {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "supplier_code"):
		return ErrSupplierCodeDuplicate
	case strings.Contains(msg, "phone"):
		return ErrSupplierPhoneDuplicate
	case strings.Contains(msg, "email"):
		return ErrSupplierEmailDuplicate
	default:
		return err
	}
}
