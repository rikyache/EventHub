package handlers

import (
	"EventHub/internal/middleware"
	"EventHub/internal/models"
	"EventHub/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req *models.EventCreateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, err := h.service.Create(ctx, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]uuid.UUID{
		"event_id": eventID,
	})
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	events, err := h.service.GetEvents(ctx)
	if err != nil {
		slog.Error("failed to get events", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(events)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(events); err != nil {
		slog.Error("failed to encode events", "error", err)
	}
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	eventID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid event ID", http.StatusBadRequest)
		return
	}

	var event *models.Event

	event, err = h.service.GetEventByID(ctx, eventID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		slog.Error("failed to encode event", "error", err)
		return
	}
}
