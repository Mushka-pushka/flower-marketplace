package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/service"
	"github.com/google/uuid"
)

type contextKey string

const UserKey contextKey = "user"
const UserIDKey contextKey = "user_id"

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// JWT — проверяет JWT-токен и добавляет пользователя в контекст
func (m *AuthMiddleware) JWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "authorization header required"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, `{"error": "invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		user, err := m.authService.GetUserFromToken(r.Context(), tokenString)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)
		ctx = context.WithValue(ctx, UserIDKey, user.ID)
		next(w, r.WithContext(ctx))
	}
}

// GetUserIDFromContext — извлекает user_id из контекста
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}