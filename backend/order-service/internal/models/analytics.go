package models

// SellerAnalytics — аналитика для продавца
type SellerAnalytics struct {
	TotalOrders     int64   `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	CompletedOrders int64   `json:"completed_orders"`
	CancelledOrders int64   `json:"cancelled_orders"`
	AverageOrderSum float64 `json:"average_order_sum"`
}

// PopularProduct — популярный товар
type PopularProduct struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalSold   int     `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
	OrdersCount int     `json:"orders_count"`
}

// OrderStatsByStatus — статистика по статусам
type OrderStatsByStatus struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// SalesDay — данные по продажам за день
type SalesDay struct {
    Date        string  `json:"date"`
    OrdersCount int     `json:"orders_count"`
    Revenue     float64 `json:"revenue"`
}