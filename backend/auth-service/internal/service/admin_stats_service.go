package service

import (
	"context"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/repository"
)

type AdminStatsService struct {
	statsRepo *repository.AdminStatsRepository
	cfg       *config.Config
}

func NewAdminStatsService(statsRepo *repository.AdminStatsRepository, cfg *config.Config) *AdminStatsService {
	return &AdminStatsService{
		statsRepo: statsRepo,
		cfg:       cfg,
	}
}

// GetAdminStats — получает общую статистику для администратора
func (s *AdminStatsService) GetAdminStats(ctx context.Context) (*models.AdminStats, error) {
	stats := &models.AdminStats{
		UsersByRole:    make(map[string]int64),
		OrdersByStatus: make(map[string]int64),
	}

	// 1. Статистика по пользователям
	totalUsers, byRole, err := s.statsRepo.GetUserStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalUsers = totalUsers
	stats.UsersByRole = byRole

	// 2. Статистика по магазинам
	totalShops, verifiedShops, err := s.statsRepo.GetShopStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalShops = totalShops
	stats.VerifiedShops = verifiedShops

	// 3. Статистика по заказам
	totalOrders, byStatus, revenue, platformRevenue, err := s.statsRepo.GetOrderStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalOrders = totalOrders
	stats.OrdersByStatus = byStatus
	stats.TotalRevenue = revenue
	stats.PlatformRevenue = platformRevenue

	// 4. Статистика по товарам
	totalProducts, activeProducts, err := s.statsRepo.GetProductStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalProducts = totalProducts
	stats.ActiveProducts = activeProducts

	// 5. Статистика по категориям
	totalCategories, err := s.statsRepo.GetCategoryStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalCategories = totalCategories

	return stats, nil
}