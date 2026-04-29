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

func (s *DashboardService) GetExpenseGraph(ctx context.Context, period string) ([]model.ChartPoint, error) {
	return s.repo.GetFinancialGraph(ctx, period, "total_expenses")
}

func (s *DashboardService) GetNetProfitGraph(ctx context.Context, period string) ([]model.ChartPoint, error) {
	return s.repo.GetFinancialGraph(ctx, period, "net_profit")
}

func (s *DashboardService) GetPopularProducts(ctx context.Context) ([]model.ProductSummary, error) {
	return s.repo.GetPopularProducts(ctx)
}

func (s *DashboardService) GetLowStockProducts(ctx context.Context) ([]model.ProductSummary, error) {
	return s.repo.GetLowStockProducts(ctx)
}