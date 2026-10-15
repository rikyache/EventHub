package repository

import (
	"EventHub/internal/models"
	"context"
	"errors"
	"fmt"

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
