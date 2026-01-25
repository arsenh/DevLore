package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arsenh/DevLore/internal/model"
	"github.com/google/uuid"
)

type DummyUserRepository struct {
	db []model.User
}

func NewDummyUserRepository() *DummyUserRepository {
	now := time.Now()

	return &DummyUserRepository{
		db: []model.User{
			{
				ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				FullName:  "Alice Johnson",
				Email:     "alice@example.com",
				Password:  "$2a$10$alicehashedpassword",
				CreatedAt: now.Add(-72 * time.Hour),
			},
			{
				ID:        uuid.MustParse("22222222-2222-2222-2222-222222222222"),
				FullName:  "Bob Smith",
				Email:     "bob@example.com",
				Password:  "$2a$10$bobhashedpassword",
				CreatedAt: now.Add(-48 * time.Hour),
			},
			{
				ID:        uuid.MustParse("33333333-3333-3333-3333-333333333333"),
				FullName:  "Carol Adams",
				Email:     "carol@example.com",
				Password:  "$2a$10$carolhashedpassword",
				CreatedAt: now.Add(-36 * time.Hour),
			},
			{
				ID:        uuid.MustParse("44444444-4444-4444-4444-444444444444"),
				FullName:  "Dave Brown",
				Email:     "dave@example.com",
				Password:  "$2a$10$davehashedpassword",
				CreatedAt: now.Add(-24 * time.Hour),
			},
			{
				ID:        uuid.MustParse("55555555-5555-5555-5555-555555555555"),
				FullName:  "Eve Thompson",
				Email:     "eve@example.com",
				Password:  "$2a$10$evehashedpassword",
				CreatedAt: now.Add(-12 * time.Hour),
			},
			{
				ID:        uuid.MustParse("66666666-6666-6666-6666-666666666666"),
				FullName:  "Frank Williams",
				Email:     "frank@example.com",
				Password:  "$2a$10$frankhashedpassword",
				CreatedAt: now.Add(-6 * time.Hour),
			},
			{
				ID:        uuid.MustParse("77777777-7777-7777-7777-777777777777"),
				FullName:  "Grace Miller",
				Email:     "grace@example.com",
				Password:  "$2a$10$gracehashedpassword",
				CreatedAt: now.Add(-3 * time.Hour),
			},
			{
				ID:        uuid.MustParse("88888888-8888-8888-8888-888888888888"),
				FullName:  "Heidi Clark",
				Email:     "heidi@example.com",
				Password:  "$2a$10$heidihashedpassword",
				CreatedAt: now.Add(-2 * time.Hour),
			},
			{
				ID:        uuid.MustParse("99999999-9999-9999-9999-999999999999"),
				FullName:  "Ivan Garcia",
				Email:     "ivan@example.com",
				Password:  "$2a$10$ivanhashpassword",
				CreatedAt: now.Add(-1 * time.Hour),
			},
			{
				ID:        uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
				FullName:  "Judy Martinez",
				Email:     "judy@example.com",
				Password:  "$2a$10$judyhashpassword",
				CreatedAt: now,
			},
		},
	}
}

func (d *DummyUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	for _, user := range d.db {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with id = %d not found", id)
}

func (d *DummyUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, user := range d.db {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}

func (d *DummyUserRepository) Create(ctx context.Context, u *model.User) error {
	d.db = append(d.db, *u)
	return nil
}

func (d *DummyUserRepository) List(ctx context.Context) ([]model.User, error) {
	return d.db, nil
}

func (d *DummyUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
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
