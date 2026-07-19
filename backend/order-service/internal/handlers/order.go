package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/service"

	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder godoc
// @Summary      Создание нового заказа
// @Description  Создаёт заказ с товарами и отправляет событие в RabbitMQ
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body models.CreateOrderRequest true "Данные заказа"
// @Success      201 {object} models.Order
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Берем user_id из заголовка (от API Gateway)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	order, err := h.orderService.CreateOrder(r.Context(), userID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "order must have at least one item" {
			status = http.StatusBadRequest
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, order)
}

// GetOrder godoc
// @Summary      Получение заказа по ID
// @Description  Возвращает заказ с позициями и историей статусов
// @Tags         orders
// @Produce      json
// @Param        id query string true "ID заказа"
// @Success      200 {object} models.OrderResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Router       /orders [get]
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	order, err := h.orderService.GetOrderByID(r.Context(), id)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "order not found" {
			status = http.StatusNotFound
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

// GetOrdersByCustomer godoc
// @Summary      Получение заказов покупателя
// @Description  Возвращает все заказы конкретного покупателя (только свои или для admin)
// @Tags         orders
// @Produce      json
// @Param        customer_id query string true "ID покупателя"
// @Success      200 {array} models.Order
// @Failure      400 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /orders/customer [get]
func (h *OrderHandler) GetOrdersByCustomer(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.URL.Query().Get("customer_id")
	if customerIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "customer_id is required")
		return
	}

	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid customer_id format")
		return
	}

	// Проверяем, что пользователь авторизован
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

	// Если пользователь не admin и запрашивает не свои заказы - запрещаем
	role := r.Header.Get("X-User-Role")
	if role != "admin" && userID != customerID {
		respondWithError(w, http.StatusForbidden, "you can only view your own orders")
		return
	}

	orders, err := h.orderService.GetOrdersByCustomer(r.Context(), customerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, orders)
}

// GetOrdersByShop godoc
// @Summary      Получение заказов магазина (для продавца)
// @Description  Возвращает все заказы конкретного магазина
// @Tags         orders
// @Produce      json
// @Param        shop_id query string true "ID магазина"
// @Success      200 {array} models.Order
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /orders/shop [get]
func (h *OrderHandler) GetOrdersByShop(w http.ResponseWriter, r *http.Request) {
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

	orders, err := h.orderService.GetOrdersByShop(r.Context(), shopID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, orders)
}

// CancelOrder godoc
// @Summary      Отмена заказа
// @Description  Покупатель или продавец отменяет заказ (только если он ещё не доставлен)
// @Tags         orders
// @Param        id query string true "ID заказа"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /orders/cancel [post]
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("id")
	if orderIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	// Берем user_id из заголовка (от API Gateway)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	// Определяем роль пользователя из заголовка
	role := r.Header.Get("X-User-Role")
	if role == "" {
		role = "customer" // по умолчанию
	}

	err = h.orderService.CancelOrder(r.Context(), orderID, userID, role)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "you can only cancel your own orders":
			status = http.StatusForbidden
		case "cannot cancel delivered order":
			status = http.StatusBadRequest
		case "order already cancelled":
			status = http.StatusBadRequest
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Заказ отменён"})
}

// UpdateOrderStatusBySeller godoc
// @Summary      Обновление статуса заказа продавцом
// @Description  Продавец изменяет статус заказа (confirmed, preparing, packing, delivery, delivered, cancelled)
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body models.UpdateOrderStatusRequest true "Новый статус"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /orders/status [put]
func (h *OrderHandler) UpdateOrderStatusBySeller(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Берем user_id из заголовка (от API Gateway)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	// Проверяем роль пользователя
	role := r.Header.Get("X-User-Role")
	if role != "seller" {
		respondWithError(w, http.StatusForbidden, "only sellers can update order status")
		return
	}

	// Получаем shop_id по user_id
	shopID, err := h.orderService.GetShopIDBySellerID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get shop")
		return
	}
	if shopID == uuid.Nil {
		respondWithError(w, http.StatusForbidden, "seller has no shop")
		return
	}

	err = h.orderService.UpdateOrderStatusBySeller(r.Context(), req.OrderID, shopID, req.Status, req.Comment)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "you can only update orders from your shop":
			status = http.StatusForbidden
		case "invalid status":
			status = http.StatusBadRequest
		case "cannot change status of delivered or cancelled order":
			status = http.StatusBadRequest
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Статус заказа обновлён"})
}

// CanReview godoc
// @Summary      Проверка возможности оставить отзыв
// @Description  Проверяет, может ли пользователь оставить отзыв на товар
// @Tags         orders
// @Produce      json
// @Param        product_id query string true "ID товара"
// @Success      200 {object} map[string]bool
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /orders/can-review [get]
func (h *OrderHandler) CanReview(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.URL.Query().Get("product_id")
	if productIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "product_id is required")
		return
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid product_id")
		return
	}

	// Берем user_id из заголовка (от API Gateway)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	canReview, err := h.orderService.CanReview(r.Context(), userID, productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]bool{"can_review": canReview})
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