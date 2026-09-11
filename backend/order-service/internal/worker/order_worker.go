package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/service"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderWorker struct {
	orderService *service.OrderService
	rabbitCh     *amqp.Channel
}

func NewOrderWorker(orderService *service.OrderService, rabbitCh *amqp.Channel) *OrderWorker {
	return &OrderWorker{
		orderService: orderService,
		rabbitCh:     rabbitCh,
	}
}

// Start — запуск воркера для обработки событий оплаты
// ВАЖНО: Воркер НЕ меняет статусы заказов автоматически.
// Он только обрабатывает результаты оплаты от Payment Service.
func (w *OrderWorker) Start(ctx context.Context) error {
	// Объявляем очереди
	queues := []string{
		"order.created",
		"order.payment_completed",
		"order.cancelled",
		"order.status_changed",
	}

	for _, queueName := range queues {
		_, err := w.rabbitCh.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
		log.Printf("Queue declared: %s", queueName)
	}

	// Подписываемся на события об оплате
	paymentMsgs, err := w.rabbitCh.Consume(
		"order.payment_completed",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register payment consumer: %w", err)
	}

	cancelledMsgs, err := w.rabbitCh.Consume(
		"order.cancelled",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register cancelled consumer: %w", err)
	}

	log.Println("Order Worker started, waiting for payment events...")

	for {
		select {
		case msg, ok := <-paymentMsgs:
			if !ok {
				return nil
			}
			w.processPaymentCompleted(ctx, msg)
		case msg, ok := <-cancelledMsgs:
			if !ok {
				return nil
			}
			w.processCancelled(ctx, msg)
		case <-ctx.Done():
			return nil
		}
	}
}

// processPaymentCompleted — обработка успешной оплаты
// Обновляет статус заказа с pending на paid
func (w *OrderWorker) processPaymentCompleted(ctx context.Context, msg amqp.Delivery) {
	log.Printf("Received payment.completed: %s", msg.Body)

	var event map[string]interface{}
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to parse message: %v", err)
		msg.Ack(false)
		return
	}

	orderIDStr, ok := event["order_id"].(string)
	if !ok {
		log.Println("Missing order_id in payment_completed event")
		msg.Ack(false)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		log.Printf("Invalid order_id: %v", err)
		msg.Ack(false)
		return
	}

	// Обновляем статус на "paid"
	err = w.orderService.UpdateOrderStatus(ctx, orderID, "paid", "system", "Оплата получена")
	if err != nil {
		log.Printf("Failed to update order status to paid: %v", err)
		msg.Nack(false, true)
		return
	}

	log.Printf("Order %s status updated to: paid", orderID)
	msg.Ack(false)
}

// processCancelled — обработка отмены заказа
func (w *OrderWorker) processCancelled(ctx context.Context, msg amqp.Delivery) {
	log.Printf("Received order.cancelled: %s", msg.Body)

	var event map[string]interface{}
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to parse message: %v", err)
		msg.Ack(false)
		return
	}

	orderIDStr, ok := event["order_id"].(string)
	if !ok {
		log.Println("Missing order_id in cancelled event")
		msg.Ack(false)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		log.Printf("Invalid order_id: %v", err)
		msg.Ack(false)
		return
	}

	// Обновляем статус на "cancelled"
	err = w.orderService.UpdateOrderStatus(ctx, orderID, "cancelled", "system", "Заказ отменён")
	if err != nil {
		log.Printf("Failed to update order status to cancelled: %v", err)
		msg.Nack(false, true)
		return
	}

	log.Printf("Order %s status updated to: cancelled", orderID)
	msg.Ack(false)
}