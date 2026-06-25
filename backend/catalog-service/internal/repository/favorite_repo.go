package repository

import (
	"context"
	"errors"

	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFavoriteNotFound = errors.New("favorite not found")
	ErrAlreadyFavorited = errors.New("product already in favorites")
)

type FavoriteRepository struct {
	db *pgxpool.Pool
}

func NewFavoriteRepository(db *pgxpool.Pool) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

// AddFavorite — добавляет товар в избранное
func (r *FavoriteRepository) AddFavorite(ctx context.Context, userID, productID uuid.UUID) error {
	// Проверяем, есть ли уже в избранном
	query := `SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND product_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, productID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyFavorited
	}

	insertQuery := `
		INSERT INTO favorites (id, user_id, product_id, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err = r.db.Exec(ctx, insertQuery, uuid.New(), userID, productID)
	return err
}

// GetFavoritesByUserID — получает все избранные товары пользователя
func (r *FavoriteRepository) GetFavoritesByUserID(ctx context.Context, userID uuid.UUID) ([]models.FavoriteWithProduct, error) {
	query := `
		SELECT 
			f.id, f.user_id, f.product_id, f.created_at,
			p.name as product_name,
			p.price as product_price,
			p.slug as product_slug,
			COALESCE((SELECT image_url FROM product_images WHERE product_id = p.id AND is_main = true LIMIT 1), '') as product_image
		FROM favorites f
		JOIN products p ON p.id = f.product_id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.FavoriteWithProduct
	for rows.Next() {
		var item models.FavoriteWithProduct
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ProductID,
			&item.CreatedAt,
			&item.ProductName,
			&item.ProductPrice,
			&item.ProductSlug,
			&item.ProductImage,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// RemoveFavorite — удаляет товар из избранного
func (r *FavoriteRepository) RemoveFavorite(ctx context.Context, userID, productID uuid.UUID) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND product_id = $2`
	_, err := r.db.Exec(ctx, query, userID, productID)
	return err
}

// IsFavorite — проверяет, находится ли товар в избранном у пользователя
func (r *FavoriteRepository) IsFavorite(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND product_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, productID).Scan(&exists)
	return exists, err
}