package handlers

import (
	"net/http"
	"strconv"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/service"

	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetSellerAnalytics — получение общей аналитики продавца
func (h *AnalyticsHandler) GetSellerAnalytics(w http.ResponseWriter, r *http.Request) {
	shopIDStr := r.URL.Query().Get("shop_id")
	if shopIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "shop_id parameter is required")
		return
	}

	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid shop_id format")
		return
	}

	analytics, err := h.analyticsService.GetSellerAnalytics(r.Context(), shopID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, analytics)
}

// GetPopularProducts — получение популярных товаров продавца
func (h *AnalyticsHandler) GetPopularProducts(w http.ResponseWriter, r *http.Request) {
	shopIDStr := r.URL.Query().Get("shop_id")
	if shopIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "shop_id parameter is required")
		return
	}

	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid shop_id format")
		return
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	products, err := h.analyticsService.GetPopularProducts(r.Context(), shopID, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

// GetOrderStatsByStatus — получение статистики по статусам заказов
func (h *AnalyticsHandler) GetOrderStatsByStatus(w http.ResponseWriter, r *http.Request) {
	shopIDStr := r.URL.Query().Get("shop_id")
	if shopIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "shop_id parameter is required")
		return
	}

	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid shop_id format")
		return
	}

	stats, err := h.analyticsService.GetOrderStatsByStatus(r.Context(), shopID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}