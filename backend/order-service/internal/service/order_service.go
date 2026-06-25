package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log" 
	"time"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/repository"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
	cfg       *config.Config
	rabbitCh  *amqp.Channel
}

func NewOrderService(orderRepo *repository.OrderRepository, cfg *config.Config, rabbitCh *amqp.Channel) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		cfg:       cfg,
		rabbitCh:  rabbitCh,
	}
}

// CreateOrder — создание заказа и отправка события в RabbitMQ
func (s *OrderService) CreateOrder(ctx context.Context, req *models.CreateOrderRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// Рассчитываем общую сумму
	var totalAmount float64
	for _, item := range req.Items {
		// Здесь можно добавить запрос к Catalog Service для получения цены товара
		// Пока используем заглушку
		totalAmount += 1000 * float64(item.Quantity) // Заглушка
	}

	now := time.Now()
	order := &models.Order{
		ID:                uuid.New(),
		CustomerID:        req.CustomerID,
		ShopID:            req.ShopID,
		DeliveryAddressID: req.DeliveryAddressID,
		PaymentTypeID:     req.PaymentTypeID,
		TotalAmount:       totalAmount,
		DeliveryTime:      req.DeliveryTime,
		Comment:           req.Comment,
		CurrentStatus:     "pending",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Парсим дату доставки
	if req.DeliveryDate != "" {
		if deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate); err == nil {
			order.DeliveryDate = &deliveryDate
		}
	}

	// Сохраняем заказ в БД
	err := s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Сохраняем позиции заказа
	for _, itemReq := range req.Items {
		item := &models.OrderItem{
			ID:        uuid.New(),
			OrderID:   order.ID,
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     1000, // Заглушка
			Total:     1000 * float64(itemReq.Quantity),
			CreatedAt: now,
		}
		err = s.orderRepo.CreateOrderItem(ctx, item)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Добавляем запись в историю статусов
	history := &models.StatusHistory{
		ID:        uuid.New(),
		OrderID:   order.ID,
		Status:    "pending",
		ChangedBy: "system",
		Comment:   "Заказ создан",
		CreatedAt: now,
	}
	err = s.orderRepo.AddStatusHistory(ctx, history)
	if err != nil {
		return nil, fmt.Errorf("failed to add status history: %w", err)
	}

	// Отправляем событие в RabbitMQ
	err = s.publishOrderCreated(ctx, order)
	if err != nil {
		// Логируем ошибку, но не отменяем создание заказа
		fmt.Printf("Failed to publish order created event: %v\n", err)
	}

	return order, nil
}

// publishOrderCreated — публикация события в RabbitMQ
func (s *OrderService) publishOrderCreated(ctx context.Context, order *models.Order) error {
	event := map[string]interface{}{
		"event":      "order.created",
		"order_id":   order.ID.String(),
		"customer_id": order.CustomerID.String(),
		"total":      order.TotalAmount,
		"status":     order.CurrentStatus,
		"timestamp":  time.Now().Unix(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = s.rabbitCh.PublishWithContext(ctx,
		"",                // exchange
		"order.created",   // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	fmt.Printf("Event published: order.created for order %s\n", order.ID)
	return nil
}

// GetOrderByID — получение заказа с деталями
func (s *OrderService) GetOrderByID(ctx context.Context, id uuid.UUID) (*models.OrderResponse, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.orderRepo.GetOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}

	statuses, err := s.orderRepo.GetStatusHistory(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.OrderResponse{
		Order:    *order,
		Items:    items,
		Statuses: statuses,
	}, nil
}

// UpdateOrderStatus — обновление статуса заказа (вызывается из воркера)
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status, changedBy, comment string) error {
	err := s.orderRepo.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return err
	}

	history := &models.StatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    status,
		ChangedBy: changedBy,
		Comment:   comment,
		CreatedAt: time.Now(),
	}
	return s.orderRepo.AddStatusHistory(ctx, history)
}

// GetOrdersByCustomer — получение заказов покупателя
func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]models.Order, error) {
	return s.orderRepo.GetOrdersByCustomer(ctx, customerID)
}

