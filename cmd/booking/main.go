package main

import (
	"EventHub/internal/database"
	"EventHub/internal/logger"
	"EventHub/internal/repository"
	"EventHub/internal/service"
	"EventHub/internal/transport/grpc/booking"
	"EventHub/internal/worker"
	"context"
	"log/slog"
	"os"
	"time"
)

//gRPC сервер! booking

func main() {
	log := logger.New("booking-service", slog.LevelInfo)
	slog.SetDefault(log)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	pool, err := database.NewPostgresPool(ctx, os.Getenv("DATABASE_URL"))
	cancel()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	reservationRepo := repository.NewReservationRepository(pool)
	reservationService := service.NewReservationService(reservationRepo)

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	expirationWorker := worker.NewReservationExpirationWorker(reservationService, 5*time.Second)
	go expirationWorker.Run(workerCtx)

	slog.Info("starting booking grpc server", "port", 50051)

	if err := booking.Run(":50051", reservationService); err != nil {
		slog.Error("grpc server stopped", "error", err)
		os.Exit(1)
	}
}
