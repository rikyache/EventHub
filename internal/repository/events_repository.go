package repository

import (
	"EventHub/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	query := `
    	INSERT INTO events (
            event_id,
            organizer_id,
            title,
            description,
            type,
            location,
            capacity,
            price,
            start_date,
            end_date,
            create_date,
            status
        )
    	VALUES ($1, $2, $3, $4, $5, $6, 
    	        $7, $8, $9, $10, $11, $12)
`

	_, err := r.db.Exec(ctx, query,
		event.EventID,
		event.OrganizerID,
		event.Title,
		event.Description,
		event.Type,
		event.Location,
		event.Capacity,
		event.Price,
		event.StartDate,
		event.EndDate,
		event.CreatedDate,
		event.Status,
	)
	return err
}

func (r *EventRepository) GetEvents(ctx context.Context) ([]*models.Event, error) {
	query := `
	SELECT 
		event_id,
		organizer_id,
		title,
		description,
		type,
		location,
		capacity,
		price,
		start_date,
		end_date,
		status
	FROM events
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*models.Event, 0)

	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.EventID,
			&event.OrganizerID,
			&event.Title,
			&event.Description,
			&event.Type,
			&event.Location,
			&event.Capacity,
			&event.Price,
			&event.StartDate,
			&event.EndDate,
			&event.CreatedDate,
			&event.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}

		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events: %w", err)
	}

	return events, nil
}

func (r *EventRepository) GetEventByID(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	query := `
	SELECT 
		event_id,
		organizer_id,
		title,
		description,
		type,
		location,
		capacity,
		price,
		start_date,
		end_date,
		create_date,
		status
	FROM events
	WHERE event_id = $1`

	var event models.Event

	err := r.db.QueryRow(ctx, query, eventID).Scan(
		&event.EventID,
		&event.OrganizerID,
		&event.Title,
		&event.Description,
		&event.Type,
		&event.Location,
		&event.Capacity,
		&event.Price,
		&event.StartDate,
		&event.EndDate,
		&event.CreatedDate,
		&event.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("get event by ID: %w", err)
	}

	return &event, nil
}
