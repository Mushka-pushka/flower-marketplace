package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/repository"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	cfg         *config.Config
	rabbitCh    *amqp.Channel
	httpClient  *http.Client
	catalogURL  string
}

// Удаляем константу, так как теперь используем значение из конфига
// const platformCommissionRate = 0.10

func NewOrderService(
	orderRepo *repository.OrderRepository,
	cfg *config.Config,
	rabbitCh *amqp.Channel,
) *OrderService {
	// Получаем URL каталога из конфига или используем значение по умолчанию
	catalogURL := cfg.CatalogServiceURL
	if catalogURL == "" {
		catalogURL = "http://localhost:8082/api/v1/catalog"
	}

	return &OrderService{
		orderRepo:  orderRepo,
		cfg:        cfg,
		rabbitCh:   rabbitCh,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		catalogURL: catalogURL,
	}
}

// GetProductPrice — получает цену товара из Catalog Service
func (s *OrderService) GetProductPrice(ctx context.Context, productID uuid.UUID) (float64, error) {
	url := fmt.Sprintf("%s/products/%s", s.catalogURL, productID.String())
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get product price: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("catalog service returned status %d: %s", resp.StatusCode, string(body))
	}

	var product struct {
		Price float64 `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return 0, fmt.Errorf("failed to decode product response: %w", err)
	}

	if product.Price <= 0 {
		return 0, fmt.Errorf("invalid price: %.2f", product.Price)
	}

	return product.Price, nil
}

// GetProductStock — получает остаток товара из Catalog Service
func (s *OrderService) GetProductStock(ctx context.Context, productID uuid.UUID) (int, error) {
	url := fmt.Sprintf("%s/products/%s", s.catalogURL, productID.String())
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get product stock: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("catalog service returned status %d: %s", resp.StatusCode, string(body))
	}

	var product struct {
		Stock int `json:"stock"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return 0, fmt.Errorf("failed to decode product response: %w", err)
	}

	if product.Stock < 0 {
		return 0, fmt.Errorf("invalid stock: %d", product.Stock)
	}

	return product.Stock, nil
}

// CreateOrder — создание заказа с проверкой наличия товаров на складе
func (s *OrderService) CreateOrder(ctx context.Context, customerID uuid.UUID, req *models.CreateOrderRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	var totalAmount float64
	var orderItems []struct {
		ProductID uuid.UUID
		Quantity  int
		Price     float64
	}

	// Проверяем наличие товаров на складе и получаем цены
	for _, item := range req.Items {
		// Получаем остаток товара
		stock, err := s.GetProductStock(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to check stock for product %s: %w", item.ProductID, err)
		}

		// Проверяем, что запрашиваемое количество есть в наличии
		if stock < item.Quantity {
			return nil, fmt.Errorf("not enough stock for product %s. Available: %d, requested: %d", 
				item.ProductID, stock, item.Quantity)
		}

		// Получаем цену товара
		price, err := s.GetProductPrice(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get price for product %s: %w", item.ProductID, err)
		}
		
		itemTotal := price * float64(item.Quantity)
		totalAmount += itemTotal
		
		orderItems = append(orderItems, struct {
			ProductID uuid.UUID
			Quantity  int
			Price     float64
		}{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
		})
		
		log.Printf("Product %s: stock=%d, price=%.2f, quantity=%d, total=%.2f", 
			item.ProductID, stock, price, item.Quantity, itemTotal)
	}

	// Используем настраиваемую комиссию из конфига
	commission := totalAmount * s.cfg.PlatformCommission
	log.Printf("Platform commission: %.2f%% (%.2f of %.2f)", 
		s.cfg.PlatformCommission*100, commission, totalAmount)

	now := time.Now()
	order := &models.Order{
		ID:                uuid.New(),
		CustomerID:        customerID,
		ShopID:            req.ShopID,
		DeliveryAddressID: req.DeliveryAddressID,
		PaymentTypeID:     req.PaymentTypeID,
		TotalAmount:       totalAmount,
		Commission:        commission,
		DeliveryTime:      req.DeliveryTime,
		Comment:           req.Comment,
		CurrentStatus:     "pending",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if req.DeliveryDate != "" {
		if deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate); err == nil {
			order.DeliveryDate = &deliveryDate
		}
	}

	err := s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Создаем позиции заказа с реальными ценами из Catalog Service
	for _, item := range orderItems {
		orderItem := &models.OrderItem{
			ID:        uuid.New(),
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Total:     item.Price * float64(item.Quantity),
			CreatedAt: now,
		}
		err = s.orderRepo.CreateOrderItem(ctx, orderItem)
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

	// Публикуем событие создания заказа
	err = s.publishOrderCreated(ctx, order)
	if err != nil {
		// Логируем ошибку, но не прерываем создание заказа
		log.Printf("Failed to publish order created event: %v", err)
	}

	return order, nil
}

func (s *OrderService) publishOrderCreated(ctx context.Context, order *models.Order) error {
	event := map[string]interface{}{
		"event":       "order.created",
		"order_id":    order.ID.String(),
		"customer_id": order.CustomerID.String(),
		"total":       order.TotalAmount,
		"commission":  order.Commission,
		"status":      order.CurrentStatus,
		"timestamp":   time.Now().Unix(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = s.rabbitCh.PublishWithContext(ctx,
		"",
		"order.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Event published: order.created for order %s", order.ID)
	return nil
}

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
	err = s.orderRepo.AddStatusHistory(ctx, history)
	if err != nil {
		return err
	}

	go s.publishOrderStatusChanged(orderID, status)

	return nil
}

func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]models.Order, error) {
	return s.orderRepo.GetOrdersByCustomer(ctx, customerID)
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, role string) error {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if role == "customer" && order.CustomerID != userID {
		return errors.New("you can only cancel your own orders")
	}

	if order.CurrentStatus == "delivered" {
		return errors.New("cannot cancel delivered order")
	}
	if order.CurrentStatus == "cancelled" {
		return errors.New("order already cancelled")
	}

	err = s.orderRepo.UpdateOrderStatus(ctx, orderID, "cancelled")
	if err != nil {
		return err
	}

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

	go s.publishOrderCancelled(orderID)

	return nil
}

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

