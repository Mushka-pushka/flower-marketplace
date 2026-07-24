package repository

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepository struct {
	db          *pgxpool.Pool
	httpClient  *http.Client
	catalogURL  string
}

func NewAnalyticsRepository(db *pgxpool.Pool, cfg *config.Config) *AnalyticsRepository {
	catalogURL := cfg.CatalogServiceURL
	if catalogURL == "" {
		catalogURL = "http://localhost:8082/api/v1/catalog"
	}

	return &AnalyticsRepository{
		db:          db,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		catalogURL:  catalogURL,
	}
}

// GetSellerAnalytics — получает общую аналитику для продавца
func (r *AnalyticsRepository) GetSellerAnalytics(ctx context.Context, shopID uuid.UUID) (*models.SellerAnalytics, error) {
	query := `
		SELECT 
			COUNT(*) as total_orders,
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) FILTER (WHERE current_status = 'delivered') as completed_orders,
			COUNT(*) FILTER (WHERE current_status = 'cancelled') as cancelled_orders,
			COALESCE(AVG(total_amount), 0) as average_order_sum
		FROM orders
		WHERE shop_id = $1
	`

	var analytics models.SellerAnalytics
	err := r.db.QueryRow(ctx, query, shopID).Scan(
		&analytics.TotalOrders,
		&analytics.TotalRevenue,
		&analytics.CompletedOrders,
		&analytics.CancelledOrders,
		&analytics.AverageOrderSum,
	)
	if err != nil {
		return nil, err
	}
	return &analytics, nil
}

// getProductNamesFromDB — получает названия товаров напрямую из БД
func (r *AnalyticsRepository) getProductNamesFromDB(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]string, error) {
    if len(productIDs) == 0 {
        return make(map[uuid.UUID]string), nil
    }

    query := `SELECT id, name FROM products WHERE id = ANY($1)`
    rows, err := r.db.Query(ctx, query, productIDs)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    nameMap := make(map[uuid.UUID]string)
    for rows.Next() {
        var id uuid.UUID
        var name string
        err := rows.Scan(&id, &name)
        if err != nil {
            return nil, err
        }
        nameMap[id] = name
    }
    return nameMap, nil
}

// GetPopularProducts — получает популярные товары (через API Catalog Service)
func (r *AnalyticsRepository) GetPopularProducts(ctx context.Context, shopID uuid.UUID, limit int) ([]models.PopularProduct, error) {
	if limit <= 0 {
		limit = 10
	}

	// Сначала получаем ID товаров с заказами
	query := `
		SELECT 
			oi.product_id,
			SUM(oi.quantity) as total_sold,
			SUM(oi.total) as total_revenue,
			COUNT(DISTINCT oi.order_id) as orders_count
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE o.shop_id = $1 AND o.current_status = 'delivered'
		GROUP BY oi.product_id
		ORDER BY total_sold DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, shopID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productIDs []uuid.UUID
	var products []models.PopularProduct
	
	for rows.Next() {
		var p models.PopularProduct
		var productID uuid.UUID
		err := rows.Scan(
			&productID,
			&p.TotalSold,
			&p.TotalRevenue,
			&p.OrdersCount,
		)
		if err != nil {
			return nil, err
		}
		
		productIDs = append(productIDs, productID)
		p.ProductID = productID.String()
		products = append(products, p)
	}

	if len(products) == 0 {
		return products, nil
	}

	// Получаем названия товаров из Catalog Service
	nameMap, err := r.getProductNamesFromDB(ctx, productIDs)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		// Используем заглушки для названий
		for i := range products {
			products[i].ProductName = fmt.Sprintf("Товар %s", products[i].ProductID[:8])
		}
		return products, nil
	}

	// Заполняем названия
	for i := range products {
		productID, err := uuid.Parse(products[i].ProductID)
		if err == nil {
			if name, ok := nameMap[productID]; ok {
				products[i].ProductName = name
			} else {
				products[i].ProductName = fmt.Sprintf("Товар %s", products[i].ProductID[:8])
			}
		}
	}

	return products, nil
}

// GetOrderStatsByStatus — получает статистику по статусам заказов
func (r *AnalyticsRepository) GetOrderStatsByStatus(ctx context.Context, shopID uuid.UUID) ([]models.OrderStatsByStatus, error) {
	query := `
		SELECT 
			current_status as status,
			COUNT(*) as count
		FROM orders
		WHERE shop_id = $1
		GROUP BY current_status
		ORDER BY count DESC
	`

	rows, err := r.db.Query(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.OrderStatsByStatus
	for rows.Next() {
		var stat models.OrderStatsByStatus
		err := rows.Scan(&stat.Status, &stat.Count)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
}

// GetSalesDynamics — получает динамику продаж по дням
func (r *AnalyticsRepository) GetSalesDynamics(ctx context.Context, shopID uuid.UUID, days int) ([]models.SalesDay, error) {
    if days <= 0 {
        days = 30
    }

    // Исправлено: DATE(created_at)::text — преобразуем дату в текст
    query := fmt.Sprintf(`
        SELECT 
            DATE(created_at)::text as day,
            COUNT(*) as orders_count,
            COALESCE(SUM(total_amount), 0) as revenue
        FROM orders
        WHERE shop_id = $1 
            AND current_status IN ('delivered', 'confirmed', 'preparing', 'packing', 'delivery')
            AND created_at >= NOW() - INTERVAL '%d days'
        GROUP BY DATE(created_at)
        ORDER BY day ASC
    `, days)

    rows, err := r.db.Query(ctx, query, shopID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var dynamics []models.SalesDay
    for rows.Next() {
        var day models.SalesDay
        err := rows.Scan(&day.Date, &day.OrdersCount, &day.Revenue)
        if err != nil {
            return nil, err
        }
        dynamics = append(dynamics, day)
    }
    return dynamics, nil
}

// GetShopIDBySellerID — возвращает shop_id продавца по его user_id
func (r *AnalyticsRepository) GetShopIDBySellerID(ctx context.Context, sellerID uuid.UUID) (uuid.UUID, error) {
    query := `SELECT id FROM shops WHERE seller_id = $1`
    var shopID uuid.UUID
    err := r.db.QueryRow(ctx, query, sellerID).Scan(&shopID)
    if err != nil {
        return uuid.Nil, err
    }
    return shopID, nil
}