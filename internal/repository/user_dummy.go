package repository

import (
	"context"
	"fmt"

	"github.com/arsenh/DevLore/internal/model"
)

type DummyUserRepository struct {
	db []model.User
}

func NewDummyUserRepository() *DummyUserRepository {
	return &DummyUserRepository{
		db: []model.User{
			{ID: 11, FullName: "Alice Johnson"},
			{ID: 22, FullName: "Bob Smith"},
			{ID: 33, FullName: "Carol Adams"},
			{ID: 44, FullName: "Dave Brown"},
			{ID: 55, FullName: "Eve Thompson"},
			{ID: 66, FullName: "Frank Williams"},
			{ID: 77, FullName: "Grace Miller"},
			{ID: 88, FullName: "Heidi Clark"},
			{ID: 99, FullName: "Ivan Garcia"},
			{ID: 111, FullName: "Judy Martinez"},
		},
	}
}

func (d DummyUserRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	for _, user := range d.db {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with id = %d not found", id)
}

func (d DummyUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, user := range d.db {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with email = %s not found", email)
}

func (d DummyUserRepository) Create(ctx context.Context, u *model.User) error {
	d.db = append(d.db, *u)
	return nil
}

func (d DummyUserRepository) List(ctx context.Context) ([]model.User, error) {
	return d.db, nil
}

func (d DummyUserRepository) Delete(ctx context.Context, id int) error {
	index := -1

	for i, user := range d.db {
		if user.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("user with id = %d not found", id)
	}

	d.db = append(d.db[:index], d.db[index+1:]...)
	return nil
}
