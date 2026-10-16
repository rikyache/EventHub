package middleware

import (
	"context"

	"github.com/google/uuid"
)

const eventIDKey contextKey = "eventID"

func EventIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(eventIDKey).(uuid.UUID)
	return id, ok
}
