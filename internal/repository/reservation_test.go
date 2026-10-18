package repository

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupReservationTest(t *testing.T) (*ReservationRepository, *pgxpool.Pool, uuid.UUID, uuid.UUID) {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("create test pool: %v", err)
	}

	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping test database: %v", err)
	}

	userID := uuid.New()
	eventID := uuid.New()

	// UUID в username/email исключает конфликты UNIQUE
	username := "test_" + userID.String()[:8]
	email := "test_" + userID.String() + "@example.com"

	_, err = pool.Exec(ctx, `
        INSERT INTO users (
            user_id,
            username,
            name,
            surname,
            email,
            password_hash
        )
        VALUES ($1, $2, 'Test', 'User', $3, 'test_hash')
    `, userID, username, email)

	if err != nil {
		t.Fatalf("create test user: %v", err)
	}

	_, err = pool.Exec(ctx, `
        INSERT INTO events (
            event_id,
            organizer_id,
            title,
            type,
            location,
            capacity,
            price,
            start_date,
            end_date,
            status
        )
        VALUES (
            $1, $2,
            'Concurrency Test',
            'concert',
            'Test Location',
            10,
            700,
            NOW() + INTERVAL '1 day',
            NOW() + INTERVAL '2 days',
            'published'
        )
    `, eventID, userID)

	if err != nil {
		t.Fatalf("create test event: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cleanupCancel()

		// Удаляем сначала дочерние записи из-за FK.
		if _, err := pool.Exec(cleanupCtx,
			`DELETE FROM reservations WHERE event_id = $1`,
			eventID,
		); err != nil {
			t.Errorf("cleanup reservations: %v", err)
		}

		if _, err := pool.Exec(cleanupCtx,
			`DELETE FROM tickets WHERE event_id = $1`,
			eventID,
		); err != nil {
			t.Errorf("cleanup tickets: %v", err)
		}

		if _, err := pool.Exec(cleanupCtx,
			`DELETE FROM events WHERE event_id = $1`,
			eventID,
		); err != nil {
			t.Errorf("cleanup event: %v", err)
		}

		if _, err := pool.Exec(cleanupCtx,
			`DELETE FROM users WHERE user_id = $1`,
			userID,
		); err != nil {
			t.Errorf("cleanup user: %v", err)
		}
	})

	repo := NewReservationRepository(pool)

	return repo, pool, eventID, userID
}

func TestConcurrentReservations(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	// Подготавливает отдельную тестовую БД,
	// пользователя и мероприятие на 10 мест.
	repo, pool, eventID, userID := setupReservationTest(t)

	const requests = 100
	const capacity = 10

	var wg sync.WaitGroup

	start := make(chan struct{})
	results := make(chan error, requests)

	for i := 0; i < requests; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Ждём общего сигнала запуска
			<-start

			_, err := repo.Create(
				ctx,
				eventID,
				userID,
			)

			results <- err
		}()
	}

	// Одновременно отпускаем все горутины
	close(start)

	wg.Wait()
	close(results)

	var success, unavailable, unexpected int

	for err := range results {
		switch {
		case err == nil:
			success++

		case errors.Is(err, ErrReservationUnavailable):
			unavailable++

		default:
			unexpected++
			t.Logf("unexpected error: %v", err)
		}
	}

	if success != capacity {
		t.Errorf(
			"expected %d reservations, got %d",
			capacity,
			success,
		)
	}

	if unavailable != requests-capacity {
		t.Errorf(
			"expected %d unavailable, got %d",
			requests-capacity,
			unavailable,
		)
	}

	if unexpected != 0 {
		t.Errorf("unexpected errors: %d", unexpected)
	}

	var reservedCount int
	var actualReservations int

	err := pool.QueryRow(ctx, `
        SELECT reserved_count
        FROM events
        WHERE event_id = $1
    `, eventID).Scan(&reservedCount)

	if err != nil {
		t.Fatal(err)
	}

	err = pool.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM reservations
        WHERE event_id = $1
          AND status = 'pending'
    `, eventID).Scan(&actualReservations)

	if err != nil {
		t.Fatal(err)
	}

	if reservedCount != capacity {
		t.Errorf(
			"reserved_count = %d, want %d",
			reservedCount,
			capacity,
		)
	}

	if actualReservations != capacity {
		t.Errorf(
			"reservations = %d, want %d",
			actualReservations,
			capacity,
		)
	}
}
