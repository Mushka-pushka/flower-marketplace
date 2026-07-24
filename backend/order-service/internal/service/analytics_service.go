package service

import (
	"context"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/repository"

	"github.com/google/uuid"
)

type AnalyticsService struct {
	analyticsRepo *repository.AnalyticsRepository
	cfg           *config.Config
}

func NewAnalyticsService(analyticsRepo *repository.AnalyticsRepository, cfg *config.Config) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
		cfg:           cfg,
	}
}

// GetSellerAnalytics — получает аналитику продавца
func (s *AnalyticsService) GetSellerAnalytics(ctx context.Context, shopID uuid.UUID) (*models.SellerAnalytics, error) {
	return s.analyticsRepo.GetSellerAnalytics(ctx, shopID)
}

// GetPopularProducts — получает популярные товары
func (s *AnalyticsService) GetPopularProducts(ctx context.Context, shopID uuid.UUID, limit int) ([]models.PopularProduct, error) {
	return s.analyticsRepo.GetPopularProducts(ctx, shopID, limit)
}

// GetOrderStatsByStatus — получает статистику по статусам
func (s *AnalyticsService) GetOrderStatsByStatus(ctx context.Context, shopID uuid.UUID) ([]models.OrderStatsByStatus, error) {
	return s.analyticsRepo.GetOrderStatsByStatus(ctx, shopID)
}

// GetSalesDynamics — получает динамику продаж
func (s *AnalyticsService) GetSalesDynamics(ctx context.Context, shopID uuid.UUID, days int) ([]models.SalesDay, error) {
    return s.analyticsRepo.GetSalesDynamics(ctx, shopID, days)
}

// GetShopIDBySellerID — возвращает shop_id продавца
func (s *AnalyticsService) GetShopIDBySellerID(ctx context.Context, sellerID uuid.UUID) (uuid.UUID, error) {
    return s.analyticsRepo.GetShopIDBySellerID(ctx, sellerID)
}