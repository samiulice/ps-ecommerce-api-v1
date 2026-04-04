package service

import (
	"context"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

type ReportService struct {
	reportRepo *repository.ReportRepo
}

func NewReportService(reportRepo *repository.ReportRepo) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

func (s *ReportService) GetPOSSalesReport(ctx context.Context, filter model.ReportFilter) (*model.POSSaleReportResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return s.reportRepo.GetPOSSalesReport(ctx, filter)
}

func (s *ReportService) GetOrdersReport(ctx context.Context, filter model.ReportFilter) (*model.OrderReportResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return s.reportRepo.GetOrdersReport(ctx, filter)
}

func (s *ReportService) GetCustomerDueReport(ctx context.Context, filter model.ReportFilter) (*model.CustomerDueReportResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return s.reportRepo.GetCustomerDueReport(ctx, filter)
}

func (s *ReportService) GetSupplierDueReport(ctx context.Context, filter model.ReportFilter) (*model.SupplierDueReportResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return s.reportRepo.GetSupplierDueReport(ctx, filter)
}

func (s *ReportService) GetLowStockReport(ctx context.Context, filter model.ReportFilter) (*model.LowStockReportResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return s.reportRepo.GetLowStockReport(ctx, filter)
}
