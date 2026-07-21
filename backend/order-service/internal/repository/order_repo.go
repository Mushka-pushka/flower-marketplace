package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder — создание заказа
func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (
			id, customer_id, shop_id, delivery_address_id, payment_type_id,
			total_amount, commission, delivery_date, delivery_time, comment, current_status,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.Exec(ctx, query,
		order.ID,
		order.CustomerID,
		order.ShopID,
		order.DeliveryAddressID,
		order.PaymentTypeID,
		order.TotalAmount,
		order.Commission,
		order.DeliveryDate,
		order.DeliveryTime,
		order.Comment,
		order.CurrentStatus,
		order.CreatedAt,
		order.UpdatedAt,
	)
	return err
}

// CreateOrderItem — создание позиции заказа
func (r *OrderRepository) CreateOrderItem(ctx context.Context, item *models.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, product_id, quantity, price, total, packaging, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		item.ID,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.Price,
		item.Total,
		item.Packaging,
		item.CreatedAt,
	)
	return err
}

// AddStatusHistory — добавление записи в историю статусов
func (r *OrderRepository) AddStatusHistory(ctx context.Context, history *models.StatusHistory) error {
	query := `
		INSERT INTO order_status_history (id, order_id, status, changed_by, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		history.ID,
		history.OrderID,
		history.Status,
		history.ChangedBy,
		history.Comment,
		history.CreatedAt,
	)
	return err
}

// GetOrderByID — получение заказа по ID
func (r *OrderRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	query := `
		SELECT id, customer_id, shop_id, delivery_address_id, payment_type_id,
			total_amount, commission, delivery_date, delivery_time, comment, current_status,
			created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var order models.Order
	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.CustomerID,
		&order.ShopID,
		&order.DeliveryAddressID,
		&order.PaymentTypeID,
		&order.TotalAmount,
		&order.Commission,
		&order.DeliveryDate,
		&order.DeliveryTime,
		&order.Comment,
		&order.CurrentStatus,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

// GetOrderItemsWithNames — получение позиций заказа с названиями товаров
func (r *OrderRepository) GetOrderItemsWithNames(ctx context.Context, orderID uuid.UUID) ([]models.OrderItemWithName, error) {
	query := `
		SELECT 
			oi.id, oi.order_id, oi.product_id, 
			COALESCE(p.name, 'Товар') as product_name,
			oi.quantity, oi.price, oi.total, oi.packaging, oi.created_at
		FROM order_items oi
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = $1
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItemWithName
	for rows.Next() {
		var item models.OrderItemWithName
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.Price,
			&item.Total,
			&item.Packaging,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetStatusHistory — получение истории статусов заказа
func (r *OrderRepository) GetStatusHistory(ctx context.Context, orderID uuid.UUID) ([]models.StatusHistory, error) {
	query := `
		SELECT id, order_id, status, changed_by, comment, created_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.StatusHistory
	for rows.Next() {
		var h models.StatusHistory
		err := rows.Scan(
			&h.ID,
			&h.OrderID,
			&h.Status,
			&h.ChangedBy,
			&h.Comment,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

// UpdateOrderStatus — обновление статуса заказа
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	query := `UPDATE orders SET current_status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, status, time.Now(), orderID)
	return err
}

// GetOrdersByCustomer — получение заказов покупателя
func (r *OrderRepository) GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]models.Order, error) {
	query := `
		SELECT id, customer_id, shop_id, delivery_address_id, payment_type_id,
			total_amount, commission, delivery_date, delivery_time, comment, current_status,
			created_at, updated_at
		FROM orders
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.ShopID,
			&order.DeliveryAddressID,
			&order.PaymentTypeID,
			&order.TotalAmount,
			&order.Commission,
			&order.DeliveryDate,
			&order.DeliveryTime,
			&order.Comment,
			&order.CurrentStatus,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// GetOrdersByShopID — получает заказы магазина
func (r *OrderRepository) GetOrdersByShopID(ctx context.Context, shopID uuid.UUID) ([]models.Order, error) {
	query := `
		SELECT id, customer_id, shop_id, delivery_address_id, payment_type_id,
			total_amount, commission, delivery_date, delivery_time, comment, current_status,
			created_at, updated_at
		FROM orders
		WHERE shop_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.ShopID,
			&order.DeliveryAddressID,
			&order.PaymentTypeID,
			&order.TotalAmount,
			&order.Commission,
			&order.DeliveryDate,
			&order.DeliveryTime,
			&order.Comment,
			&order.CurrentStatus,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// GetShopIDBySellerID — возвращает shop_id продавца по его user_id
func (r *OrderRepository) GetShopIDBySellerID(ctx context.Context, sellerID uuid.UUID) (uuid.UUID, error) {
	query := `SELECT id FROM shops WHERE seller_id = $1`
	var shopID uuid.UUID
	err := r.db.QueryRow(ctx, query, sellerID).Scan(&shopID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, nil // продавец не имеет магазина
		}
		return uuid.Nil, err
	}
	return shopID, nil
}

// CanReview — проверяет, может ли пользователь оставить отзыв на товар
func (r *OrderRepository) CanReview(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM orders o
			JOIN order_items oi ON oi.order_id = o.id
			WHERE o.customer_id = $1 
			  AND oi.product_id = $2 
			  AND o.current_status = 'delivered'
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, productID).Scan(&exists)
	return exists, err
}

// GetOrderItemsByCustomer — получает все позиции заказов пользователя
func (r *OrderRepository) GetOrderItemsByCustomer(ctx context.Context, customerID uuid.UUID) ([]models.OrderItemWithStatus, error) {
	query := `
		SELECT 
			oi.id,
			oi.order_id,
			oi.product_id,
			p.name as product_name,
			p.price as product_price,
			oi.quantity,
			oi.total,
			o.current_status as order_status,
			o.created_at,
			o.updated_at,
			o.shop_id,
			o.delivery_date,
			o.delivery_time,
			o.comment
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN products p ON p.id = oi.product_id
		WHERE o.customer_id = $1
		ORDER BY o.created_at DESC, oi.id
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItemWithStatus
	for rows.Next() {
		var item models.OrderItemWithStatus
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.ProductPrice,
			&item.Quantity,
			&item.Total,
			&item.OrderStatus,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ShopID,
			&item.DeliveryDate,
			&item.DeliveryTime,
			&item.Comment,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}