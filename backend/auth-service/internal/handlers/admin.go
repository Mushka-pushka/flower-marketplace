package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/service"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GetSellers — получение списка продавцов
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

// VerifySeller — верификация продавца
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

// UpdateUserStatus — блокировка/разблокировка пользователя
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

// GetUsersList — получение списка пользователей
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