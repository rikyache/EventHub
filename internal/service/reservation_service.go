package service

import (
	"EventHub/internal/repository"
	"context"

	"github.com/google/uuid"
)

type ReservationService struct {
	repo *repository.ReservationRepository
}

func NewReservationService(repo *repository.ReservationRepository) *ReservationService {
	return &ReservationService{
		repo: repo,
	}
}

func (s *ReservationService) Create(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (uuid.UUID, error) {
	var reservationID uuid.UUID

	reservationID, err := s.repo.Create(ctx, eventID, userID)

	return reservationID, err
}
