package repository

import (
	"context"
	"errors"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
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

// GetUsersListWithFilters — получает список пользователей с фильтрацией и поиском
func (r *AdminRepository) GetUsersListWithFilters(ctx context.Context, req *models.UsersListRequest) ([]models.UserDetails, int64, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	query := `
		SELECT 
			u.id, u.email, u.phone, u.first_name, u.last_name, u.role, u.is_active, u.created_at, u.updated_at,
			s.id as shop_id, s.name as shop_name, s.is_verified, s.rating
		FROM users u
		LEFT JOIN shops s ON s.seller_id = u.id
		WHERE 1=1
	`

	countQuery := `SELECT COUNT(*) FROM users u WHERE 1=1`
	args := []interface{}{}
	countArgs := []interface{}{}
	argIndex := 1
	countArgIndex := 1

	// Фильтр по роли
	if req.Role != "" {
		query += " AND u.role = $" + string(rune('0'+argIndex))
		args = append(args, req.Role)
		argIndex++
		countQuery += " AND u.role = $" + string(rune('0'+countArgIndex))
		countArgs = append(countArgs, req.Role)
		countArgIndex++
	}

	// Поиск по email или имени
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query += " AND (u.email ILIKE $" + string(rune('0'+argIndex)) + " OR u.first_name ILIKE $" + string(rune('0'+argIndex+1)) + " OR u.last_name ILIKE $" + string(rune('0'+argIndex+2)) + ")"
		args = append(args, searchTerm, searchTerm, searchTerm)
		argIndex += 3

		countQuery += " AND (u.email ILIKE $" + string(rune('0'+countArgIndex)) + " OR u.first_name ILIKE $" + string(rune('0'+countArgIndex+1)) + " OR u.last_name ILIKE $" + string(rune('0'+countArgIndex+2)) + ")"
		countArgs = append(countArgs, searchTerm, searchTerm, searchTerm)
		countArgIndex += 3
	}

	// Фильтр по статусу
	if req.IsActive != nil {
		query += " AND u.is_active = $" + string(rune('0'+argIndex))
		args = append(args, *req.IsActive)
		argIndex++
		countQuery += " AND u.is_active = $" + string(rune('0'+countArgIndex))
		countArgs = append(countArgs, *req.IsActive)
		countArgIndex++
	}

	query += " ORDER BY u.created_at DESC LIMIT $" + string(rune('0'+argIndex)) + " OFFSET $" + string(rune('0'+argIndex+1))
	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.UserDetails
	for rows.Next() {
		var u models.UserDetails
		var shopID *uuid.UUID
		var shopName *string
		var shopVerified *bool
		var shopRating *float64

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
			&shopID,
			&shopName,
			&shopVerified,
			&shopRating,
		)
		if err != nil {
			return nil, 0, err
		}

		if shopID != nil {
			u.Shop = &models.ShopInfo{
				ID:         *shopID,
				Name:       *shopName,
				IsVerified: *shopVerified,
				Rating:     *shopRating,
			}
		}

		users = append(users, u)
	}

	// Подсчёт общего количества
	var total int64
	err = r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByIDForAdmin — получает детальную информацию о пользователе для админа
func (r *AdminRepository) GetUserByIDForAdmin(ctx context.Context, userID uuid.UUID) (*models.UserDetails, error) {
	query := `
		SELECT 
			u.id, u.email, u.phone, u.first_name, u.last_name, u.role, u.is_active, u.created_at, u.updated_at,
			s.id as shop_id, s.name as shop_name, s.is_verified, s.rating
		FROM users u
		LEFT JOIN shops s ON s.seller_id = u.id
		WHERE u.id = $1
	`

	var u models.UserDetails
	var shopID *uuid.UUID
	var shopName *string
	var shopVerified *bool
	var shopRating *float64

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&u.ID,
		&u.Email,
		&u.Phone,
		&u.FirstName,
		&u.LastName,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
		&shopID,
		&shopName,
		&shopVerified,
		&shopRating,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if shopID != nil {
		u.Shop = &models.ShopInfo{
			ID:         *shopID,
			Name:       *shopName,
			IsVerified: *shopVerified,
			Rating:     *shopRating,
		}
	}

	return &u, nil
}

// GetShopByID — получает магазин по ID
func (r *AdminRepository) GetShopByID(ctx context.Context, shopID uuid.UUID) (*models.Shop, error) {
    query := `SELECT id, name, description, seller_id, is_verified, rating, created_at, updated_at FROM shops WHERE id = $1`
    var shop models.Shop
    err := r.db.QueryRow(ctx, query, shopID).Scan(
        &shop.ID,
        &shop.Name,
        &shop.Description,
        &shop.SellerID,
        &shop.IsVerified,
        &shop.Rating,
        &shop.CreatedAt,
        &shop.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &shop, nil
}

// UpdateShopName — обновляет название магазина
func (r *AdminRepository) UpdateShopName(ctx context.Context, shopID uuid.UUID, name string) error {
    query := `UPDATE shops SET name = $1, updated_at = NOW() WHERE id = $2`
    _, err := r.db.Exec(ctx, query, name, shopID)
    return err
}

// GetShopIDBySellerID — возвращает shop_id продавца
func (r *AdminRepository) GetShopIDBySellerID(ctx context.Context, sellerID uuid.UUID) (uuid.UUID, error) {
    query := `SELECT id FROM shops WHERE seller_id = $1`
    var shopID uuid.UUID
    err := r.db.QueryRow(ctx, query, sellerID).Scan(&shopID)
    if err != nil {
        return uuid.Nil, err
    }
    return shopID, nil
}
