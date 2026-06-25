package models

import (
	"time"

	"github.com/google/uuid"
)

// Favorite — избранный товар
type Favorite struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// FavoriteWithProduct — избранный товар с данными о товаре
type FavoriteWithProduct struct {
	Favorite
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	ProductImage string  `json:"product_image,omitempty"`
	ProductSlug  string  `json:"product_slug"`
}

// AddFavoriteRequest — запрос на добавление в избранное
type AddFavoriteRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
}