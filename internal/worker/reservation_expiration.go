package worker

import (
	"context"
	"log/slog"
	"time"
)

type ReservationExpirer interface {
	Expire(ctx context.Context) (bool, error)
}

type ReservationExpirationWorker struct {
	expirer  ReservationExpirer
	interval time.Duration
}

func NewReservationExpirationWorker(expirer ReservationExpirer, interval time.Duration) *ReservationExpirationWorker {
	return &ReservationExpirationWorker{
		expirer:  expirer,
		interval: interval,
	}
}

// Run just to check
func (w *ReservationExpirationWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("reservation expiration worker stopped")
			return

		case <-ticker.C:
			w.processExpired(ctx)
		}
	}
}

// Process, continue working while having expired reservations, otherwise wait for next tick
func (w *ReservationExpirationWorker) processExpired(
	ctx context.Context,
) {
	for {
		if ctx.Err() != nil {
			return
		}

		processed, err := w.expirer.Expire(ctx)
		if err != nil {
			slog.Error("failed to expire reservation", "error", err)
			return
		}

		if !processed {
			return
		}
	}
}
