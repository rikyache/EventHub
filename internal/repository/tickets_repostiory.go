package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TicketRepository struct {
	db *pgxpool.Pool
}

func NewTicketRepository(db *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) BuyTicket(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	var price int64
	var validUntil time.Time

	query := `
	UPDATE events
	SET sold_count = sold_count + 1
	WHERE event_id = $1
		AND sold_count < capacity
		AND status = 'published'
		AND start_date > NOW()
	RETURNING price, end_date
`
	err = tx.QueryRow(ctx, query, eventID).Scan(
		&price,
		&validUntil,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrTicketUnavailable
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not buy ticket: %w", err)
	}

	ticketID := uuid.New()

	query = `
    INSERT INTO tickets (
        ticket_id,
        event_id,
        user_id,
        price,
        status,
        valid_until
    )
    VALUES ($1, $2, $3, $4, 'active', $5)
`

	_, err = tx.Exec(ctx, query,
		ticketID,
		eventID,
		userID,
		price,
		validUntil,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("could not buy ticket: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("could not commit: %w", err)
	}

	return ticketID, nil
}
