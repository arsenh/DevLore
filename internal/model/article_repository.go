package model

import "context"

type ArticleRepository interface {
	FindByID(ctx context.Context, id int) (*Article, error)
	FindByTitle(ctx context.Context, email string) (*Article, error)
	Create(ctx context.Context, u *Article) error
	List(ctx context.Context) ([]Article, error)
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, query string) ([]Article, error)
}
