package service

import (
	"EventHub/internal/auth"
	"context"
	"errors"
	"fmt"
	"time"

	"EventHub/internal/models"
	"EventHub/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   *repository.UserRepository
	tokens *auth.TokenManager
}

func NewAuthService(
	repo *repository.UserRepository,
	tokens *auth.TokenManager,
) *AuthService {
	return &AuthService{
		repo:   repo,
		tokens: tokens,
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

func (s *AuthService) Login(
	ctx context.Context,
	req models.LoginRequest,
) (string, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)

	if errors.Is(err, repository.ErrUserNotFound) {
		return "", ErrInvalidCredentials
	}

	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(user.UserID)
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {

	return nil
}
