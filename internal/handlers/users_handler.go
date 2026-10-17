package handlers

import (
	"EventHub/internal/DTO"
	"EventHub/internal/middleware"
	"EventHub/internal/models"
	"EventHub/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	user, err := h.service.GetByID(ctx, userID)
	if err != nil {
		slog.Error("failed to get  user", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	response := models.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		Role:      user.Role,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func (h *UserHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.DepositRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("failed to decode request", "error", err)
		return
	}

	if req.Amount <= 0 {
		slog.Error("invalid amount", "amount", req.Amount)
		return
	}

	ctx := r.Context()

	err = h.service.Deposit(ctx, userID, req.Amount)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
