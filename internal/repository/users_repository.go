package repository

import (
	"EventHub/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
	INSERT INTO users(
	    user_id, username,name, surname, email, 
	    password_hash, role, balance, created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

	_, err := r.db.Exec(
		ctx,
		query,
		user.UserID,
		user.Username,
		user.Name,
		user.Surname,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Balance,
		user.CreatedAt,
	)

	if err != nil {
		return err
	}

	if err != nil {

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_unique":
				return ErrEmailTaken
			case "users_username_unique":
				return ErrUsernameTaken
			}
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
	SELECT
		user_id,
		username,
		name,
		surname,
		email,
		password_hash,
		role,
		balance,
		created_at
	FROM users
	WHERE email = $1
`
	var user models.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.UserID,
		&user.Username,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Balance,
		&user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `
	SELECT
		user_id,
		username,
		name,
		surname,
		email,
		password_hash,
		role,
		balance,
		created_at
	FROM users
	WHERE user_id = $1
`
	var user models.User

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Balance,
		&user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetRoleByID(ctx context.Context, userID uuid.UUID) (models.UserRole, error) {
	query := `
        SELECT role
        FROM users
        WHERE user_id = $1
    `

	var role models.UserRole

	err := r.db.QueryRow(ctx, query, userID).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}
