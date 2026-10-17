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

type TicketHandler struct {
	service *service.TicketService
}

func NewTicketHandler(service *service.TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

func (h *TicketHandler) BuyTicket(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid event ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ticketID, err := h.service.CreateTicket(ctx, eventID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTicketUnavailable) {
			http.Error(w, "ticket unavailable", http.StatusConflict)
			return
		}

		slog.Error("failed to buy ticket", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]uuid.UUID{
		"ticket_id": ticketID}); err != nil {
		slog.Error("failed to encode ticket response", "error", err)
	}
}
