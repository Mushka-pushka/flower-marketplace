package models

import (
	"time"

	"github.com/google/uuid"
)

// SellerWithUser — продавец с данными пользователя
type SellerWithUser struct {
	ShopID     uuid.UUID `json:"shop_id"`
	ShopName   string    `json:"shop_name"`
	ShopDesc   string    `json:"shop_description,omitempty"`
	IsVerified bool      `json:"is_verified"`
	Rating     float64   `json:"rating"`
	UserID     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

// VerifySellerRequest — запрос на верификацию продавца
type VerifySellerRequest struct {
	ShopID uuid.UUID `json:"shop_id" binding:"required"`
	Verify bool      `json:"verify"` // true — подтвердить, false — отклонить
}

// UpdateUserStatusRequest — запрос на блокировку/разблокировку пользователя
type UpdateUserStatusRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	IsActive bool      `json:"is_active"`
}

// UsersListRequest — запрос на получение списка пользователей
type UsersListRequest struct {
	Role     string `json:"role"`
	Search   string `json:"search"` 
	IsActive *bool  `json:"is_active"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

// UserDetails — детальная информация о пользователе
type UserDetails struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Shop      *ShopInfo `json:"shop,omitempty"`
}

// ShopInfo — информация о магазине (для продавцов)
type ShopInfo struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	IsVerified bool      `json:"is_verified"`
	Rating     float64   `json:"rating"`
}

// UsersListResponse — ответ со списком пользователей
type UsersListResponse struct {
	Users      []UserDetails `json:"users"`
	Total      int64         `json:"total"`
	Limit      int           `json:"limit"`
	Offset     int           `json:"offset"`
	HasMore    bool          `json:"has_more"`
}

// UsersListRequestFull — запрос на получение списка пользователей с поиском
type UsersListRequestFull struct {
	Role     string `json:"role"`
	Search   string `json:"search"`
	IsActive *bool  `json:"is_active"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

// AdminStats — общая статистика для администратора
type AdminStats struct {
	TotalUsers      int64            `json:"total_users"`
	UsersByRole     map[string]int64 `json:"users_by_role"`
	TotalShops      int64            `json:"total_shops"`
	VerifiedShops   int64            `json:"verified_shops"`
	TotalOrders     int64            `json:"total_orders"`
	OrdersByStatus  map[string]int64 `json:"orders_by_status"`
	TotalRevenue    float64          `json:"total_revenue"`
	PlatformRevenue float64          `json:"platform_revenue"`
	TotalProducts   int64            `json:"total_products"`
	ActiveProducts  int64            `json:"active_products"`
	TotalCategories int64            `json:"total_categories"`
}