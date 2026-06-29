package models

import (
	"time"

	"github.com/google/uuid"
)

// DeliveryAddress — адрес доставки
type DeliveryAddress struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Name       string    `json:"name" db:"name"`
	Address    string    `json:"address" db:"address"`
	Entrance   string    `json:"entrance,omitempty" db:"entrance"`
	Floor      string    `json:"floor,omitempty" db:"floor"`
	Intercom   string    `json:"intercom,omitempty" db:"intercom"`
	Comment    string    `json:"comment,omitempty" db:"comment"`
	IsDefault  bool      `json:"is_default" db:"is_default"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateAddressRequest — запрос на создание адреса
type CreateAddressRequest struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Entrance  string `json:"entrance"`
	Floor     string `json:"floor"`
	Intercom  string `json:"intercom"`
	Comment   string `json:"comment"`
	IsDefault bool   `json:"is_default"`
}

// UpdateAddressRequest — запрос на обновление адреса
type UpdateAddressRequest struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Entrance  string `json:"entrance"`
	Floor     string `json:"floor"`
	Intercom  string `json:"intercom"`
	Comment   string `json:"comment"`
	IsDefault bool   `json:"is_default"`
}