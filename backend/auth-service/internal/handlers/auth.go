package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/service"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary      Регистрация нового пользователя
// @Description  Создаёт нового пользователя с ролью customer или seller
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "Данные для регистрации"
// @Success      201 {object} models.UserResponse
// @Failure      400 {object} ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	response := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login godoc
// @Summary      Вход в систему
// @Description  Авторизация пользователя с выдачей JWT-токенов
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Email и пароль"
// @Success      200 {object} models.LoginResponse
// @Failure      401 {object} ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	user, accessToken, refreshToken, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	response := models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		User: &models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Phone:     user.Phone,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Me godoc
// @Summary      Получение профиля текущего пользователя
// @Description  Возвращает данные авторизованного пользователя
// @Tags         auth
// @Produce      json
// @Security     Bearer
// @Success      200 {object} models.UserResponse
// @Failure      401 {object} ErrorResponse
// @Router       /auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	response := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateProfile godoc
// @Summary      Обновление профиля пользователя
// @Description  Изменяет имя, фамилию и телефон пользователя
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body models.UpdateProfileRequest true "Данные для обновления"
// @Success      200 {object} models.UserResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /auth/profile [put]
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

	user, err := h.authService.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	})
}

// ChangePassword godoc
// @Summary      Смена пароля пользователя
// @Description  Меняет пароль при условии корректного старого пароля
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body models.ChangePasswordRequest true "Старый и новый пароль"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /auth/password [put]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := uuid.MustParse("6b75b13b-2b7b-4df1-b700-b39ac0bc1d45")

	err := h.authService.ChangePassword(r.Context(), userID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invalid old password" {
			status = http.StatusBadRequest
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Пароль успешно изменён"})
}

// ============================================================
// ВАЛИДАЦИЯ JWT (ДЛЯ API GATEWAY)
// ============================================================

// ValidateToken godoc
// @Summary      Проверка JWT-токена и получение user_id
// @Description  Проверяет валидность JWT-токена и возвращает user_id
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} map[string]string
// @Failure      401 {object} ErrorResponse
// @Router       /auth/validate [post]
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		respondWithError(w, http.StatusUnauthorized, "missing token")
		return
	}

	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}

	userID, err := h.authService.ValidateTokenAndGetUserID(tokenStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"user_id": userID.String()})
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

// RefreshToken — обновление access токена
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    accessToken, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "token_type":    "Bearer",
    })
}