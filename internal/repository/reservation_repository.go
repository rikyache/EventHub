package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepository struct {
	db *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{
		db: db,
	}
}

func (h *ReservationRepository) Create(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (uuid.UUID, error) {
	query := `
	UPDATE events
	SET reserved_count = reserved_count + 1
		WHERE event_id = $1
  			AND sold_count + reserved_count < capacity
  			AND status = 'published'
  			AND start_date > NOW()
	RETURNING price;
`

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	var price int64

	err = tx.QueryRow(ctx, query, eventID).Scan(&price)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrReservationUnavailable
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("reserve event: %w", err)
	}

	query = `
		INSERT INTO reservations (
			reservation_id,
			event_id,
			user_id,
			price,
			status,
			expires_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			'pending',
			NOW() + INTERVAL '10 minutes'
		);
`
	reservationID := uuid.New()

	tag, err := tx.Exec(ctx, query, reservationID, eventID, userID, price)
	if err != nil {
		return uuid.Nil, fmt.Errorf("reserve event: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return uuid.Nil, ErrReservationUnavailable
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("reserve event: %w", err)
	}

	return reservationID, nil
}
