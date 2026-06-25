package models

import (
	"time"

	"github.com/google/uuid"
)

// CartItem — товар в корзине
type CartItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CartItemWithProduct — товар в корзине с данными о товаре
type CartItemWithProduct struct {
	CartItem
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	ProductImage string  `json:"product_image,omitempty"`
	TotalPrice   float64 `json:"total_price"` // quantity * price
}

// AddToCartRequest — запрос на добавление в корзину
type AddToCartRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

// UpdateCartRequest — запрос на обновление количества
type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}