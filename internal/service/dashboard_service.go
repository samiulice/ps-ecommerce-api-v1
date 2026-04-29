package service

import (
	"context"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetStatsCards(ctx context.Context, period string) (*model.StatsCards, error) {
	return s.repo.GetStatsCards(ctx, period)
}

func (s *DashboardService) GetSaleComparison(ctx context.Context, period string) ([]model.SaleComparisonData, error) {
	return s.repo.GetSaleComparisonChart(ctx, period)
}

func (s *DashboardService) GetOperationalExpenseGraph(ctx context.Context, period string) ([]model.ChartPoint, error) {
	return s.repo.GetFinancialGraph(ctx, period, "operational_expenses")
}

func (s *DashboardService) GetNetCashFlowGraph(ctx context.Context, period string) ([]model.ChartPoint, error) {
	return s.repo.GetFinancialGraph(ctx, period, "net_cash_flow")
}

func (s *DashboardService) GetPopularProducts(ctx context.Context) ([]model.ProductSummary, error) {
	return s.repo.GetPopularProducts(ctx)
}

func (s *DashboardService) GetLowStockProducts(ctx context.Context) ([]model.ProductSummary, error) {
	return s.repo.GetLowStockProducts(ctx)
}
