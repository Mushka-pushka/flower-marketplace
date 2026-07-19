package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

// Start — запуск воркера
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

	// Начинаем слушать очередь order.created
	msgs, err := w.rabbitCh.Consume(
		"order.created",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Println("Order Worker started, waiting for messages...")

	for msg := range msgs {
		w.processMessage(ctx, msg)
	}

	return nil
}

// processMessage — обработка одного сообщения
func (w *OrderWorker) processMessage(ctx context.Context, msg amqp.Delivery) {
	log.Printf("Received message: %s", msg.Body)

	var event map[string]interface{}
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to parse message: %v", err)
		msg.Ack(false)
		return
	}

	eventType, ok := event["event"].(string)
	if !ok {
		log.Println("Missing event field")
		msg.Ack(false)
		return
	}

	switch eventType {
	case "order.created":
		w.handleOrderCreated(ctx, event, msg)
	case "order.payment_completed":
		w.handlePaymentCompleted(ctx, event, msg)
	default:
		log.Printf("Unknown event type: %s", eventType)
		msg.Ack(false)
	}
}

// handlePaymentCompleted — обработка события order.payment_completed
func (w *OrderWorker) handlePaymentCompleted(ctx context.Context, event map[string]interface{}, msg amqp.Delivery) {
	orderIDStr, ok := event["order_id"].(string)
	if !ok {
		log.Println("Missing order_id in payment_completed event")
		msg.Ack(false)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		log.Printf("Invalid order_id in payment_completed event: %v", err)
		msg.Ack(false)
		return
	}

	status, ok := event["status"].(string)
	if !ok {
		status = "paid" // статус по умолчанию
	}

	log.Printf("Payment completed for order %s, updating status to %s", orderID, status)

	// Обновляем статус заказа на "paid"
	err = w.orderService.UpdateOrderStatus(ctx, orderID, "paid", "system", "Оплата получена")
	if err != nil {
		log.Printf("Failed to update order status to paid: %v", err)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
	log.Printf("Order %s status updated to: paid", orderID)
}

// handleOrderCreated — обработка события order.created (исправленная версия с таймерами)
func (w *OrderWorker) handleOrderCreated(ctx context.Context, event map[string]interface{}, msg amqp.Delivery) {
	orderIDStr, ok := event["order_id"].(string)
	if !ok {
		log.Println("Missing order_id")
		msg.Ack(false)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		log.Printf("Invalid order_id: %v", err)
		msg.Ack(false)
		return
	}

	log.Printf("Processing order %s", orderID)

	// Проверяем статус заказа
	order, err := w.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		msg.Nack(false, true)
		return
	}

	currentStatus := order.Order.CurrentStatus

	if currentStatus == "cancelled" || currentStatus == "delivered" {
		log.Printf("Order %s already %s, skipping", orderID, currentStatus)
		msg.Ack(false)
		return
	}

	// ШАГ 1: Ждём оплаты (pending → paid)
	if currentStatus == "pending" {
		log.Printf("Order %s waiting for payment...", orderID)

		// Используем ticker вместо sleep в цикле
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		timeout := time.After(5 * time.Minute) // Максимальное время ожидания оплаты

		for {
			select {
			case <-ticker.C:
				order, err = w.orderService.GetOrderByID(ctx, orderID)
				if err != nil {
					log.Printf("Failed to get order: %v", err)
					msg.Nack(false, true)
					return
				}

				if order.Order.CurrentStatus == "paid" {
					log.Printf("Order %s paid, proceeding...", orderID)
					currentStatus = "paid"
					goto waitConfirmation
				}

				if order.Order.CurrentStatus == "cancelled" {
					log.Printf("Order %s cancelled, stopping worker", orderID)
					msg.Ack(false)
					return
				}

			case <-timeout:
				log.Printf("Order %s payment timeout, cancelling", orderID)
				err = w.orderService.UpdateOrderStatus(ctx, orderID, "cancelled", "system", "Превышено время ожидания оплаты")
				if err != nil {
					log.Printf("Failed to cancel order: %v", err)
				}
				msg.Ack(false)
				return

			case <-ctx.Done():
				log.Printf("Context cancelled for order %s", orderID)
				msg.Ack(false)
				return
			}
		}
	}

waitConfirmation:
	// ШАГ 2: Ждём подтверждения продавца (paid → confirmed)
	if currentStatus == "paid" {
		log.Printf("Order %s waiting for seller confirmation...", orderID)

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		timeout := time.After(10 * time.Minute) // Максимальное время ожидания подтверждения

		for {
			select {
			case <-ticker.C:
				order, err = w.orderService.GetOrderByID(ctx, orderID)
				if err != nil {
					log.Printf("Failed to get order: %v", err)
					msg.Nack(false, true)
					return
				}

				if order.Order.CurrentStatus == "confirmed" {
					log.Printf("Order %s confirmed by seller, proceeding...", orderID)
					currentStatus = "confirmed"
					goto processFlow
				}

				if order.Order.CurrentStatus == "cancelled" || order.Order.CurrentStatus == "delivered" {
					log.Printf("Order %s %s, stopping worker", orderID, order.Order.CurrentStatus)
					msg.Ack(false)
					return
				}

			case <-timeout:
				log.Printf("Order %s confirmation timeout, cancelling", orderID)
				err = w.orderService.UpdateOrderStatus(ctx, orderID, "cancelled", "system", "Превышено время ожидания подтверждения продавцом")
				if err != nil {
					log.Printf("Failed to cancel order: %v", err)
				}
				msg.Ack(false)
				return

			case <-ctx.Done():
				log.Printf("Context cancelled for order %s", orderID)
				msg.Ack(false)
				return
			}
		}
	}

processFlow:
	// ШАГ 3: Автоматическое обновление статусов с задержками
	statusFlow := []struct {
		from    string
		to      string
		delay   time.Duration
		comment string
	}{
		{"confirmed", "preparing", 3 * time.Second, "Заказ готовится"},
		{"preparing", "packing", 2 * time.Second, "Заказ упаковывается"},
		{"packing", "delivery", 2 * time.Second, "Заказ передан курьеру"},
		{"delivery", "delivered", 3 * time.Second, "Заказ доставлен получателю"},
	}

	for _, step := range statusFlow {
		// Проверяем, что заказ всё ещё в нужном статусе
		order, err = w.orderService.GetOrderByID(ctx, orderID)
		if err != nil {
			log.Printf("Failed to get order: %v", err)
			msg.Nack(false, true)
			return
		}

		currentStatus = order.Order.CurrentStatus

		if currentStatus == "cancelled" || currentStatus == "delivered" {
			log.Printf("Order %s %s, stopping worker", orderID, currentStatus)
			msg.Ack(false)
			return
		}

		if currentStatus != step.from {
			log.Printf("Order %s status changed to %s, skipping step %s -> %s", orderID, currentStatus, step.from, step.to)
			continue
		}

		// Используем time.After вместо time.Sleep
		select {
		case <-time.After(step.delay):
			// Проверяем статус перед обновлением
			order, err = w.orderService.GetOrderByID(ctx, orderID)
			if err != nil {
				log.Printf("Failed to get order: %v", err)
				msg.Nack(false, true)
				return
			}

			if order.Order.CurrentStatus != step.from {
				log.Printf("Order %s status changed to %s, skipping step", orderID, order.Order.CurrentStatus)
				continue
			}

			err = w.orderService.UpdateOrderStatus(ctx, orderID, step.to, "system", step.comment)
			if err != nil {
				log.Printf("Failed to update status to %s: %v", step.to, err)
				msg.Nack(false, true)
				return
			}
			log.Printf("Order %s status updated to: %s", orderID, step.to)

		case <-ctx.Done():
			log.Printf("Context cancelled for order %s", orderID)
			msg.Ack(false)
			return
		}
	}

	msg.Ack(false)
	log.Printf("Order %s processed successfully!", orderID)
}