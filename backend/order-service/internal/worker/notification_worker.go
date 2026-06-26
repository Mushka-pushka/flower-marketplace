package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationWorker struct {
	notifService *service.NotificationService
	rabbitCh     *amqp.Channel
}

func NewNotificationWorker(notifService *service.NotificationService, rabbitCh *amqp.Channel) *NotificationWorker {
	return &NotificationWorker{
		notifService: notifService,
		rabbitCh:     rabbitCh,
	}
}

// Start — запуск воркера уведомлений
func (w *NotificationWorker) Start(ctx context.Context) error {
	queue, err := w.rabbitCh.QueueDeclare(
		"notification.email",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
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
		return err
	}

	log.Println("Notification Worker started, waiting for emails...")

	for msg := range msgs {
		var notification models.EmailNotification
		if err := json.Unmarshal(msg.Body, &notification); err != nil {
			log.Printf("Failed to parse notification: %v", err)
			msg.Ack(false)
			continue
		}

		err := w.notifService.SendEmailNotification(ctx, &notification)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			msg.Nack(false, true)
		} else {
			msg.Ack(false)
		}
	}

	return nil
}