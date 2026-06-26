package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreatePayment — создаёт платёж
func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (id, order_id, amount, status, payment_method, transaction_id, payment_url, completed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query,
		payment.ID,
		payment.OrderID,
		payment.Amount,
		payment.Status,
		payment.PaymentMethod,
		payment.TransactionID,
		payment.PaymentURL,
		payment.CompletedAt,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	return err
}

// GetPaymentByID — получает платёж по ID
func (r *PaymentRepository) GetPaymentByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	query := `
		SELECT id, order_id, amount, status, payment_method, transaction_id, payment_url, completed_at, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	var payment models.Payment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.PaymentURL,
		&payment.CompletedAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// GetPaymentByOrderID — получает платёж по ID заказа
func (r *PaymentRepository) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error) {
	query := `
		SELECT id, order_id, amount, status, payment_method, transaction_id, payment_url, completed_at, created_at, updated_at
		FROM payments
		WHERE order_id = $1
	`

	var payment models.Payment
	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.PaymentURL,
		&payment.CompletedAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// UpdatePaymentStatus — обновляет статус платежа
func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status string, completedAt *time.Time) error {
	query := `UPDATE payments SET status = $1, completed_at = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, status, completedAt, time.Now(), id)
	return err
}