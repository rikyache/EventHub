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

type ReservationRepository struct {
	db *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{
		db: db,
	}
}

func (r *ReservationRepository) Create(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (uuid.UUID, error) {
	query := `
	UPDATE events
	SET reserved_count = reserved_count + 1
		WHERE event_id = $1
  			AND sold_count + reserved_count < capacity
  			AND status = 'published'
  			AND start_date > NOW()
	RETURNING price;
`

	tx, err := r.db.Begin(ctx)
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

func (r *ReservationRepository) Confirm(ctx context.Context, reservationID, userID uuid.UUID) (uuid.UUID, error) {
	query := `
	SELECT event_id, price, status, expires_at
		FROM reservations
		WHERE reservation_id = $1
  			AND user_id = $2
	FOR UPDATE;
`
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	var eventID uuid.UUID
	var price int64
	var status string
	var expiresAt time.Time

	err = tx.QueryRow(
		ctx,
		query,
		reservationID,
		userID,
	).Scan(
		&eventID,
		&price,
		&status,
		&expiresAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrReservationNotFound
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("get reservation: %w", err)
	}

	var expired bool

	err = tx.QueryRow(ctx, `
		SELECT expires_at <= clock_timestamp()
		FROM reservations
		WHERE reservation_id = $1
	`, reservationID).Scan(&expired)

	if err != nil {
		return uuid.Nil, err
	}

	if expired {
		return uuid.Nil, ErrReservationExpired
	}

	if status != "pending" {
		return uuid.Nil, ErrReservationNotPending
	}

	query = `
		UPDATE users
		SET balance = balance - $1
		WHERE user_id = $2
		AND balance >= $1
	`

	tag, err := tx.Exec(ctx, query, price, userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("reserve user: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return uuid.Nil, ErrInsufficientFunds
	}

	query = `
	UPDATE events
		SET reserved_count = reserved_count - 1,
		    sold_count = sold_count + 1
		WHERE event_id = $1
			AND reserved_count > 0
`

	tag, err = tx.Exec(ctx, query, eventID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("update event counters: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return uuid.Nil, ErrReservationUnavailable
	}

	query = `
    SELECT end_date
    FROM events
    WHERE event_id = $1
`

	var validUntil time.Time

	err = tx.QueryRow(ctx, query, eventID).Scan(&validUntil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get event end date: %w", err)
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

	_, err = tx.Exec(
		ctx,
		query,
		ticketID,
		eventID,
		userID,
		price,
		validUntil,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create ticket: %w", err)
	}

	query = `
    UPDATE reservations
    SET status = 'confirmed'
    WHERE reservation_id = $1
      AND user_id = $2
      AND status = 'pending'
`

	tag, err = tx.Exec(ctx, query, reservationID, userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("confirm reservation: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return uuid.Nil, ErrReservationNotPending
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit reservation: %w", err)
	}

	return ticketID, nil
}

// find and set reservation status to expired
func (r *ReservationRepository) Expire(ctx context.Context) (bool, error) {
	query := `
    SELECT reservation_id, event_id
    FROM reservations
    WHERE status = 'pending'
      AND expires_at <= clock_timestamp()
    ORDER BY expires_at, reservation_id
    LIMIT 1
    FOR UPDATE SKIP LOCKED
`
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	var reservationID uuid.UUID
	var eventID uuid.UUID

	err = tx.QueryRow(ctx, query).Scan(
		&reservationID,
		&eventID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("find expired reservation: %w", err)
	}

	query = `
	UPDATE reservations
	SET status = 'expired'
	WHERE reservation_id = $1
	  AND status = 'pending'
	  AND expires_at <= clock_timestamp();
`

	tag, err := tx.Exec(ctx, query, reservationID)
	if err != nil {
		return false, fmt.Errorf("update reservation status: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return false, fmt.Errorf(
			"expected to expire one reservation, affected %d rows",
			tag.RowsAffected(),
		)
	}

	query = `
	UPDATE events
	SET reserved_count = reserved_count - 1
	WHERE event_id = $1
	  AND reserved_count > 0;
`

	tag, err = tx.Exec(ctx, query, eventID)
	if err != nil {
		return false, fmt.Errorf("release reserved seat: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return false, fmt.Errorf(
			"expected to release one seat, affected %d rows",
			tag.RowsAffected(),
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("commit expiration: %w", err)
	}

	return true, nil
}
