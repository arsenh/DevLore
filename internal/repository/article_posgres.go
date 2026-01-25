package repository

import (
	"context"
	"errors"

	"github.com/arsenh/DevLore/internal/database"
	"github.com/arsenh/DevLore/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrArticleNotFound = errors.New("article not found")

const (
	CommandArticleFindByID    = "SELECT id, title, content, user_id, created_at, updated_at FROM articles WHERE id = $1;"
	CommandArticleFindByTitle = "SELECT id, title, content, user_id, created_at, updated_at FROM articles WHERE title = $1;"
	CommandArticleCreate      = "INSERT INTO articles (title, content, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at;"
	CommandArticleList        = "SELECT * FROM articles;"
	CommandArticleDelete      = "DELETE FROM articles WHERE id = $1;"
	CommandArticleSearch      = "SELECT * FROM articles WHERE title ILIKE '%' || $1 || '%';"
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
	var article model.Article

	err := d.db.QueryRow(ctx, CommandArticleFindByID, id).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.UserId,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrArticleNotFound
	}
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (d *PostgresArticleRepository) FindByTitle(ctx context.Context, title string) (*model.Article, error) {
	var article model.Article

	err := d.db.QueryRow(ctx, CommandArticleFindByTitle, title).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.UserId,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrArticleNotFound
	}
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (d *PostgresArticleRepository) Create(ctx context.Context, a *model.Article) error {
	return d.db.QueryRow(ctx, CommandArticleCreate,
		a.Title,
		a.Content,
		a.UserId,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (d *PostgresArticleRepository) List(ctx context.Context) ([]model.Article, error) {
	rows, err := d.db.QueryRows(ctx, CommandArticleList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []model.Article{}

	for rows.Next() {
		var a model.Article
		if err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Content,
			&a.UserId,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}

func (d *PostgresArticleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := d.db.Exec(ctx, CommandArticleDelete, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (d *PostgresArticleRepository) Search(ctx context.Context, query string) ([]model.Article, error) {
	rows, err := d.db.QueryRows(ctx, CommandArticleSearch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []model.Article{}

	for rows.Next() {
		var a model.Article
		if err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Content,
			&a.UserId,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}
