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
	eventID       uuid.UUID         `json:"eventID"`
	userID        uuid.UUID         `json:"userID"`
	price         int64             `json:"price"`
	status        ReservationStatus `json:"status"`
	createdAt     time.Time         `json:"createdAt"`
	ExpiresAt     time.Time         `json:"expiresAt"`
}
