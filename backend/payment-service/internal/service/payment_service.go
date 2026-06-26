package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/repository"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	cfg         *config.Config
	rabbitCh    *amqp.Channel
}

func NewPaymentService(paymentRepo *repository.PaymentRepository, cfg *config.Config, rabbitCh *amqp.Channel) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		cfg:         cfg,
		rabbitCh:    rabbitCh,
	}
}

// CreatePayment — создаёт платёж и эмулирует оплату
func (s *PaymentService) CreatePayment(ctx context.Context, req *models.CreatePaymentRequest) (*models.Payment, error) {
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

	err := s.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Эмулируем оплату в фоне
	go s.processPayment(payment.ID)

	return payment, nil
}

// processPayment — эмуляция обработки платежа
func (s *PaymentService) processPayment(paymentID uuid.UUID) {
	time.Sleep(5 * time.Second) // Имитация задержки

	ctx := context.Background()
	payment, err := s.paymentRepo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		log.Printf("Failed to get payment: %v", err)
		return
	}

	// 80% успешных платежей, 20% ошибок
	success := time.Now().UnixNano()%10 < 8

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

	// Отправляем событие в RabbitMQ
	s.publishPaymentStatusChanged(payment.OrderID, status)
}

// publishPaymentStatusChanged — публикация события об изменении статуса платежа
func (s *PaymentService) publishPaymentStatusChanged(orderID uuid.UUID, status string) {
	event := map[string]interface{}{
		"event":     "payment.status_changed",
		"order_id":  orderID.String(),
		"status":    status,
		"timestamp": time.Now().Unix(),
	}

	body, _ := json.Marshal(event)
	err := s.rabbitCh.Publish(
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