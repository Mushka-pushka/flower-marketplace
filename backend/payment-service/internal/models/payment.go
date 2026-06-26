package models

import (
	"time"

	"github.com/google/uuid"
)

// Payment — платёж
type Payment struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	OrderID       uuid.UUID  `json:"order_id" db:"order_id"`
	Amount        float64    `json:"amount" db:"amount"`
	Status        string     `json:"status" db:"status"` // pending, completed, failed, refunded
	PaymentMethod string     `json:"payment_method" db:"payment_method"`
	TransactionID string     `json:"transaction_id" db:"transaction_id"`
	PaymentURL    string     `json:"payment_url" db:"payment_url"`
	CompletedAt   *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// CreatePaymentRequest — запрос на создание платежа
type CreatePaymentRequest struct {
	OrderID       uuid.UUID `json:"order_id" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,gt=0"`
	PaymentMethod string    `json:"payment_method" binding:"required"` // card, cash, online
}

// ConfirmPaymentRequest — запрос на подтверждение платежа
type ConfirmPaymentRequest struct {
	PaymentID uuid.UUID `json:"payment_id" binding:"required"`
	Status    string    `json:"status" binding:"required"` // completed, failed
}