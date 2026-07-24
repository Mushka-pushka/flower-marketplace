package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

    "github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/middleware"  
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/service"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GetSellers godoc
// @Summary      Получение списка продавцов
// @Description  Возвращает список всех продавцов с фильтром по верификации
// @Tags         admin
// @Produce      json
// @Security     Bearer
// @Param        verified query bool false "Фильтр по верификации"
// @Success      200 {array} models.SellerWithUser
// @Failure      500 {object} ErrorResponse
// @Router       /admin/sellers [get]
func (h *AdminHandler) GetSellers(w http.ResponseWriter, r *http.Request) {
	var verified *bool
	if v := r.URL.Query().Get("verified"); v != "" {
		b := v == "true"
		verified = &b
	}

	sellers, err := h.adminService.GetSellers(r.Context(), verified)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, sellers)
}

// VerifySeller godoc
// @Summary      Верификация продавца
// @Description  Подтверждает или отклоняет продавца
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body models.VerifySellerRequest true "Данные для верификации"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /admin/sellers/verify [put]
func (h *AdminHandler) VerifySeller(w http.ResponseWriter, r *http.Request) {
	var req models.VerifySellerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.adminService.VerifySeller(r.Context(), req.ShopID, req.Verify)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := "Продавец подтверждён"
	if !req.Verify {
		message = "Продавец отклонён"
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": message})
}

// UpdateUserStatus godoc
// @Summary      Блокировка/разблокировка пользователя
// @Description  Изменяет статус пользователя (активен/заблокирован)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body models.UpdateUserStatusRequest true "Данные для обновления"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /admin/users/status [put]
func (h *AdminHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateUserStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.adminService.UpdateUserStatus(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := "Пользователь разблокирован"
	if !req.IsActive {
		message = "Пользователь заблокирован"
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": message})
}

// GetUsersList godoc
// @Summary      Получение списка пользователей
// @Description  Возвращает список пользователей с фильтрацией по роли и статусу
// @Tags         admin
// @Produce      json
// @Security     Bearer
// @Param        role query string false "Роль пользователя"
// @Param        is_active query bool false "Статус активности"
// @Param        limit query int false "Количество записей" default(20)
// @Param        offset query int false "Смещение" default(0)
// @Success      200 {array} models.UserResponse
// @Failure      500 {object} ErrorResponse
// @Router       /admin/users [get]
func (h *AdminHandler) GetUsersList(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	var isActive *bool
	if v := r.URL.Query().Get("is_active"); v != "" {
		b := v == "true"
		isActive = &b
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	users, err := h.adminService.GetUsersList(r.Context(), role, isActive, limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

// GetUsersListWithFilters godoc
// @Summary      Получение списка пользователей с расширенной фильтрацией
// @Description  Возвращает список пользователей с фильтрацией по роли, статусу и поиску
// @Tags         admin
// @Produce      json
// @Security     Bearer
// @Param        role query string false "Роль пользователя"
// @Param        is_active query bool false "Статус активности"
// @Param        search query string false "Поиск по email или имени"
// @Param        limit query int false "Количество записей" default(20)
// @Param        offset query int false "Смещение" default(0)
// @Success      200 {object} models.UsersListResponse
// @Failure      500 {object} ErrorResponse
// @Router       /admin/users/list [get]
func (h *AdminHandler) GetUsersListWithFilters(w http.ResponseWriter, r *http.Request) {
	req := models.UsersListRequest{
		Role:   r.URL.Query().Get("role"),
		Search: r.URL.Query().Get("search"),
	}

	if v := r.URL.Query().Get("is_active"); v != "" {
		b := v == "true"
		req.IsActive = &b
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			req.Limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			req.Offset = val
		}
	}

	resp, err := h.adminService.GetUsersListWithFilters(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// GetUserByIDForAdmin godoc
// @Summary      Получение детальной информации о пользователе
// @Description  Возвращает полную информацию о пользователе по ID
// @Tags         admin
// @Produce      json
// @Security     Bearer
// @Param        id query string true "ID пользователя"
// @Success      200 {object} models.UserDetails
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /admin/users/details [get]
func (h *AdminHandler) GetUserByIDForAdmin(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	user, err := h.adminService.GetUserByIDForAdmin(r.Context(), userID)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// UpdateShopName — обновление названия магазина
func (h *AdminHandler) UpdateShopName(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name string `json:"name"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    if req.Name == "" {
        respondWithError(w, http.StatusBadRequest, "shop name is required")
        return
    }

    // Получаем user_id из контекста
    userID, ok := middleware.GetUserIDFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "user not authenticated")
        return
    }

    // Получаем shop_id продавца
    shopID, err := h.adminService.GetShopIDBySellerID(r.Context(), userID)
    if err != nil || shopID == uuid.Nil {
        respondWithError(w, http.StatusForbidden, "seller has no shop")
        return
    }

    err = h.adminService.UpdateShopName(r.Context(), shopID, req.Name)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    // Возвращаем обновлённое название
    shop, err := h.adminService.GetShopByID(r.Context(), shopID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{
        "shop_id":   shopID.String(),
        "shop_name": shop.Name,
    })
}

// GetShopInfo — получение информации о магазине продавца
func (h *AdminHandler) GetShopInfo(w http.ResponseWriter, r *http.Request) {
    // Получаем user_id из контекста
    userID, ok := middleware.GetUserIDFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "user not authenticated")
        return
    }

    // Получаем shop_id продавца
    shopID, err := h.adminService.GetShopIDBySellerID(r.Context(), userID)
    if err != nil || shopID == uuid.Nil {
        respondWithError(w, http.StatusForbidden, "seller has no shop")
        return
    }

    // Получаем информацию о магазине
    shop, err := h.adminService.GetShopByID(r.Context(), shopID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]interface{}{
        "id":          shop.ID,
        "name":        shop.Name,
        "is_verified": shop.IsVerified,
        "rating":      shop.Rating,
    })
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