package middleware

import (
    "context"
    "net/http"

    "github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// AuthMiddleware — проверяет X-User-ID заголовок
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
    userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
    return userID, ok
}