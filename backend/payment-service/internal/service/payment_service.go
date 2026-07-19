package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/repository"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Order — структура заказа для получения из Order Service
type Order struct {
	ID          uuid.UUID `json:"id"`
	CustomerID  uuid.UUID `json:"customer_id"`
	ShopID      uuid.UUID `json:"shop_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
}

type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	cfg         *config.Config
	rabbitCh    *amqp.Channel
	httpClient  *http.Client
	orderURL    string
}

func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	cfg *config.Config,
	rabbitCh *amqp.Channel,
) *PaymentService {
	// Получаем URL сервиса заказов из конфига
	orderURL := cfg.OrderServiceURL
	if orderURL == "" {
		orderURL = "http://localhost:8083/api/v1/orders"
	}

	return &PaymentService{
		paymentRepo: paymentRepo,
		cfg:         cfg,
		rabbitCh:    rabbitCh,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		orderURL:    orderURL,
	}
}

// GetOrderByID — получает заказ из Order Service
func (s *PaymentService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*Order, error) {
	url := fmt.Sprintf("%s?id=%s", s.orderURL, orderID.String())
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Добавляем заголовки для аутентификации (если нужны)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("order service returned status %d: %s", resp.StatusCode, string(body))
	}

	var order Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	return &order, nil
}

// CreatePayment — создаёт платёж и эмулирует оплату
func (s *PaymentService) CreatePayment(ctx context.Context, req *models.CreatePaymentRequest) (*models.Payment, error) {
	// Проверяем, существует ли уже платеж для этого заказа
	existing, err := s.paymentRepo.GetPaymentByOrderID(ctx, req.OrderID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("payment already exists for this order")
	}

	now := time.Now()

	// Генерируем транзакцию
	transactionID := fmt.Sprintf("TXN-%d-%s", time.Now().UnixNano(), uuid.New().String()[:8])

	payment := &models.Payment{
		ID:            uuid.New(),
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		Status:        "pending",
		PaymentMethod: req.PaymentMethod,
		TransactionID: transactionID,
		PaymentURL:    fmt.Sprintf("/payment/checkout/%s", uuid.New().String()),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = s.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Публикуем событие создания платежа
	go s.publishPaymentCreated(payment)

	// Эмулируем оплату в фоне
	go s.processPayment(payment.ID)

	return payment, nil
}

// processPayment — эмуляция обработки платежа (исправленная версия)
func (s *PaymentService) processPayment(paymentID uuid.UUID) {
	// Используем time.After вместо time.Sleep
	select {
	case <-time.After(5 * time.Second):
		// Обработка платежа
		ctx := context.Background()
		payment, err := s.paymentRepo.GetPaymentByID(ctx, paymentID)
		if err != nil {
			log.Printf("Failed to get payment: %v", err)
			return
		}

		// Используем настраиваемую вероятность успеха
		successRate := s.cfg.PaymentSuccessRate
		if successRate <= 0 {
			successRate = 0.8 // По умолчанию 80%
		}
		
		success := time.Now().UnixNano()%100 < int64(successRate*100)

		var status string
		var completedAt *time.Time
		if success {
			status = "completed"
			now := time.Now()
			completedAt = &now
			log.Printf("Payment %s completed successfully", paymentID)
		} else {
			status = "failed"
			log.Printf("Payment %s failed", paymentID)
		}

		err = s.paymentRepo.UpdatePaymentStatus(ctx, paymentID, status, completedAt)
		if err != nil {
			log.Printf("Failed to update payment status: %v", err)
			return
		}

		s.publishPaymentStatusChanged(payment.OrderID, status)

		if status == "completed" {
			s.updateOrderStatus(payment.OrderID, "paid")
		}
	}
}

// updateOrderStatus — обновляет статус заказа через RabbitMQ
func (s *PaymentService) updateOrderStatus(orderID uuid.UUID, status string) {
	// Отправляем событие в RabbitMQ для Order Service
	event := map[string]interface{}{
		"event":     "order.payment_completed",
		"order_id":  orderID.String(),
		"status":    status,
		"timestamp": time.Now().Unix(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal order.payment_completed event: %v", err)
		return
	}

	err = s.rabbitCh.Publish(
		"",
		"order.payment_completed",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish order.payment_completed: %v", err)
	} else {
		log.Printf("Event published: order.payment_completed for order %s", orderID)
	}
}

// publishPaymentCreated — публикация события создания платежа
func (s *PaymentService) publishPaymentCreated(payment *models.Payment) {
	event := map[string]interface{}{
		"event":          "payment.created",
		"payment_id":     payment.ID.String(),
		"order_id":       payment.OrderID.String(),
		"amount":         payment.Amount,
		"payment_method": payment.PaymentMethod,
		"status":         payment.Status,
		"transaction_id": payment.TransactionID,
		"timestamp":      time.Now().Unix(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal payment.created event: %v", err)
		return
	}

	err = s.rabbitCh.Publish(
		"",
		"payment.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish payment.created: %v", err)
	} else {
		log.Printf("Event published: payment.created for payment %s", payment.ID)
	}
}

// publishPaymentStatusChanged — публикация события об изменении статуса платежа
func (s *PaymentService) publishPaymentStatusChanged(orderID uuid.UUID, status string) {
	event := map[string]interface{}{
		"event":     "payment.status_changed",
		"order_id":  orderID.String(),
		"status":    status,
		"timestamp": time.Now().Unix(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal payment.status_changed event: %v", err)
		return
	}

	err = s.rabbitCh.Publish(
		"",
		"payment.status_changed",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish payment.status_changed: %v", err)
	} else {
		log.Printf("Event published: payment.status_changed for order %s -> %s", orderID, status)
	}
}

// GetPaymentStatus — получает статус платежа
func (s *PaymentService) GetPaymentStatus(ctx context.Context, paymentID uuid.UUID) (*models.Payment, error) {
	return s.paymentRepo.GetPaymentByID(ctx, paymentID)
}

// GetPaymentByOrderID — получает платёж по ID заказа
func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error) {
	return s.paymentRepo.GetPaymentByOrderID(ctx, orderID)
}