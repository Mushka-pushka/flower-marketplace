package repository

import (
	"context"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepository struct {
	db *pgxpool.Pool
}

func NewAnalyticsRepository(db *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
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

// GetPopularProducts — получает популярные товары продавца
func (r *AnalyticsRepository) GetPopularProducts(ctx context.Context, shopID uuid.UUID, limit int) ([]models.PopularProduct, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT 
			p.id as product_id,
			p.name as product_name,
			SUM(oi.quantity) as total_sold,
			SUM(oi.total) as total_revenue,
			COUNT(DISTINCT oi.order_id) as orders_count
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN products p ON p.id = oi.product_id
		WHERE o.shop_id = $1 AND o.current_status = 'delivered'
		GROUP BY p.id, p.name
		ORDER BY total_sold DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, shopID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.PopularProduct
	for rows.Next() {
		var product models.PopularProduct
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.TotalSold,
			&product.TotalRevenue,
			&product.OrdersCount,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
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