package main

import (
	"EventHub/internal/auth"
	"EventHub/internal/database"
	"EventHub/internal/handlers"
	"EventHub/internal/logger"
	"EventHub/internal/middleware"
	"EventHub/internal/repository"
	"EventHub/internal/service"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// TODO:
// реализовать три сервиса event service/psotgresql
// реализовать booking service/postgresql, redis, kafka -> noti worker
// paymant service -> postgresql

func main() {
	// logger
	log := logger.New("event-service", slog.LevelInfo)
	slog.SetDefault(log)

	// cfg

	dsn := os.Getenv("DATABASE_URL")

	// db

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := database.NewPostgresPool(ctx, dsn)
	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	tokenManager, err := auth.NewTokenManager(os.Getenv("JWT_SECRET"))
	if err != nil {
		slog.Error("failed to initialize JWT token manager", "error", err)
		os.Exit(1)
	}

	// repo
	userRepo := repository.NewUsersRepository(pool)
	// service
	authService := service.NewAuthService(userRepo, tokenManager)
	//	handler
	authHandler := handlers.NewAuthHandler(authService)

	// router
	mux := http.NewServeMux()
	// routes
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Healthy")
	})
	// USERS
	mux.Handle("GET /users/me", middleware.Auth(tokenManager)(http.HandlerFunc(userHandler.Profile)))
	mux.HandleFunc("GET /users/tickets", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /users/deposit", func(w http.ResponseWriter, r *http.Request) {})
	// AUTH
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)
	// mux.HandleFunc("POST /auth/refresh", func(w http.ResponseWriter, r *http.Request) {})
	// EVENTS
	mux.HandleFunc("POST /events", func(w http.ResponseWriter, r *http.Request) {}) // создание мероприятия
	mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /events/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("PATCH /events/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /events/{id}/publish", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /events/{id}/cancel", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /events/{id}/tickets", func(w http.ResponseWriter, r *http.Request) {})
	//TICKETS
	mux.HandleFunc("POST /events/{id}/tickets", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /tickets/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /tickets/{id}/refund", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /tickets/{id}/check-in", func(w http.ResponseWriter, r *http.Request) {})
	//http server

	slog.Info("server started", "port", 8080)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("HTTP server stopped", "error", err)
		os.Exit(1)
	}
}
