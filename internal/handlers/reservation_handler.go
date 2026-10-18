package handlers

import (
	pb "EventHub/gen/booking/v1"
	"EventHub/internal/middleware"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//через HTTP API мы вызываем функции на удаленном сервере

type ReservationHandler struct {
	bookingClient pb.BookingServiceClient
}

func NewReservationHandler(bookingClient pb.BookingServiceClient) *ReservationHandler {
	return &ReservationHandler{
		bookingClient: bookingClient,
	}
}

func (h *ReservationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Missing UserID", http.StatusUnauthorized)
		return
	}

	eventID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Event ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.bookingClient.CreateReservation(ctx, &pb.CreateReservationRequest{
		EventId: eventID.String(),
		UserId:  userID.String(),
	},
	)
	if err != nil {
		handleBookingError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"reservation_id": resp.GetReservationId(),
	}); err != nil {
		slog.Error("failed to encode reservation response", "error", err)
	}
}

func (h *ReservationHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Missing UserID", http.StatusUnauthorized)
		return
	}

	reservationID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Reservation ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.bookingClient.ConfirmReservation(ctx, &pb.ConfirmReservationRequest{
		ReservationId: reservationID.String(),
		UserId:        userID.String(),
	})
	if err != nil {
		handleBookingError(w, err)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]string{"ticket_id": resp.GetTicketId()}); err != nil {
		slog.Error("failed to encode reservation response", "error", err)
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

	_, err = h.bookingClient.CancelReservation(
		ctx,
		&pb.CancelReservationRequest{
			ReservationId: reservationID.String(),
			UserId:        userID.String(),
		},
	)
	if err != nil {
		handleBookingError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func handleBookingError(w http.ResponseWriter, err error) {
	switch status.Code(err) {
	case codes.InvalidArgument:
		http.Error(w, "Invalid argument", http.StatusBadRequest)

	case codes.NotFound:
		http.Error(w, "Reservation not found", http.StatusNotFound)

	case codes.FailedPrecondition:
		http.Error(w, "Reservation unavailable", http.StatusConflict)

	case codes.Unavailable:
		slog.Error("booking service unavailable", "error", err)
		http.Error(w, "Booking service unavailable", http.StatusServiceUnavailable)

	case codes.DeadlineExceeded:
		slog.Error("booking service timeout", "error", err)
		http.Error(w, "Booking service timeout", http.StatusGatewayTimeout)

	case codes.Canceled:
		slog.Warn("booking request canceled", "error", err)
		http.Error(w, "Request canceled", http.StatusRequestTimeout)

	default:
		slog.Error("booking grpc request failed", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
