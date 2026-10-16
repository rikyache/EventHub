package handlers

import (
	"EventHub/internal/middleware"
	"EventHub/internal/models"
	"EventHub/internal/service"
	"encoding/json"
	"net/http"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.EventCreateRequest

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

	err = h.service.Create(ctx, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
