package repository

import (
	"context"
	"log"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// LogEmail — логирует отправку email (эмуляция)
func (r *NotificationRepository) LogEmail(ctx context.Context, notification *models.EmailNotification) error {
	log.Printf("EMAIL NOTIFICATION:")
	log.Printf("   To: %s", notification.To)
	log.Printf("   Subject: %s", notification.Subject)
	log.Printf("   Type: %s", notification.Type)
	log.Printf("   Body: %s", notification.Body)
	log.Println("   --- Email sent (emulated) ---")
	return nil
}