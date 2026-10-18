package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	pb "EventHub/gen/booking/v1"
	"EventHub/internal/auth"
	"EventHub/internal/database"
	"EventHub/internal/handlers"
	"EventHub/internal/logger"
	"EventHub/internal/middleware"
	"EventHub/internal/models"
	"EventHub/internal/repository"
	"EventHub/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Logger
	log := logger.New("event-service", slog.LevelInfo)
	slog.SetDefault(log)

	// Config
	dsn := os.Getenv("DATABASE_URL")

	// PostgreSQL
	dbCtx, dbCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	pool, err := database.NewPostgresPool(dbCtx, dsn)
	dbCancel()

	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		return
	}
	defer pool.Close()

	// JWT
	tokenManager, err := auth.NewTokenManager(os.Getenv("JWT_SECRET"))
	if err != nil {
		slog.Error("failed to initialize JWT token manager", "error", err)
		return
	}

	// Repositories
	userRepo := repository.NewUsersRepository(pool)
	eventRepo := repository.NewEventRepository(pool)
	ticketRepo := repository.NewTicketRepository(pool)

	// Services
	authService := service.NewAuthService(userRepo, tokenManager)
	userService := service.NewUserService(userRepo)
	eventService := service.NewEventService(eventRepo)
	ticketService := service.NewTicketService(ticketRepo)

	// gRPC client
	conn, err := grpc.NewClient(
		"booking:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("failed to create grpc client", "error", err)
		return
	}
	defer conn.Close()

	bookingClient := pb.NewBookingServiceClient(conn)

	// HTTP handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	eventHandler := handlers.NewEventHandler(eventService)
	ticketHandler := handlers.NewTicketHandler(ticketService)
	reservationHandler := handlers.NewReservationHandler(bookingClient)

	// Router
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Healthy")
	})

	// Users
	mux.Handle(
		"GET /users/me",
		middleware.Auth(tokenManager)(
			http.HandlerFunc(userHandler.Profile),
		),
	)

	mux.HandleFunc("GET /users/tickets", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /users/deposit", func(w http.ResponseWriter, r *http.Request) {})

	// Auth
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /auth/refresh", func(w http.ResponseWriter, r *http.Request) {})

	// Events
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

	// Reservations — через gRPC
	mux.Handle(
		"POST /events/{id}/reservations",
		middleware.Auth(tokenManager)(
			http.HandlerFunc(reservationHandler.Create),
		),
	)

	mux.Handle(
		"POST /reservations/{id}/confirm",
		middleware.Auth(tokenManager)(
			http.HandlerFunc(reservationHandler.Confirm),
		),
	)

	mux.Handle(
		"POST /reservations/{id}/cancel",
		middleware.Auth(tokenManager)(
			http.HandlerFunc(reservationHandler.Cancel),
		),
	)

	// Tickets
	mux.Handle(
		"POST /events/{id}/tickets",
		middleware.Auth(tokenManager)(
			http.HandlerFunc(ticketHandler.BuyTicket),
		),
	)

	mux.HandleFunc("GET /tickets/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /tickets/{id}/refund", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /tickets/{id}/check-in", func(w http.ResponseWriter, r *http.Request) {})

	// HTTP server
	slog.Info("HTTP server started", "port", 8080)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("HTTP server stopped", "error", err)
	}
}
