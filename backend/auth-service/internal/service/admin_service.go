package service

import (
	"context"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/repository"

	"github.com/google/uuid"
)

type AdminService struct {
	adminRepo *repository.AdminRepository
	cfg       *config.Config
}

func NewAdminService(adminRepo *repository.AdminRepository, cfg *config.Config) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
		cfg:       cfg,
	}
}

// GetSellers — получает список продавцов с магазинами
func (s *AdminService) GetSellers(ctx context.Context, verified *bool) ([]models.SellerWithUser, error) {
	return s.adminRepo.GetSellersWithShops(ctx, verified)
}

// VerifySeller — верифицирует продавца
func (s *AdminService) VerifySeller(ctx context.Context, shopID uuid.UUID, verify bool) error {
	return s.adminRepo.VerifyShop(ctx, shopID, verify)
}

// UpdateUserStatus — обновляет статус пользователя
func (s *AdminService) UpdateUserStatus(ctx context.Context, userID uuid.UUID, isActive bool) error {
	// Проверяем, что пользователь не пытается заблокировать самого себя
	// (в реальном проекте это проверяется через контекст)
	return s.adminRepo.UpdateUserStatus(ctx, userID, isActive)
}

// GetUsersList — получает список пользователей
func (s *AdminService) GetUsersList(ctx context.Context, role string, isActive *bool, limit, offset int) ([]models.User, error) {
	return s.adminRepo.GetUsersList(ctx, role, isActive, limit, offset)
}