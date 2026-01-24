package repository

import (
	"context"

	"github.com/arsenh/DevLore/internal/database"
	"github.com/arsenh/DevLore/internal/model"
	"github.com/google/uuid"
)

type PostgresArticleRepository struct {
	db *database.DbContext
}

func NewPostgresArticleRepository(db *database.DbContext) *PostgresArticleRepository {
	return &PostgresArticleRepository{
		db: db,
	}
}

func (d *PostgresArticleRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Article, error) {
	return nil, nil
}

func (d *PostgresArticleRepository) FindByTitle(ctx context.Context, title string) (*model.Article, error) {
	return nil, nil
}

func (d *PostgresArticleRepository) Create(ctx context.Context, u *model.Article) error {
	return nil
}

func (d *PostgresArticleRepository) List(ctx context.Context) ([]model.Article, error) {
	return nil, nil
}

func (d *PostgresArticleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (d *PostgresArticleRepository) Search(ctx context.Context, query string) ([]model.Article, error) {
	return nil, nil
}
