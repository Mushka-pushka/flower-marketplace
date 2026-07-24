package handlers

import (
	"encoding/json"
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

// GetSellerAnalytics godoc
// @Summary      Получение аналитики продавца
// @Description  Возвращает общую аналитику по заказам продавца (количество заказов, выручка, средний чек)
// @Tags         analytics
// @Produce      json
// @Security     Bearer
// @Param        shop_id query string true "ID магазина"
// @Success      200 {object} models.SellerAnalytics
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /analytics/seller [get]
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

// GetPopularProducts godoc
// @Summary      Получение популярных товаров
// @Description  Возвращает самые продаваемые товары продавца с количеством продаж и выручкой
// @Tags         analytics
// @Produce      json
// @Security     Bearer
// @Param        shop_id query string true "ID магазина"
// @Param        limit query int false "Количество товаров" default(10)
// @Success      200 {array} models.PopularProduct
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /analytics/popular [get]
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

// GetOrderStatsByStatus godoc
// @Summary      Получение статистики по статусам
// @Description  Возвращает количество заказов по каждому статусу для магазина
// @Tags         analytics
// @Produce      json
// @Security     Bearer
// @Param        shop_id query string true "ID магазина"
// @Success      200 {array} models.OrderStatsByStatus
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /analytics/statuses [get]
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

// GetSalesDynamics — получение динамики продаж
func (h *AnalyticsHandler) GetSalesDynamics(w http.ResponseWriter, r *http.Request) {
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

    // Проверяем, что пользователь имеет доступ к этому магазину
    userIDStr := r.Header.Get("X-User-ID")
    if userIDStr == "" {
        respondWithError(w, http.StatusUnauthorized, "user not authenticated")
        return
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "invalid user id")
        return
    }

    // Проверяем, что продавец владеет этим магазином
    sellerShopID, err := h.analyticsService.GetShopIDBySellerID(r.Context(), userID)
    if err != nil || sellerShopID != shopID {
        respondWithError(w, http.StatusForbidden, "you don't have access to this shop")
        return
    }

    days := 30
    if d := r.URL.Query().Get("days"); d != "" {
        if val, err := strconv.Atoi(d); err == nil && val > 0 {
            days = val
        }
    }

    dynamics, err := h.analyticsService.GetSalesDynamics(r.Context(), shopID, days)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, dynamics)
}

// ============================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ============================================================

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message, Code: code})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}