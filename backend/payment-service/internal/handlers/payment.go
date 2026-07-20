package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/service"

	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreatePayment godoc
// @Summary      Создание платежа
// @Description  Создаёт новый платеж для заказа
// @Tags         payments
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body models.CreatePaymentRequest true "Данные платежа"
// @Success      201 {object} models.Payment
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /payments [post]
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
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

	ctx := context.WithValue(r.Context(), "user_id", userID)

	// Проверяем, что платеж принадлежит пользователю
	order, err := h.paymentService.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "order not found")
		return
	}

	if order.CustomerID != userID {
		respondWithError(w, http.StatusForbidden, "you can only pay for your own orders")
		return
	}

	payment, err := h.paymentService.CreatePayment(ctx, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, payment)
}

// GetPaymentStatus godoc
// @Summary      Получение статуса платежа
// @Description  Возвращает статус платежа по ID
// @Tags         payments
// @Produce      json
// @Security     Bearer
// @Param        id query string true "ID платежа"
// @Success      200 {object} models.Payment
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /payments [get]
func (h *PaymentHandler) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
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

	payment, err := h.paymentService.GetPaymentStatus(r.Context(), id)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "payment not found" {
			status = http.StatusNotFound
		}
		respondWithError(w, status, err.Error())
		return
	}

	// Проверяем, что платеж принадлежит пользователю
	order, err := h.paymentService.GetOrderByID(r.Context(), payment.OrderID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	role := r.Header.Get("X-User-Role")
	if role != "admin" && order.CustomerID != userID {
		respondWithError(w, http.StatusForbidden, "you can only view your own payments")
		return
	}

	respondWithJSON(w, http.StatusOK, payment)
}

// GetPaymentByOrderID godoc
// @Summary      Получение платежа по ID заказа
// @Description  Возвращает платеж по ID заказа
// @Tags         payments
// @Produce      json
// @Security     Bearer
// @Param        order_id query string true "ID заказа"
// @Success      200 {object} models.Payment
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /payments/order [get]
func (h *PaymentHandler) GetPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("order_id")
	if orderIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "order_id parameter is required")
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid order_id format")
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

	payment, err := h.paymentService.GetPaymentByOrderID(r.Context(), orderID)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "payment not found" {
			status = http.StatusNotFound
		}
		respondWithError(w, status, err.Error())
		return
	}

	// Проверяем, что платеж принадлежит пользователю
	order, err := h.paymentService.GetOrderByID(r.Context(), payment.OrderID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	role := r.Header.Get("X-User-Role")
	if role != "admin" && order.CustomerID != userID {
		respondWithError(w, http.StatusForbidden, "you can only view your own payments")
		return
	}

	respondWithJSON(w, http.StatusOK, payment)
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
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}