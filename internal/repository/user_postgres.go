package repository

import (
	"context"

	"github.com/arsenh/DevLore/internal/database"
	"github.com/arsenh/DevLore/internal/model"
	"github.com/google/uuid"
)

const ()

type PostgresUserRepository struct {
	db *database.DbContext
}

func NewPostgresUserRepository(db *database.DbContext) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (d *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return nil, nil
}

func (d *PostgresUserRepository) FindByEmail(ctx context.Context, email string) *model.User {
	return nil
}

func (d *PostgresUserRepository) Create(ctx context.Context, u *model.User) error {
	return nil
}

func (d *PostgresUserRepository) List(ctx context.Context) ([]model.User, error) {
	return nil, nil
}

func (d *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
