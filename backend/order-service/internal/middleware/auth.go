package middleware

import (
	"context"
	"net/http"

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/config"
	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type AuthMiddleware struct {
	cfg *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

// AuthMiddleware — проверяет X-User-ID заголовок
func (m *AuthMiddleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			http.Error(w, `{"error": "user not authenticated"}`, http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, `{"error": "invalid user id"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

// GetUserIDFromContext — извлекает user_id из контекста
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}