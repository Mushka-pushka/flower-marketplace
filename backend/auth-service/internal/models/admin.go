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
	IsActive *bool  `json:"is_active"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}