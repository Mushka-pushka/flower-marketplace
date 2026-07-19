package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type AdminStatsRepository struct {
    httpClient  *http.Client
    catalogURL  string
    orderURL    string
}

func NewAdminStatsRepository() *AdminStatsRepository {
    return &AdminStatsRepository{
        httpClient: &http.Client{Timeout: 5 * time.Second},
        catalogURL: "http://localhost:8082/api/v1",
        orderURL:   "http://localhost:8083/api/v1",
    }
}

// GetUserStats — статистика по пользователям (локальная БД)
func (r *AdminStatsRepository) GetUserStats(ctx context.Context) (total int64, byRole map[string]int64, err error) {
    // Этот запрос остается в Auth Service (таблица users)
    // ... существующий код ...
}

// GetShopStats — статистика по магазинам (локальная БД)
func (r *AdminStatsRepository) GetShopStats(ctx context.Context) (total, verified int64, err error) {
    // Этот запрос остается в Auth Service (таблица shops)
    // ... существующий код ...
}

// GetOrderStats — статистика по заказам (через API Order Service)
func (r *AdminStatsRepository) GetOrderStats(ctx context.Context) (total int64, byStatus map[string]int64, revenue, platformRevenue float64, err error) {
    url := fmt.Sprintf("%s/admin/stats/orders", r.orderURL)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return 0, nil, 0, 0, err
    }

    resp, err := r.httpClient.Do(req)
    if err != nil {
        return 0, nil, 0, 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return 0, nil, 0, 0, fmt.Errorf("order service error: %s", string(body))
    }

    var stats struct {
        Total           int64             `json:"total"`
        ByStatus        map[string]int64  `json:"by_status"`
        Revenue         float64           `json:"revenue"`
        PlatformRevenue float64           `json:"platform_revenue"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
        return 0, nil, 0, 0, err
    }

    return stats.Total, stats.ByStatus, stats.Revenue, stats.PlatformRevenue, nil
}

// GetProductStats — статистика по товарам (через API Catalog Service)
func (r *AdminStatsRepository) GetProductStats(ctx context.Context) (total, active int64, err error) {
    url := fmt.Sprintf("%s/admin/stats/products", r.catalogURL)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return 0, 0, err
    }

    resp, err := r.httpClient.Do(req)
    if err != nil {
        return 0, 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return 0, 0, fmt.Errorf("catalog service error: %s", string(body))
    }

    var stats struct {
        Total  int64 `json:"total"`
        Active int64 `json:"active"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
        return 0, 0, err
    }

    return stats.Total, stats.Active, nil
}

// GetCategoryStats — статистика по категориям (через API Catalog Service)
func (r *AdminStatsRepository) GetCategoryStats(ctx context.Context) (total int64, err error) {
    url := fmt.Sprintf("%s/admin/stats/categories", r.catalogURL)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return 0, err
    }

    resp, err := r.httpClient.Do(req)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return 0, fmt.Errorf("catalog service error: %s", string(body))
    }

    var stats struct {
        Total int64 `json:"total"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
        return 0, err
    }

    return stats.Total, nil
}