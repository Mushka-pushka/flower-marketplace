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

    req.Header.Set("Content-Type", "application/json")

    if userID, ok := ctx.Value("user_id").(uuid.UUID); ok {
        req.Header.Set("X-User-ID", userID.String())
    }

    resp, err := s.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get order: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("order service returned status %d: %s", resp.StatusCode, string(body))
    }

    var orderResponse struct {
        Order Order `json:"order"`
    }
    if err := json.Unmarshal(body, &orderResponse); err != nil {
        return nil, fmt.Errorf("failed to decode order response: %w", err)
    }

    return &orderResponse.Order, nil
}

// CreatePayment — создаёт платёж (сразу завершённый)
func (s *PaymentService) CreatePayment(ctx context.Context, req *models.CreatePaymentRequest) (*models.Payment, error) {
	// Проверяем, существует ли уже платеж для этого заказа
	existing, err := s.paymentRepo.GetPaymentByOrderID(ctx, req.OrderID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("payment already exists for this order")
	}

	now := time.Now()

	// Генерируем транзакцию
	transactionID := fmt.Sprintf("TXN-%d-%s", time.Now().UnixNano(), uuid.New().String()[:8])

	// Создаём платёж сразу со статусом "completed"
	payment := &models.Payment{
		ID:            uuid.New(),
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		Status:        "completed", 
		PaymentMethod: req.PaymentMethod,
		TransactionID: transactionID,
		PaymentURL:    fmt.Sprintf("/payment/checkout/%s", uuid.New().String()),
		CompletedAt:   &now, 
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = s.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Сразу отправляем событие об успешной оплате
	s.publishPaymentStatusChanged(payment.OrderID, "completed")
	s.updateOrderStatus(payment.OrderID, "Оплачен")

	log.Printf("Payment %s completed immediately for order %s", payment.ID, payment.OrderID)

	return payment, nil
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