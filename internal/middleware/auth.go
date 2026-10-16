package middleware

import (
	"EventHub/internal/auth"
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "userID"

func Auth(tokens *auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			scheme, token, ok := strings.Cut(header, " ")

			if !ok || !strings.EqualFold(scheme, "Bearer") ||
				token == "" || strings.ContainsAny(token, " \t") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			userID, err := tokens.Validate(token)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
