package repository

import (
	"context"
	"errors"

	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrReviewNotFound = errors.New("review not found")
	ErrReviewExists   = errors.New("review already exists for this order")
)

type ReviewRepository struct {
	db *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// CreateReview — создаёт отзыв
func (r *ReviewRepository) CreateReview(ctx context.Context, review *models.Review) error {
	// Проверяем, есть ли уже отзыв для этого заказа
	if review.OrderID != nil {
		query := `SELECT EXISTS(SELECT 1 FROM reviews WHERE order_id = $1)`
		var exists bool
		err := r.db.QueryRow(ctx, query, review.OrderID).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			return ErrReviewExists
		}
	}

	query := `
		INSERT INTO reviews (id, product_id, user_id, order_id, rating, comment, is_approved, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		review.ID,
		review.ProductID,
		review.UserID,
		review.OrderID,
		review.Rating,
		review.Comment,
		review.IsApproved,
		review.CreatedAt,
		review.UpdatedAt,
	)
	return err
}

// GetReviewsByProductID — получает все отзывы на товар
func (r *ReviewRepository) GetReviewsByProductID(ctx context.Context, productID uuid.UUID) ([]models.ReviewWithUser, error) {
	query := `
		SELECT 
			r.id, r.product_id, r.user_id, r.order_id, r.rating, r.comment, r.is_approved, r.created_at, r.updated_at,
			u.first_name || ' ' || u.last_name as user_name,
			u.email as user_email
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		WHERE r.product_id = $1 AND r.is_approved = true
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ReviewWithUser
	for rows.Next() {
		var rev models.ReviewWithUser
		err := rows.Scan(
			&rev.ID,
			&rev.ProductID,
			&rev.UserID,
			&rev.OrderID,
			&rev.Rating,
			&rev.Comment,
			&rev.IsApproved,
			&rev.CreatedAt,
			&rev.UpdatedAt,
			&rev.UserName,
			&rev.UserEmail,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}

// GetReviewsByUserID — получает все отзывы пользователя
func (r *ReviewRepository) GetReviewsByUserID(ctx context.Context, userID uuid.UUID) ([]models.ReviewWithUser, error) {
	query := `
		SELECT 
			r.id, r.product_id, r.user_id, r.order_id, r.rating, r.comment, r.is_approved, r.created_at, r.updated_at,
			u.first_name || ' ' || u.last_name as user_name,
			u.email as user_email
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ReviewWithUser
	for rows.Next() {
		var rev models.ReviewWithUser
		err := rows.Scan(
			&rev.ID,
			&rev.ProductID,
			&rev.UserID,
			&rev.OrderID,
			&rev.Rating,
			&rev.Comment,
			&rev.IsApproved,
			&rev.CreatedAt,
			&rev.UpdatedAt,
			&rev.UserName,
			&rev.UserEmail,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}

// GetReviewByID — получает отзыв по ID
func (r *ReviewRepository) GetReviewByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	query := `
		SELECT id, product_id, user_id, order_id, rating, comment, is_approved, created_at, updated_at
		FROM reviews
		WHERE id = $1
	`

	var rev models.Review
	err := r.db.QueryRow(ctx, query, id).Scan(
		&rev.ID,
		&rev.ProductID,
		&rev.UserID,
		&rev.OrderID,
		&rev.Rating,
		&rev.Comment,
		&rev.IsApproved,
		&rev.CreatedAt,
		&rev.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}
	return &rev, nil
}

// UpdateReview — обновляет отзыв
func (r *ReviewRepository) UpdateReview(ctx context.Context, review *models.Review) error {
	query := `
		UPDATE reviews 
		SET rating = $1, comment = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query, review.Rating, review.Comment, review.UpdatedAt, review.ID)
	return err
}

// DeleteReview — удаляет отзыв
func (r *ReviewRepository) DeleteReview(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// ApproveReview — одобряет отзыв (для админа)
func (r *ReviewRepository) ApproveReview(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE reviews SET is_approved = true, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// GetAverageRating — получает средний рейтинг товара
func (r *ReviewRepository) GetAverageRating(ctx context.Context, productID uuid.UUID) (float64, int, error) {
	query := `
		SELECT COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as count
		FROM reviews
		WHERE product_id = $1 AND is_approved = true
	`
	var avgRating float64
	var count int
	err := r.db.QueryRow(ctx, query, productID).Scan(&avgRating, &count)
	return avgRating, count, err
}

// GetReviewByUserAndProduct — получает отзыв пользователя на товар
func (r *ReviewRepository) GetReviewByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) (*models.Review, error) {
    query := `
        SELECT id, product_id, user_id, order_id, rating, comment, is_approved, created_at, updated_at
        FROM reviews
        WHERE user_id = $1 AND product_id = $2
    `

    var review models.Review
    err := r.db.QueryRow(ctx, query, userID, productID).Scan(
        &review.ID,
        &review.ProductID,
        &review.UserID,
        &review.OrderID,
        &review.Rating,
        &review.Comment,
        &review.IsApproved,
        &review.CreatedAt,
        &review.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil // отзыва нет
        }
        return nil, err
    }
    return &review, nil
}