func (s *OrderService) GetOrdersByShop(ctx context.Context, shopID uuid.UUID) ([]models.Order, error) {
	return s.orderRepo.GetOrdersByShopID(ctx, shopID)
}

func (s *OrderService) UpdateOrderStatusBySeller(ctx context.Context, orderID, shopID uuid.UUID, status, comment string) error {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.ShopID != shopID {
		return errors.New("you can only update orders from your shop")
	}

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

	if order.CurrentStatus == "delivered" || order.CurrentStatus == "cancelled" {
		return errors.New("cannot change status of delivered or cancelled order")
	}

	err = s.orderRepo.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return err
	}

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

	go s.publishOrderStatusChanged(orderID, status)

	return nil
}

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

// AssignCourier — назначает курьера на заказ (ВРЕМЕННО ОТКЛЮЧЕНО)
func (s *OrderService) AssignCourier(ctx context.Context, orderID uuid.UUID) (*models.Courier, error) {
	// Временно отключаем, чтобы не мешал
	return nil, nil
}

// CompleteDelivery — завершает доставку (ВРЕМЕННО ОТКЛЮЧЕНО)
func (s *OrderService) CompleteDelivery(ctx context.Context, orderID uuid.UUID) error {
	// Временно отключаем
	return nil
}

// publishDeliveryAssigned — публикация события назначения курьера (ВРЕМЕННО ОТКЛЮЧЕНО)
func (s *OrderService) publishDeliveryAssigned(orderID uuid.UUID, courier *models.Courier) {
	// Временно отключаем
}

// publishDeliveryCompleted — публикация события завершения доставки (ВРЕМЕННО ОТКЛЮЧЕНО)
func (s *OrderService) publishDeliveryCompleted(orderID uuid.UUID) {
	// Временно отключаем
}

func (s *OrderService) CanReview(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	return s.orderRepo.CanReview(ctx, userID, productID)
}

// GetShopIDBySellerID — возвращает shop_id продавца по его user_id
func (s *OrderService) GetShopIDBySellerID(ctx context.Context, sellerID uuid.UUID) (uuid.UUID, error) {
	return s.orderRepo.GetShopIDBySellerID(ctx, sellerID)
}