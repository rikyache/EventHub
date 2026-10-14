package handlers

import (
	"EventHub/internal/models"
	"EventHub/internal/repository"
	"EventHub/internal/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"uuid"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Service: service,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := service.RegisterInput{
		Username: req.Username,
		Name:     req.Name,
		Surname:  req.Surname,
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := r.Context()

	userID, err := h.Service.Register(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPassword):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, repository.ErrEmailTaken),
			errors.Is(err, repository.ErrUsernameTaken):
			http.Error(w, "email or username already taken", http.StatusConflict)

		default:
			slog.ErrorContext(r.Context(), "failed to register user", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]uuid.UUID{"user_id": userID}); err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

}
