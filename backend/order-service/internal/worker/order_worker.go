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
	queue, err := w.rabbitCh.QueueDeclare(
		"order.created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := w.rabbitCh.Consume(
		queue.Name,
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

	log.Printf("Processing order %s", orderID)

	// Получаем текущий статус заказа
	order, err := w.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		msg.Nack(false, true)
		return
	}

	currentStatus := order.Order.CurrentStatus

	// Если заказ уже отменён или доставлен — завершаем
	if currentStatus == "cancelled" || currentStatus == "delivered" {
		log.Printf("Order %s already %s, skipping", orderID, currentStatus)
		msg.Ack(false)
		return
	}

	// ШАГ 1: Ждём подтверждения продавца (pending → confirmed)
	if currentStatus == "pending" {
		log.Printf("Order %s waiting for seller confirmation...", orderID)
		for {
			time.Sleep(2 * time.Second)
			order, err = w.orderService.GetOrderByID(ctx, orderID)
			if err != nil {
				log.Printf("Failed to get order: %v", err)
				msg.Nack(false, true)
				return
			}
			if order.Order.CurrentStatus == "confirmed" {
				log.Printf("Order %s confirmed by seller, proceeding...", orderID)
				break
			}
			if order.Order.CurrentStatus == "cancelled" || order.Order.CurrentStatus == "delivered" {
				log.Printf("Order %s %s, stopping worker", orderID, order.Order.CurrentStatus)
				msg.Ack(false)
				return
			}
		}
	}

	// Проверяем статус снова
	currentStatus = order.Order.CurrentStatus
	if currentStatus == "cancelled" || currentStatus == "delivered" {
		log.Printf("Order %s %s, skipping", orderID, currentStatus)
		msg.Ack(false)
		return
	}

	// ШАГ 2: Подготовка (confirmed → preparing)
	if currentStatus == "confirmed" {
		time.Sleep(3 * time.Second)
		order, err = w.orderService.GetOrderByID(ctx, orderID)
		if err != nil {
			log.Printf("Failed to get order: %v", err)
			msg.Nack(false, true)
			return
		}
		if order.Order.CurrentStatus == "confirmed" {
			err = w.orderService.UpdateOrderStatus(ctx, orderID, "preparing", "system", "Заказ готовится")
			if err != nil {
				log.Printf("Failed to update status to preparing: %v", err)
				msg.Nack(false, true)
				return
			}
			log.Printf("Order %s status updated to: preparing", orderID)
		} else {
			log.Printf("Order %s status changed by seller, skipping", orderID)
		}
	}

	// Проверяем статус снова
	order, err = w.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		msg.Nack(false, true)
		return
	}
	currentStatus = order.Order.CurrentStatus
	if currentStatus == "cancelled" || currentStatus == "delivered" {
		log.Printf("Order %s %s, skipping", orderID, currentStatus)
		msg.Ack(false)
		return
	}

	// ШАГ 3: Упаковка (preparing → packing)
	if currentStatus == "preparing" {
		time.Sleep(2 * time.Second)
		order, err = w.orderService.GetOrderByID(ctx, orderID)
		if err != nil {
			log.Printf("Failed to get order: %v", err)
			msg.Nack(false, true)
			return
		}
		if order.Order.CurrentStatus == "preparing" {
			err = w.orderService.UpdateOrderStatus(ctx, orderID, "packing", "system", "Заказ упаковывается")
			if err != nil {
				log.Printf("Failed to update status to packing: %v", err)
				msg.Nack(false, true)
				return
			}
			log.Printf("Order %s status updated to: packing", orderID)
		} else {
			log.Printf("Order %s status changed by seller, skipping", orderID)
		}
	}

	// Проверяем статус снова
	order, err = w.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		msg.Nack(false, true)
		return
	}
	currentStatus = order.Order.CurrentStatus
	if currentStatus == "cancelled" || currentStatus == "delivered" {
		log.Printf("Order %s %s, skipping", orderID, currentStatus)
		msg.Ack(false)
		return
	}

	// ШАГ 4: Доставка (packing → delivery)
	if currentStatus == "packing" {
		time.Sleep(2 * time.Second)
		order, err = w.orderService.GetOrderByID(ctx, orderID)
		if err != nil {
			log.Printf("Failed to get order: %v", err)
			msg.Nack(false, true)
			return
		}
		if order.Order.CurrentStatus == "packing" {
			// Обновляем статус на delivery
			err = w.orderService.UpdateOrderStatus(ctx, orderID, "delivery", "system", "Заказ передан курьеру")
			if err != nil {
				log.Printf("Failed to update status to delivery: %v", err)
				msg.Nack(false, true)
				return
			}
			log.Printf("Order %s status updated to: delivery", orderID)

			// Назначаем курьера и отправляем событие
			courier, err := w.orderService.AssignCourier(ctx, orderID)
			if err != nil {
				log.Printf("Failed to assign courier: %v", err)
			} else {
				log.Printf("Courier %s assigned to order %s", courier.Name, orderID)
			}
		} else {
			log.Printf("Order %s status changed by seller, skipping", orderID)
		}
	}

	// Проверяем статус снова
	order, err = w.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		msg.Nack(false, true)
		return
	}
	currentStatus = order.Order.CurrentStatus
	if currentStatus == "cancelled" || currentStatus == "delivered" {
		log.Printf("Order %s %s, skipping", orderID, currentStatus)
		msg.Ack(false)
		return
	}

	// ШАГ 5: Доставлен (delivery → delivered)
	if currentStatus == "delivery" {
		time.Sleep(3 * time.Second)
		order, err = w.orderService.GetOrderByID(ctx, orderID)
		if err != nil {
			log.Printf("Failed to get order: %v", err)
			msg.Nack(false, true)
			return
		}
		if order.Order.CurrentStatus == "delivery" {
			// Отмечаем завершение доставки
			err = w.orderService.CompleteDelivery(ctx, orderID)
			if err != nil {
				log.Printf("Failed to complete delivery: %v", err)
				msg.Nack(false, true)
				return
			}

			// Обновляем статус на delivered
			err = w.orderService.UpdateOrderStatus(ctx, orderID, "delivered", "system", "Заказ доставлен получателю")
			if err != nil {
				log.Printf("Failed to update status to delivered: %v", err)
				msg.Nack(false, true)
				return
			}
			log.Printf("Order %s status updated to: delivered", orderID)
		} else {
			log.Printf("Order %s status changed by seller, skipping", orderID)
		}
	}

	// Подтверждаем успешную обработку
	msg.Ack(false)
	log.Printf("Order %s processed successfully!", orderID)
}