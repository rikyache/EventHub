package handlers

import (
	"EventHub/internal/middleware"
	"EventHub/internal/repository"
	"EventHub/internal/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type ReservationHandler struct {
	service *service.ReservationService
}

func NewReservationHandler(service *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		service: service,
	}
}

func (h *ReservationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Missing UserID", http.StatusUnauthorized)
		return
	}

	eventID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Event ID", http.StatusBadRequest)
		return
	}

	reservationID, err := h.service.Create(ctx, eventID, userID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReservationUnavailable):
			http.Error(w, "reservation unavailable", http.StatusConflict)

		default:
			slog.Error("failed to create reservation", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]uuid.UUID{
		"reservation_id": reservationID,
	}); err != nil {
		slog.Error("failed to encode reservation response", "error", err)
	}
}
