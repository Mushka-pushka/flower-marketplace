package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminStatsRepository struct {
	db *pgxpool.Pool
}

func NewAdminStatsRepository(db *pgxpool.Pool) *AdminStatsRepository {
	return &AdminStatsRepository{db: db}
}

// GetUserStats — статистика по пользователям
func (r *AdminStatsRepository) GetUserStats(ctx context.Context) (total int64, byRole map[string]int64, err error) {
	// Общее количество пользователей
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return 0, nil, err
	}

	// Количество по ролям
	rows, err := r.db.Query(ctx, `SELECT role, COUNT(*) FROM users GROUP BY role`)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	byRole = make(map[string]int64)
	for rows.Next() {
		var role string
		var count int64
		err := rows.Scan(&role, &count)
		if err != nil {
			return 0, nil, err
		}
		byRole[role] = count
	}
	return total, byRole, nil
}

// GetShopStats — статистика по магазинам
func (r *AdminStatsRepository) GetShopStats(ctx context.Context) (total, verified int64, err error) {
	err = r.db.QueryRow(ctx, `SELECT COUNT(*), COUNT(*) FILTER (WHERE is_verified = true) FROM shops`).Scan(&total, &verified)
	return total, verified, err
}

// GetOrderStats — статистика по заказам (с комиссией)
func (r *AdminStatsRepository) GetOrderStats(ctx context.Context) (total int64, byStatus map[string]int64, revenue, platformRevenue float64, err error) {
    // Общее количество заказов, общая выручка (оборот) и комиссия платформы
    err = r.db.QueryRow(ctx, `
        SELECT 
            COUNT(*), 
            COALESCE(SUM(total_amount), 0),
            COALESCE(SUM(commission), 0)
        FROM orders
    `).Scan(&total, &revenue, &platformRevenue)
    if err != nil {
        return 0, nil, 0, 0, err
    }

    // Количество по статусам
    rows, err := r.db.Query(ctx, `SELECT current_status, COUNT(*) FROM orders GROUP BY current_status`)
    if err != nil {
        return 0, nil, 0, 0, err
    }
    defer rows.Close()

    byStatus = make(map[string]int64)
    for rows.Next() {
        var status string
        var count int64
        err := rows.Scan(&status, &count)
        if err != nil {
            return 0, nil, 0, 0, err
        }
        byStatus[status] = count
    }
    return total, byStatus, revenue, platformRevenue, nil
}

// GetProductStats — статистика по товарам (из Catalog Service)
func (r *AdminStatsRepository) GetProductStats(ctx context.Context) (total, active int64, err error) {
	err = r.db.QueryRow(ctx, `SELECT COUNT(*), COUNT(*) FILTER (WHERE is_active = true) FROM products`).Scan(&total, &active)
	return total, active, err
}

// GetCategoryStats — статистика по категориям
func (r *AdminStatsRepository) GetCategoryStats(ctx context.Context) (total int64, err error) {
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM categories`).Scan(&total)
	return total, err
}