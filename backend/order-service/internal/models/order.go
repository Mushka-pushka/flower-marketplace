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
	Commission        float64    `json:"commission" db:"commission"`
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
	Order    Order                `json:"order"`
	Items    []OrderItemWithName  `json:"items"`   
	Statuses []StatusHistory      `json:"statuses,omitempty"`
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

// UpdateOrderStatusRequest — запрос на обновление статуса заказа
type UpdateOrderStatusRequest struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
	Status  string    `json:"status" binding:"required"`
	Comment string    `json:"comment"`
}

// Courier — курьер
type Courier struct {
    ID          uuid.UUID `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Phone       string    `json:"phone" db:"phone"`
    IsAvailable bool      `json:"is_available" db:"is_available"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// DeliveryAssignment — назначение курьера на заказ
type DeliveryAssignment struct {
    ID          uuid.UUID  `json:"id" db:"id"`
    OrderID     uuid.UUID  `json:"order_id" db:"order_id"`
    CourierID   uuid.UUID  `json:"courier_id" db:"courier_id"`
    AssignedAt  time.Time  `json:"assigned_at" db:"assigned_at"`
    CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
    Status      string     `json:"status" db:"status"` // assigned, completed, failed
}

// OrderItemWithStatus — позиция заказа с данными о товаре и статусе
type OrderItemWithStatus struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	OrderID      uuid.UUID  `json:"order_id" db:"order_id"`
	ProductID    uuid.UUID  `json:"product_id" db:"product_id"`
	ProductName  string     `json:"product_name" db:"product_name"`
	ProductPrice float64    `json:"product_price" db:"product_price"`
	Quantity     int        `json:"quantity" db:"quantity"`
	Total        float64    `json:"total" db:"total"`
	OrderStatus  string     `json:"order_status" db:"order_status"`
	ShopID       uuid.UUID  `json:"shop_id" db:"shop_id"`
	DeliveryDate *time.Time `json:"delivery_date" db:"delivery_date"`
	DeliveryTime string     `json:"delivery_time" db:"delivery_time"`
	Comment      string     `json:"comment" db:"comment"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// OrderItemWithName — позиция заказа с названием товара
type OrderItemWithName struct {
	ID          uuid.UUID `json:"id" db:"id"`
	OrderID     uuid.UUID `json:"order_id" db:"order_id"`
	ProductID   uuid.UUID `json:"product_id" db:"product_id"`
	ProductName string    `json:"product_name" db:"product_name"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Price       float64   `json:"price" db:"price"`
	Total       float64   `json:"total" db:"total"`
	Packaging   string    `json:"packaging" db:"packaging"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}