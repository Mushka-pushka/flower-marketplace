package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

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

// CreateProduct — POST /api/v1/catalog/products
func (h *CatalogHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	product, err := h.catalogService.CreateProduct(r.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "required") {
			status = http.StatusBadRequest
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, product)
}

// GetProductByID — GET /api/v1/catalog/products?id={uuid}
func (h *CatalogHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
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

	product, err := h.catalogService.GetProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			respondWithError(w, http.StatusNotFound, "product not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, product)
}

// GetProductBySlug — GET /api/v1/catalog/products/slug/{slug}
func (h *CatalogHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		respondWithError(w, http.StatusBadRequest, "slug is required")
		return
	}

	product, err := h.catalogService.GetProductBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			respondWithError(w, http.StatusNotFound, "product not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, product)
}

// SearchProducts — GET /api/v1/catalog/search
// Параметры:
//   - q         — поисковый запрос (текст)
//   - category  — slug категории
//   - tags      — теги через запятую (романтика,свадьба)
//   - min_price — минимальная цена
//   - max_price — максимальная цена
//   - sort_by   — поле сортировки (price_asc, price_desc, rating, relevance, newest)
//   - limit     — количество записей (по умолчанию 24)
//   - offset    — смещение (по умолчанию 0)
func (h *CatalogHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	req := models.SearchRequest{
		Query:    r.URL.Query().Get("q"),
		Category: r.URL.Query().Get("category"),
		SortBy:   r.URL.Query().Get("sort_by"),
	}

	// Парсинг тегов (разделитель — запятая)
	if tagsStr := r.URL.Query().Get("tags"); tagsStr != "" {
		tags := strings.Split(tagsStr, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		req.Tags = tags
	}

	// Парсинг цен
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			req.MinPrice = &val
		}
	}
	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			req.MaxPrice = &val
		}
	}

	// Парсинг пагинации
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			req.Limit = val
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			req.Offset = val
		}
	}

	// Поиск по категории (фильтр)
	if req.Category == "" {
		req.Category = r.URL.Query().Get("category")
	}

	resp, err := h.catalogService.SearchProducts(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// GetCategories — GET /api/v1/catalog/categories
// Параметр: with_count=true — добавить количество товаров в категориях
func (h *CatalogHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	withProducts := r.URL.Query().Get("with_count") == "true"

	categories, err := h.catalogService.GetCategories(r.Context(), withProducts)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

// UpdateProduct — PUT /api/v1/catalog/products/{id}
func (h *CatalogHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	product, err := h.catalogService.UpdateProduct(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			respondWithError(w, http.StatusNotFound, "product not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, product)
}

// DeleteProduct — DELETE /api/v1/catalog/products/{id}
func (h *CatalogHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	err = h.catalogService.DeleteProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			respondWithError(w, http.StatusNotFound, "product not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// ============================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ДЛЯ ОТВЕТОВ
// ============================================================

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

// ============================================================
// КОРЗИНА (CART)
// ============================================================

// AddToCart — добавление товара в корзину
func (h *CatalogHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	// Временно: получаем user_id из запроса (позже будет из JWT)
	var req models.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Временно используем тестового пользователя
	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

	err := h.catalogService.AddToCart(r.Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Товар добавлен в корзину"})
}

// GetCart — получение корзины пользователя
func (h *CatalogHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

	items, err := h.catalogService.GetCart(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Считаем общую сумму
	var total float64
	for _, item := range items {
		total += item.TotalPrice
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
		"total": total,
	})
}

// UpdateCartItem — обновление количества товара в корзине
func (h *CatalogHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса (не из Path)
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

	err = h.catalogService.UpdateCartItem(r.Context(), cartItemID, req.Quantity)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Количество обновлено"})
}

// RemoveFromCart — удаление товара из корзины
func (h *CatalogHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса (не из Path)
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

	err = h.catalogService.RemoveFromCart(r.Context(), cartItemID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Товар удалён из корзины"})
}

// ИЗБРАННОЕ (FAVORITES)

// AddFavorite — добавление товара в избранное
func (h *CatalogHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	var req models.AddFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

	err := h.catalogService.AddFavorite(r.Context(), userID, req.ProductID)
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

// GetFavorites — получение списка избранных товаров
func (h *CatalogHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

	items, err := h.catalogService.GetFavorites(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

// RemoveFavorite — удаление товара из избранного
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

	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

	err = h.catalogService.RemoveFavorite(r.Context(), userID, productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Товар удалён из избранного"})
}

// CheckFavorite — проверка, находится ли товар в избранном
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

	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

	isFavorite, err := h.catalogService.IsFavorite(r.Context(), userID, productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]bool{"is_favorite": isFavorite})
}

// ОТЗЫВЫ (REVIEWS)

// CreateReview — создание отзыва
func (h *CatalogHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req models.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

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

// GetProductReviews — получение отзывов на товар
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

// GetMyReviews — получение моих отзывов
func (h *CatalogHandler) GetMyReviews(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

	reviews, err := h.catalogService.GetMyReviews(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, reviews)
}

// UpdateReview — обновление отзыва
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

	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

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

// DeleteReview — удаление отзыва
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

	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

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

// ApproveReview — одобрение отзыва (для админа)
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

// GetAutocompleteSuggestions — получение подсказок для поиска
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