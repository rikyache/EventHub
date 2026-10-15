package service

import (
	"context"
	"fmt"
	"time"

	"EventHub/internal/models"
	"EventHub/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

type RegisterInput struct {
	Username string
	Name     string
	Surname  string
	Email    string
	Password string
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (uuid.UUID, error) {
	if len(input.Password) < 8 || len([]byte(input.Password)) > 72 {
		return uuid.Nil, ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("hash password: %w", err)
	}

	user := models.User{
		UserID:       uuid.New(),
		Username:     input.Username,
		Name:         input.Name,
		Surname:      input.Surname,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         models.RoleUser,
		Balance:      0,
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	return user.UserID, nil
}
