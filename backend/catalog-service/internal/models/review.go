package models

import (
	"time"

	"github.com/google/uuid"
)

// Review — отзыв на товар
type Review struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ProductID  uuid.UUID `json:"product_id" db:"product_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	OrderID    *uuid.UUID `json:"order_id,omitempty" db:"order_id"`
	Rating     int       `json:"rating" db:"rating"`
	Comment    string    `json:"comment" db:"comment"`
	IsApproved bool      `json:"is_approved" db:"is_approved"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// ReviewWithUser — отзыв с данными о пользователе
type ReviewWithUser struct {
	Review
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	UserAvatar *string `json:"user_avatar,omitempty"`
}

// CreateReviewRequest — запрос на создание отзыва
type CreateReviewRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	OrderID   uuid.UUID `json:"order_id,omitempty"`
	Rating    int       `json:"rating" binding:"required,min=1,max=5"`
	Comment   string    `json:"comment" binding:"required"`
}

// UpdateReviewRequest — запрос на обновление отзыва
type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"min=1,max=5"`
	Comment string `json:"comment"`
}