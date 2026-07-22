package handlers

import (
	"encoding/json"
	"log"
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
// @Security     Bearer
// @Param        request body models.CreateOrderRequest true "Данные заказа"
// @Success      201 {object} models.Order
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
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
// @Description  Возвращает все заказы конкретного покупателя
// @Tags         orders
// @Produce      json
// @Security     Bearer
// @Param        customer_id query string true "ID покупателя"
// @Success      200 {array} models.Order
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
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

	// Проверка прав
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
// @Security     Bearer
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
// @Security     Bearer
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
// @Security     Bearer
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

	userIDStr := r.Header.Get("X-User-ID")
	log.Printf("UpdateOrderStatusBySeller: X-User-ID: %s", userIDStr)
	
	role := r.Header.Get("X-User-Role")
	log.Printf("UpdateOrderStatusBySeller: X-User-Role: %s", role)

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
	if role != "seller" {
		log.Printf("User is not a seller: role=%s", role)
		respondWithError(w, http.StatusForbidden, "only sellers can update order status")
		return
	}

	// Получаем shop_id по user_id
	shopID, err := h.orderService.GetShopIDBySellerID(r.Context(), userID)
	if err != nil || shopID == uuid.Nil {
		log.Printf("Seller has no shop: userID=%s, err=%v", userID, err)
		respondWithError(w, http.StatusForbidden, "seller has no shop")
		return
	}
	log.Printf("UpdateOrderStatusBySeller: seller shop_id: %s", shopID)

	// Проверяем, что заказ принадлежит магазину продавца
	order, err := h.orderService.GetOrderByIDSimple(r.Context(), req.OrderID)
	if err != nil {
		log.Printf("Order not found: orderID=%s, err=%v", req.OrderID, err)
		respondWithError(w, http.StatusNotFound, "order not found")
		return
	}
	log.Printf("UpdateOrderStatusBySeller: order shop_id: %s", order.ShopID)

	if order.ShopID != shopID {
		log.Printf("Shop mismatch: order.shop_id=%s, seller.shop_id=%s", order.ShopID, shopID)
		respondWithError(w, http.StatusForbidden, "you can only update orders from your shop")
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
		log.Printf("Failed to update order status: %v", err)
		respondWithError(w, status, err.Error())
		return
	}

	log.Printf("Order status updated: orderID=%s, newStatus=%s", req.OrderID, req.Status)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Статус заказа обновлён"})
}

// CanReview godoc
// @Summary      Проверка возможности оставить отзыв
// @Description  Проверяет, может ли пользователь оставить отзыв на товар
// @Tags         orders
// @Produce      json
// @Security     Bearer
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

// GetOrderItemsByCustomer — получение всех товаров пользователя (как отдельные позиции)
func (h *OrderHandler) GetOrderItemsByCustomer(w http.ResponseWriter, r *http.Request) {
    // Берем user_id из заголовка
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

    // Проверяем, что пользователь запрашивает свои данные
    // (админ может смотреть всех)
    role := r.Header.Get("X-User-Role")
    if role != "admin" {
        // Только свои товары
        items, err := h.orderService.GetOrderItemsByCustomer(r.Context(), userID)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, err.Error())
            return
        }
        respondWithJSON(w, http.StatusOK, items)
        return
    }

    // Админ может указать customer_id
    customerIDStr := r.URL.Query().Get("customer_id")
    if customerIDStr == "" {
        respondWithError(w, http.StatusBadRequest, "customer_id is required for admin")
        return
    }

    customerID, err := uuid.Parse(customerIDStr)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "invalid customer_id format")
        return
    }

    items, err := h.orderService.GetOrderItemsByCustomer(r.Context(), customerID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, items)
}
