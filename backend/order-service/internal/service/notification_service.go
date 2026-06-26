package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationService struct {
	notifRepo *repository.NotificationRepository
	cfg       *config.Config
	rabbitCh  *amqp.Channel
}

func NewNotificationService(
	notifRepo *repository.NotificationRepository,
	cfg *config.Config,
	rabbitCh *amqp.Channel,
) *NotificationService {
	return &NotificationService{
		notifRepo: notifRepo,
		cfg:       cfg,
		rabbitCh:  rabbitCh,
	}
}

// SendEmailNotification — отправляет email-уведомление (эмуляция)
func (s *NotificationService) SendEmailNotification(ctx context.Context, notif *models.EmailNotification) error {
	return s.notifRepo.LogEmail(ctx, notif)
}

// PublishOrderEmail — публикация события для отправки email
func (s *NotificationService) PublishOrderEmail(orderID, userEmail, userName, status, eventType string) {
	subject := "Статус вашего заказа обновлён"
	body := fmt.Sprintf(
		"Здравствуйте, %s!\n\nСтатус вашего заказа #%s изменён на: %s.\n\nСпасибо, что выбираете наш маркетплейс!",
		userName, orderID[:8], status,
	)

	notification := &models.EmailNotification{
		To:      userEmail,
		Subject: subject,
		Body:    body,
		Type:    eventType,
	}

	data, _ := json.Marshal(notification)
	err := s.rabbitCh.Publish(
		"",
		"notification.email",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		log.Printf("Failed to publish email notification: %v", err)
	} else {
		log.Printf("Email notification queued for %s", userEmail)
	}
}