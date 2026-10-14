package models

import (
	"time"
	"uuid"
)

type TicketStatus string

const (
	TicketStatusActive    TicketStatus = "active"
	TicketStatusUsed      TicketStatus = "used"
	TicketStatusCancelled TicketStatus = "cancelled"
	TicketStatusRefunded  TicketStatus = "refunded"
)

type Ticket struct {
	TicketID uuid.UUID `json:"ticket_id"`
	EventID  uuid.UUID `json:"event_id"`
	UserID   uuid.UUID `json:"user_id"`

	Price  int64  `json:"price"`
	Status string `json:"status"`

	BoughtAt   time.Time  `json:"bought_at"`
	ValidUntil time.Time  `json:"valid_until"`
	UsedAt     *time.Time `json:"used_at"`
}
