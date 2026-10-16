package service

import (
	"EventHub/internal/models"
	"EventHub/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) Create(ctx context.Context, userID uuid.UUID, req *models.EventCreateRequest) (uuid.UUID, error) {
	if req == nil {
		return uuid.Nil, errors.New("event request is nil")
	}

	if strings.TrimSpace(req.Title) == "" {
		return uuid.Nil, errors.New("event title cannot be empty")
	}

	if req.Price < 0 {
		return uuid.Nil, errors.New("event price cannot be negative")
	}

	if req.Capacity < 1 {
		return uuid.Nil, errors.New("event capacity must be greater than zero")
	}

	if !req.StartDate.Before(req.EndDate) {
		return uuid.Nil, errors.New("event start date must be before end date")
	}

	if !req.StartDate.After(time.Now()) {
		return uuid.Nil, errors.New("event start date must be in the future")
	}

	switch req.Type {
	case models.EventTypeConcert,
		models.EventTypeTheatre,
		models.EventTypeFestival,
		models.EventTypeCinema,
		models.EventTypeConference,
		models.EventTypeSport,
		models.EventTypeExhibition:
	default:
		return uuid.Nil, errors.New("invalid event type")
	}

	event := models.Event{
		EventID:     uuid.New(),
		OrganizerID: userID,

		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Location:    req.Location,
		Capacity:    req.Capacity,
		Price:       req.Price,

		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CreatedDate: time.Now().UTC(),
		Status:      models.EventStatusDraft,
	}

	err := s.repo.CreateEvent(ctx, &event)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create event: %w", err)
	}

	return event.EventID, nil
}

func (s *EventService) GetEvents(ctx context.Context) ([]*models.Event, error) {
	events, err := s.repo.GetEvents(ctx)

	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}

	return events, nil
}
