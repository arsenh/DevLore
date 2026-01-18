package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/arsenh/DevLore/internal/model"
	"github.com/google/uuid"
)

type DummyArticleRepository struct {
	db []model.Article
}

func NewDummyArticleRepository() *DummyArticleRepository {
	now := time.Now()
	return &DummyArticleRepository{
		db: []model.Article{
			{
				ID:        uuid.MustParse("0d3f5c6a-1f6a-4b8f-9b57-1b5f9b9a0001"),
				Title:     "Getting Started with Go",
				Content:   "Notes about installing Go, setting GOPATH, and writing the first program.",
				UserId:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				CreatedAt: now.AddDate(0, 0, -10),
				UpdatedAt: now.AddDate(0, 0, -10),
			},
			{
				ID:        uuid.MustParse("0d3f5c6a-1f6a-4b8f-9b57-1b5f9b9a0002"),
				Title:     "Understanding Pointers in Go",
				Content:   "Pointers let you share and modify data without copying. & and * operators explained.",
				UserId:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
				CreatedAt: now.AddDate(0, 0, -9),
				UpdatedAt: now.AddDate(0, 0, -9),
			},
			{
				ID:        uuid.MustParse("0d3f5c6a-1f6a-4b8f-9b57-1b5f9b9a0003"),
				Title:     "HTTP Server Basics",
				Content:   "Minimal net/http example and how handlers work.",
				UserId:    uuid.MustParse("33333333-3333-3333-3333-333333333333"),
				CreatedAt: now.AddDate(0, 0, -8),
				UpdatedAt: now.AddDate(0, 0, -8),
			},
			{
				ID:        uuid.MustParse("0d3f5c6a-1f6a-4b8f-9b57-1b5f9b9a0004"),
				Title:     "Working with Structs",
				Content:   "Structs group related data together and are the building blocks of Go programs.",
				UserId:    uuid.MustParse("44444444-4444-4444-4444-444444444444"),
				CreatedAt: now.AddDate(0, 0, -7),
				UpdatedAt: now.AddDate(0, 0, -7),
			},
			{
				ID:        uuid.MustParse("0d3f5c6a-1f6a-4b8f-9b57-1b5f9b9a0005"),
				Title:     "Interfaces — Duck Typing in Go",
				Content:   "Interfaces describe behavior, not data. Any type that implements the methods satisfies it.",
				UserId:    uuid.MustParse("55555555-5555-5555-5555-555555555555"),
				CreatedAt: now.AddDate(0, 0, -6),
				UpdatedAt: now.AddDate(0, 0, -6),
			},
			{
				ID:        uuid.MustParse("0d3f5c6a-1f6a-4b8f-9b57-1b5f9b9a0006"),
				Title:     "Basic SQL with Go",
				Content:   "Connecting to a database, preparing statements, and scanning rows.",
				UserId:    uuid.MustParse("66666666-6666-6666-6666-666666666666"),
				CreatedAt: now.AddDate(0, 0, -5),
				UpdatedAt: now.AddDate(0, 0, -5),
			},
			{
				ID:        uuid.MustParse("0d3f5c6a-1f6a-4b8f-9b57-1b5f9b9a0007"),
				Title:     "Error Handling Patterns",
				Content:   "Why Go uses explicit errors and common patterns like wrapping and sentinel errors.",
				UserId:    uuid.MustParse("77777777-7777-7777-7777-777777777777"),
				CreatedAt: now.AddDate(0, 0, -4),
				UpdatedAt: now.AddDate(0, 0, -4),
			},
			{
				ID:        uuid.MustParse("0d3f5c6a-1f6a-4b8f-9b57-1b5f9b9a0008"),
				Title:     "Goroutines and Channels",
				Content:   "Concurrency basics: launching goroutines and communicating safely with channels.",
				UserId:    uuid.MustParse("88888888-8888-8888-8888-888888888888"),
				CreatedAt: now.AddDate(0, 0, -3),
				UpdatedAt: now.AddDate(0, 0, -3),
			},
			{
				ID:        uuid.MustParse("0d3f5c6a-1f6a-4b8f-9b57-1b5f9b9a0009"),
				Title:     "Template Rendering in Go",
				Content:   "Using html/template, parsing templates once, and passing data structs.",
				UserId:    uuid.MustParse("99999999-9999-9999-9999-999999999999"),
				CreatedAt: now.AddDate(0, 0, -2),
				UpdatedAt: now.AddDate(0, 0, -2),
			},
			{
				ID:        uuid.MustParse("0d3f5c6a-1f6a-4b8f-9b57-1b5f9b9a0010"),
				Title:     "Project Structure Best Practices",
				Content:   "Separating cmd/, internal/, handlers/, models/, templates/, and static/ folders.",
				UserId:    uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
				CreatedAt: now.AddDate(0, 0, -1),
				UpdatedAt: now.AddDate(0, 0, -1),
			},
		},
	}
}

func (d *DummyArticleRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Article, error) {
	for _, article := range d.db {
		if article.ID == id {
			return &article, nil
		}
	}

	return nil, fmt.Errorf("article with id = %d not found", id)
}

func (d *DummyArticleRepository) FindByTitle(ctx context.Context, title string) (*model.Article, error) {
	for _, article := range d.db {
		if article.Title == title {
			return &article, nil
		}
	}

	return nil, fmt.Errorf("article with title = %s not found", title)
}

func (d *DummyArticleRepository) Create(ctx context.Context, u *model.Article) error {
	d.db = append(d.db, *u)
	return nil
}

func (d *DummyArticleRepository) List(ctx context.Context) ([]model.Article, error) {
	return d.db, nil
}

func (d *DummyArticleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	index := -1

	for i, article := range d.db {
		if article.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("article with id = %d not found", id)
	}

	d.db = append(d.db[:index], d.db[index+1:]...)
	return nil
}

func (d *DummyArticleRepository) Search(ctx context.Context, query string) ([]model.Article, error) {

	var finds []model.Article

	for _, article := range d.db {
		if strings.Contains(strings.ToLower(article.Title), strings.ToLower(query)) {
			finds = append(finds, article)
		}
	}

	return finds, nil
}
