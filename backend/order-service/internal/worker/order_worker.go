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
	// Объявляем очередь
	queue, err := w.rabbitCh.QueueDeclare(
		"order.created", // имя очереди
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Начинаем потребление сообщений
	msgs, err := w.rabbitCh.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack (отключаем, чтобы подтверждать вручную)
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Println("Order Worker started, waiting for messages...")

	// Обрабатываем сообщения
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
	default:
		log.Printf("Unknown event type: %s", eventType)
		msg.Ack(false)
	}
}

// handleOrderCreated — обработка события order.created
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

	// Имитация обработки заказа
	log.Printf("Processing order %s", orderID)

	// Шаг 1: Подтверждение заказа
	time.Sleep(2 * time.Second)
	err = w.orderService.UpdateOrderStatus(ctx, orderID, "confirmed", "system", "Заказ подтверждён")
	if err != nil {
		log.Printf("Failed to update status to confirmed: %v", err)
		msg.Nack(false, true) // retry
		return
	}
	log.Printf("Order %s status updated to: confirmed", orderID)

	// Шаг 2: Подготовка заказа
	time.Sleep(3 * time.Second)
	err = w.orderService.UpdateOrderStatus(ctx, orderID, "preparing", "system", "Заказ готовится")
	if err != nil {
		log.Printf("Failed to update status to preparing: %v", err)
		msg.Nack(false, true)
		return
	}
	log.Printf("Order %s status updated to: preparing", orderID)

	// Шаг 3: Упаковка
	time.Sleep(2 * time.Second)
	err = w.orderService.UpdateOrderStatus(ctx, orderID, "packing", "system", "Заказ упаковывается")
	if err != nil {
		log.Printf("Failed to update status to packing: %v", err)
		msg.Nack(false, true)
		return
	}
	log.Printf("Order %s status updated to: packing", orderID)

	// Шаг 4: Доставка
	time.Sleep(2 * time.Second)
	err = w.orderService.UpdateOrderStatus(ctx, orderID, "delivery", "system", "Заказ передан курьеру")
	if err != nil {
		log.Printf("Failed to update status to delivery: %v", err)
		msg.Nack(false, true)
		return
	}
	log.Printf("Order %s status updated to: delivery", orderID)

	// Шаг 5: Доставлен
	time.Sleep(3 * time.Second)
	err = w.orderService.UpdateOrderStatus(ctx, orderID, "delivered", "system", "Заказ доставлен получателю")
	if err != nil {
		log.Printf("Failed to update status to delivered: %v", err)
		msg.Nack(false, true)
		return
	}
	log.Printf("Order %s status updated to: delivered", orderID)

	// Подтверждаем успешную обработку
	msg.Ack(false)
	log.Printf("Order %s processed successfully!", orderID)
}