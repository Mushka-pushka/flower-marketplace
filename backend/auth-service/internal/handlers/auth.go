package handlers

import (
	"encoding/json"
	"log" 
	"fmt" 
	"io"  
	"net/http"
	"os"   
	"path/filepath"
	"strings"
	"time"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/middleware"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/service"
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
		AvatarURL: user.AvatarURL,
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
		log.Printf("Login: invalid request body: %v", err)
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt for email: %s", req.Email)

	user, accessToken, refreshToken, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		log.Printf("Login failed for %s: %v", req.Email, err) 
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	log.Printf("Login successful for %s", req.Email) 

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
			AvatarURL: user.AvatarURL,
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
		AvatarURL: user.AvatarURL,
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

	// Берем user_id из контекста
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

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
		AvatarURL: user.AvatarURL,
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

	// Берем user_id из контекста
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

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

// RefreshToken godoc
// @Summary      Обновление JWT-токенов
// @Description  Обновляет access и refresh токены
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshRequest true "Refresh токен"
// @Success      200 {object} RefreshResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /auth/refresh [post]
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

// ValidateToken godoc
// @Summary      Проверка JWT-токена
// @Description  Проверяет валидность JWT-токена и возвращает user_id
// @Tags         auth
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

// UploadAvatar — загрузка аватара пользователя
func (h *AuthHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
    // Получаем user_id из контекста
    userID, ok := middleware.GetUserIDFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "user not authenticated")
        return
    }

    // Ограничиваем размер файла (10MB)
    err := r.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "file too large (max 10MB)")
        return
    }

    // Получаем файл из запроса
    file, header, err := r.FormFile("avatar")
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "failed to get avatar file")
        return
    }
    defer file.Close()

    // Проверяем тип файла
    contentType := header.Header.Get("Content-Type")
    if !strings.HasPrefix(contentType, "image/") {
        respondWithError(w, http.StatusBadRequest, "file must be an image")
        return
    }

    // Создаём директорию для загрузок
    uploadDir := "uploads/avatars"
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        respondWithError(w, http.StatusInternalServerError, "failed to create upload directory")
        return
    }

    // Генерируем уникальное имя файла
    ext := filepath.Ext(header.Filename)
    filename := fmt.Sprintf("%s-%d%s", userID.String(), time.Now().UnixNano(), ext)
    filePath := filepath.Join(uploadDir, filename)

    // Сохраняем файл
    dst, err := os.Create(filePath)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "failed to save file")
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        respondWithError(w, http.StatusInternalServerError, "failed to save file")
        return
    }

    // Формируем URL аватара
    avatarURL := "/uploads/avatars/" + filename

    // Обновляем URL аватара в БД
    err = h.authService.UpdateAvatar(r.Context(), userID, &avatarURL)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    // Возвращаем обновлённого пользователя
    user, err := h.authService.GetUserByID(r.Context(), userID)
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
        AvatarURL: user.AvatarURL,
        CreatedAt: user.CreatedAt,
    })
}

// ============================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ (используются из admin.go)
// ============================================================

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}