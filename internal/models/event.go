package models

import (
	"time"
	"uuid"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCanceled  EventStatus = "canceled"
	EventStatusCompleted EventStatus = "completed"
)

type Event struct {
	EventID     uuid.UUID `json:"event_id"`
	OrganizerID uuid.UUID `json:"organizer_id"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Location    string `json:"location"`

	Capacity int `json:"capacity"`
	Price    int `json:"price"`

	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedDate time.Time `json:"created_date"`

	Status string `json:"status"`
}
