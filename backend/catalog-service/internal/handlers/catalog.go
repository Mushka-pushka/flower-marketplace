package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/middleware"
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/repository"
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/service"

	"github.com/google/uuid"
)

type CatalogHandler struct {
	catalogService *service.CatalogService
}

func NewCatalogHandler(catalogService *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalogService: catalogService}
}

// ... (CreateProduct, GetProductByID, GetProductBySlug, SearchProducts, GetCategories, UpdateProduct, DeleteProduct остаются без изменений) ...

// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
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

// getUserIDFromContext - вспомогательная функция для получения userID из контекста
func getUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		return uuid.Nil, errors.New("user not authenticated")
	}
	return userID, nil
}

// КОРЗИНА (CART)

// AddToCart godoc
// @Summary      Добавление товара в корзину
// @Description  Увеличивает количество товара в корзине или добавляет новый
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request body models.AddToCartRequest true "Товар и количество"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/cart [post]
func (h *CatalogHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req models.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.catalogService.AddToCart(r.Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Товар добавлен в корзину"})
}

// GetCart godoc
// @Summary      Получение корзины пользователя
// @Description  Возвращает все товары в корзине с общей суммой
// @Tags         cart
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/cart [get]
func (h *CatalogHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	items, err := h.catalogService.GetCart(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var total float64
	for _, item := range items {
		total += item.TotalPrice
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
		"total": total,
	})
}

// UpdateCartItem godoc
// @Summary      Обновление количества товара в корзине
// @Description  Изменяет количество конкретного товара в корзине
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        id query string true "ID позиции в корзине"
// @Param        request body models.UpdateCartRequest true "Новое количество"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/cart [put]
func (h *CatalogHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	cartItemID, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	var req models.UpdateCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.catalogService.UpdateCartItem(r.Context(), cartItemID, userID, req.Quantity)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Количество обновлено"})
}

// RemoveFromCart godoc
// @Summary      Удаление товара из корзины
// @Description  Удаляет конкретный товар из корзины пользователя
// @Tags         cart
// @Param        id query string true "ID позиции в корзине"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/cart [delete]
func (h *CatalogHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	cartItemID, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.catalogService.RemoveFromCart(r.Context(), cartItemID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Товар удалён из корзины"})
}

// ИЗБРАННОЕ (FAVORITES)

// AddFavorite godoc
// @Summary      Добавление товара в избранное
// @Description  Сохраняет товар в список избранных
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param        request body models.AddFavoriteRequest true "ID товара"
// @Success      201 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/favorites [post]
func (h *CatalogHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	var req models.AddFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.catalogService.AddFavorite(r.Context(), userID, req.ProductID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "product already in favorites" {
			status = http.StatusConflict
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Товар добавлен в избранное"})
}

// GetFavorites godoc
// @Summary      Получение списка избранных товаров
// @Description  Возвращает все товары, добавленные в избранное
// @Tags         favorites
// @Produce      json
// @Success      200 {array} models.FavoriteWithProduct
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/favorites [get]
func (h *CatalogHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	items, err := h.catalogService.GetFavorites(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

// RemoveFavorite godoc
// @Summary      Удаление товара из избранного
// @Description  Удаляет товар из списка избранных
// @Tags         favorites
// @Param        product_id query string true "ID товара"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/favorites [delete]
func (h *CatalogHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.URL.Query().Get("product_id")
	if productIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "product_id parameter is required")
		return
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid product_id format")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.catalogService.RemoveFavorite(r.Context(), userID, productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Товар удалён из избранного"})
}

// CheckFavorite godoc
// @Summary      Проверка наличия товара в избранном
// @Description  Проверяет, находится ли товар в избранном у пользователя
// @Tags         favorites
// @Produce      json
// @Param        product_id query string true "ID товара"
// @Success      200 {object} map[string]bool
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/favorites/check [get]
func (h *CatalogHandler) CheckFavorite(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.URL.Query().Get("product_id")
	if productIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "product_id parameter is required")
		return
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid product_id format")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	isFavorite, err := h.catalogService.IsFavorite(r.Context(), userID, productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]bool{"is_favorite": isFavorite})
}

// ОТЗЫВЫ (REVIEWS)

// CreateReview godoc
// @Summary      Создание отзыва
// @Description  Добавляет отзыв на товар с рейтингом
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        request body models.CreateReviewRequest true "Данные отзыва"
// @Success      201 {object} models.Review
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/reviews [post]
func (h *CatalogHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req models.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	review, err := h.catalogService.CreateReview(r.Context(), &req, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "review already exists for this order" {
			status = http.StatusConflict
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, review)
}

// GetProductReviews godoc
// @Summary      Получение отзывов на товар
// @Description  Возвращает все одобренные отзывы для конкретного товара
// @Tags         reviews
// @Produce      json
// @Param        product_id query string true "ID товара"
// @Success      200 {array} models.ReviewWithUser
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/reviews [get]
func (h *CatalogHandler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.URL.Query().Get("product_id")
	if productIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "product_id parameter is required")
		return
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid product_id format")
		return
	}

	reviews, err := h.catalogService.GetProductReviews(r.Context(), productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, reviews)
}

// GetMyReviews godoc
// @Summary      Получение моих отзывов
// @Description  Возвращает все отзывы текущего пользователя
// @Tags         reviews
// @Produce      json
// @Success      200 {array} models.ReviewWithUser
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/reviews/me [get]
func (h *CatalogHandler) GetMyReviews(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	reviews, err := h.catalogService.GetMyReviews(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, reviews)
}

// UpdateReview godoc
// @Summary      Обновление отзыва
// @Description  Изменяет текст или рейтинг отзыва
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        id query string true "ID отзыва"
// @Param        request body models.UpdateReviewRequest true "Новые данные"
// @Success      200 {object} models.Review
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/reviews [put]
func (h *CatalogHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	reviewIDStr := r.URL.Query().Get("id")
	if reviewIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	var req models.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	review, err := h.catalogService.UpdateReview(r.Context(), reviewID, &req, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "you can only update your own reviews" {
			status = http.StatusForbidden
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, review)
}

// DeleteReview godoc
// @Summary      Удаление отзыва
// @Description  Удаляет отзыв пользователя
// @Tags         reviews
// @Param        id query string true "ID отзыва"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/reviews [delete]
func (h *CatalogHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewIDStr := r.URL.Query().Get("id")
	if reviewIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.catalogService.DeleteReview(r.Context(), reviewID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "you can only delete your own reviews" {
			status = http.StatusForbidden
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Отзыв удалён"})
}

// ApproveReview godoc
// @Summary      Одобрение отзыва (админ)
// @Description  Администратор одобряет отзыв для публикации
// @Tags         admin
// @Param        id query string true "ID отзыва"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /admin/reviews/approve [put]
func (h *CatalogHandler) ApproveReview(w http.ResponseWriter, r *http.Request) {
	reviewIDStr := r.URL.Query().Get("id")
	if reviewIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	err = h.catalogService.ApproveReview(r.Context(), reviewID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Отзыв одобрен"})
}

// GetAutocompleteSuggestions godoc
// @Summary      Автодополнение поиска
// @Description  Возвращает подсказки по товарам, категориям и тегам
// @Tags         catalog
// @Produce      json
// @Param        q query string true "Поисковый запрос"
// @Param        limit query int false "Количество подсказок" default(10)
// @Success      200 {array} models.AutocompleteSuggestion
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/autocomplete [get]
func (h *CatalogHandler) GetAutocompleteSuggestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithJSON(w, http.StatusOK, []models.AutocompleteSuggestion{})
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	suggestions, err := h.catalogService.GetAutocompleteSuggestions(r.Context(), query, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, suggestions)
}

// АДРЕСА ДОСТАВКИ

// CreateAddress godoc
// @Summary      Создание адреса доставки
// @Description  Добавляет новый адрес для пользователя
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        request body models.CreateAddressRequest true "Данные адреса"
// @Success      201 {object} models.DeliveryAddress
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/addresses [post]
func (h *CatalogHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	address, err := h.catalogService.CreateAddress(r.Context(), userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, address)
}

// GetAddresses godoc
// @Summary      Получение адресов пользователя
// @Description  Возвращает все адреса доставки текущего пользователя
// @Tags         addresses
// @Produce      json
// @Success      200 {array} models.DeliveryAddress
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/addresses [get]
func (h *CatalogHandler) GetAddresses(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	addresses, err := h.catalogService.GetAddresses(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, addresses)
}

// UpdateAddress godoc
// @Summary      Обновление адреса доставки
// @Description  Изменяет данные существующего адреса
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        id query string true "ID адреса"
// @Param        request body models.UpdateAddressRequest true "Новые данные"
// @Success      200 {object} models.DeliveryAddress
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/addresses [put]
func (h *CatalogHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	addressIDStr := r.URL.Query().Get("id")
	if addressIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	var req models.UpdateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	address, err := h.catalogService.UpdateAddress(r.Context(), addressID, userID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "you can only update your own addresses" {
			status = http.StatusForbidden
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, address)
}

// DeleteAddress godoc
// @Summary      Удаление адреса доставки
// @Description  Удаляет адрес пользователя
// @Tags         addresses
// @Param        id query string true "ID адреса"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/addresses [delete]
func (h *CatalogHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	addressIDStr := r.URL.Query().Get("id")
	if addressIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.catalogService.DeleteAddress(r.Context(), addressID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "you can only delete your own addresses" {
			status = http.StatusForbidden
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Адрес удалён"})
}

// SetDefaultAddress godoc
// @Summary      Установка адреса по умолчанию
// @Description  Делает выбранный адрес основным для пользователя
// @Tags         addresses
// @Param        id query string true "ID адреса"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /catalog/addresses/default [post]
func (h *CatalogHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	addressIDStr := r.URL.Query().Get("id")
	if addressIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.catalogService.SetDefaultAddress(r.Context(), addressID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "you can only set your own addresses as default" {
			status = http.StatusForbidden
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Адрес установлен по умолчанию"})
}