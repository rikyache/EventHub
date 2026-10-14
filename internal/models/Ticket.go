package models

import (
	"time"
	"uuid"
)

type Ticket struct {
	TicketID uuid.UUID `json:"ticket_id"`
	EventID  uuid.UUID `json:"event_id"`
	UserID   uuid.UUID `json:"user_id"`

	Price  int    `json:"price"`
	Status string `json:"status"`

	BoughtAt   time.Time  `json:"bought_at"`
	ValidUntil time.Time  `json:"valid_until"`
	UsedAt     *time.Time `json:"used_at"`
}
