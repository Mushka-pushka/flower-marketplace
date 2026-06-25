package models

import (
	"time"

	"github.com/google/uuid"
)

// Order — заказ
type Order struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	CustomerID        uuid.UUID  `json:"customer_id" db:"customer_id"`
	ShopID            uuid.UUID  `json:"shop_id" db:"shop_id"`
	DeliveryAddressID uuid.UUID  `json:"delivery_address_id" db:"delivery_address_id"`
	PaymentTypeID     int        `json:"payment_type_id" db:"payment_type_id"`
	TotalAmount       float64    `json:"total_amount" db:"total_amount"`
	DeliveryDate      *time.Time `json:"delivery_date" db:"delivery_date"`
	DeliveryTime      string     `json:"delivery_time" db:"delivery_time"`
	Comment           string     `json:"comment" db:"comment"`
	CurrentStatus     string     `json:"current_status" db:"current_status"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// OrderItem — позиция в заказе
type OrderItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	OrderID   uuid.UUID `json:"order_id" db:"order_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
	Total     float64   `json:"total" db:"total"`
	Packaging string    `json:"packaging" db:"packaging"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateOrderRequest — запрос на создание заказа
type CreateOrderRequest struct {
	CustomerID        uuid.UUID `json:"customer_id"`
	ShopID            uuid.UUID `json:"shop_id"`
	DeliveryAddressID uuid.UUID `json:"delivery_address_id"`
	PaymentTypeID     int       `json:"payment_type_id"`
	DeliveryDate      string    `json:"delivery_date"`
	DeliveryTime      string    `json:"delivery_time"`
	Comment           string    `json:"comment"`
	Items             []OrderItemRequest `json:"items"`
}

// OrderItemRequest — позиция в заказе (запрос)
type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

// OrderStatusUpdate — обновление статуса заказа
type OrderStatusUpdate struct {
	OrderID   uuid.UUID `json:"order_id"`
	Status    string    `json:"status"`
	ChangedBy string    `json:"changed_by"`
	Comment   string    `json:"comment"`
}

// OrderResponse — ответ с заказом
type OrderResponse struct {
	Order     Order        `json:"order"`
	Items     []OrderItem  `json:"items"`
	Statuses  []StatusHistory `json:"statuses,omitempty"`
}

// StatusHistory — история статусов
type StatusHistory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	OrderID   uuid.UUID `json:"order_id" db:"order_id"`
	Status    string    `json:"status" db:"status"`
	ChangedBy string    `json:"changed_by" db:"changed_by"`
	Comment   string    `json:"comment" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}