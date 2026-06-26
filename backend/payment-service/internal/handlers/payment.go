package handlers

import (
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

// CreatePayment — создание платежа
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	payment, err := h.paymentService.CreatePayment(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, payment)
}

// GetPaymentStatus — получение статуса платежа
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

	payment, err := h.paymentService.GetPaymentStatus(r.Context(), id)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "payment not found" {
			status = http.StatusNotFound
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, payment)
}

// GetPaymentByOrderID — получение платежа по ID заказа
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

	payment, err := h.paymentService.GetPaymentByOrderID(r.Context(), orderID)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "payment not found" {
			status = http.StatusNotFound
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, payment)
}

// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ

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