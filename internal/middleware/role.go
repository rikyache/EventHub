package middleware

import (
	"EventHub/internal/models"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type RoleLookup func(ctx context.Context, userID uuid.UUID) (models.UserRole, error)

func RequireRole(getRole RoleLookup, allowedRoles ...models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			role, err := getRole(r.Context(), userID)
			if err != nil {
				http.Error(w, "failed to check permissions", http.StatusInternalServerError)
				return
			}
			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}