// CancelOrder — отмена заказа
func (s *OrderService) CancelOrder(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, role string) error {
	// Получаем заказ
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Проверяем, что заказ принадлежит пользователю или пользователь — продавец/админ
	if role == "customer" && order.CustomerID != userID {
		return errors.New("you can only cancel your own orders")
	}

	// Проверяем, что заказ ещё не доставлен и не отменён
	if order.CurrentStatus == "delivered" {
		return errors.New("cannot cancel delivered order")
	}
	if order.CurrentStatus == "cancelled" {
		return errors.New("order already cancelled")
	}

	// Обновляем статус
	err = s.orderRepo.UpdateOrderStatus(ctx, orderID, "cancelled")
	if err != nil {
		return err
	}

	// Добавляем запись в историю
	history := &models.StatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    "cancelled",
		ChangedBy: userID.String(),
		Comment:   "Заказ отменён пользователем",
		CreatedAt: time.Now(),
	}
	err = s.orderRepo.AddStatusHistory(ctx, history)
	if err != nil {
		return err
	}

	// Отправляем событие в RabbitMQ
	go s.publishOrderCancelled(orderID)

	return nil
}

// publishOrderCancelled — публикация события об отмене заказа
func (s *OrderService) publishOrderCancelled(orderID uuid.UUID) {
	event := map[string]interface{}{
		"event":     "order.cancelled",
		"order_id":  orderID.String(),
		"timestamp": time.Now().Unix(),
	}

	body, _ := json.Marshal(event)
	s.rabbitCh.Publish(
		"",
		"order.cancelled",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	log.Printf("Event published: order.cancelled for order %s", orderID)
}

// GetOrdersByShop — получает заказы магазина
func (s *OrderService) GetOrdersByShop(ctx context.Context, shopID uuid.UUID) ([]models.Order, error) {
	return s.orderRepo.GetOrdersByShopID(ctx, shopID)
}

// UpdateOrderStatusBySeller — обновление статуса заказа продавцом
func (s *OrderService) UpdateOrderStatusBySeller(ctx context.Context, orderID, shopID uuid.UUID, status, comment string) error {
	// Получаем заказ
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Проверяем, что заказ принадлежит магазину продавца
	if order.ShopID != shopID {
		return errors.New("you can only update orders from your shop")
	}

	// Проверяем допустимость статуса
	validStatuses := map[string]bool{
		"confirmed": true,
		"preparing": true,
		"packing":   true,
		"delivery":  true,
		"delivered": true,
		"cancelled": true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	// Нельзя изменить статус доставленного или отменённого заказа
	if order.CurrentStatus == "delivered" || order.CurrentStatus == "cancelled" {
		return errors.New("cannot change status of delivered or cancelled order")
	}

	// Обновляем статус
	err = s.orderRepo.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return err
	}

	// Добавляем запись в историю
	history := &models.StatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    status,
		ChangedBy: "seller",
		Comment:   comment,
		CreatedAt: time.Now(),
	}
	err = s.orderRepo.AddStatusHistory(ctx, history)
	if err != nil {
		return err
	}

	// Отправляем событие в RabbitMQ
	go s.publishOrderStatusChanged(orderID, status)

	return nil
}

// publishOrderStatusChanged — публикация события об изменении статуса
func (s *OrderService) publishOrderStatusChanged(orderID uuid.UUID, status string) {
	event := map[string]interface{}{
		"event":     "order.status_changed",
		"order_id":  orderID.String(),
		"status":    status,
		"timestamp": time.Now().Unix(),
	}

	body, _ := json.Marshal(event)
	s.rabbitCh.Publish(
		"",
		"order.status_changed",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	log.Printf("Event published: order.status_changed for order %s -> %s", orderID, status)
}