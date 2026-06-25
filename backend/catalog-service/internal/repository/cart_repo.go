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
	ErrCartItemNotFound = errors.New("cart item not found")
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

// AddToCart — добавляет товар в корзину (или увеличивает количество)
func (r *CartRepository) AddToCart(ctx context.Context, userID, productID uuid.UUID, quantity int) error {
	// Проверяем, есть ли уже такой товар в корзине
	query := `SELECT id, quantity FROM cart_items WHERE user_id = $1 AND product_id = $2`
	var id uuid.UUID
	var existingQuantity int
	err := r.db.QueryRow(ctx, query, userID, productID).Scan(&id, &existingQuantity)

	if err == nil {
		// Товар уже есть — обновляем количество
		newQuantity := existingQuantity + quantity
		updateQuery := `UPDATE cart_items SET quantity = $1, updated_at = NOW() WHERE id = $2`
		_, err = r.db.Exec(ctx, updateQuery, newQuantity, id)
		return err
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	// Товара нет — добавляем новый
	insertQuery := `
		INSERT INTO cart_items (id, user_id, product_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`
	_, err = r.db.Exec(ctx, insertQuery, uuid.New(), userID, productID, quantity)
	return err
}

// GetCartByUserID — получает корзину пользователя с данными о товарах
func (r *CartRepository) GetCartByUserID(ctx context.Context, userID uuid.UUID) ([]models.CartItemWithProduct, error) {
	query := `
		SELECT 
			ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
			p.name as product_name, 
			p.price as product_price,
			COALESCE((SELECT image_url FROM product_images WHERE product_id = p.id AND is_main = true LIMIT 1), '') as product_image
		FROM cart_items ci
		JOIN products p ON p.id = ci.product_id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItemWithProduct
	for rows.Next() {
		var item models.CartItemWithProduct
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ProductName,
			&item.ProductPrice,
			&item.ProductImage,
		)
		if err != nil {
			return nil, err
		}
		item.TotalPrice = float64(item.Quantity) * item.ProductPrice
		items = append(items, item)
	}
	return items, nil
}

// UpdateCartItemQuantity — обновляет количество товара в корзине
func (r *CartRepository) UpdateCartItemQuantity(ctx context.Context, cartItemID uuid.UUID, quantity int) error {
	query := `UPDATE cart_items SET quantity = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, quantity, cartItemID)
	return err
}

// RemoveFromCart — удаляет товар из корзины
func (r *CartRepository) RemoveFromCart(ctx context.Context, cartItemID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := r.db.Exec(ctx, query, cartItemID)
	return err
}

// ClearCart — очищает корзину пользователя
func (r *CartRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

// GetCartItemByID — получает элемент корзины по ID
func (r *CartRepository) GetCartItemByID(ctx context.Context, id uuid.UUID) (*models.CartItem, error) {
	query := `
		SELECT id, user_id, product_id, quantity, created_at, updated_at
		FROM cart_items
		WHERE id = $1
	`

	var item models.CartItem
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.UserID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCartItemNotFound
		}
		return nil, err
	}
	return &item, nil
}