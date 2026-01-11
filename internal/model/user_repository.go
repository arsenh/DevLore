package model

import "context"

type UserRepository interface {
	FindByID(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) *User
	Create(ctx context.Context, u *User) error
	List(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id int) error
}
