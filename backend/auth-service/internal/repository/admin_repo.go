package repository

import (
	"context"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

// GetSellersWithShops — получает всех продавцов с данными магазинов
func (r *AdminRepository) GetSellersWithShops(ctx context.Context, verified *bool) ([]models.SellerWithUser, error) {
	query := `
		SELECT 
			s.id as shop_id,
			s.name as shop_name,
			s.description as shop_description,
			s.is_verified,
			s.rating,
			u.id as user_id,
			u.email,
			u.phone,
			u.first_name,
			u.last_name,
			u.is_active,
			u.created_at
		FROM shops s
		JOIN users u ON u.id = s.seller_id
		WHERE u.role = 'seller'
	`

	if verified != nil {
		if *verified {
			query += " AND s.is_verified = true"
		} else {
			query += " AND s.is_verified = false"
		}
	}

	query += " ORDER BY s.created_at DESC"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sellers []models.SellerWithUser
	for rows.Next() {
		var s models.SellerWithUser
		err := rows.Scan(
			&s.ShopID,
			&s.ShopName,
			&s.ShopDesc,
			&s.IsVerified,
			&s.Rating,
			&s.UserID,
			&s.Email,
			&s.Phone,
			&s.FirstName,
			&s.LastName,
			&s.IsActive,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sellers = append(sellers, s)
	}
	return sellers, nil
}

// VerifyShop — верифицирует магазин
func (r *AdminRepository) VerifyShop(ctx context.Context, shopID uuid.UUID, verified bool) error {
	query := `UPDATE shops SET is_verified = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, verified, shopID)
	return err
}

// UpdateUserStatus — обновляет статус пользователя
func (r *AdminRepository) UpdateUserStatus(ctx context.Context, userID uuid.UUID, isActive bool) error {
	query := `UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, isActive, userID)
	return err
}

// GetUsersList — получает список пользователей с фильтрацией
func (r *AdminRepository) GetUsersList(ctx context.Context, role string, isActive *bool, limit, offset int) ([]models.User, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, email, phone, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if role != "" {
		query += " AND role = $" + string(rune('0'+argIndex))
		args = append(args, role)
		argIndex++
	}
	if isActive != nil {
		query += " AND is_active = $" + string(rune('0'+argIndex))
		args = append(args, *isActive)
		argIndex++
	}

	query += " ORDER BY created_at DESC LIMIT $" + string(rune('0'+argIndex)) + " OFFSET $" + string(rune('0'+argIndex+1))
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Phone,
			&u.FirstName,
			&u.LastName,
			&u.Role,
			&u.IsActive,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}