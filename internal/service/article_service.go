package service

import (
	"context"
	"fmt"
	"time"

	"github.com/arsenh/DevLore/internal/database"
	"github.com/arsenh/DevLore/internal/logger"
	"github.com/arsenh/DevLore/internal/model"
	"github.com/arsenh/DevLore/internal/repository"
	"github.com/arsenh/DevLore/internal/views"
	"github.com/google/uuid"
)

const articleTimeFormat = "Jan 2, 2006 • 15:04"

type ArticleService struct {
	articleRepository model.ArticleRepository
	userRepository    model.UserRepository
}

func NewArticleService(db *database.DbContext) *ArticleService {
	return &ArticleService{
		articleRepository: repository.NewDummyArticleRepository(),
		userRepository:    repository.NewDummyUserRepository(),
	}
}

func (s *ArticleService) GetDashboardData(ctx context.Context) ([]views.ArticleShortItem, error) {
	articles, err := s.articleRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	articleItems := make([]views.ArticleShortItem, 0)

	for _, article := range articles {
		articleItem := views.ArticleShortItem{
			ID:        article.ID.String(),
			Title:     article.Title,
			UpdatedAt: article.UpdatedAt.Local().Format(articleTimeFormat),
		}

		articleItems = append(articleItems, articleItem)
	}

	return articleItems, nil
}

func (s *ArticleService) GetArticleById(ctx context.Context, id uuid.UUID) (*views.ArticleView, error) {
	article, err := s.articleRepository.FindByID(ctx, id)
	if err != nil {
		logger.L().Warnf("the article id = %d not found", id)
		return nil, err
	}

	view := &views.ArticleView{
		BaseView: views.BaseView{
			UserName: "",
		},
		Article: views.ArticleFullViewItem{
			ID:        article.ID.String(),
			Title:     article.Title,
			Content:   article.Content,
			UserId:    article.UserId.String(),
			CreatedAt: article.CreatedAt.Local().Format(articleTimeFormat),
			UpdatedAt: article.UpdatedAt.Local().Format(articleTimeFormat),
		},
		CreatedBy: "",
	}

	return view, nil
}

func (s *ArticleService) SaveArticle(ctx context.Context, title string, content string, userID uuid.UUID) (uuid.UUID, error) {
	if title == "" || content == "" {
		return uuid.Nil, fmt.Errorf("the title and content of article must be not empty")
	}

	id := uuid.New() // generate for new article

	newArticle := &model.Article{
		ID:        id,
		Title:     title,
		Content:   content,
		UserId:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.articleRepository.Create(ctx, newArticle); err != nil {
		return uuid.Nil, logger.LogAndErr("failed to store new article in database")
	}
	return id, nil
}

func (s *ArticleService) DeleteArticleById(ctx context.Context, id uuid.UUID) error {
	if err := s.articleRepository.Delete(ctx, id); err != nil {
		return logger.LogAndErr("failed to delete article from database")
	}
	return nil
}

func (s *ArticleService) EditArticleById(ctx context.Context, id uuid.UUID, title string, content string) (*views.ArticleView, error) {
	article, err := s.articleRepository.FindByID(ctx, id)
	if err != nil {
		return nil, logger.LogAndErr("failed to get article by id = %d", id)
	}

	newArticle := model.Article{
		ID:        article.ID,
		Title:     title,
		Content:   content,
		UserId:    article.UserId,
		CreatedAt: article.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if err := s.DeleteArticleById(ctx, id); err != nil {
		return nil, logger.LogAndErr("failed to delete article by id = %d", id)
	}

	if err := s.articleRepository.Create(ctx, &newArticle); err != nil {
		return nil, logger.LogAndErr("failed to add updated article by id = %d", id)
	}

	view := &views.ArticleView{
		BaseView: views.BaseView{
			UserName: "",
		},
		Article: views.ArticleFullViewItem{
			ID:        newArticle.ID.String(),
			Title:     newArticle.Title,
			Content:   newArticle.Content,
			CreatedAt: newArticle.CreatedAt.Local().Format(articleTimeFormat),
			UpdatedAt: newArticle.UpdatedAt.Local().Format(articleTimeFormat),
		},
	}

	return view, nil
}

func (s *ArticleService) SearchArticlesByQuery(ctx context.Context, query string) (*views.SearchView, error) {
	view := &views.SearchView{
		DashboardView: views.DashboardView{
			BaseView: views.BaseView{
				UserName: "",
			},
		},
		Query: query,
	}

	if query == "" {
		return view, nil
	}

	articles, err := s.articleRepository.Search(ctx, query)
	if err != nil {
		return nil, logger.LogAndErr("failed to search articles with query = %s", query)
	}

	for _, article := range articles {
		articleShortView := views.ArticleShortItem{
			ID:        article.ID.String(),
			Title:     article.Title,
			UpdatedAt: article.UpdatedAt.Local().Format(articleTimeFormat),
		}

		view.Articles = append(view.Articles, articleShortView)
	}

	return view, nil
}
