package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arsenh/DevLore/internal/model"
)

type DummyUserRepository struct {
	db []model.User
}

func NewDummyUserRepository() *DummyUserRepository {
	now := time.Now()

	return &DummyUserRepository{
		db: []model.User{
			{
				ID:        11,
				FullName:  "Alice Johnson",
				Email:     "alice@example.com",
				Password:  "$2a$10$alicehashedpassword",
				CreatedAt: now.Add(-72 * time.Hour),
			},
			{
				ID:        22,
				FullName:  "Bob Smith",
				Email:     "bob@example.com",
				Password:  "$2a$10$bobhashedpassword",
				CreatedAt: now.Add(-48 * time.Hour),
			},
			{
				ID:        33,
				FullName:  "Carol Adams",
				Email:     "carol@example.com",
				Password:  "$2a$10$carolhashedpassword",
				CreatedAt: now.Add(-36 * time.Hour),
			},
			{
				ID:        44,
				FullName:  "Dave Brown",
				Email:     "dave@example.com",
				Password:  "$2a$10$davehashedpassword",
				CreatedAt: now.Add(-24 * time.Hour),
			},
			{
				ID:        55,
				FullName:  "Eve Thompson",
				Email:     "eve@example.com",
				Password:  "$2a$10$evehashedpassword",
				CreatedAt: now.Add(-12 * time.Hour),
			},
			{
				ID:        66,
				FullName:  "Frank Williams",
				Email:     "frank@example.com",
				Password:  "$2a$10$frankhashedpassword",
				CreatedAt: now.Add(-6 * time.Hour),
			},
			{
				ID:        77,
				FullName:  "Grace Miller",
				Email:     "grace@example.com",
				Password:  "$2a$10$gracehashedpassword",
				CreatedAt: now.Add(-3 * time.Hour),
			},
			{
				ID:        88,
				FullName:  "Heidi Clark",
				Email:     "heidi@example.com",
				Password:  "$2a$10$heidihashedpassword",
				CreatedAt: now.Add(-2 * time.Hour),
			},
			{
				ID:        99,
				FullName:  "Ivan Garcia",
				Email:     "ivan@example.com",
				Password:  "$2a$10$ivanhashpassword",
				CreatedAt: now.Add(-1 * time.Hour),
			},
			{
				ID:        111,
				FullName:  "Judy Martinez",
				Email:     "judy@example.com",
				Password:  "$2a$10$judyhashpassword",
				CreatedAt: now,
			},
		},
	}
}

func (d *DummyUserRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	for _, user := range d.db {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with id = %d not found", id)
}

func (d *DummyUserRepository) FindByEmail(ctx context.Context, email string) *model.User {
	for _, user := range d.db {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

func (d *DummyUserRepository) Create(ctx context.Context, u *model.User) error {
	d.db = append(d.db, *u)
	return nil
}

func (d *DummyUserRepository) List(ctx context.Context) ([]model.User, error) {
	return d.db, nil
}

func (d *DummyUserRepository) Delete(ctx context.Context, id int) error {
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
