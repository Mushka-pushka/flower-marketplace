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

    log.Println("Order Worker started, waiting for messages...")

    for {
        select {
        case msg, ok := <-msgs:
            if !ok {
                return nil
            }
            w.processMessage(ctx, msg)
        case msg, ok := <-paymentMsgs:
            if !ok {
                return nil
            }
            w.processMessage(ctx, msg)
        case msg, ok := <-cancelledMsgs:
            if !ok {
                return nil
            }
            w.processMessage(ctx, msg)
        case <-ctx.Done():
            return nil
        }
    }
}

// processMessage — обработка одного сообщения с retry
func (w *OrderWorker) processMessage(ctx context.Context, msg amqp.Delivery) {
    log.Printf("Received message: %s", msg.Body)

    go func() {
        retryCount := 0
        if msg.Headers != nil {
            if retry, ok := msg.Headers["x-retry-count"].(int64); ok {
                retryCount = int(retry)
            }
        }

        maxRetries := 3

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

        var err error
        switch eventType {
        case "order.created":
            err = w.handleOrderCreatedWithRetry(ctx, event, retryCount)
        case "order.payment_completed":
            err = w.handlePaymentCompletedWithRetry(ctx, event, retryCount)
        default:
            log.Printf("Unknown event type: %s", eventType)
            msg.Ack(false)
            return
        }

        if err != nil {
            retryCount++
            if retryCount < maxRetries {
                log.Printf("Retrying message (attempt %d/%d): %v", retryCount+1, maxRetries, err)

                headers := amqp.Table{}
                if msg.Headers != nil {
                    for k, v := range msg.Headers {
                        headers[k] = v
                    }
                }
                headers["x-retry-count"] = int64(retryCount)

                msg.Nack(false, true)
                return
            }

            log.Printf("Failed after %d retries: %v", maxRetries, err)
            msg.Ack(false)
            return
        }

        msg.Ack(false)
    }()
}

// handlePaymentCompletedWithRetry — обработка события order.payment_completed с retry
func (w *OrderWorker) handlePaymentCompletedWithRetry(ctx context.Context, event map[string]interface{}, retryCount int) error {
	orderIDStr, ok := event["order_id"].(string)
	if !ok {
		return fmt.Errorf("missing order_id in payment_completed event")
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return fmt.Errorf("invalid order_id in payment_completed event: %w", err)
	}

	status, ok := event["status"].(string)
	if !ok {
		status = "paid" // статус по умолчанию
	}

	log.Printf("Payment completed for order %s, updating status to %s (retry %d)", orderID, status, retryCount)

	// Обновляем статус заказа на "paid"
	err = w.orderService.UpdateOrderStatus(ctx, orderID, "paid", "system", "Оплата получена")
	if err != nil {
		return fmt.Errorf("failed to update order status to paid: %w", err)
	}

	log.Printf("Order %s status updated to: paid", orderID)
	return nil
}

// handleOrderCreatedWithRetry — обработка события order.created с retry
func (w *OrderWorker) handleOrderCreatedWithRetry(ctx context.Context, event map[string]interface{}, retryCount int) error {
	orderIDStr, ok := event["order_id"].(string)
	if !ok {
		return fmt.Errorf("missing order_id")
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return fmt.Errorf("invalid order_id: %w", err)
	}

	log.Printf("Processing order %s (retry %d)", orderID, retryCount)

	// Проверяем статус заказа
	order, err := w.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	currentStatus := order.Order.CurrentStatus

	if currentStatus == "cancelled" || currentStatus == "delivered" {
		log.Printf("Order %s already %s, skipping", orderID, currentStatus)
		return nil
	}

	// ШАГ 1: Ждём оплаты (pending → paid)
	if currentStatus == "pending" {
		log.Printf("Order %s waiting for payment...", orderID)

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		timeout := time.After(2 * time.Minute)

		for {
			select {
			case <-ticker.C:
				order, err = w.orderService.GetOrderByID(ctx, orderID)
				if err != nil {
					return fmt.Errorf("failed to get order: %w", err)
				}

				if order.Order.CurrentStatus == "paid" {
					log.Printf("Order %s paid, auto-confirming...", orderID)
					
					
					err = w.orderService.UpdateOrderStatus(ctx, orderID, "confirmed", "system", "Автоматическое подтверждение после оплаты")
					if err != nil {
						return fmt.Errorf("failed to confirm order: %w", err)
					}
					currentStatus = "confirmed"
					goto processFlow
				}

				if order.Order.CurrentStatus == "cancelled" {
					log.Printf("Order %s cancelled, stopping worker", orderID)
					return nil
				}

			case <-timeout:
				log.Printf("Order %s payment timeout, cancelling", orderID)
				err = w.orderService.UpdateOrderStatus(ctx, orderID, "cancelled", "system", "Превышено время ожидания оплаты")
				if err != nil {
					return fmt.Errorf("failed to cancel order: %w", err)
				}
				return nil

			case <-ctx.Done():
				log.Printf("Context cancelled for order %s", orderID)
				return ctx.Err()
			}
		}
	}

processFlow:
	// ШАГ 2: Автоматическое обновление статусов с задержками
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
		order, err = w.orderService.GetOrderByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}

		currentStatus = order.Order.CurrentStatus

		if currentStatus == "cancelled" || currentStatus == "delivered" {
			log.Printf("Order %s %s, stopping worker", orderID, currentStatus)
			return nil
		}

		if currentStatus != step.from {
			log.Printf("Order %s status changed to %s, skipping step %s -> %s", orderID, currentStatus, step.from, step.to)
			continue
		}

		select {
		case <-time.After(step.delay):
			order, err = w.orderService.GetOrderByID(ctx, orderID)
			if err != nil {
				return fmt.Errorf("failed to get order: %w", err)
			}

			if order.Order.CurrentStatus != step.from {
				log.Printf("Order %s status changed to %s, skipping step", orderID, order.Order.CurrentStatus)
				continue
			}

			err = w.orderService.UpdateOrderStatus(ctx, orderID, step.to, "system", step.comment)
			if err != nil {
				return fmt.Errorf("failed to update status to %s: %w", step.to, err)
			}
			log.Printf("Order %s status updated to: %s", orderID, step.to)

		case <-ctx.Done():
			log.Printf("Context cancelled for order %s", orderID)
			return ctx.Err()
		}
	}

	log.Printf("Order %s processed successfully!", orderID)
	return nil
}