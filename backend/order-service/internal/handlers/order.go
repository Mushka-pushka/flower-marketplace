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

// CreateOrder — создание заказа
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.orderService.CreateOrder(r.Context(), &req)
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

// GetOrder — получение заказа по ID
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

// GetOrdersByCustomer — получение заказов покупателя
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

	orders, err := h.orderService.GetOrdersByCustomer(r.Context(), customerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, orders)
}

// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ

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

// CancelOrder — отмена заказа
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

	// Временно: получаем user_id из запроса (позже будет из JWT)
	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")
	role := "customer" // временно

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

// GetOrdersByShop — получение заказов магазина
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

// UpdateOrderStatusBySeller — обновление статуса заказа продавцом
func (h *OrderHandler) UpdateOrderStatusBySeller(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Временно используем тестовый shop_id (позже будет из JWT)
	shopID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	err := h.orderService.UpdateOrderStatusBySeller(r.Context(), req.OrderID, shopID, req.Status, req.Comment)
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