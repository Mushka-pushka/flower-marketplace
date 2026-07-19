package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// getProductNamesFromCatalog — получает названия товаров из Catalog Service
func (r *AnalyticsRepository) getProductNamesFromCatalog(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	if len(productIDs) == 0 {
		return make(map[uuid.UUID]string), nil
	}

	// Формируем URL с параметрами
	url := fmt.Sprintf("%s/products/batch", r.catalogURL)
	
	// Создаем запрос с списком ID
	requestBody := map[string]interface{}{
		"product_ids": productIDs,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewReader(body))

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get product names: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("catalog service returned status %d: %s", resp.StatusCode, string(body))
	}

	var products []struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("failed to decode product names: %w", err)
	}

	// Создаем мапу для быстрого доступа
	nameMap := make(map[uuid.UUID]string)
	for _, p := range products {
		nameMap[p.ID] = p.Name
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
	nameMap, err := r.getProductNamesFromCatalog(ctx, productIDs)
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