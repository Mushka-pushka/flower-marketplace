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
	ErrAddressNotFound = errors.New("address not found")
)

type AddressRepository struct {
	db *pgxpool.Pool
}

func NewAddressRepository(db *pgxpool.Pool) *AddressRepository {
	return &AddressRepository{db: db}
}

// CreateAddress — создаёт адрес
func (r *AddressRepository) CreateAddress(ctx context.Context, address *models.DeliveryAddress) error {
	// Если адрес по умолчанию — снимаем флаг с других адресов пользователя
	if address.IsDefault {
		_, err := r.db.Exec(ctx, `UPDATE delivery_addresses SET is_default = false WHERE user_id = $1`, address.UserID)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO delivery_addresses (id, user_id, name, address, entrance, floor, intercom, comment, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query,
		address.ID,
		address.UserID,
		address.Name,
		address.Address,
		address.Entrance,
		address.Floor,
		address.Intercom,
		address.Comment,
		address.IsDefault,
		address.CreatedAt,
		address.UpdatedAt,
	)
	return err
}

// GetAddressesByUserID — получает все адреса пользователя
func (r *AddressRepository) GetAddressesByUserID(ctx context.Context, userID uuid.UUID) ([]models.DeliveryAddress, error) {
	query := `
		SELECT id, user_id, name, address, entrance, floor, intercom, comment, is_default, created_at, updated_at
		FROM delivery_addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.DeliveryAddress
	for rows.Next() {
		var addr models.DeliveryAddress
		err := rows.Scan(
			&addr.ID,
			&addr.UserID,
			&addr.Name,
			&addr.Address,
			&addr.Entrance,
			&addr.Floor,
			&addr.Intercom,
			&addr.Comment,
			&addr.IsDefault,
			&addr.CreatedAt,
			&addr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

// GetAddressByID — получает адрес по ID
func (r *AddressRepository) GetAddressByID(ctx context.Context, id uuid.UUID) (*models.DeliveryAddress, error) {
	query := `
		SELECT id, user_id, name, address, entrance, floor, intercom, comment, is_default, created_at, updated_at
		FROM delivery_addresses
		WHERE id = $1
	`

	var addr models.DeliveryAddress
	err := r.db.QueryRow(ctx, query, id).Scan(
		&addr.ID,
		&addr.UserID,
		&addr.Name,
		&addr.Address,
		&addr.Entrance,
		&addr.Floor,
		&addr.Intercom,
		&addr.Comment,
		&addr.IsDefault,
		&addr.CreatedAt,
		&addr.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}
	return &addr, nil
}

// UpdateAddress — обновляет адрес
func (r *AddressRepository) UpdateAddress(ctx context.Context, address *models.DeliveryAddress) error {
	// Если адрес по умолчанию — снимаем флаг с других адресов пользователя
	if address.IsDefault {
		_, err := r.db.Exec(ctx, `UPDATE delivery_addresses SET is_default = false WHERE user_id = $1 AND id != $2`, address.UserID, address.ID)
		if err != nil {
			return err
		}
	}

	query := `
		UPDATE delivery_addresses SET
			name = $1, address = $2, entrance = $3, floor = $4, intercom = $5,
			comment = $6, is_default = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.Exec(ctx, query,
		address.Name,
		address.Address,
		address.Entrance,
		address.Floor,
		address.Intercom,
		address.Comment,
		address.IsDefault,
		address.UpdatedAt,
		address.ID,
	)
	return err
}

// DeleteAddress — удаляет адрес
func (r *AddressRepository) DeleteAddress(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM delivery_addresses WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// SetDefaultAddress — устанавливает адрес по умолчанию
func (r *AddressRepository) SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID) error {
	// Снимаем флаг со всех адресов пользователя
	_, err := r.db.Exec(ctx, `UPDATE delivery_addresses SET is_default = false WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// Устанавливаем флаг на нужный адрес
	query := `UPDATE delivery_addresses SET is_default = true WHERE id = $1 AND user_id = $2`
	_, err = r.db.Exec(ctx, query, addressID, userID)
	return err
}