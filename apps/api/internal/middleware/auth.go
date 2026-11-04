package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/farrelarea/budgeting-app/apps/api/internal/auth"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// RequireAuth validates the Bearer token and injects userID into context.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		claims, err := auth.ParseToken(strings.TrimPrefix(header, "Bearer "), auth.SecretKey())
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext returns the authenticated user ID if present.
func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDKey).(string)
	return v, ok
}
