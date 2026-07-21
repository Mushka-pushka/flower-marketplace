package middleware

import (
	"context"
	"log" 
	"net/http"
	"strings"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/service"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserKey   contextKey = "user"
	UserIDKey contextKey = "user_id"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// JWT — проверяет JWT-токен и добавляет пользователя в контекст
func (m *AuthMiddleware) JWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("JWT middleware called")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("JWT: no Authorization header")
			http.Error(w, `{"error": "authorization header required"}`, http.StatusUnauthorized)
			return
		}

		log.Printf("JWT: Authorization header: %s", authHeader)

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Println("JWT: invalid header format")
			http.Error(w, `{"error": "invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		log.Printf("JWT: token: %s", tokenString)

		user, err := m.authService.GetUserFromToken(r.Context(), tokenString)
		if err != nil {
			log.Printf("JWT: GetUserFromToken error: %v", err)
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}
		log.Printf("JWT: user found: %s", user.Email)

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

// GetUserFromContext — извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(UserKey).(*models.User)
	return user, ok
}