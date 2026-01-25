package model

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
	List(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
