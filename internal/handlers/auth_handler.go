package handlers

import (
	"EventHub/internal/service"
	"net/http"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Service: service,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request)    {}
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request)   {}
