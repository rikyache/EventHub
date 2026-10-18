package models

import (
	"time"

	"github.com/google/uuid"
)

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusExpired   ReservationStatus = "expired"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusRejected  ReservationStatus = "rejected"
)

type Reservation struct {
	ReservationID uuid.UUID         `json:"reservationID"`
	EventID       uuid.UUID         `json:"eventID"`
	UserID        uuid.UUID         `json:"userID"`
	Price         int64             `json:"price"`
	Status        ReservationStatus `json:"status"`
	CreatedAt     time.Time         `json:"createdAt"`
	ExpiresAt     time.Time         `json:"expiresAt"`
}
