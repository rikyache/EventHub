package service

import (
	"EventHub/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type TicketService struct {
	repo *repository.TicketRepository
}

func NewTicketService(repo *repository.TicketRepository) *TicketService {
	return &TicketService{
		repo: repo,
	}
}

func (r *TicketService) CreateTicket(ctx context.Context, eventID, userID uuid.UUID) (uuid.UUID, error) {
	ticketID, err := r.repo.BuyTicket(ctx, eventID, userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("buy ticket: %w", err)
	}
	return ticketID, err
}
