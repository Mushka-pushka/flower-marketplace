package repository

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

type AdminStatsRepository struct {
    db *pgxpool.Pool
}

func NewAdminStatsRepository(db *pgxpool.Pool) *AdminStatsRepository {
    return &AdminStatsRepository{
        db: db,
    }
}

// GetUserStats — статистика по пользователям (из БД Auth Service)
func (r *AdminStatsRepository) GetUserStats(ctx context.Context) (total int64, byRole map[string]int64, err error) {
    byRole = make(map[string]int64)
    
    // Общее количество пользователей
    err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
    if err != nil {
        return 0, nil, fmt.Errorf("failed to get total users: %w", err)
    }
    
    // Количество по ролям
    rows, err := r.db.Query(ctx, `
        SELECT role, COUNT(*) 
        FROM users 
        GROUP BY role
    `)
    if err != nil {
        return 0, nil, fmt.Errorf("failed to get users by role: %w", err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var role string
        var count int64
        if err := rows.Scan(&role, &count); err != nil {
            return 0, nil, fmt.Errorf("failed to scan role stats: %w", err)
        }
        byRole[role] = count
    }
    
    return total, byRole, nil
}

// GetShopStats — статистика по магазинам (из БД Auth Service)
func (r *AdminStatsRepository) GetShopStats(ctx context.Context) (total, verified int64, err error) {
    // Общее количество магазинов
    err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM shops`).Scan(&total)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to get total shops: %w", err)
    }
    
    // Верифицированные магазины
    err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM shops WHERE is_verified = true`).Scan(&verified)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to get verified shops: %w", err)
    }
    
    return total, verified, nil
}

// GetOrderStats — статистика по заказам (из БД Auth Service)
// ВНИМАНИЕ: предполагается, что таблица orders есть в этой же БД
func (r *AdminStatsRepository) GetOrderStats(ctx context.Context) (total int64, byStatus map[string]int64, revenue, platformRevenue float64, err error) {
    byStatus = make(map[string]int64)
    
    // Общее количество заказов и общая выручка
    err = r.db.QueryRow(ctx, `
        SELECT 
            COUNT(*) as total,
            COALESCE(SUM(total_amount), 0) as revenue,
            COALESCE(SUM(commission), 0) as platform_revenue
        FROM orders
    `).Scan(&total, &revenue, &platformRevenue)
    if err != nil {
        // Если таблицы orders нет в этой БД, логируем и возвращаем нули
        log.Printf("Warning: failed to get order stats: %v", err)
        return 0, make(map[string]int64), 0, 0, nil
    }
    
    // Количество по статусам
    rows, err := r.db.Query(ctx, `
        SELECT current_status, COUNT(*) 
        FROM orders 
        GROUP BY current_status
    `)
    if err != nil {
        return 0, nil, 0, 0, fmt.Errorf("failed to get orders by status: %w", err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var status string
        var count int64
        if err := rows.Scan(&status, &count); err != nil {
            return 0, nil, 0, 0, fmt.Errorf("failed to scan order status stats: %w", err)
        }
        byStatus[status] = count
    }
    
    return total, byStatus, revenue, platformRevenue, nil
}

// GetProductStats — статистика по товарам (из БД Auth Service)
func (r *AdminStatsRepository) GetProductStats(ctx context.Context) (total, active int64, err error) {
    // Общее количество товаров
    err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&total)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to get total products: %w", err)
    }
    
    // Активные товары
    err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE is_active = true`).Scan(&active)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to get active products: %w", err)
    }
    
    return total, active, nil
}

// GetCategoryStats — статистика по категориям (из БД Auth Service)
func (r *AdminStatsRepository) GetCategoryStats(ctx context.Context) (total int64, err error) {
    err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM categories`).Scan(&total)
    if err != nil {
        return 0, fmt.Errorf("failed to get total categories: %w", err)
    }
    return total, nil
}

// GetDailyStats — статистика по дням для графиков (из БД Auth Service)
func (r *AdminStatsRepository) GetDailyStats(ctx context.Context, days int) ([]DailyStat, error) {
    if days <= 0 {
        days = 30
    }
    
    query := fmt.Sprintf(`
        SELECT 
            DATE(created_at) as date,
            COUNT(*) as orders_count,
            COALESCE(SUM(total_amount), 0) as revenue
        FROM orders
        WHERE created_at >= NOW() - INTERVAL '%d days'
        GROUP BY DATE(created_at)
        ORDER BY date ASC
    `, days)
    
    rows, err := r.db.Query(ctx, query)
    if err != nil {
        // Если таблицы orders нет, возвращаем пустой массив
        log.Printf("Warning: failed to get daily stats: %v", err)
        return []DailyStat{}, nil
    }
    defer rows.Close()
    
    var stats []DailyStat
    for rows.Next() {
        var stat DailyStat
        err := rows.Scan(&stat.Date, &stat.OrdersCount, &stat.Revenue)
        if err != nil {
            return nil, fmt.Errorf("failed to scan daily stats: %w", err)
        }
        stats = append(stats, stat)
    }
    
    return stats, nil
}

// DailyStat — структура для дневной статистики
type DailyStat struct {
    Date        string  `json:"date"`
    OrdersCount int     `json:"orders_count"`
    Revenue     float64 `json:"revenue"`
}