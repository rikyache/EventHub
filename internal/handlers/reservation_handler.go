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

func (h *ReservationHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Missing UserID", http.StatusUnauthorized)
		return
	}

	reservationID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Reservation ID", http.StatusBadRequest)
		return
	}

	ticketID, err := h.service.Confirm(ctx, reservationID, userID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReservationNotFound):
			http.Error(w, "Reservation not found", http.StatusNotFound)

		case errors.Is(err, repository.ErrReservationExpired):
			http.Error(w, "Reservation expired", http.StatusConflict)

		case errors.Is(err, repository.ErrReservationNotPending):
			http.Error(w, "Reservation is not pending", http.StatusConflict)

		case errors.Is(err, repository.ErrInsufficientFunds):
			http.Error(w, "Insufficient funds", http.StatusConflict)

		default:
			slog.Error("failed to confirm reservation", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]uuid.UUID{
		"ticket_id": ticketID,
	}); err != nil {
		slog.Error("failed to encode confirmation response", "error", err)
	}
}

func (h *ReservationHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Missing UserID", http.StatusUnauthorized)
		return
	}

	reservationID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Reservation ID", http.StatusBadRequest)
		return
	}

	err = h.service.Cancel(ctx, userID, reservationID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReservationNotFound):
			http.Error(w, "Reservation not found", http.StatusNotFound)

		case errors.Is(err, repository.ErrReservationNotPending):
			http.Error(w, "Reservation is not pending", http.StatusConflict)

		case errors.Is(err, repository.ErrReservationUnavailable):
			http.Error(w, "Reservation unavailable", http.StatusConflict)

		default:
			slog.Error("failed to cancel reservation",
				"reservation_id", reservationID,
				"error", err,
			)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
