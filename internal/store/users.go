package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jesusthecreator017/fswithgo/internal/store/dbsqlc"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Permissions  int32     `json:"permissions"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserStore struct {
	queries *dbsqlc.Queries
}

func (u *UserStore) Create(ctx context.Context, user *User) error {
	row, err := u.queries.CreateUser(ctx, dbsqlc.CreateUserParams{
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
	})

	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	user.ID = row.ID
	user.Permissions = row.Permissions
	user.CreatedAt = row.CreatedAt.Time
	user.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (u *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	row, err := u.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting user: %w", err)
	}

	return &User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Name:         row.Name,
		Permissions:  row.Permissions,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}, nil
}

func (u *UserStore) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row, err := u.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting user: %w", err)
	}

	return &User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Name:         row.Name,
		Permissions:  row.Permissions,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}, nil
}
