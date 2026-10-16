package repository

import (
	"EventHub/internal/models"
	"context"

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
