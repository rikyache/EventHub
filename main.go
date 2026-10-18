package main

import (
	"EventHub/internal/auth"
	"EventHub/internal/database"
	"EventHub/internal/handlers"
	"EventHub/internal/logger"
	"EventHub/internal/middleware"
	"EventHub/internal/models"
	"EventHub/internal/repository"
	"EventHub/internal/service"
	"EventHub/internal/worker"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// TODO:
// TODO: реализовать три сервиса event service/postgresql
// TODO: реализовать booking service/postgresql, redis, kafka -> noti worker
// TODO: payment service -> postgresql

func main() {
	// Logger
	log := logger.New("event-service", slog.LevelInfo)
	slog.SetDefault(log)

	// Config
	dsn := os.Getenv("DATABASE_URL")

	// Create context related to application life cycle
	dbCtx, dbCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	// Open database for connections
	pool, err := database.NewPostgresPool(dbCtx, dsn)
	dbCancel()

	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Initiate token manager for JWT logic
	tokenManager, err := auth.NewTokenManager(os.Getenv("JWT_SECRET"))
	if err != nil {
		slog.Error("failed to initialize JWT token manager", "error", err)
		os.Exit(1)
	}

	// Repository init
	userRepo := repository.NewUsersRepository(pool)
	eventRepo := repository.NewEventRepository(pool)
	ticketRepo := repository.NewTicketRepository(pool)
	reservationRepo := repository.NewReservationRepository(pool)

	// Service init
	authService := service.NewAuthService(userRepo, tokenManager)
	userService := service.NewUserService(userRepo)
	eventService := service.NewEventService(eventRepo)
	ticketService := service.NewTicketService(ticketRepo)
	reservationService := service.NewReservationService(reservationRepo)

	//	Handlers init
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	eventHandler := handlers.NewEventHandler(eventService)
	ticketHandler := handlers.NewTicketHandler(ticketService)
	reservationHandler := handlers.NewReservationHandler(reservationService)

	// Background worker
	workerCtx, workerCancel := context.WithCancel(
		context.Background(),
	)
	defer workerCancel()

	expirationWorker := worker.NewReservationExpirationWorker(reservationService, 5*time.Second)
	go expirationWorker.Run(workerCtx)

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
	mux.HandleFunc("POST /auth/refresh", func(w http.ResponseWriter, r *http.Request) {})

	// EVENTS
	mux.Handle(
		"POST /events",
		middleware.Auth(tokenManager)(
			middleware.RequireRole(
				userService.GetRoleByID,
				models.RoleOrganizer,
				models.RoleAdmin,
			)(
				http.HandlerFunc(eventHandler.Create),
			),
		),
	)
	mux.HandleFunc("GET /events", eventHandler.GetEvents)
	mux.HandleFunc("GET /events/{id}", eventHandler.GetEventByID)
	mux.HandleFunc("PATCH /events/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /events/{id}/publish", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /events/{id}/cancel", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /events/{id}/tickets", func(w http.ResponseWriter, r *http.Request) {})
	mux.Handle("POST /events/{id}/reservations",
		middleware.Auth(tokenManager)(
			http.HandlerFunc(reservationHandler.Create),
		))
	mux.HandleFunc("POST /reservations/{id}/confirm", reservationHandler.Confirm)

	//TICKETS
	mux.Handle("POST /events/{id}/tickets", middleware.Auth(tokenManager)(http.HandlerFunc(ticketHandler.BuyTicket)))
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
