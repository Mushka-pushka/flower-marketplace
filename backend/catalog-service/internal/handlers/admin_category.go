package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/models"

	"github.com/google/uuid"
)

// ============================================================
// АДМИН: УПРАВЛЕНИЕ КАТЕГОРИЯМИ
// ============================================================

// AdminCreateCategory — создание категории (админ)
func (h *CatalogHandler) AdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	category, err := h.catalogService.AdminCreateCategory(r.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category with this slug already exists" {
			status = http.StatusConflict
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, category)
}

// AdminGetAllCategories — получение всех категорий (админ)
func (h *CatalogHandler) AdminGetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.catalogService.AdminGetAllCategories(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

// AdminGetCategoryByID — получение категории по ID (админ)
func (h *CatalogHandler) AdminGetCategoryByID(w http.ResponseWriter, r *http.Request) {
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

	category, err := h.catalogService.AdminGetCategoryByID(r.Context(), id)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, category)
}

// AdminUpdateCategory — обновление категории (админ)
func (h *CatalogHandler) AdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
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

	var req models.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	category, err := h.catalogService.AdminUpdateCategory(r.Context(), id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		}
		if err.Error() == "category with this slug already exists" {
			status = http.StatusConflict
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, category)
}

// AdminDeleteCategory — удаление категории (админ)
func (h *CatalogHandler) AdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
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

	err = h.catalogService.AdminDeleteCategory(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		}
		if err.Error() == "cannot delete category with existing products" ||
			err.Error() == "cannot delete category with subcategories" {
			status = http.StatusBadRequest
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Категория удалена"})
}