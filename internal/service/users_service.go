package service

import (
	"EventHub/internal/models"
	"EventHub/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	return user, nil
}

func (s *UserService) GetRoleByID(ctx context.Context, userID uuid.UUID) (models.UserRole, error) {
	role, err := s.repo.GetRoleByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get user role: %w", err)
	}

	return role, nil
}
